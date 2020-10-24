package controllers

import (
	"github.com/gogearbox/gearbox"
	"net/http"

	"github.com/muhammednagy/PR-519-Software-development-project-backend/api/responses"
)

func (server *Server) Home(ctx gearbox.Context) {
	responses.JSON(ctx, http.StatusOK, "Welcome To This Awesome API")

}
