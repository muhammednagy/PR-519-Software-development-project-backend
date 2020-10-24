package responses

import (
	"fmt"
	"github.com/gogearbox/gearbox"
	"net/http"
)

func JSON(ctx gearbox.Context, statusCode int, data interface{}) {
	ctx.Status(statusCode)
	err := ctx.SendJSON(data)
	if err != nil {
		fmt.Fprintf(ctx.Context().Response.BodyWriter(), "%s", err.Error())
	}
}

func ERROR(ctx gearbox.Context, statusCode int, err error) {
	if err != nil {
		JSON(ctx, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(ctx, http.StatusBadRequest, nil)
}
