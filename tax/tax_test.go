package tax

import (
	"errors"
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
		{499999.99, 34999.999},
		{500000, 35000},
		{750000, 72500},
		{999999.99, 109999.9985},
		{1000000, 110000},
		{1500000, 210000},
		{1999999.99, 309999.998},
		{2000000, 310000},
		{3000000, 660000},
		{4000000, 1010000},
		{10000000, 3110000},
	}

	for _, tc := range testCases {
		got := calculateGraduatedTax(tc.taxableIncome)
		if got != tc.want {
			t.Errorf("calculateGraduatedTax(%f) returned %f, expected %f", tc.taxableIncome, got, tc.want)
		}
	}
}

func TestCalculateNetTaxAndRefund(t *testing.T) {
	testCases := []struct {
		tax      float64
		wht      float64
		expected []float64
	}{
		{
			tax:      100.0,
			wht:      20.0,
			expected: []float64{80.0, 0.0},
		},
		{
			tax:      50.0,
			wht:      50.0,
			expected: []float64{0.0, 0.0},
		},
	}

	for _, tc := range testCases {
		netTax, taxRefund := calculateNetTaxAndRefund(tc.tax, tc.wht)
		if netTax != tc.expected[0] || taxRefund != tc.expected[1] {
			t.Errorf("Expected %v, %v but got %v, %v", tc.expected[0], tc.expected[1], netTax, taxRefund)
		}
	}
}

func TestValidateTotalIncome(t *testing.T) {
	testCases := []struct {
		request       Request
		expectedError error
	}{
		{Request{TotalIncome: 100000.0}, nil},
		{Request{TotalIncome: 0.0}, errors.New("totalIncome must be greater than zero")},
		{Request{TotalIncome: -50000.0}, errors.New("totalIncome must be greater than zero")},
	}

	for _, tc := range testCases {
		err := validateTotalIncome(&tc.request)
		if (err == nil && tc.expectedError != nil) || (err != nil && err.Error() != tc.expectedError.Error()) {
			t.Errorf("validateTotalIncome(%+v) returned error %v, expected %v", tc.request, err, tc.expectedError)
		}
	}
}

func TestValidateWHT(t *testing.T) {
	testCases := []struct {
		request       Request
		expectedError error
	}{
		{Request{
			WHT:         100,
			TotalIncome: 1000,
		}, nil},
		{Request{
			WHT:         -100,
			TotalIncome: 1000,
		}, errors.New("invalid WHT must be greater than zero and morn than TotalIncome")},
		{Request{
			WHT:         2000,
			TotalIncome: 1000,
		}, errors.New("invalid WHT must be greater than zero and morn than TotalIncome")},
	}

	for _, tc := range testCases {
		err := validateWHT(&tc.request)
		if (err == nil && tc.expectedError != nil) || (err != nil && err.Error() != tc.expectedError.Error()) {
			t.Errorf("validateWHT(%+v) returned error %v, expected %v", tc.request, err, tc.expectedError)
		}
	}
}
