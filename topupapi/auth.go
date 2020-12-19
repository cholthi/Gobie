package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/cholthi/topupapi/model"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/uuid"
	auth "github.com/netlify/gotrue/api"
	"github.com/netlify/gotrue/models"
)

var ConfigFile string = "../conf.env"

type contextkey string

const userContextKey = contextkey("user")

func getAuthHandler() (http.Handler, error) {
	api, _, err := auth.NewAPIFromConfigFile(ConfigFile, "Unknown version")
	if err != nil {
		return nil, err
	}
	return api.GetHandler(), nil
}

func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("The Token is required")
	}
	matches := bearerRegexp.FindStringSubmatch(authHeader)
	if len(matches) != 2 {
		return "", errors.New("Invalid Token format")
	}

	return matches[1], nil
}

func parseJTWToken(bearer string) (*jwt.Token, error) {
	p := jwt.Parser{ValidMethods: []string{jwt.SigningMethodHS256.Name}}
	token, err := p.ParseWithClaims(bearer, &auth.GoTrueClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("djkdhjehfejhr973"), nil
	})

	if err != nil {
		return nil, err
	}
	return token, err
}

func addUserToContext(ctx context.Context, token *jwt.Token) (context.Context, error) {
	claims := token.Claims.(*auth.GoTrueClaims)
	userid, err := uuid.FromString(claims.Subject)
	if err != nil {
		//panic(err)
		return nil, err
	}
	user, err := model.FindUserByID(userid)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, userContextKey, user)
	return ctx, nil
}

func getUser(ctx context.Context) *models.User {
	val := ctx.Value(userContextKey)
	if val == nil {
		return nil
	}
	return val.(*models.User)
}
