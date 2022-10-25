package main

import "fmt"

type CheckingAccount struct {
	holder string
	agencyNumber int
	accountNumber int
	balance float64
}

func main() {
	account1 := CheckingAccount{holder: "John",
		agencyNumber: 589, accountNumber: 123456, balance: 125.5}

	account2 := CheckingAccount{"Erika", 222, 111222, 200}

	fmt.Println(account1)
	fmt.Println(account2)
}