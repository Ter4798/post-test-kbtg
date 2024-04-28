package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/Ter4798/post-test-kbtg/admin"
	"github.com/Ter4798/post-test-kbtg/tax"
	"github.com/labstack/echo/v4"
)

// const (
// 	host     = "postgres"
// 	portp    = 5432
// 	user     = "postgres"
// 	password = "postgres"
// 	dbname   = "ktaxes"
// )

func main() {

	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, portp, user, password, dbname)
	connStr := "user=postgres password=postgres host=localhost port=5432 sslmode=disable dbname=ktaxes"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS taxdeduction (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        amount FLOAT8 NOT NULL
    )`)
	if err != nil {
		panic(err)
	}

	fmt.Print("Connect DB")

	e := echo.New()
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))

	e.POST("/tax/calculations", func(c echo.Context) error {
		req := new(tax.Request)

		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := tax.ValidateRequest(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		t, taxRefund, taxLevels, err := tax.CalculateTax(db, req.TotalIncome, req.WHT, req.Allowances)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		resp := &tax.Response{
			Tax:       t,
			TaxLevels: taxLevels,
		}
		if taxRefund > 0 {
			resp.TaxRefund = taxRefund
		}

		return c.JSON(http.StatusOK, resp)
	})

	e.POST("/admin/deductions/personal", admin.UpdatePersonalAllowance(db))

	go func() {
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("ListenAndServe error: ", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	fmt.Println("Shutting down the server")
}
