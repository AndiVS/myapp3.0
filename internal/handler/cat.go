package handler

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"myapp3.0/internal/model"
	"myapp3.0/internal/repository"
	"net/http"
)

type CatHandler struct {
	Pool *pgxpool.Pool
}

//curl -d '{"name":"A","type":"B"}'  10.1.0.1:8080/records
func (h *CatHandler) Add(c echo.Context) error {

	rec := new(model.Record)

	if err := c.Bind(rec); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return c.JSON(http.StatusBadRequest, "")
	}

	rep := repository.New(h.Pool)
	id, err := rep.Insert(rec, c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, id)
}

//curl 10.1.0.1:8080/records/6
func (h *CatHandler) Get(c echo.Context) error {

	rec := new(model.Record)

	if err := c.Bind(rec); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return c.JSON(http.StatusBadRequest, "")
	}

	rep := repository.New(h.Pool)

	r, err := rep.Select(rec, c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, r)
}

//SelectAll provides all cats
func (h *CatHandler) GetAll(c echo.Context) error {

	var rec []model.Record
	rep := repository.New(h.Pool)
	rec, err := rep.SelectAll(c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rec)
}

//curl -XPUT -H 'Content-Type: application/json' -d '{"name":"AAAA","type":"BBBB"}'  127.0.0.1:8080/records/5
func (h *CatHandler) Update(c echo.Context) error {
	rec := new(model.Record)

	if err := c.Bind(rec); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return c.JSON(http.StatusBadRequest, "")
	}

	rep := repository.New(h.Pool)
	err := rep.Update(rec, c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "completed successfully")
}

//curl -XDELETE -H 'Content-Type: application/json'  127.0.0.1:8080/records/6
func (h *CatHandler) Delete(c echo.Context) error {

	rec := new(model.Record)

	if err := c.Bind(rec); err != nil {
		log.Errorf("Bind fail : %v\n", err)
		return c.JSON(http.StatusBadRequest, "")
	}

	rep := repository.New(h.Pool)
	err := rep.Delete(rec, c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "completed successfully")
}
