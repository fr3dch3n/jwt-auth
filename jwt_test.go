package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

type TestClaims struct {
	jwt.StandardClaims
}

var jwksConfigPath string

func init()  {
	logrus.SetLevel(logrus.FatalLevel)
	jwksConfigPath = path.Join("test-resources", "mock-jwks.json")
}

func TestDecodeToken_(t *testing.T) {
	t.Run("no valid auth header", func(t *testing.T) {
		// given
		token := ""

		// when
		_, err := DecodeToken(token, &TestClaims{})

		// then
		assert.Error(t, err, "error getting token from authorization header")
	})
	t.Run("token decode fails", func(t *testing.T) {
		// given
		token := givenAuthToken(givenToken("some-key"))  + ".123"

		// when
		_, err := DecodeToken(token, &TestClaims{})

		// then
		assert.Error(t, err)
	})
}

func TestGetKey(t *testing.T) {
	t.Run("no kid in token", func(t *testing.T) {
		// given
		token := jwt.Token{}

		// when
		_, err := getKey(&token)

		// then
		assert.Error(t, err)
	})
	t.Run("jwks key lookup broken", func(t *testing.T) {
		// given
		token := givenToken("invalid-key")
		NewAuth(FetchJwksConfigurationFromFS, jwksConfigPath, 0)

		// when
		_, err := getKey(token)
		defer StopReloadingJWKS()

		// then
		assert.Error(t, err)
	})
	t.Run("jwks key is fine", func(t *testing.T) {
		// given
		token := givenToken("some-key")
		NewAuth(FetchJwksConfigurationFromFS, jwksConfigPath, 0)

		// when
		res, err := getKey(token)
		defer StopReloadingJWKS()

		// then
		assert.NotNil(t, res)
		assert.NoError(t, err)
	})
}

func givenToken(keyId string) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		Id:        "some-id",
		IssuedAt:  time.Now().Unix()-2000,
		ExpiresAt: time.Now().Unix()-1000,
	})
	token.Header["kid"] = keyId
	return token
}

func givenAuthToken(token *jwt.Token) string {
	bs, _ := os.Open(path.Join("test-resources", "mockPrivKey.priv"))
	bytes, _ := ioutil.ReadAll(bs)
	var privateKey, _ = jwt.ParseRSAPrivateKeyFromPEM(bytes)
	tokenString, _ := token.SignedString(privateKey)
	return "Bearer " + tokenString
}
