package arn

import (
	"strconv"

	"github.com/aerogo/nano"
)

// PayPalPayment is an approved and exeucted PayPal payment.
type PayPalPayment struct {
	UserID   string `json:"userId"`
	PayerID  string `json:"payerId"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	Method   string `json:"method"`
	Created  string `json:"created"`

	hasID
}

// NewPayPalPayment creates a new PayPalPayment object with the paypal provided ID.
func NewPayPalPayment(paymentID, payerID, userID, method, amount, currency string) *PayPalPayment {
	return &PayPalPayment{
		hasID: hasID{
			ID: paymentID,
		},
		PayerID:  payerID,
		UserID:   userID,
		Method:   method,
		Amount:   amount,
		Currency: currency,
		Created:  DateTimeUTC(),
	}
}

// Gems returns the total amount of gems.
func (payment *PayPalPayment) Gems() int {
	amount, err := strconv.ParseFloat(payment.Amount, 64)

	if err != nil {
		return 0
	}

	return int(amount)
}

// User returns the user who made the payment.
func (payment *PayPalPayment) User() *User {
	user, _ := GetUser(payment.UserID)
	return user
}

// Save saves the paypal payment in the database.
func (payment *PayPalPayment) Save() {
	DB.Set("PayPalPayment", payment.ID, payment)
}

// StreamPayPalPayments returns a stream of all paypal payments.
func StreamPayPalPayments() <-chan *PayPalPayment {
	channel := make(chan *PayPalPayment, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("PayPalPayment") {
			channel <- obj.(*PayPalPayment)
		}

		close(channel)
	}()

	return channel
}

// AllPayPalPayments returns a slice of all paypal payments.
func AllPayPalPayments() ([]*PayPalPayment, error) {
	all := make([]*PayPalPayment, 0, DB.Collection("PayPalPayment").Count())

	for obj := range StreamPayPalPayments() {
		all = append(all, obj)
	}

	return all, nil
}

// FilterPayPalPayments filters all paypal payments by a custom function.
func FilterPayPalPayments(filter func(*PayPalPayment) bool) ([]*PayPalPayment, error) {
	var filtered []*PayPalPayment

	for obj := range StreamPayPalPayments() {
		if filter(obj) {
			filtered = append(filtered, obj)
		}
	}

	return filtered, nil
}
