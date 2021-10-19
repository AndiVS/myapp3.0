// Package handler contain function for handling request
package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"myapp3.0/internal/model"
	"myapp3.0/internal/service"

	"net/http"
)

// CatHandler struct that contain repository linc
type CatHandler struct {
	Service *service.Service
}

// NewC function for customization handler
func NewC(Service *service.Service) *CatHandler {
	return &CatHandler{Service: Service}
}

// AddC record about cat
func (h *CatHandler) AddC(c echo.Context) error {
	rec := new(model.Record)

	if err := c.Bind(rec); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err := h.Service.AddC(c.Request().Context(), rec)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusCreated, rec.ID)
}

// GetC provides cat
func (h *CatHandler) GetC(c echo.Context) error {
	id := c.Param("_id")
	_id, err1 := uuid.Parse(id)
	if err1 != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad id")
	}
	r, err := h.Service.GetC(c.Request().Context(), _id)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusOK, r)
}

// GetAllC provides all cats
func (h *CatHandler) GetAllC(c echo.Context) error {
	var rec []*model.Record

	rec, err := h.Service.GetAllC(c.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusOK, rec)
}

// UpdateC updating record about cat
func (h *CatHandler) UpdateC(c echo.Context) error {
	rec := new(model.Record)

	if err := c.Bind(rec); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err := h.Service.UpdateC(c.Request().Context(), rec)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusOK, "completed successfully")
}

// DeleteC record about cat
func (h *CatHandler) DeleteC(c echo.Context) error {
	id := c.Param("_id")
	_id, err1 := uuid.Parse(id)
	if err1 != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad id ")
	}

	err := h.Service.DeleteC(c.Request().Context(), _id)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusOK, "completed successfully")
}
