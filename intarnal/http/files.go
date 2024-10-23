package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/magmaheat/cache-service/intarnal/entity"
	"github.com/magmaheat/cache-service/intarnal/service"
	log "github.com/sirupsen/logrus"
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
	handler.DELETE("/docs/:id", f.deleteDocument)
}

// @Summary Create document
// @Description Create document
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} http.Response
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /docs [post]
func (f *cacheRouter) createDocument(c echo.Context) error {
	var meta entity.Meta

	jsonMeta := c.FormValue("meta")
	if err := json.Unmarshal([]byte(jsonMeta), &meta); jsonMeta == "" || err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid body request")
		return fmt.Errorf("meta json invalid")
	}

	jsonStr := c.FormValue("json")
	var jsonField map[string]interface{}
	if jsonStr != "" {
		if err := json.Unmarshal([]byte(jsonStr), &jsonField); err != nil {
			newErrorResponse(c, http.StatusBadRequest, "invalid json field")
			return err

		}
	}

	file, err := c.FormFile("file")
	if err != nil && jsonStr == "" {
		newErrorResponse(c, http.StatusBadRequest, "missing file in body request")
		return err
	}

	err = f.cacheRouter.SaveData(c.Request().Context(), meta, jsonStr, file)
	if err != nil {
		if errors.Is(err, service.ErrFileAlreadyExists) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}

		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Json map[string]interface{} `json:"json"`
		File string                 `json:"file"`
	}

	return c.JSON(http.StatusOK, Response{
		Data: response{
			Json: jsonField,
			File: meta.Name,
		},
	})
}

// @Summary Get document
// @Description Get document. The HEAD method is also supported and returns only the headers.
// @Tags files
// @Produce json
// @Param id path string true "ID to get document"
// @Success 200 {object} http.Response
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /docs/{id} [get]
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

// @Summary Get documents
// @Description Get documents. The HEAD method is also supported and returns only the headers.
// @Tags files
// @Accept json
// @Produce json
// @Success 200 {object} http.Response
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /docs [get]
func (f *cacheRouter) getDocuments(c echo.Context) error {
	var input entity.SearchDocuments

	if err := c.Bind(&input); err != nil {
		log.Errorf("http - files - c.Bind: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "invalid body request")
		return err
	}

	if input.Login == "" {
		input.Login = c.Get("login").(string)
	}

	metaList, err := f.cacheRouter.GetDocuments(c.Request().Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrFileNotFound) {
			newErrorResponse(c, http.StatusNotFound, "files not found")
			return err
		}

		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Docs entity.MetaSlice `json:"docs"`
	}

	return c.JSON(http.StatusOK, Response{
		Data: response{
			Docs: metaList,
		},
	})
}

// @Summary Delete document
// @Description Delete document
// @Tags files
// @Produce json
// @Param id path string true "ID to delete document"
// @Success 200 {object} http.Response
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /docs/{id} [delete]
func (f *cacheRouter) deleteDocument(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		log.Errorf("http - files - c.Param: param id is empty")
		newErrorResponse(c, http.StatusBadRequest, "param id is empty")
		return fmt.Errorf("param id is empty")
	}

	err := f.cacheRouter.DeleteDocument(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrFileNotFound) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}

		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, Response{
		Response: map[string]bool{
			id: true,
		},
	})
}
