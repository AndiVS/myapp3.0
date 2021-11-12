package handler

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/AndiVS/myapp3.0/internal/repository"
	"github.com/AndiVS/myapp3.0/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	mockContext  = mock.Anything
	errSomeError = errors.New("some error")
	firstCat     = model.Cat{
		ID:   uuid.New(),
		Name: "firstCat",
		Type: "firstType",
	}
	secondCat = model.Cat{
		ID:   uuid.New(),
		Name: "secondCat",
		Type: "secondType",
	}
	cats = []*model.Cat{&firstCat, &secondCat}
)

func TestAddCat(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := uuid.New()
	req := model.Cat{ID: uuid.Nil, Name: "thirdCat", Type: "thirdType"}
	s.On("AddCat", mockContext, &req).Return(id, nil)
	ctx, rec := setup(http.MethodPost, &req)

	// Act
	handle := NewHandlerCat(s)
	err := handle.AddCat(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, mustEncodeJSON(id.String()), rec.Body.String())
}

func TestAddCatServiceFailed(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := uuid.New()
	req := model.Cat{ID: id, Name: "thirdCat", Type: "thirdType"}
	s.On("AddCat", mockContext, &req).Return(nil, errSomeError)
	ctx, _ := setup(http.MethodPost, &req)

	// Act
	handle := NewHandlerCat(s)
	err := handle.AddCat(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, errSomeError.Error()), err)
}

func TestGetCat(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := firstCat.ID
	s.On("GetCat", mockContext, id).Return(firstCat, nil)
	ctx, rec := setup(http.MethodGet, nil)
	ctx.SetParamNames("_id")
	ctx.SetParamValues(id.String())

	// Act
	handle := NewHandlerCat(s)
	err := handle.GetCat(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, mustEncodeJSON(firstCat), rec.Body.String())
}

func TestGetCatMalformedId(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	ctx, _ := setup(http.MethodGet, nil)
	ctx.SetParamNames("_id")
	ctx.SetParamValues("malformed-uuid")

	// Act
	handle := NewHandlerCat(s)
	err := handle.GetCat(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusBadRequest), err)
}

func TestGetCatNotFound(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := uuid.New().String()
	s.On("GetCat", mockContext, mock.AnythingOfType("uuid.UUID")).Return(model.Cat{}, repository.ErrNotFound)
	ctx, _ := setup(http.MethodGet, nil)
	ctx.SetParamNames("_id")
	ctx.SetParamValues(id)

	// Act
	handle := NewHandlerCat(s)
	err := handle.GetCat(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusNotFound), err)
}

func TestGetCatRepositoryFailed(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := uuid.New().String()
	s.On("GetCat", mockContext, mock.AnythingOfType("uuid.UUID")).Return(model.Cat{}, errSomeError)
	ctx, _ := setup(http.MethodGet, nil)
	ctx.SetParamNames("_id")
	ctx.SetParamValues(id)

	// Act
	handle := NewHandlerCat(s)
	err := handle.GetCat(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, errSomeError.Error()), err)
}

func TestGetAllCat(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	s.On("GetAllCat", mockContext).Return(cats, nil)
	ctx, rec := setup(http.MethodGet, nil)

	// Act
	handle := NewHandlerCat(s)
	err := handle.GetAllCat(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, mustEncodeJSON(cats), rec.Body.String())
}

func TestGetAllCatRepositoryFailed(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	s.On("GetAllCat", mockContext).Return(nil, errSomeError)
	ctx, _ := setup(http.MethodGet, nil)

	// Act
	handle := NewHandlerCat(s)
	err := handle.GetAllCat(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, errSomeError.Error()), err)
}

func TestDeleteCat(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := firstCat.ID
	s.On("DeleteCat", mockContext, id).Return(nil)
	ctx, rec := setup(http.MethodDelete, nil)
	ctx.SetParamNames("_id")
	ctx.SetParamValues(id.String())

	// Act
	handle := NewHandlerCat(s)
	err := handle.DeleteCat(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "", rec.Body.String())
}

func TestDeleteCatFailed(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	ctx, _ := setup(http.MethodDelete, nil)
	ctx.SetParamNames("_id")
	ctx.SetParamValues("malformed-uuid")

	// Act
	handle := NewHandlerCat(s)
	err := handle.DeleteCat(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusBadRequest), err)
}

func TestDeleteCatNotFound(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := uuid.New().String()
	s.On("DeleteCat", mockContext, mock.AnythingOfType("uuid.UUID")).Return(repository.ErrNotFound)
	ctx, _ := setup(http.MethodDelete, nil)
	ctx.SetParamNames("_id")
	ctx.SetParamValues(id)

	// Act
	handle := NewHandlerCat(s)
	err := handle.DeleteCat(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusNotFound), err)
}

func TestUpdateCat(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := firstCat.ID
	req := model.Cat{ID: id, Name: "thirdCat", Type: "thirdType"}
	s.On("UpdateCat", mockContext, &req).Return(nil)
	ctx, rec := setup(http.MethodPut, &req)
	ctx.SetParamNames("_id")
	ctx.SetParamValues(id.String())

	// Act
	handle := NewHandlerCat(s)
	err := handle.UpdateCat(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "", rec.Body.String())
}

func TestUpdateCatMalformedId(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	req := model.Cat{ID: uuid.Nil, Name: "thirdCat", Type: "thirdType"}
	ctx, _ := setup(http.MethodPut, &req)
	ctx.SetParamNames("_id")
	ctx.SetParamValues("malformed-uuid")

	// Act
	handle := NewHandlerCat(s)
	err := handle.UpdateCat(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusBadRequest), err)
}

func TestUpdateCatNotFound(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := uuid.New()
	req := model.Cat{ID: id, Name: "thirdCat", Type: "thirdType"}
	s.On("UpdateCat", mockContext, &req).Return(repository.ErrNotFound)
	ctx, _ := setup(http.MethodPut, &req)
	ctx.SetParamNames("_id")
	ctx.SetParamValues(id.String())

	// Act
	handle := NewHandlerCat(s)
	err := handle.UpdateCat(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusNotFound), err)
}

func setup(method string, body interface{}) (echo.Context, *httptest.ResponseRecorder) {
	jsonBody := ""
	if body != nil {
		jsonBody = mustEncodeJSON(body)
	}
	request := httptest.NewRequest(method, "/", strings.NewReader(jsonBody))
	if body != nil {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)
	return c, recorder
}

func mustEncodeJSON(data interface{}) string {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	err := enc.Encode(data)
	if err != nil {
		panic(err)
	}
	return buffer.String()
}
