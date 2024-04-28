package admin

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UpdatePersonalAllowance(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req personalAllowanceRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM taxdeduction WHERE name = 'personalAllowance')").Scan(&exists)
		if err != nil {
			return err
		}

		if exists {
			_, err = db.Exec("UPDATE taxdeduction SET amount = $1 WHERE name = 'personalAllowance'", req.Amount)
			if err != nil {
				return err
			}
		} else {
			_, err = db.Exec("INSERT INTO taxdeduction (name, amount) VALUES ('personalAllowance', $1)", req.Amount)
			if err != nil {
				return err
			}
		}

		resp := personalAllowanceResponse{
			PersonalDeduction: req.Amount,
		}
		return c.JSON(http.StatusOK, resp)
	}
}
