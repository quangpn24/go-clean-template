package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"go-clean-template/pkg/apperror"
	"go-clean-template/pkg/config"
	"go-clean-template/pkg/constant"

	"github.com/go-jose/go-jose/v3"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

type Authentication struct {
	SkipPaths         []string
	CognitoUserPoolId string
	jcache            *jwkCache
	cfg               *config.Config
}
type jwkCache struct {
	jwks      *jose.JSONWebKeySet
	timestamp time.Time
}

func NewAuthentication(userPoolId string, skipPaths []string, cfg *config.Config) *Authentication {
	return &Authentication{
		CognitoUserPoolId: userPoolId,
		SkipPaths:         skipPaths,
		jcache:            &jwkCache{},
		cfg:               cfg,
	}
}

func (a *Authentication) Middleware() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper:    a.Skipper,
		Validator:  a.ValidateAccessToken,
		AuthScheme: "Bearer",
		KeyLookup:  "header:Authorization",
	})
}

func (a *Authentication) Skipper(c echo.Context) bool {
	for _, p := range a.SkipPaths {
		if strings.HasPrefix(c.Path(), p) {
			return true
		}
	}
	return false
}

type Claims struct {
	Sub      string `json:"sub"`
	ClientId string `json:"client_id"`
	UserName string `json:"username"`
	TokenUse string `json:"token_use"`
	jwt.StandardClaims
}

func (a *Authentication) ValidateAccessToken(token string, c echo.Context) (bool, error) {
	var claims Claims

	_, err := jwt.ParseWithClaims(token, &claims, a.lookupKey)
	if err != nil {
		return false, apperror.ErrUnauthorized(err)
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return false, apperror.ErrUnauthorized(errors.New("Token expired"))
	}

	if claims.Issuer != a.cfg.CognitoIssuer {
		return false, apperror.ErrUnauthorized(errors.New("Invalid issuer"))
	}

	if claims.TokenUse != "access" {
		return false, apperror.ErrUnauthorized(errors.New("No access"))
	}

	c.Set(constant.UserIDKey, claims.Sub)
	return true, nil
}

func (a *Authentication) lookupKey(token *jwt.Token) (interface{}, error) {
	kid := token.Header["kid"].(string)

	//get JWKS(Json Web Key Sets)
	jwks, err := a.getJWKS()
	if err != nil {
		return nil, err
	}

	var jwk jose.JSONWebKey
	for _, key := range jwks.Keys {
		if key.KeyID == kid {
			jwk = key
		}
	}
	return jwk.Key, nil
}

func (a *Authentication) getJWKS() (*jose.JSONWebKeySet, error) {
	if a.jcache.jwks != nil && time.Since(a.jcache.timestamp).Minutes() < 10.0 {
		return a.jcache.jwks, nil
	}

	jwks := &jose.JSONWebKeySet{}
	resp, err := http.Get(a.cfg.CognitoURLGetJWKS)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&jwks); err != nil {
		return nil, err
	}
	a.jcache.jwks = jwks
	a.jcache.timestamp = time.Now()
	return jwks, nil
}
