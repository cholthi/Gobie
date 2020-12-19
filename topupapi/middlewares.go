package main

import (
	"log"
	"net/http"
	"regexp"
)

var bearerRegexp = regexp.MustCompile(`^(?:B|b)earer (\S+$)`)

type Middleware func(http.Handler) http.Handler

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
			/*body, err := ioutil.ReadAll(r.Body)// Got ya mutherfucker. io.Reader is only readable once.
			if err != nil {
				logger.Print(err)
			}*/
			logger.Printf("%s %s", r.Method, r.RemoteAddr)
			h.ServeHTTP(res, r)
		})

	}
}

func requireAthenticationMiddleware(logger log.Logger) Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
			bearer, err := extractBearerToken(r)
			if err != nil {
				errorResponse(res, http.StatusUnauthorized, err)
				logger.Println(err)
				return
			}
			token, err := parseJTWToken(bearer)

			if err != nil {
				errorResponse(res, http.StatusUnauthorized, err)
				logger.Println(err)
				return
			}
			ctx, err := addUserToContext(r.Context(), token)
			if err != nil {
				errorResponse(res, http.StatusUnauthorized, err)
				logger.Println(err)
				return
			}

			r = r.WithContext(ctx)
			handler.ServeHTTP(res, r)
			return

		})
	}
}

func AuthMiddleware(logger log.Logger) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {

			errRes := TopUpResponse{
				StatusCode:    20,
				StatusMessage: "AUTHENTICATION_REQUIRED",
			}
			apikey := r.Header.Get("api-key")

			if apikey == "" {
				h.ServeHTTP(res, r)
				return
			}
			body, err := encodeTopupResponse(errRes)
			if err != nil {
				logger.Print(err)
			}
			res.Header().Set("Content-Type", "application/json")
			_, err = res.Write(body)
			if err != nil {
				logger.Print(err)
			}
			return
		})
	}
}
