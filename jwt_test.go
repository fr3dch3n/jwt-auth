package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestClaims struct {
	jwt.StandardClaims
}

func TestDecodeToken_noAuthHeader(t *testing.T) {
	// given
	token := ""

	// when
	_, err := DecodeToken(token, &TestClaims{})

	// then
	assert.Error(t, err, "Error getting token from authorization header")
}
