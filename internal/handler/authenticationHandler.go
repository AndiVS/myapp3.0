package handler

import (
	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/AndiVS/myapp3.0/internal/service"
	"github.com/AndiVS/myapp3.0/protocol"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"net/http"
	"time"
)

// AuthenticationHandler handler for aunt
type AuthenticationHandler struct {
	Service service.Authentication
	protocol.UnimplementedUserServiceServer
}

// NewHandlerAuthentication create AuthenticationHandler
func NewHandlerAuthentication(Service service.Authentication) *AuthenticationHandler {
	return &AuthenticationHandler{Service: Service}
}

// SignUp User about cat
func (h *AuthenticationHandler) SignUp(c echo.Context) error {
	user := new(model.User)

	if err := c.Bind(user); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err := h.Service.SignUp(c.Request().Context(), user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusCreated)
}

// SignIn generate token
func (h *AuthenticationHandler) SignIn(c echo.Context) error {
	user := new(model.User)

	if err := c.Bind(user); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	accessToken, refreshToken, err := h.Service.SignIn(c.Request().Context(), user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	service.SetUserCookie(user, time.Now().Add(20*time.Second), c)

	service.SetTokenCookie("refreshToken", refreshToken, time.Now().Add(1000*time.Second), c)

	return c.JSON(http.StatusOK, echo.Map{
		"token": accessToken,
	})
}
