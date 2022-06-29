package main

import (
	"TakeoffHomework/Data"
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const StartingCash float64 = 10000

var workingCash float64
var workingAccount *Data.Account

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	log.SetOutput(os.Stdout)

	workingCash = StartingCash

	for scanner.Scan() {

		userInput := scanner.Text()

		command := strings.Split(userInput, " ")

		switch command[0] {

		case "authorize":
			//Ignore extra args I guess?
			if len(command) >= 3 {
				err := handleAuthorize(command[1], command[2])
				if err != nil {
					log.Println(err)
				}
			} else {
				log.Println("Not enough arguments to complete command")
			}

		case "withdraw":
			//Ignore extra args I guess?
			if len(command) >= 2 {
				if isAccountAuthorized(workingAccount) {
					resetAuthTime(workingAccount)
					amount, err := strconv.Atoi(command[1])
					if err != nil {
						log.Println(err)
					}
					if amount%20 > 0 {
						log.Println("Withdrawal amount must be divisible evenly by 20")
					} else {
						floatAmount := float64(amount)
						err = handleWithdraw(floatAmount, workingAccount)
						if err != nil {
							log.Println(err)
						}
					}
				}
			}

		case "deposit":
			//Ignore extra args I guess?
			if len(command) >= 2 {
				if isAccountAuthorized(workingAccount) {
					resetAuthTime(workingAccount)
					amount, err := strconv.ParseFloat(command[1], 64)
					if err != nil {
						log.Println(err)
					}
					handleDeposit(amount, workingAccount)
					addHistoryRecord(workingAccount, amount)
				}
			}

		case "balance":
			if isAccountAuthorized(workingAccount) {
				resetAuthTime(workingAccount)
				handleBalance(workingAccount)
			}

		case "history":
			handleHistory(workingAccount)

		case "logout":
			handleLogout(workingAccount)

		case "help":
			handleHelp()

		case "end":
			syscall.Exit(0)

		default:
			log.Println(fmt.Sprintf("Unrecognized command %s", command[0]))
			handleHelp()
		}
	}

}

func resetAuthTime(account *Data.Account) {
	account.AuthorizationTime = time.Now()
}

func isAccountAuthorized(account *Data.Account) bool {
	if account.AuthorizationTime.Add(2 * time.Minute).Before(time.Now()) {
		account.AuthorizationTime = time.Unix(0, 0)
		log.Println("Authorization Required")
		return false
	}

	return true
}

func handleAuthorize(accountId string, pin string) error {
	//Check to see if we know about the account
	if accountToAuthorize, ok := Data.Accounts[accountId]; ok {
		if accountToAuthorize.Pin == pin {
			accountToAuthorize.AuthorizationTime = time.Now()
			workingAccount = accountToAuthorize
			log.Println(fmt.Sprintf("%s successfully authorized", accountToAuthorize.AccountId))
		} else {
			return errors.New("authorization failed")
		}
	} else {
		return errors.New("authorization failed")
	}
	return nil
}

func handleWithdraw(value float64, account *Data.Account) error {
	//I'm going to refuse to dispense money if their account is at 0 as well because otherwise that would be cruel...
	if account.Balance <= 0 {
		return errors.New("your account is overdrawn! You may not make withdrawals at this time")
	} else {
		if workingCash == 0 {
			return errors.New("unable to process your withdrawal at this time")
		} else if workingCash < value {
			//How do I adjust value here best?
			twentiesLeft := workingCash / 20
			twentiesLeft = math.Trunc(twentiesLeft)
			value = twentiesLeft * 20
			log.Println("unable to dispense full amount requested at this time")
		}
		account.Balance = account.Balance - value
		log.Println(fmt.Sprintf("Amount dispensed: $%.2f", value))
		if account.Balance < 0 {
			account.Balance = account.Balance - 5
			log.Println(fmt.Sprintf("You have been charged an overdraft fee of $5. Current balance: $%.2f", account.Balance))
		}

		addHistoryRecord(workingAccount, -value)

	}
	return nil
}

func handleDeposit(amount float64, account *Data.Account) {
	account.Balance = account.Balance + amount
	workingCash = workingCash + amount
	log.Println(fmt.Sprintf("Current balance: %.2f", account.Balance))
}

func handleBalance(account *Data.Account) {
	log.Println(fmt.Sprintf("Current balance: %.2f", account.Balance))
}

func addHistoryRecord(account *Data.Account, amount float64) {
	account.AccountHistory = append(account.AccountHistory, fmt.Sprintf("%s %.2f %.2f",
		time.Now().Format(time.RFC3339), amount, account.Balance))
}

func handleHistory(account *Data.Account) {
	if len(account.AccountHistory) > 0 {
		for i := len(account.AccountHistory) - 1; i >= 0; i-- {
			fmt.Println(account.AccountHistory[i])
		}
	} else {
		log.Println(fmt.Sprintf("No history found for account: %s", account.AccountId))
	}
}

func handleLogout(account *Data.Account) {
	if account != nil && account.AccountId != "" {
		//Doesn't matter if the account in no longer authorized so no need to check
		account.AuthorizationTime = time.Unix(0, 0)
		log.Println(fmt.Sprintf("Account: %s logged out.", account.AccountId))
		//Need to reset working account here
		workingAccount = &Data.Account{}
	} else {
		log.Println("No account is currently authorized")
	}
}

func handleHelp() {
	//Not using the logger here because it puts ugly timestamps in front of the lines
	fmt.Fprintln(os.Stdout, "Welcome to my example ATM program, you may use the following commands: ")
	fmt.Fprintln(os.Stdout, "authorize <account_id> <pin>")
	fmt.Fprintln(os.Stdout, "Authorizes an account locally until they are logged out. "+
		"Will be logged out if there is no activity for 2 minutes")
	fmt.Fprintln(os.Stdout, "-------------------------------------------------------------------------")
	fmt.Fprintln(os.Stdout, "withdraw <value>")
	fmt.Fprintln(os.Stdout, "Removes value from the authorized account. Must be a multiple of 20")
	fmt.Fprintln(os.Stdout, "-------------------------------------------------------------------------")
	fmt.Fprintln(os.Stdout, "deposit <value>")
	fmt.Fprintln(os.Stdout, "Adds value to the authorized account. The deposited amount does not need to be a multiple of 20.")
	fmt.Fprintln(os.Stdout, "-------------------------------------------------------------------------")
	fmt.Fprintln(os.Stdout, "balance")
	fmt.Fprintln(os.Stdout, "Returns the account’s current balance.")
	fmt.Fprintln(os.Stdout, "-------------------------------------------------------------------------")
	fmt.Fprintln(os.Stdout, "history")
	fmt.Fprintln(os.Stdout, "Returns the account’s transaction history.")
	fmt.Fprintln(os.Stdout, "-------------------------------------------------------------------------")
	fmt.Fprintln(os.Stdout, "logout")
	fmt.Fprintln(os.Stdout, "Removes authorization from the current account")
	fmt.Fprintln(os.Stdout, "-------------------------------------------------------------------------")
	fmt.Fprintln(os.Stdout, "end")
	fmt.Fprintln(os.Stdout, "Ends the example ATM program, thanks for trying it out!")
	fmt.Fprintln(os.Stdout, "-------------------------------------------------------------------------")
}
