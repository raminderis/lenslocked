package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Invalid number of args. hash takes 1 argument. compare takes 2.")
		return
	}
	switch os.Args[1] {
	case "hash":
		if len(os.Args) == 3 {
			passwordInPlainText := os.Args[2]
			passwordInHash, err := hash(passwordInPlainText)
			if err != nil {
				fmt.Println("Hashing failed")
				panic(err)
			}
			fmt.Printf("hash of the password %v is %v\n", passwordInPlainText, passwordInHash)
		} else {
			fmt.Println("Invalid number of args. hash takes 1 argument.")
		}
		// fmt.Println(compare(passwordInPlainText, passwordInHash))
	case "compare":
		if len(os.Args) == 4 {
			if compare(os.Args[2], os.Args[3]) {
				fmt.Println("Password is correct!")
			} else {
				fmt.Println("Password is invalid.")
			}
		} else {
			fmt.Println("Invalid number of args. compare takes 2 argument.")
		}
	default:
		fmt.Printf("Invalid command: %v. We only support hash and compare\n", os.Args[1])
	}
}

func hash(passwordInPlainText string) (string, error) {
	passwordInHash, err := bcrypt.GenerateFromPassword([]byte(passwordInPlainText), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("error  hashing the password %s: %v\n", passwordInPlainText, err)
		return "", err
	}
	return string(passwordInHash), nil

}

func compare(passwordInPlainText, passwordInHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordInHash), []byte(passwordInPlainText))
	if err != nil {
		return false
	}
	return true
}
