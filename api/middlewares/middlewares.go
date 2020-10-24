package middlewares

import (
	"errors"
	"github.com/gogearbox/gearbox"
	"net/http"

	"github.com/muhammednagy/PR-519-Software-development-project-backend/api/auth"
	"github.com/muhammednagy/PR-519-Software-development-project-backend/api/responses"
)

func SetMiddlewareAuthentication(ctx gearbox.Context) {
	err := auth.TokenValid(ctx)
	if err != nil {
		responses.ERROR(ctx, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	ctx.Next()
}
