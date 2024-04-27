package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Ter4798/post-test-kbtg/tax"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))

	e.POST("/tax/calculations", func(c echo.Context) error {
		req := new(tax.Request)
		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		t := tax.CalculateTax(req.TotalIncome)
		resp := &tax.Response{
			Tax: t,
		}

		return c.JSON(http.StatusOK, resp)
	})

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
