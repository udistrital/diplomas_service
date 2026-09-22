package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/udistrital/diplomas_mid/services"
)

type DiplomaPreviewController struct {
	APIController
	Service services.DiplomaPreviewService
}

func (c *DiplomaPreviewController) URLMapping() {}

func (c *DiplomaPreviewController) Generar() {
	documentoID, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.writeError(http.StatusBadRequest, errors.New("invalid documento id"))
		return
	}

	result, err := c.Service.Generar(c.Ctx.Request.Context(), documentoID)
	if err != nil {
		c.writeError(mapServiceError(err), err)
		return
	}

	c.writeJSON(http.StatusOK, result)
}
