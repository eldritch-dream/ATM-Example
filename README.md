# Example ATM Program

This program aims to demonstrate the basic features of an atm that can only dispense $20 bills and charge overdraft 
fees if you would over draw your account.

## Supported Commands

You can run the program and call the help command or refer to the section below

| Command                        | Description                                                                                                       |
|--------------------------------|-------------------------------------------------------------------------------------------------------------------|
| authorize <account_id\> <pin\> | Authorizes an account locally until they are logged out. Will be logged out if there is no activity for 2 minutes |
| withdraw <value\>              | Removes value from the authorized account. Must be a multiple of 20                                               |
| deposit <value\>               | Adds value to the authorized account. The deposited amount does not need to be a multiple of 20.                  |
| balance                        | Returns the account’s current balance.                                                                            |
| history                        | Returns the account’s transaction history.                                                                        | 
| logout                         | Removes authorization from the current account                                                                    | 
| end                            | Ends the example ATM program, thanks for trying it out!                                                           |
