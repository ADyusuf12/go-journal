package main

import (
	"errors"
	"fmt"
)

// PaymentProcessor defines the contract for charging a customer.
type PaymentProcessor interface {
	ProcessPayment(amountCents int64, customerID string) (string, error)
}

// PaystackGateway handles production payments.
type PaystackGateway struct {
	APIKey string
}

func (p *PaystackGateway) ProcessPayment(amountCents int64, customerID string) (string, error) {
	if amountCents <= 0 {
		return "", errors.New("invalid payment amount")
	}

	// Simulated API Call
	txnID := fmt.Sprintf("pstk_tx_%d", amountCents)
	return txnID, nil
}

func (p *PaystackGateway) Refund(txnID string) error {
	fmt.Printf("Refunding transaction %s via Paystack API...\n", txnID)
	return nil
}

func ProcessRefundIfSupported(processor PaymentProcessor, txnID string) {
	// Type assertion: "Is processor internally holding a *PaystackGateway?"
	paystack, ok := processor.(*PaystackGateway)
	if !ok {
		fmt.Println("Refund skipped: Provided processor does not support Paystack refunds.")
		return
	}

	// We now have access to Paystack-specific methods!
	paystack.Refund(txnID)
}

func InspectGateway(processor PaymentProcessor) {
	switch v := processor.(type) {
	case *PaystackGateway:
		fmt.Printf("Gateway Type: Paystack (API Key configured: %s)\n", v.APIKey)
	case MockGateway:
		fmt.Println("Gateway Type: Local Mock (No network activity)")
	default:
		fmt.Println("Gateway Type: Unknown processor implementation")
	}
}

// MockGateway handles testing or local dev payments.
type MockGateway struct{}

func (m MockGateway) ProcessPayment(amountCents int64, customerID string) (string, error) {
	return "mock_tx_12345", nil
}

// CheckoutService depends entirely on the PaymentProcessor abstraction.
type CheckoutService struct {
	Processor PaymentProcessor
}

func (cs CheckoutService) CompleteOrder(amountCents int64, customerID string) error {
	txnID, err := cs.Processor.ProcessPayment(amountCents, customerID)
	if err != nil {
		return fmt.Errorf("checkout failed: %w", err)
	}

	fmt.Printf("Order succeeded! Transaction ID: %s\n", txnID)
	return nil
}

func main() {
	// 1. Production setup using Paystack
	paystack := &PaystackGateway{APIKey: "sk_live_12345"}
	prodService := CheckoutService{Processor: paystack}

	fmt.Println("--- Running Production Checkout ---")
	err := prodService.CompleteOrder(50000, "cust_001")
	if err != nil {
		fmt.Println("Error:", err)
	}

	// 2. Local/Test setup using Mock
	mock := MockGateway{}
	testService := CheckoutService{Processor: mock}

	fmt.Println("\n--- Running Test Checkout ---")
	err = testService.CompleteOrder(50000, "cust_001")
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("\n--- Testing Type Assertions ---")
	ProcessRefundIfSupported(paystack, "pstk_tx_50000") // Should succeed!
	ProcessRefundIfSupported(mock, "mock_tx_12345")     // Should fail gracefully!

	fmt.Println("\n--- Testing Type Switches ---")
	InspectGateway(paystack)
	InspectGateway(mock)
}
