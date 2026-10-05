package controllers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/udistrital/diplomas_mid/services"
)

type CeremoniaController struct {
	APIController
	Service services.CeremoniaService
}

func (c *CeremoniaController) URLMapping() {}

func (c *CeremoniaController) requestContext() context.Context {
	ctx := c.Ctx.Request.Context()
	if token := c.Ctx.Request.Header.Get("Authorization"); token != "" {
		ctx = context.WithValue(ctx, "Authorization", token)
	}
	return ctx
}

func (c *CeremoniaController) Listar() {
	documento, _ := strconv.ParseInt(c.GetString("documento_identidad_secretario"), 10, 64)
	result, err := c.Service.Listar(c.requestContext(), &services.ListarCeremoniasInput{DocumentoIdentidadSecretario: documento})
	if err != nil {
		c.writeError(mapServiceError(err), err)
		return
	}
	c.writeJSON(http.StatusOK, result)
}

func (c *CeremoniaController) Actualizar() {
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil || id <= 0 {
		c.writeError(http.StatusBadRequest, services.ErrInvalidInput)
		return
	}
	payload, err := decodeBody[services.ActualizarCeremoniaInput](&c.Controller)
	if err != nil {
		c.writeError(http.StatusBadRequest, err)
		return
	}
	result, err := c.Service.Actualizar(c.requestContext(), id, payload)
	if err != nil {
		c.writeError(mapServiceError(err), err)
		return
	}
	c.writeJSON(http.StatusOK, result)
}

func (c *CeremoniaController) Eliminar() {
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil || id <= 0 {
		c.writeError(http.StatusBadRequest, services.ErrInvalidInput)
		return
	}
	documento, _ := strconv.ParseInt(c.GetString("documento_identidad_secretario"), 10, 64)
	err = c.Service.Eliminar(c.requestContext(), id, &services.EliminarCeremoniaInput{DocumentoIdentidadSecretario: documento})
	if err != nil {
		c.writeError(mapServiceError(err), err)
		return
	}
	c.writeJSON(http.StatusOK, map[string]bool{"deleted": true})
}

func (c *CeremoniaController) Crear() {
	payload, err := decodeBody[services.CrearCeremoniaInput](&c.Controller)
	if err != nil {
		c.writeError(http.StatusBadRequest, err)
		return
	}
	result, err := c.Service.Crear(c.requestContext(), payload)
	if err != nil {
		c.writeError(mapServiceError(err), err)
		return
	}
	c.writeJSON(http.StatusCreated, result)
}

func (c *CeremoniaController) Asignar() {
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil || id <= 0 {
		c.writeError(http.StatusBadRequest, services.ErrInvalidInput)
		return
	}
	payload, err := decodeBody[services.AsignarEstudiantesCeremoniaInput](&c.Controller)
	if err != nil {
		c.writeError(http.StatusBadRequest, err)
		return
	}
	result, err := c.Service.AsignarEstudiantes(c.requestContext(), id, payload)
	if err != nil {
		c.writeError(mapServiceError(err), err)
		return
	}
	c.writeJSON(http.StatusOK, result)
}

func (c *CeremoniaController) QuitarPrograma() {
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil || id <= 0 {
		c.writeError(http.StatusBadRequest, services.ErrInvalidInput)
		return
	}
	payload, err := decodeBody[services.QuitarProgramaCeremoniaInput](&c.Controller)
	if err != nil {
		c.writeError(http.StatusBadRequest, err)
		return
	}
	result, err := c.Service.QuitarPrograma(c.requestContext(), id, payload)
	if err != nil {
		c.writeError(mapServiceError(err), err)
		return
	}
	c.writeJSON(http.StatusOK, result)
}

func (c *CeremoniaController) EstudiantesPrograma() {
	programaID, err := strconv.ParseInt(c.Ctx.Input.Param(":programa_id"), 10, 64)
	if err != nil || programaID <= 0 {
		c.writeError(http.StatusBadRequest, services.ErrInvalidInput)
		return
	}
	documento, _ := strconv.ParseInt(c.GetString("documento_identidad_secretario"), 10, 64)
	result, err := c.Service.EstudiantesPrograma(c.requestContext(), &services.EstudiantesProgramaCeremoniaInput{DocumentoIdentidadSecretario: documento, ProgramaAcademicoID: programaID})
	if err != nil {
		c.writeError(mapServiceError(err), err)
		return
	}
	c.writeJSON(http.StatusOK, result)
}
