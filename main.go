package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	storedHash := "$2a$10$BnvrDFy5KlF5OyFJuYMiEePGjKat2dBL7LsGaw8sLGtFDSUS04klC"
	testPassword := "securepassword123"

	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(testPassword))
	if err != nil {
		fmt.Printf("Password does not match: %v\n", err)
		return
	}

	fmt.Println("Password matches!")
}
