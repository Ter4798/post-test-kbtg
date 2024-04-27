package tax

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type Request struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type Response struct {
	Tax       float64 `json:"tax"`
	TaxRefund float64 `json:"taxRefund,omitempty"`
}
