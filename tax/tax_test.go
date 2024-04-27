package tax

import (
	"testing"
)

func TestCalculateDeductions(t *testing.T) {
	want := 60000.0
	got := calculateDeductions()

	if got != want {
		t.Errorf("calculateDeductions() returned %f, expected %f", got, want)
	}
}

func TestCalculateTaxableIncome(t *testing.T) {
	testCases := []struct {
		totalIncome    float64
		totalDeduction float64
		want           float64
	}{
		{100000.0, 60000.0, 40000.0},
		{500000.0, 60000.0, 440000.0},
		{1500000.0, 60000.0, 1440000.0},
	}

	for _, tc := range testCases {
		got := calculateTaxableIncome(tc.totalIncome, tc.totalDeduction)
		if got != tc.want {
			t.Errorf("calculateTaxableIncome(%f, %f) returned %f, expected %f", tc.totalIncome, tc.totalDeduction, got, tc.want)
		}
	}
}

func TestCalculateGraduatedTax(t *testing.T) {
	testCases := []struct {
		taxableIncome float64
		want          float64
	}{
		{0, 0},
		{149999.99, 0},
		{150000, 0},
		{325000, 17500},
		{499999.99, 34999.999000},
		{500000, 35000},
		{750000, 72500},
		{999999.99, 109999.9985},
		{1000000, 110000},
		{1500000, 235000},
		{1999999.99, 334999.998},
		{2000000, 335000},
		{3000000, 885000},
		{4000000, 1235000},
		{10000000, 3335000},
	}

	for _, tc := range testCases {
		got := calculateGraduatedTax(tc.taxableIncome)
		if got != tc.want {
			t.Errorf("calculateGraduatedTax(%f) returned %f, expected %f", tc.taxableIncome, got, tc.want)
		}
	}
}
