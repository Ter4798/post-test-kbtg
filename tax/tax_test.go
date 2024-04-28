package tax

import (
	"errors"
	"testing"
)

func TestCalculateDeductions(t *testing.T) {
	testCases := []struct {
		name           string
		allowances     []Allowance
		maxDonation    float64
		maxKReceipt    float64
		expectedResult float64
	}{
		{
			name:           "No allowances",
			allowances:     []Allowance{},
			maxDonation:    100000,
			maxKReceipt:    50000,
			expectedResult: 0,
		},
		{
			name: "Single donation within limit",
			allowances: []Allowance{
				{AllowanceType: "donation", Amount: 50000},
			},
			maxDonation:    100000,
			maxKReceipt:    50000,
			expectedResult: 50000,
		},
		{
			name: "Single donation exceeding limit",
			allowances: []Allowance{
				{AllowanceType: "donation", Amount: 150000},
			},
			maxDonation:    100000,
			maxKReceipt:    50000,
			expectedResult: 100000,
		},
		{
			name: "Single k-receipt within limit",
			allowances: []Allowance{
				{AllowanceType: "k-receipt", Amount: 30000},
			},
			maxDonation:    100000,
			maxKReceipt:    50000,
			expectedResult: 30000,
		},
		{
			name: "Single k-receipt exceeding limit",
			allowances: []Allowance{
				{AllowanceType: "k-receipt", Amount: 60000},
			},
			maxDonation:    100000,
			maxKReceipt:    50000,
			expectedResult: 50000,
		},
		{
			name: "Multiple allowances",
			allowances: []Allowance{
				{AllowanceType: "donation", Amount: 80000},
				{AllowanceType: "k-receipt", Amount: 40000},
				{AllowanceType: "donation", Amount: 30000},
			},
			maxDonation:    100000,
			maxKReceipt:    50000,
			expectedResult: 150000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := calculateDeductions(tc.allowances, tc.maxDonation, tc.maxKReceipt)
			if result != tc.expectedResult {
				t.Errorf("Expected %f, got %f", tc.expectedResult, result)
			}
		})
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

func TestValidateAllowanceTypes(t *testing.T) {
	testCases := []struct {
		request       Request
		expectedError error
	}{
		{Request{
			Allowances: []Allowance{
				{AllowanceType: "donation"},
				{AllowanceType: "k-receipt"},
			},
		}, nil},
		{Request{
			Allowances: []Allowance{
				{AllowanceType: "invalid"},
			},
		}, errors.New("invalid allowance type type should be donation, k-receipt")},
		{Request{
			Allowances: []Allowance{
				{AllowanceType: "donation"},
				{AllowanceType: "invalid"},
				{AllowanceType: "k-receipt"},
			},
		}, errors.New("invalid allowance type type should be donation, k-receipt")},
	}

	for _, tc := range testCases {
		err := validateAllowanceTypes(&tc.request)
		if (err == nil && tc.expectedError != nil) || (err != nil && err.Error() != tc.expectedError.Error()) {
			t.Errorf("validateWHT(%+v) returned error %v, expected %v", tc.request, err, tc.expectedError)
		}
	}
}

func TestCalculateTaxLevels(t *testing.T) {
	testCases := []struct {
		name          string
		taxableIncome float64
		expected      []TaxLevel
	}{
		{
			name:          "No tax",
			taxableIncome: 100000,
			expected: []TaxLevel{
				{"0-150,000", 0.0},
				{"150,001-500,000", 0.0},
				{"500,001-1,000,000", 0.0},
				{"1,000,001-2,000,000", 0.0},
				{"2,000,001 ขึ้นไป", 0.0},
			},
		},
		{
			name:          "First tax bracket",
			taxableIncome: 300000,
			expected: []TaxLevel{
				{"0-150,000", 0.0},
				{"150,001-500,000", 15000.0},
				{"500,001-1,000,000", 0.0},
				{"1,000,001-2,000,000", 0.0},
				{"2,000,001 ขึ้นไป", 0.0},
			},
		},
		{
			name:          "Second tax bracket",
			taxableIncome: 750000,
			expected: []TaxLevel{
				{"0-150,000", 0.0},
				{"150,001-500,000", 35000.0},
				{"500,001-1,000,000", 37500.0},
				{"1,000,001-2,000,000", 0.0},
				{"2,000,001 ขึ้นไป", 0.0},
			},
		},
		{
			name:          "Third tax bracket",
			taxableIncome: 1500000,
			expected: []TaxLevel{
				{"0-150,000", 0.0},
				{"150,001-500,000", 35000.0},
				{"500,001-1,000,000", 75000.0},
				{"1,000,001-2,000,000", 100000.0},
				{"2,000,001 ขึ้นไป", 0.0},
			},
		},
		{
			name:          "Fourth tax bracket",
			taxableIncome: 3000000,
			expected: []TaxLevel{
				{"0-150,000", 0.0},
				{"150,001-500,000", 35000.0},
				{"500,001-1,000,000", 75000.0},
				{"1,000,001-2,000,000", 200000.0},
				{"2,000,001 ขึ้นไป", 350000.0},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := calculateTaxLevels(tc.taxableIncome)
			for i, level := range result {
				if level.Level != tc.expected[i].Level || level.Tax != tc.expected[i].Tax {
					t.Errorf("Incorrect tax level calculation. Expected: %v, Got: %v", tc.expected[i], level)
				}
			}
		})
	}
}
