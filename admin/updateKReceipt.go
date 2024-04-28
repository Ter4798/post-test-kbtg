package admin

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UpdateKReceiptAllowance(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req kReceiptAllowanceRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		if err := validateKReceiptAllowance(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM taxdeduction WHERE name = 'kReceiptAllowance')").Scan(&exists)
		if err != nil {
			return err
		}

		if exists {
			_, err = db.Exec("UPDATE taxdeduction SET amount = $1 WHERE name = 'kReceiptAllowance'", req.Amount)
			if err != nil {
				return err
			}
		} else {
			_, err = db.Exec("INSERT INTO taxdeduction (name, amount) VALUES ('kReceiptAllowance', $1)", req.Amount)
			if err != nil {
				return err
			}
		}

		resp := kReceiptAllowanceResponse{
			KReceiptDeduction: req.Amount,
		}
		return c.JSON(http.StatusOK, resp)
	}
}
