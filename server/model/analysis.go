package model

type OrderSummary struct {
	Date           string  `json:"date"`
	OrderTotal     int64   `json:"order_total"`
	IncomeTotal    float64 `json:"income_total"`
	GeneralTotal   int64   `json:"general_total"`
	RechargeTotal  int64   `json:"recharge_total"`
	SubscribeTotal int64   `json:"subscribe_total"`
}
type UserSummary struct {
	Date          string `json:"date"`
	RegisterTotal int64  `json:"register_total"`
}
