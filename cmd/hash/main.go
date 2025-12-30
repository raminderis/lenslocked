package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {
	// Step1 Think about the Key
	secretKeyForHash := "MyKeyAsStringWillGoAsBytesInHashFunc"
	// Step2 Think about the hash function
	hash_func := hmac.New(sha256.New, []byte(secretKeyForHash))
	// Step3 Think about the password
	passwordToBeHashed := "My super secret password"
	// Step4 Think about writing the password to the hash function
	hash_func.Write([]byte(passwordToBeHashed))
	// Step5 Think about reading the hashed password from the hash function
	result := hash_func.Sum(nil)
	// Step6 Think about encoding to string to be able to print the hashed valued
	result_encoded_toPrint := hex.EncodeToString(result)
	fmt.Println(result_encoded_toPrint)
	// Step7 Also may think about Reseting the hash function
	hash_func.Reset()
}
