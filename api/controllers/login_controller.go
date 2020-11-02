package controllers

import (
	"github.com/gogearbox/gearbox"
	"net/http"

	"github.com/muhammednagy/PR-519-Software-development-project-backend/api/auth"
	"github.com/muhammednagy/PR-519-Software-development-project-backend/api/models"
	"github.com/muhammednagy/PR-519-Software-development-project-backend/api/responses"
	"github.com/muhammednagy/PR-519-Software-development-project-backend/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) LoginOptions(ctx gearbox.Context) {
	ctx.Set("Access-Control-Allow-Origin", "*")
	ctx.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func (server *Server) Login(ctx gearbox.Context) {
	ctx.Set("Access-Control-Allow-Origin", "*")
	user := models.User{}
	err := ctx.ParseBody(&user)
	if err != nil {
		responses.ERROR(ctx, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(ctx, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(ctx, http.StatusOK, token)
}

func (server *Server) SignIn(email, password string) (string, error) {

	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}
