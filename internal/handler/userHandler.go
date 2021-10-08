// Package handler contain function for handling request
package handler

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"myapp3.0/internal/model"
	"myapp3.0/internal/service"

	"net/http"
)

// UserHandler struct that contain repository linc
type UserHandler struct {
	Service *service.Service
}

// AddU record about cat
func (h *UserHandler) AddU(c echo.Context) error {
	rec := new(model.User)

	if err := c.Bind(rec); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err := h.Service.AddU(c.Request().Context(), rec)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusCreated, rec.Username)
}

// GetU provides cat
func (h *UserHandler) GetU(c echo.Context) error {
	id := c.Param("")

	r, err := h.Service.GetU(c.Request().Context(), id)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusOK, r)
}

// GetAllU provides all cats
func (h *UserHandler) GetAllU(c echo.Context) error {
	var rec []*model.User

	rec, err := h.Service.GetAllU(c.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusOK, rec)
}

// UpdateU updating record about cat
func (h *UserHandler) UpdateU(c echo.Context) error {
	rec := new(model.User)

	if err := c.Bind(rec); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err := h.Service.UpdateU(c.Request().Context(), rec.Username, rec.IsAdmin)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusOK, "completed successfully")
}

// DeleteU record about cat
func (h *UserHandler) DeleteU(c echo.Context) error {
	username := c.Param("username")

	err := h.Service.DeleteU(c.Request().Context(), username)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusOK, "completed successfully")
}
