package middleware

import (
	"knowledge-capsule-api/utils"
	"log"
	"net/http"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				utils.ErrorResponse(w, http.StatusInternalServerError, err.(error))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
