// hash-generator.go
package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// <-- UNESITE ŽELJENU LOZINKU ZA ADMINA OVDE
	passwordToHash := "123"

	// Heširanje lozinke
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwordToHash), 14)
	if err != nil {
		log.Fatal(err)
	}

	hashedPassword := string(bytes)
	fmt.Println("Vaša heširana lozinka je:")
	fmt.Println(hashedPassword)
}
