package admin

type personalAllowanceRequest struct {
	Amount float64 `json:"amount"`
}

type personalAllowanceResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}
