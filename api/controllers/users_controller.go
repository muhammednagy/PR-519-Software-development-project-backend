package controllers

import (
	"errors"
	"fmt"
	"github.com/gogearbox/gearbox"
	"net/http"
	"strconv"

	"github.com/muhammednagy/PR-519-Software-development-project-backend/api/auth"
	"github.com/muhammednagy/PR-519-Software-development-project-backend/api/models"
	"github.com/muhammednagy/PR-519-Software-development-project-backend/api/responses"
	"github.com/muhammednagy/PR-519-Software-development-project-backend/api/utils/formaterror"
)

func (server *Server) CreateUser(ctx gearbox.Context) {
	ctx.Set("Access-Control-Allow-Origin", "*")
	user := models.User{}
	err := ctx.ParseBody(&user)
	if err != nil {
		responses.ERROR(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	user.Prepare()
	err = user.Validate("")
	if err != nil {
		responses.ERROR(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	userCreated, err := user.SaveUser(server.DB)

	if err != nil {
		responses.ERROR(ctx, http.StatusInternalServerError, formaterror.FormatError(err.Error()))
		return
	}
	token, err := auth.CreateToken(userCreated.ID)
	if err != nil {
		responses.ERROR(ctx, http.StatusInternalServerError, formaterror.FormatError(err.Error()))
		return
	}
	responses.JSON(ctx, http.StatusCreated, token)
}

func (server *Server) UpdateUser(ctx gearbox.Context) {
	ctx.Set("Access-Control-Allow-Origin", "*")
	uid, err := strconv.ParseUint(ctx.Query("id"), 10, 32)
	if err != nil {
		responses.ERROR(ctx, http.StatusBadRequest, err)
		return
	}
	user := models.User{}
	err = ctx.ParseBody(&user)
	if err != nil {
		responses.ERROR(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(ctx)
	if err != nil {
		responses.ERROR(ctx, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	if uint64(tokenID) != uid {
		responses.ERROR(ctx, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		responses.ERROR(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	updatedUser, err := user.UpdateAUser(server.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(ctx, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(ctx, http.StatusOK, updatedUser)
}



func (server *Server) DeleteUser(ctx gearbox.Context) {
	ctx.Set("Access-Control-Allow-Origin", "*")
	user := models.User{}

	uid, err := strconv.ParseUint(ctx.Query("id"), 10, 32)
	if err != nil {
		responses.ERROR(ctx, http.StatusBadRequest, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(ctx)
	if err != nil {
		responses.ERROR(ctx, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	if tokenID != 0 && uint64(tokenID) != uid {
		responses.ERROR(ctx, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(ctx, http.StatusNoContent, "")
}
