package handler

import (
	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/AndiVS/myapp3.0/internal/repository"
	"github.com/AndiVS/myapp3.0/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

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

/*
func TestSingUp(t *testing.T) {
	// Arrange
	s := new(service.MockUsers)
	req := model.User{Username: "thirdUser", Password: "thirdUser", IsAdmin: false}
	s.On("SignUp", mockContext, &req).Return(nil)
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
	req := model.User{Username: "thirdUser", Password: "thirdUser", IsAdmin: false}
	s.On("SignUp", mockContext, &req).Return(errSomeError)
	ctx, _ := setup(http.MethodPost, &req)

	// Act
	handl := NewHandlerUser(s)
	err := handl.AddU(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, errSomeError.Error()), err)
}
*/
func TestGetAllUser(t *testing.T) {
	// Arrange
	s := new(service.MockUsers)
	s.On("GetAllUser", mockContext).Return(users, nil)
	ctx, rec := setup(http.MethodGet, nil)

	// Act
	handl := NewHandlerUser(s)
	err := handl.GetAllUser(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, mustEncodeJSON(users), rec.Body.String())
}

func TestGetAllUserRepositoryFailed(t *testing.T) {
	// Arrange
	s := new(service.MockUsers)
	s.On("GetAllUser", mockContext).Return(nil, errSomeError)
	ctx, _ := setup(http.MethodGet, nil)

	// Act
	handl := NewHandlerUser(s)
	err := handl.GetAllUser(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, errSomeError.Error()), err)
}

func TestDeleteUser(t *testing.T) {
	// Arrange
	s := new(service.MockUsers)
	username := firstU.Username
	s.On("DeleteUser", mockContext, username).Return(nil)
	ctx, rec := setup(http.MethodDelete, nil)
	ctx.SetParamNames("username")
	ctx.SetParamValues(username)

	// Act
	handle := NewHandlerUser(s)
	err := handle.DeleteUser(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "", rec.Body.String())
}

func TestDeleteUserNotFound(t *testing.T) {
	// Arrange
	s := new(service.MockUsers)
	username := " "
	s.On("DeleteUser", mockContext, username).Return(repository.ErrNotFound)
	ctx, _ := setup(http.MethodDelete, nil)
	ctx.SetParamNames("username")
	ctx.SetParamValues(username)

	// Act
	handle := NewHandlerUser(s)
	err := handle.DeleteUser(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusNotFound), err)
}

func TestDeleteUserRepositoryFailed(t *testing.T) {
	// Arrange
	s := new(service.MockUsers)
	username := " "
	s.On("DeleteUser", mockContext, username).Return(errSomeError)
	ctx, _ := setup(http.MethodDelete, nil)
	ctx.SetParamNames("username")
	ctx.SetParamValues(username)

	// Act
	handle := NewHandlerUser(s)
	err := handle.DeleteUser(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, errSomeError.Error()), err)
}

func TestUpdateU(t *testing.T) {
	// Arrange
	s := new(service.MockUsers)
	req := model.User{Username: "thirdUser", Password: "thirdUser", IsAdmin: false}
	s.On("UpdateUser", mockContext, req.Username, req.IsAdmin).Return(nil)
	ctx, rec := setup(http.MethodPut, &req)
	ctx.SetParamNames("username")
	ctx.SetParamValues(req.Username)

	// Act
	handle := NewHandlerUser(s)
	err := handle.UpdateUser(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "", rec.Body.String())
}

func TestUpdateUserNotFound(t *testing.T) {
	// Arrange
	s := new(service.MockUsers)
	req := model.User{Username: "thirdUser", Password: "thirdUser", IsAdmin: false}
	s.On("UpdateUser", mockContext, req.Username, req.IsAdmin).Return(repository.ErrNotFound)
	ctx, _ := setup(http.MethodPut, nil)
	ctx.SetParamNames("username")
	ctx.SetParamValues(req.Username)

	// Act
	handle := NewHandlerUser(s)
	err := handle.UpdateUser(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusNotFound), err)
}

func TestUpdateUserRepositoryFailed(t *testing.T) {
	// Arrange
	s := new(service.MockUsers)
	req := model.User{Username: "thirdUser", Password: "thirdUser", IsAdmin: false}
	s.On("UpdateUser", mockContext, req.Username, req.IsAdmin).Return(errSomeError)
	ctx, _ := setup(http.MethodPut, nil)
	ctx.SetParamNames("username")
	ctx.SetParamValues(req.Username)

	// Act
	handle := NewHandlerUser(s)
	err := handle.UpdateUser(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, errSomeError.Error()), err)
}
