package tax

import (
	"errors"
)

func ValidateRequest(req *Request) error {

	if req.TotalIncome <= 0 {
		return errors.New("totalIncome must be greater than zero")
	}

	return nil
}
