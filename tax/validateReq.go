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
	if err := validateAllowanceTypes(req); err != nil {
		return err
	}
	if err := validateAllowanceAmounts(req); err != nil {
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

func validateAllowanceTypes(req *Request) error {
	validAllowanceTypes := map[string]bool{
		"donation":  true,
		"k-receipt": true,
	}

	for _, allowance := range req.Allowances {
		if _, ok := validAllowanceTypes[allowance.AllowanceType]; !ok {
			return errors.New("invalid allowance type type should be donation, k-receipt")
		}
	}
	return nil
}

func validateAllowanceAmounts(req *Request) error {
	for _, allowance := range req.Allowances {
		if allowance.Amount < 0 {
			return errors.New("invalid allowance amount must be greater than 0")
		}
	}
	return nil
}
