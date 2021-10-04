package handler

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"myapp3.0/internal/model"
	"net/http"
	"strconv"
)

type CatHandler struct {
	Pool *pgxpool.Pool
}

//curl -d '{"name":"A","type":"B"}'  10.1.0.1:8080/records
func (h *CatHandler) Insert( c echo.Context) error {
	rec := model.Record{}

	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body for Insert: %s\n", err)
		return c.String(http.StatusBadRequest, "")
	}

	err = json.Unmarshal(b, &rec)
	if err != nil {
		log.Printf("Failed unmarshaling in addCats: %s\n", err)
		return c.String(http.StatusBadRequest, "")
	}

	conn, err := h.Pool.Acquire(c.Request().Context())
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v\n", err)
		return c.String(http.StatusInternalServerError, "")
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		"INSERT INTO catsbase (name, type) VALUES ($1, $2) RETURNING id", rec.Name, rec.Type)

	var id uint64
	err = row.Scan(&id)
	if err != nil {
		log.Errorf("Unable to INSERT: %v", err)
		return c.String(http.StatusInternalServerError, "")
	}

	rs := make(map[string]string, 1)
	rs["id"] = strconv.FormatUint(id, 10)

	resp,err := json.Marshal(rs)
	if err != nil {
		log.Printf("Failed marshaling in addCats: %s\n", err)
		return c.String(http.StatusInternalServerError, "")//500
	}


	log.Printf("this is your cat: %#v\n", resp)
	return c.JSON(http.StatusOK, id)
}

//curl 10.1.0.1:8080/records/6
func (h *CatHandler) Select(c echo.Context) error {

	idp := c.Param("id")

	id, err := strconv.ParseUint(idp, 10, 64)
	if err != nil {
		log.Printf("Failed reading the request body for Select: %s\n", err)
		return c.String(http.StatusBadRequest, "")
	}

	conn, err := h.Pool.Acquire(c.Request().Context())
	if err != nil {
		log.Printf("Failed reading the request body for addCats: %s\n", err)
		return c.String(http.StatusInternalServerError, "")
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		"SELECT id, name, type FROM catsbase WHERE id = $1", id)

	var rec model.Record
	err = row.Scan(&rec.Id, &rec.Name, &rec.Type)
	if err == pgx.ErrNoRows {
		log.Errorf("No such row: %v", err)
		return c.String(http.StatusNotFound, "")
	}
	if err != nil {
		log.Errorf("Unable to SELECT: %v", err)
		return c.String(http.StatusInternalServerError, "")

	}

	resp,err := json.Marshal(rec)
	if err != nil {
		log.Printf("Failed unmarshaling in addCats: %s\n", err)
		return c.String(http.StatusInternalServerError, "")//500
	}

	log.Printf("this is your cat: %#v\n", resp)
	return c.JSON(http.StatusOK,rec)

}

//curl 10.1.0.1:8080/records
func (h *CatHandler) SelectAll(c echo.Context) error  {
	conn, err := h.Pool.Acquire(c.Request().Context())
	if err != nil {
		log.Errorf("Unable to acquire a database connection for SelectAll: %v\n", err)
		return c.String(http.StatusInternalServerError, "")
	}
	defer conn.Release()

	row, err := conn.Query(context.Background(),
		"SELECT * FROM catsbase")

	rec := []model.Record{}
	for row.Next() {
		var rc model.Record
		err = row.Scan(&rc.Id, &rc.Name, &rc.Type)
		if err == pgx.ErrNoRows {
			return c.String(http.StatusNotFound, "")
		}
		if err != nil {
			log.Errorf("Unable to SELECT: %v", err)
			return c.String(http.StatusInternalServerError, "")

		}
		rec = append(rec, rc)
	}

	/*resp,err := json.Marshal(rec)
	if err != nil {
		log.Errorf("Unable to encode json: %v", err)
		return c.String(http.StatusInternalServerError, "")//500
	}*/

	//log.Printf("this is your cat: %#v\n", rec)
	return c.JSON(http.StatusOK, rec)
}

//curl -XPUT  -d '{"name":"AAA","type":"BBB"}'  10.1.0.1:8080/records/10
func (h *CatHandler) Update(c echo.Context) error {
	idp := c.Param("id")

	id, err := strconv.ParseUint(idp, 10, 64)
	if err != nil {
		log.Printf("Failed reading the request body for Updating: %s\n", err)
		return c.String(http.StatusInternalServerError, "")
	}

	rec := model.Record{}

	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body for Updating: %s\n", err)
		return c.String(http.StatusBadRequest, "")
	}

	err = json.Unmarshal(b, &rec)
	if err != nil {
		log.Printf("Failed unmarshaling in Updating: %s\n", err)
		return c.String(http.StatusBadRequest, "")
	}


	conn, err := h.Pool.Acquire(c.Request().Context())
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v", err)
		return c.String(http.StatusInternalServerError, "")

	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(),
		"UPDATE catsbase SET name = $2, type = $3 WHERE id = $1", id, rec.Name, rec.Type)

	if err != nil {
		log.Printf("Failed reading the request body for addCats: %s\n", err)
		return c.String(http.StatusInternalServerError, "")
	}

	if ct.RowsAffected() == 0 {
		log.Printf("Failed reading the request body for addCats: %s\n", err)
		return c.String(http.StatusNotFound, "")
	}


	return c.String(http.StatusOK, "updating complit ")
}

//curl -XDELETE  10.1.0.1:8080/records/9
func (h *CatHandler) Delete(c echo.Context) error {
	idp := c.Param("id")

	id, err := strconv.ParseUint(idp, 10, 64)
	if err != nil {
		log.Printf("Failed reading the request body for addCats: %s\n", err)
		return c.String(http.StatusBadRequest, "")
	}

	conn, err := h.Pool.Acquire(c.Request().Context())
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v", err)
		return c.String(http.StatusInternalServerError, "")
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), "DELETE FROM catsbase WHERE id = $1", id)

	if err != nil {
		log.Errorf("Unable to DELETE: %v", err)
		return c.String(http.StatusInternalServerError, "")
	}

	if ct.RowsAffected() == 0 {
		return c.String(http.StatusNotFound, "")
	}

	return c.String(http.StatusOK, "Delet secses ")
}

