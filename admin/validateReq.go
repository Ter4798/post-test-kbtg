package admin

import (
	"errors"
)

func ValidateRequest(req *personalAllowanceRequest) error {

	if err := validatePersonalAllowance(req); err != nil {
		return err
	}
	return nil
}

func validatePersonalAllowance(req *personalAllowanceRequest) error {
	if req.Amount < 10000 || req.Amount > 100000 {
		return errors.New("amount must be between 10,000 and 100,000.")
	}
	return nil
}
