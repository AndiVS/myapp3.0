package handler

import (
	echo "github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"myapp3.0/internal/model"
	"myapp3.0/internal/repository"
	"net/http"
)

type CatHandler struct {
	Rep *repository.Repository
}

//Add add record about cat
func (h *CatHandler) Add(c echo.Context) error {

	rec := new(model.Record)

	if err := c.Bind(rec); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err := h.Rep.Insert(rec, c.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusCreated, rec.ID)
}

//Get provides cat
func (h *CatHandler) Get(c echo.Context) error {

	rec := new(model.Record)

	if err := c.Bind(rec); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	r, err := h.Rep.Select(&rec.ID, c.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusOK, r)
}

//GetAll provides all cats
func (h *CatHandler) GetAll(c echo.Context) error {

	var rec []*model.Record

	rec, err := h.Rep.SelectAll(c.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusOK, rec)
}

//Update updating record about cat
func (h *CatHandler) Update(c echo.Context) error {
	rec := new(model.Record)

	if err := c.Bind(rec); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err := h.Rep.Update(rec, c.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusOK, "completed successfully")
}

//Delete delete record about cat
func (h *CatHandler) Delete(c echo.Context) error {

	rec := new(model.Record)

	if err := c.Bind(rec); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err := h.Rep.Delete(&rec.ID, c.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusOK, "completed successfully")
}
