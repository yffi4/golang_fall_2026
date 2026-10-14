package main

import (
	"errors"
	"fmt"
)

type BankAccount struct {
	Owner   string
	balance float64
}

func (a *BankAccount) Deposit(amount float64) {
	a.balance += amount
}
func (a *BankAccount) Withdraw(amount float64) error {
	if amount > a.balance {
		return errors.New("not enough funds")

	}
	a.balance -= amount
	return nil
}

func (a BankAccount) Balance() float64 {
	return a.balance
}

func main() {
	accounts := []BankAccount{
		{Owner: "Alice", balance: 1000},
		{Owner: "Bob", balance: 2000},
		{Owner: "Carol", balance: 500},
	}

	for i := range accounts {
		accounts[i].Deposit(accounts[i].Balance() * 0.05)
	}

	for _, acc := range accounts {
		fmt.Printf("%s: %.2f\n", acc.Owner, acc.Balance())
	}

	err := accounts[2].Withdraw(10000)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
