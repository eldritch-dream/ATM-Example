package main

import (
	"TakeoffHomework/Data"
	"strings"
	"testing"
	"time"
)

func TestResetAuthTime(t *testing.T) {
	epochTime := time.Unix(0, 0)
	testAccount := Data.Account{AuthorizationTime: epochTime}

	resetAuthTime(&testAccount)
	if testAccount.AuthorizationTime == epochTime {
		t.Errorf("Authorization time should have been close to time.Now() but was epoch time instead")
	}
}

func TestIsAccountAuthorized(t *testing.T) {
	testAccount := Data.Account{}
	accountAuthorized := isAccountAuthorized(&testAccount)
	if accountAuthorized {
		t.Errorf("Empty account should NOT be authorized")
	}

	testAccount = Data.Account{AuthorizationTime: time.Now()}
	accountAuthorized = isAccountAuthorized(&testAccount)
	if !accountAuthorized {
		t.Errorf("Account with auth time of now should be authorized")
	}
}

func TestHandleAuthorize(t *testing.T) {
	err := handleAuthorize("2", "2")
	if err == nil {
		t.Errorf("Account 2 with pin 2 should not exist")
	}

	err = handleAuthorize("1", "2")
	if err == nil {
		t.Errorf("Account 1 does not have pin 2")
	}

	err = handleAuthorize("1", "1")
	if err != nil {
		t.Errorf("Account 1 has pin 1 and should not error")
	}
}

func TestHandleWithdraw(t *testing.T) {
	testAccount := &Data.Account{Balance: 0}
	err := handleWithdraw(0, testAccount)
	if err == nil {
		t.Errorf("Withdraw on 0 should cause an error")
	}

	testAccount = &Data.Account{Balance: 10}
	workingCash = 0
	err = handleWithdraw(0, testAccount)
	if err == nil {
		t.Errorf("Withdraw action when ATM is empty should error")
	}

	testAccount = &Data.Account{Balance: 40, AccountHistory: []string{}}
	workingAccount = testAccount
	workingCash = 20
	_ = handleWithdraw(40, testAccount)
	if testAccount.Balance != 20 {
		t.Errorf("Balance should only be 20 since machine can only dispense 1 bill but was %f", testAccount.Balance)
	}

	testAccount = &Data.Account{Balance: 20}
	workingCash = 40
	_ = handleWithdraw(40, testAccount)
	if testAccount.Balance != -25 {
		t.Errorf("Balance should be -25 due to overdraft but was %f", testAccount.Balance)
	}

	testAccount = &Data.Account{Balance: 30}
	workingCash = 40
	_ = handleWithdraw(20, testAccount)
	if testAccount.Balance != 10 {
		t.Errorf("Balance should be 10 but was %f", testAccount.Balance)
	}
}

func TestHandleDeposit(t *testing.T) {
	testAccount := &Data.Account{Balance: 30}

	handleDeposit(13.37, testAccount)
	if testAccount.Balance != 43.37 {
		t.Errorf("Balance should be 43.37 but was %f", testAccount.Balance)
	}
}

func TestAddHistoryRecord(t *testing.T) {
	testAccount := &Data.Account{Balance: 9000}
	addHistoryRecord(testAccount, 0)
	if !strings.Contains(testAccount.AccountHistory[0], "9000") {
		t.Errorf("Account history should contain a trasaction with 9000 in it")
	}
}

func TestHandleLogoout(t *testing.T) {
	testAccount := &Data.Account{AccountId: "1", AuthorizationTime: time.Now()}
	handleLogout(testAccount)
	if testAccount.AuthorizationTime == time.Now() {
		t.Errorf("Object should have been reset as a result of logout")
	}
}
