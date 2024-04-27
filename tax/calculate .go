package tax

func calculateDeductions(allowances []Allowance) float64 {
	var totalDeduction float64
	personalAllowance := 60000.0
	totalDeduction += personalAllowance
	for _, allowance := range allowances {
		switch allowance.AllowanceType {
		case "donation":
			deduction := allowance.Amount
			if deduction > 100000 {
				deduction = 100000
			}
			totalDeduction += deduction
		case "k-receipt":
			deduction := allowance.Amount
			if deduction > 50000 {
				deduction = 50000
			}
			totalDeduction += deduction
		}
	}
	return totalDeduction
}

func calculateTaxableIncome(totalIncome, totalDeduction float64) float64 {
	return totalIncome - totalDeduction
}

func calculateGraduatedTax(taxableIncome float64) float64 {
	var tax float64

	switch {
	case taxableIncome <= 150000:
		tax = 0
	case taxableIncome <= 500000:
		tax = 0.1 * (taxableIncome - 150000)
	case taxableIncome <= 1000000:
		tax = 35000 + 0.15*(taxableIncome-500000)
	case taxableIncome <= 2000000:
		tax = 135000 + 0.2*(taxableIncome-1000000)
	default:
		tax = 535000 + 0.35*(taxableIncome-2000000)
	}

	return tax
}

func calculateNetTaxAndRefund(tax, wht float64) (float64, float64) {
	taxRefund := 0.0
	if tax < wht {
		taxRefund = wht - tax
		tax = 0
	} else {
		tax = tax - wht
	}

	return tax, taxRefund
}

func CalculateTax(totalIncome float64, wht float64, allowances []Allowance) (float64, float64) {
	totalDeduction := calculateDeductions(allowances)
	taxableIncome := calculateTaxableIncome(totalIncome, totalDeduction)
	tax := calculateGraduatedTax(taxableIncome)
	netTax, taxRefund := calculateNetTaxAndRefund(tax, wht)
	return netTax, taxRefund
}
