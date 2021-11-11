// Package handler contain function for handling request
package handler

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"myapp3.0/internal/model"
	"myapp3.0/internal/repository"
	"myapp3.0/internal/service"
	"net/http"
)

// UserHandler struct that contain repository linc
type UserHandler struct {
	Service service.Users
}

// NewHandlerUser add new user handler
func NewHandlerUser(Service service.Users) *UserHandler {
	return &UserHandler{Service: Service}
}

// GetUser provides user
func (h *UserHandler) GetUser(c echo.Context) error {
	user := new(model.User)

	if err := c.Bind(user); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	user, err := h.Service.GetUser(c.Request().Context(), user.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

// GetAllUser provides all users
func (h *UserHandler) GetAllUser(c echo.Context) error {
	var user []*model.User

	user, err := h.Service.GetAllUser(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateUser updating record about user
func (h *UserHandler) UpdateUser(c echo.Context) error {
	user := new(model.User)

	if err := c.Bind(user); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err := h.Service.UpdateUser(c.Request().Context(), user.Username, user.IsAdmin)
	if err != nil {
		if err.Error() == repository.ErrNotFound.Error() {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// DeleteUser delete record about user
func (h *UserHandler) DeleteUser(c echo.Context) error {
	username := c.Param("username")

	err := h.Service.DeleteUser(c.Request().Context(), username)
	if err != nil {
		if err.Error() == "not found" {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
