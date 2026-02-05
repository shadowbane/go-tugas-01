package dto

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}
