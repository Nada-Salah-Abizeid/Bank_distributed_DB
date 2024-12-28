
package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
)

type Command struct {
	Instruction   string
	Acc_no        int
	NId           int
	AccountHolder string
	Balance       float64
}

func InsertAccount(account_no int, nationalId int, accountHolder string, balance float64, conn net.Conn) {
	myStruct := Command{
		Instruction:   "insert",
		Acc_no:        account_no,
		NId:           nationalId,
		AccountHolder: accountHolder,
		Balance:       balance,
	}

	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(myStruct)
	if err != nil {
		log.Fatal("Failed to encode struct:", err)
	}

}

func UpdateAccountBalance(accountNumber int, newBalance float64, conn net.Conn) {
	myStruct := Command{
		Instruction:   "updateB",
		Acc_no:        accountNumber,
		NId:           0,
		AccountHolder: "",
		Balance:       newBalance,
	}
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(myStruct)
	if err != nil {
		log.Fatal("Failed to encode struct:", err)
	}
}

func UpdateAccountInfo(accountNumber int, accountHolder string, conn net.Conn) {
	myStruct := Command{
		Instruction:   "updateI",
		Acc_no:        accountNumber,
		NId:           0,
		AccountHolder: accountHolder,
		Balance:       0.0,
	}
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(myStruct)
	if err != nil {
		log.Fatal("Failed to encode struct:", err)
	}
}

func DeleteAccount(accountNumber int, conn net.Conn) {
	myStruct := Command{
		Instruction:   "delete",
		Acc_no:        accountNumber,
		NId:           0,
		AccountHolder: "",
		Balance:       0.0,
	}
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(myStruct)
	if err != nil {
		log.Fatal("Failed to encode struct:", err)
	}
}

func QueryAccount(accountNumber int, conn net.Conn) (string, float64) {
	myStruct := Command{
		Instruction:   "queryAccount",
		Acc_no:        accountNumber,
		NId:           0,
		AccountHolder: "",
		Balance:       0.0,
	}
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(myStruct)
	if err != nil {
		log.Fatal("Failed to encode struct:", err)
	}
	decoder := gob.NewDecoder(conn)
	//var myStruct Command
	err = decoder.Decode(&myStruct)
	if err != nil {
		log.Println("Failed to decode struct:", err)
	}

	return myStruct.AccountHolder, myStruct.Balance
}

func Leave(conn net.Conn) {
	myStruct := Command{
		Instruction:   "leaving",
		Acc_no:        0,
		NId:           0,
		AccountHolder: "",
		Balance:       0.0,
	}
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(myStruct)
	if err != nil {
		log.Fatal("Failed to encode struct:", err)
	}
}

func main() {

	// Connect to the server
	conn, err := net.Dial("tcp", "192.168.1.32:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("Connected!")

	var command int
	for {
		fmt.Println("Press 1 for creading an account.")
		fmt.Println("Press 2 for removing an account.")
		fmt.Println("Press 3 for updating an account balance.")
		fmt.Println("Press 4 for updating an account holder.")
		fmt.Println("Press 5 for searching by account no. an account.")
		fmt.Println("Press 6 to exit.")

		n, err := fmt.Scanln(&command)
		if n < 1 || err != nil {
			panic(err)
		}

		switch command {
		case 1:
			var account_no int
			var national_id int
			var Name string
			var balance float64
			fmt.Println("Please Enter your account number.")
			_, err = fmt.Scanln(&account_no)
			if err != nil {
				panic(err)
			}
			fmt.Println("Please Enter your national_id.")
			_, err = fmt.Scanln(&national_id)
			if err != nil {
				panic(err)
			}
			fmt.Println("Please Enter your Name .")
			_, err = fmt.Scanln(&Name)
			if err != nil {
				panic(err)
			}
			fmt.Println("Please Enter your balance .")
			_, err = fmt.Scanln(&balance)
			if err != nil {
				panic(err)
			}

			InsertAccount(account_no, national_id, Name, balance, conn)
			fmt.Println("Done")
		case 2:
			fmt.Println("Please Enter the account number.")
			var account_no int
			_, err = fmt.Scanln(&account_no)
			if err != nil {
				panic(err)
			}
			DeleteAccount(account_no, conn)
			fmt.Println("Done")

		case 3:
			fmt.Println("Please Enter the account number.")
			var account_no int
			_, err = fmt.Scanln(&account_no)
			if err != nil {
				panic(err)
			}
			fmt.Println("Please Enter the new balance.")
			var new_balance float64
			_, err = fmt.Scanln(&new_balance)
			if err != nil {
				panic(err)
			}
			UpdateAccountBalance(account_no, new_balance, conn)
			fmt.Println("Done")

		case 4:
			fmt.Println("Please Enter the account number.")
			var account_no int
			_, err = fmt.Scanln(&account_no)
			if err != nil {
				panic(err)
			}
			fmt.Println("Please Enter the name of the new holder.")
			var new_holder string
			_, err = fmt.Scanln(&new_holder)
			if err != nil {
				panic(err)
			}
			UpdateAccountInfo(account_no, new_holder, conn)
			fmt.Println("Done")

		case 5:
			fmt.Println("Please Enter the account number.")
			var account_no int
			_, err = fmt.Scanln(&account_no)
			if err != nil {
				panic(err)
			}
			accountHolder, balance := QueryAccount(account_no, conn)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Account holder: %s, Balance: %f\n", accountHolder, balance)
			fmt.Println("Done")

		case 6:
			Leave(conn)
			os.Exit(2)
		default:
			fmt.Println("Invalid choice")
		}
	}

}
