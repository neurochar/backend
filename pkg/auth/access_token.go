package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

func ParseAccessToken(tokenStr string, validate bool, secret []byte) (*UserTenantSessionAccessClaims, error) {
	ops := []jwt.ParserOption{
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	}

	if !validate {
		ops = append(ops, jwt.WithoutClaimsValidation())
	}

	token, err := jwt.ParseWithClaims(tokenStr, &UserTenantSessionAccessClaims{}, func(token *jwt.Token) (any, error) {
		return secret, nil
	}, ops...)
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*UserTenantSessionAccessClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func IssueAccessJWT(access *UserTenantSessionAccessClaims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, access)

	return token.SignedString(secret)
}
