package handler

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"myapp3.0/internal/model"
	"myapp3.0/internal/repository"
	"myapp3.0/internal/service"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	mockContext  = mock.Anything
	errSomeError = errors.New("some error")
	firstC       = model.Record{
		ID:   uuid.New(),
		Name: "firstCat",
		Type: "firstType",
	}
	secondC = model.Record{
		ID:   uuid.New(),
		Name: "secondCat",
		Type: "secondType",
	}
	cats = []*model.Record{&firstC, &secondC}
)

func TestAddC(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := uuid.New()
	req := model.Record{ID: uuid.Nil, Name: "thirdCat", Type: "thirdType"}
	s.On("AddC", mockContext, &req).Return(id, nil)
	ctx, rec := setup(http.MethodPost, &req)

	// Act
	handl := NewHandlerCat(s)
	err := handl.AddC(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, mustEncodeJSON(id.String()), rec.Body.String())
}

func TestAddCServiceFailed(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := uuid.New()
	req := model.Record{ID: id, Name: "thirdCat", Type: "thirdType"}
	s.On("AddC", mockContext, &req).Return(nil, errSomeError)
	ctx, _ := setup(http.MethodPost, &req)

	// Act
	handl := NewHandlerCat(s)
	err := handl.AddC(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, errSomeError.Error()), err)
}

func TestGetC(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := firstC.ID
	s.On("GetC", mockContext, id).Return(firstC, nil)
	ctx, rec := setup(http.MethodGet, nil)
	ctx.SetParamNames("_id")
	ctx.SetParamValues(id.String())

	// Act
	handl := NewHandlerCat(s)
	err := handl.GetC(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, mustEncodeJSON(firstC), rec.Body.String())
}

func TestGetCMalformedId(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	ctx, _ := setup(http.MethodGet, nil)
	ctx.SetParamNames("_id")
	ctx.SetParamValues("malformed-uuid")

	// Act
	handl := NewHandlerCat(s)
	err := handl.GetC(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusBadRequest), err)
}

func TestGetCNotFound(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := uuid.New().String()
	s.On("GetC", mockContext, mock.AnythingOfType("uuid.UUID")).Return(model.Record{}, repository.ErrNotFound)
	ctx, _ := setup(http.MethodGet, nil)
	ctx.SetParamNames("_id")
	ctx.SetParamValues(id)

	// Act
	handl := NewHandlerCat(s)
	err := handl.GetC(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusNotFound), err)
}

func TestGetCRepositoryFailed(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := uuid.New().String()
	s.On("GetC", mockContext, mock.AnythingOfType("uuid.UUID")).Return(model.Record{}, errSomeError)
	ctx, _ := setup(http.MethodGet, nil)
	ctx.SetParamNames("_id")
	ctx.SetParamValues(id)

	// Act
	handl := NewHandlerCat(s)
	err := handl.GetC(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, errSomeError.Error()), err)
}

func TestGetAllC(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	s.On("GetAllC", mockContext).Return(cats, nil)
	ctx, rec := setup(http.MethodGet, nil)

	// Act
	handl := NewHandlerCat(s)
	err := handl.GetAllC(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, mustEncodeJSON(cats), rec.Body.String())
}

func TestGetAllCRepositoryFailed(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	s.On("GetAllC", mockContext).Return(nil, errSomeError)
	ctx, _ := setup(http.MethodGet, nil)

	// Act
	handl := NewHandlerCat(s)
	err := handl.GetAllC(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, errSomeError.Error()), err)
}

func TestDeleteCat(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := firstC.ID
	s.On("DeleteC", mockContext, id).Return(nil)
	ctx, rec := setup(http.MethodDelete, nil)
	ctx.SetParamNames("_id")
	ctx.SetParamValues(id.String())

	// Act
	handl := NewHandlerCat(s)
	err := handl.DeleteC(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "", rec.Body.String())
}

func TestDeleteCFailed(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	ctx, _ := setup(http.MethodDelete, nil)
	ctx.SetParamNames("_id")
	ctx.SetParamValues("malformed-uuid")

	// Act
	handl := NewHandlerCat(s)
	err := handl.DeleteC(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusBadRequest), err)
}

func TestDeleteCNotFound(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := uuid.New().String()
	s.On("DeleteC", mockContext, mock.AnythingOfType("uuid.UUID")).Return(repository.ErrNotFound)
	ctx, _ := setup(http.MethodDelete, nil)
	ctx.SetParamNames("_id")
	ctx.SetParamValues(id)

	// Act
	handl := NewHandlerCat(s)
	err := handl.DeleteC(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusNotFound), err)
}

func TestUpdateC(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := firstC.ID
	req := model.Record{ID: id, Name: "thirdCat", Type: "thirdType"}
	s.On("UpdateC", mockContext, &req).Return(nil)
	ctx, rec := setup(http.MethodPut, &req)
	ctx.SetParamNames("_id")
	ctx.SetParamValues(id.String())

	// Act
	handl := NewHandlerCat(s)
	err := handl.UpdateC(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "", rec.Body.String())
}

func TestUpdateCMalformedId(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	req := model.Record{ID: uuid.Nil, Name: "thirdCat", Type: "thirdType"}
	ctx, _ := setup(http.MethodPut, &req)
	ctx.SetParamNames("_id")
	ctx.SetParamValues("malformed-uuid")

	// Act
	handl := NewHandlerCat(s)
	err := handl.UpdateC(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusBadRequest), err)
}

func TestUpdateCNotFound(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := uuid.New()
	req := model.Record{ID: id, Name: "thirdCat", Type: "thirdType"}
	s.On("UpdateC", mockContext, &req).Return(repository.ErrNotFound)
	ctx, _ := setup(http.MethodPut, &req)
	ctx.SetParamNames("_id")
	ctx.SetParamValues(id.String())

	// Act
	handl := NewHandlerCat(s)
	err := handl.UpdateC(ctx)

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
