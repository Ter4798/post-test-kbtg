package tax

import (
	"database/sql"
)

func getPersonalAllowance(db *sql.DB) (float64, error) {
	var personalAllowance float64 = 60000.0
	var exists bool

	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM taxdeduction WHERE name = 'personalAllowance')").Scan(&exists)
	if err != nil {
		return 0, err
	}

	if exists {
		err = db.QueryRow("SELECT amount FROM taxdeduction WHERE name = 'personalAllowance'").Scan(&personalAllowance)
		if err != nil {
			return 0, err
		}
	}

	return personalAllowance, nil
}
