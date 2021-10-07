// Package handler contain function for handling request
package handler

import (
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

// New function for customization handler
func New(Service *service.Service) CatHandler {
	return CatHandler{Service: Service}
}

// Add record about cat
func (h *CatHandler) Add(c echo.Context) error {
	rec := new(model.Record)

	if err := c.Bind(rec); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err := h.Service.Add(c.Request().Context(), rec)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusCreated, rec.ID)
}

// Get provides cat
func (h *CatHandler) Get(c echo.Context) error {
	id := c.Param("id")

	r, err := h.Service.Get(c.Request().Context(), id)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusOK, r)
}

// GetAll provides all cats
func (h *CatHandler) GetAll(c echo.Context) error {
	var rec []*model.Record

	rec, err := h.Service.GetAll(c.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusOK, rec)
}

// Update updating record about cat
func (h *CatHandler) Update(c echo.Context) error {
	rec := new(model.Record)

	if err := c.Bind(rec); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err := h.Service.Update(c.Request().Context(), rec)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusOK, "completed successfully")
}

// Delete record about cat
func (h *CatHandler) Delete(c echo.Context) error {
	id := c.Param("id")

	err := h.Service.Delete(c.Request().Context(), id)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusOK, "completed successfully")
}
