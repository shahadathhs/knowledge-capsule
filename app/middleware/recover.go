package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"knowledge-capsule/pkg/logger"
	"knowledge-capsule/pkg/utils"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				var err error
				switch e := rec.(type) {
				case error:
					err = e
				default:
					err = errors.New(fmt.Sprint(e))
				}
				logger.ErrorRequest(r, logger.EventPanic, err)
				utils.ErrorResponse(w, r, http.StatusInternalServerError, err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
