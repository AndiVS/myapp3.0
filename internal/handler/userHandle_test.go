package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"myapp3.0/internal/model"
	"myapp3.0/internal/service"
	"net/http"
	"testing"
)

var (
	firstU = model.User{
		Username: "firstUser",
		Password: "firstUser",
		IsAdmin:  false,
	}
	secondU = model.User{
		Username: "secondUser",
		Password: "secondUser",
		IsAdmin:  false,
	}
	users = []*model.User{&firstU, &secondU}
)

func TestAddU(t *testing.T) {
	// Arrange
	s := new(service.MockUsers)
	req := model.User{"thirdUser", "thirdUser", false}
	s.On("AddU", mockContext, &req).Return(nil)
	ctx, rec := setup(http.MethodPost, &req)

	// Act
	handl := NewHandlerUser(s)
	err := handl.AddU(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, "", rec.Body.String())
}

func TestAddUServiceFailed(t *testing.T) {
	// Arrange
	s := new(service.MockUsers)
	req := model.User{"thirdUser", "thirdUser", false}
	s.On("AddU", mockContext, &req).Return(errSomeError)
	ctx, _ := setup(http.MethodPost, &req)

	// Act
	handl := NewHandlerUser(s)
	err := handl.AddU(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, errSomeError.Error()), err)
}

func TestGetAllU(t *testing.T) {
	// Arrange
	s := new(service.MockUsers)
	s.On("GetAllU", mockContext).Return(users, nil)
	ctx, rec := setup(http.MethodGet, nil)

	// Act
	handl := NewHandlerUser(s)
	err := handl.GetAllU(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, mustEncodeJSON(users), rec.Body.String())
}

func TestGetAllURepositoryFailed(t *testing.T) {
	// Arrange
	s := new(service.MockUsers)
	s.On("GetAllU", mockContext).Return(nil, errSomeError)
	ctx, _ := setup(http.MethodGet, nil)

	// Act
	handl := NewHandlerUser(s)
	err := handl.GetAllU(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, errSomeError.Error()), err)
}

func TestDeleteU(t *testing.T) {
	// Arrange
	s := new(service.MockUsers)
	usrename := firstU.Username
	s.On("DeleteU", mockContext, usrename).Return(nil)
	ctx, rec := setup(http.MethodDelete, nil)
	ctx.SetParamNames("username")
	ctx.SetParamValues(usrename)

	// Act
	handl := NewHandlerUser(s)
	err := handl.DeleteU(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "", rec.Body.String())
}

func TestDeleteUNotFound(t *testing.T) {
	// Arrange
	s := new(service.MockUsers)
	usrename := " "
	s.On("DeleteU", mockContext, usrename).Return(errors.New("not found"))
	ctx, _ := setup(http.MethodDelete, nil)
	ctx.SetParamNames("username")
	ctx.SetParamValues(usrename)

	// Act
	handl := NewHandlerUser(s)
	err := handl.DeleteU(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusNotFound), err)
}

func TestDeleteURepositoryFailed(t *testing.T) {
	// Arrange
	s := new(service.MockUsers)
	usrename := " "
	s.On("DeleteU", mockContext, usrename).Return(errSomeError)
	ctx, _ := setup(http.MethodDelete, nil)
	ctx.SetParamNames("username")
	ctx.SetParamValues(usrename)

	// Act
	handl := NewHandlerUser(s)
	err := handl.DeleteU(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, errSomeError.Error()), err)
}

func TestUpdateU(t *testing.T) {
	// Arrange
	s := new(service.MockUsers)
	req := model.User{"thirdUser", "thirdUser", false}
	s.On("UpdateU", mockContext, req.Username, req.IsAdmin).Return(nil)
	ctx, rec := setup(http.MethodPut, &req)
	ctx.SetParamNames("username")
	ctx.SetParamValues(req.Username)

	// Act
	handl := NewHandlerUser(s)
	err := handl.UpdateU(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "", rec.Body.String())
}

func TestUpdateUNotFound(t *testing.T) {
	// Arrange
	s := new(service.MockUsers)
	req := model.User{"thirdUser", "thirdUser", false}
	s.On("UpdateU", mockContext, req.Username, req.IsAdmin).Return(errors.New("not found"))
	ctx, _ := setup(http.MethodPut, nil)
	ctx.SetParamNames("username")
	ctx.SetParamValues(req.Username)

	// Act
	handl := NewHandlerUser(s)
	err := handl.UpdateU(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusNotFound), err)
}

func TestUpdateURepositoryFailed(t *testing.T) {
	// Arrange
	s := new(service.MockUsers)
	req := model.User{"thirdUser", "thirdUser", false}
	s.On("UpdateU", mockContext, req.Username, req.IsAdmin).Return(errSomeError)
	ctx, _ := setup(http.MethodPut, nil)
	ctx.SetParamNames("username")
	ctx.SetParamValues(req.Username)

	// Act
	handl := NewHandlerUser(s)
	err := handl.UpdateU(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, errSomeError.Error()), err)
}
