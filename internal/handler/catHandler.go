// Package handler contain function for handling request
package handler

import (
	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/AndiVS/myapp3.0/internal/repository"
	"github.com/AndiVS/myapp3.0/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"net/http"
)

// CatHandler struct that contain repository linc
type CatHandler struct {
	Service service.Cats
}

// NewHandlerCat function for customization handler
func NewHandlerCat(Service service.Cats) *CatHandler {
	return &CatHandler{Service: Service}
}

// AddCat record about cat
func (h *CatHandler) AddCat(c echo.Context) error {
	cat := new(model.Cat)

	if err := c.Bind(cat); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	cat.ID = uuid.New()

	id, err := h.Service.AddCat(c.Request().Context(), cat)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, id)
}

// GetCat provides cat
func (h *CatHandler) GetCat(c echo.Context) error {
	id := c.Param("_id")
	_id, err1 := uuid.Parse(id)
	if err1 != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	r, err := h.Service.GetCat(c.Request().Context(), _id)
	if err != nil {
		if err.Error() == repository.ErrNotFound.Error() {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, r)
}

// GetAllCat provides all cats
func (h *CatHandler) GetAllCat(c echo.Context) error {
	var cat []*model.Cat

	cat, err := h.Service.GetAllCat(c.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, cat)
}

// UpdateCat updating record about cat
func (h *CatHandler) UpdateCat(c echo.Context) error {
	cat := new(model.Cat)

	if err := c.Bind(cat); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err := h.Service.UpdateCat(c.Request().Context(), cat)
	if err != nil {
		if err.Error() == repository.ErrNotFound.Error() {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// DeleteCat record about cat
func (h *CatHandler) DeleteCat(c echo.Context) error {
	id := c.Param("_id")
	_id, err1 := uuid.Parse(id)
	if err1 != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err := h.Service.DeleteCat(c.Request().Context(), _id)
	if err != nil {
		if err.Error() == repository.ErrNotFound.Error() {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
