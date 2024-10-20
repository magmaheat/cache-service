package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/magmaheat/cache-service/intarnal/entity"
	"github.com/magmaheat/cache-service/intarnal/service"
	"net/http"
)

type cacheRouter struct {
	cacheRouter service.Cache
}

func NewFilesRouter(handler *echo.Group, services *service.Services) {
	f := &cacheRouter{
		cacheRouter: services.Cache,
	}

	handler.POST("/docs", f.createDocument)
	handler.GET("/docs/:id", f.getDocument)
	handler.HEAD("/docs/:id", f.getDocument)
	handler.GET("/docs", f.getDocuments)
	handler.HEAD("/docs", f.getDocuments)
}

func (f *cacheRouter) createDocument(c echo.Context) error {
	var meta entity.Meta

	jsonMeta := c.FormValue("meta")
	if err := json.Unmarshal([]byte(jsonMeta), &meta); jsonMeta == "" || err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid body request")
		return fmt.Errorf("meta json invalid")
	}

	jsonField := c.FormValue("json")

	file, err := c.FormFile("file")
	if err != nil && jsonField == "" {
		newErrorResponse(c, http.StatusBadRequest, "missing file in body request")
		return err
	}

	err = f.cacheRouter.SaveData(c.Request().Context(), meta, jsonField, file)
	if err != nil {
		if errors.Is(err, service.ErrFileAlreadyExists) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}

		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Json string `json:"json"` //TODO transform?
		File string `json:"file"`
	}

	return c.JSON(http.StatusOK, Response{
		Data: response{
			Json: jsonField,
			File: meta.Name,
		},
	})
}

func (f *cacheRouter) getDocument(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		newErrorResponse(c, http.StatusBadRequest, "param id not found")
		return fmt.Errorf("param id not found")
	}

	document, err := f.cacheRouter.GetDocument(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrFileNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return err
		}

		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	if document.JsonBody == "" {
		c.Response().Header().Set("Content-Type", document.Mime)
		c.Response().Header().Set("Content-Disposition", "attachment; filename="+document.Name)
	}

	if c.Request().Method == http.MethodGet {
		if document.JsonBody == "" {
			return c.Blob(http.StatusOK, document.Mime, document.Body)
		}

		type response struct {
			Json string `json:"json"`
		}
		c.JSON(http.StatusOK, Response{
			Data: response{
				Json: document.JsonBody,
			},
		})
	}

	return nil
}

func (f *cacheRouter) getDocuments(c echo.Context) error {
	return nil
}
