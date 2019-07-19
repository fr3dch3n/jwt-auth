# jwt-auth
> Authenticate requests with jwt, extensible claims and a specific jwks-set.

[![Build Status][travis-image]][travis-url]

This library provides a way to decode a jwt-token with self-defined claims.
Additionally one can specify a function to retrieve the JWKS.
See [here](https://auth0.com/docs/jwks) for further information about JWKS.

## Installation

* with go-get: `go get github.com/fr3dch3n/jwt-auth`
* with dep: `dep ensure --add github.com/fr3dch3n/jwt-auth`


## Usage example

**Initialization**

First the JWKS has to be initialized via one of the following ways:

_Initialize jwk-set from json in AWS-SSM_
```go
jwt.NewAuth("/path/in/ssm", auth.FetchJwksConfigurationFromSSM)
```

_Initialize jwk-set from local json-file_
```go
jwt.NewAuth("file/to/local/jwks.json", auth.FetchJwksConfigurationFromFS)
```

**Authorizing requests**

Then the jwt-component can be used to authenticate requests.
To use specific claims, just extend the jwtgo.StandardClaims.

```go
package handler

import (
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/fr3dch3n/jwt-auth"
	"net/http"
)

func MyAuthenticatedHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	token, err := jwt.DecodeToken(authHeader, &jwtgo.StandardClaims{})
	if err != nil {
	    fmt.Println(err)
	}
	fmt.Println(token)
}
```

## Release History

* 0.0.1
    * initial release

## Meta

[@fr3dch3n](https://twitter.com/fr3dch3n)

Distributed under the Apache 2.0 license. See ``LICENSE`` for more information.

## Contributing

1. Fork it (<https://github.com/fr3dch3n/jwt-auth/fork>)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request

<!-- Markdown link & img dfn's -->
[travis-image]: https://img.shields.io/travis/fr3dch3n/jwt-auth/master.svg?style=flat-square
[travis-url]: https://travis-ci.org/fr3dch3n/jwt-auth
