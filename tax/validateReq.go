package tax

import (
	"errors"
)

func ValidateRequest(req *Request) error {

	if err := validateTotalIncome(req); err != nil {
		return err
	}
	if err := validateWHT(req); err != nil {
		return err
	}
	return nil
}

func validateTotalIncome(req *Request) error {
	if req.TotalIncome <= 0 {
		return errors.New("totalIncome must be greater than zero")
	}
	return nil
}

func validateWHT(req *Request) error {
	if req.WHT < 0 || req.WHT > req.TotalIncome {
		return errors.New("invalid WHT must be greater than zero and morn than TotalIncome")
	}
	return nil
}
