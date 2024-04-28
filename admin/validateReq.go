package admin

import (
	"errors"
)

func validatePersonalAllowance(req *personalAllowanceRequest) error {
	if req.Amount < 10000 || req.Amount > 100000 {
		return errors.New("amount must be between 10000 and 100000")
	}
	return nil
}

func validateKReceiptAllowance(req *kReceiptAllowanceRequest) error {
	if req.Amount < 0 || req.Amount > 100000 {
		return errors.New("amount must be between 0 and 100000")
	}
	return nil
}
