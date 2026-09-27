package main

import (
	"errors"
	"fmt"
)

func getMarketplaceCommission(amount int) (int, error) {
	if amount < 0 {
		return 0, errors.New("amount cannot be negative")
	}
	commission := amount * 10 / 100 // 10% commission
	return commission, nil
}

func main() {
	commission, err := getMarketplaceCommission(-230)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Marketplace Commission:", commission)
}
