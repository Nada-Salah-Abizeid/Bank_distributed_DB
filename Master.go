package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// type BankDatabase struct {
// 	db *sql.DB
// }

type Command struct {
	Instruction   string
	Acc_no        int
	NId           int
	AccountHolder string
	Balance       float64
}

// func (db *BankDatabase) CreateDatabase(dbName string) error {
// 	_, err := db.db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
// 	return err
// }

// func (db *BankDatabase) CreateTables() error {
// 	_, err := db.db.Exec(`
// 		CREATE TABLE IF NOT EXISTS accounts (
// 			account_number INT(255) AUTO_INCREMENT,
// 			national_id INT NOT NULL ,
// 			account_holder VARCHAR(255) NOT NULL,
// 			balance DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
// 			PRIMARY KEY (account_number)
// 		);

// 	`)
// 	return err
// }

func main() {

	var SlaveNum int = 1
	conns := []net.Conn{}

	//Listen for incoming connections
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	// Accept connections
	for {

		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		conns = append(conns, conn)
		defer conn.Close()

		if len(conns) == SlaveNum {
			break
		}

	}
	fmt.Println("Connected to slave!")

	// Get a database handle.
	db, err := sql.Open("mysql", "root:rootroot@tcp(localhost:3306)/testing")
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected to database!")

	//BankDatabase := &BankDatabase{db}
	for {
		decoder := gob.NewDecoder(conns[0])
		var myStruct Command
		err = decoder.Decode(&myStruct)
		if err != nil {
			log.Println("Failed to decode struct:", err)
			return
		}
		if myStruct.Instruction == "insert" {
			stmt, err := db.Prepare("INSERT INTO accounts (account_number, national_id, account_holder, balance) VALUES (?,?,?,?)")
			if err != nil {
				panic(err)
			}
			defer stmt.Close()
			_, err = stmt.Exec(myStruct.Acc_no, myStruct.NId, myStruct.AccountHolder, myStruct.Balance)
			if err != nil {
				panic(err)
			}
			fmt.Println("Done")
		} else if myStruct.Instruction == "updateB" {
			stmt, err := db.Prepare("UPDATE accounts SET balance =? WHERE account_number =?")
			if err != nil {
				panic(err)
			}
			defer stmt.Close()
			_, err = stmt.Exec(myStruct.Balance, myStruct.Acc_no)
			if err != nil {
				panic(err)
			}

		} else if myStruct.Instruction == "updateI" {
			stmt, err := db.Prepare("UPDATE accounts SET account_holder =? WHERE account_number =?")
			if err != nil {
				panic(err)
			}
			defer stmt.Close()
			_, err = stmt.Exec(myStruct.AccountHolder, myStruct.Acc_no)
			if err != nil {
				panic(err)
			}
		} else if myStruct.Instruction == "delete" {
			stmt, err := db.Prepare("DELETE FROM accounts WHERE account_number =?")
			if err != nil {
				panic(err)
			}
			defer stmt.Close()
			_, err = stmt.Exec(myStruct.Acc_no)
			if err != nil {
				panic(err)
			}
		} else if myStruct.Instruction == "queryAccount" {
			stmt, err := db.Prepare("SELECT account_holder, balance FROM accounts WHERE account_number =?")
			if err != nil {
				panic(err)
			}
			defer stmt.Close()
			err = stmt.QueryRow(myStruct.Acc_no).Scan(&myStruct.AccountHolder, &myStruct.Balance)
			if err != nil {
				panic(err)
			}
			encoder := gob.NewEncoder(conns[0])
			err = encoder.Encode(myStruct)
			if err != nil {
				log.Fatal("Failed to encode struct:", err)
			}
		} else {
			os.Exit(2)
			//break
		}

	}

}
