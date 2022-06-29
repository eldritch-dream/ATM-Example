package Data

import "time"

type Account struct {
	AccountId         string
	Pin               string
	Balance           float64
	AuthorizationTime time.Time
	AccountHistory    []string
}

var account1 = Account{
	AccountId:      "2859459814",
	Pin:            "7386",
	Balance:        10.24,
	AccountHistory: []string{},
}

var account2 = Account{
	AccountId:      "1434597300",
	Pin:            "4557",
	Balance:        90000.55,
	AccountHistory: []string{},
}

var account3 = Account{
	AccountId:      "7089382418",
	Pin:            "0075",
	Balance:        0.00,
	AccountHistory: []string{},
}

var account4 = Account{
	AccountId:      "2001377812",
	Pin:            "5950",
	Balance:        60.00,
	AccountHistory: []string{},
}

var account5 = Account{
	AccountId:      "1",
	Pin:            "1",
	Balance:        100.00,
	AccountHistory: []string{},
}

var Accounts = map[string]*Account{account1.AccountId: &account1,
	account2.AccountId: &account2,
	account3.AccountId: &account3,
	account4.AccountId: &account4,
	account5.AccountId: &account5}
