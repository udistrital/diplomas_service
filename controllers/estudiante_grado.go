package controllers

import (
	"net/http"

	"github.com/udistrital/diplomas_mid/services"
)

type EstudianteGradoController struct {
	APIController
	Service services.EstudianteGradoService
}

func (c *EstudianteGradoController) URLMapping() {}

func (c *EstudianteGradoController) AprobadosPorFacultad() {
	result, err := c.Service.ListarAprobadosPorFacultad(c.Ctx.Request.Context())
	if err != nil {
		c.writeError(mapServiceError(err), err)
		return
	}

	c.writeJSON(http.StatusOK, result)
}

func (c *EstudianteGradoController) CrearDocumentosAprobados() {
	payload, err := decodeBody[services.CrearDocumentosAprobadosInput](&c.Controller)
	if err != nil {
		c.writeError(http.StatusBadRequest, err)
		return
	}

	result, err := c.Service.CrearDocumentosAprobados(c.Ctx.Request.Context(), payload)
	if err != nil {
		c.writeError(mapServiceError(err), err)
		return
	}

	status := http.StatusCreated
	if result.DryRun {
		status = http.StatusOK
	}
	if result.Errores > 0 {
		status = http.StatusMultiStatus
	}
	c.writeJSON(status, result)
}
