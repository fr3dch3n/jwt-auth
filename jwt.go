package jwt

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat/go-jwx/jwk"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

var jwksSet *jwk.Set = nil

func NewAuth(ssmPath string, jwksFetcher func(string) (*jwk.Set, error)) {
	log.Info("Initializing jwt-auth")
	newSet, err := jwksFetcher(ssmPath)
	if err != nil {
		log.Error(err)
	}
	jwksSet = newSet
}

func FetchJwksConfigurationFromSSM(ssmPath string) (*jwk.Set, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String("eu-central-1")},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}

	ssmsvc := ssm.New(sess)
	withDecryption := true
	param, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           &ssmPath,
		WithDecryption: &withDecryption,
	})
	if err != nil {
		return nil, err
	}

	value := *param.Parameter.Value
	return jwk.Parse([]byte(value))
}

func FetchJwksConfigurationFromFS(jwksURL string) (*jwk.Set, error) {
	jsonFile, err := os.Open(jwksURL)
	if err != nil {
		log.Error(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	return jwk.Parse(byteValue)
}

func getKey(token *jwt.Token) (interface{}, error) {
	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("expecting JWT header to have string kid")
	}

	if key := jwksSet.LookupKeyID(keyID); len(key) == 1 {
		return key[0].Materialize()
	}

	return nil, errors.New("unable to find key")
}

func DecodeToken(bearerToken string, claims jwt.Claims) (*jwt.Token, error) {
	extractedToken := strings.Split(bearerToken, "Bearer ")
	if len(extractedToken) != 2 {
		return nil, errors.New("error getting token from authorization header")
	}
	tokenString := extractedToken[1]

	token, err := jwt.ParseWithClaims(tokenString, claims, getKey)
	if err != nil {
		log.Error("Error decoding token: ", err)
		return nil, err
	}

	return token, nil
}
