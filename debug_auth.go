package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Test case - we know these credentials should work
	email := "admin@example.com"
	password := "securepassword123"

	// This is the password hash from the database (replace with the actual hash from your database)
	// To get the actual hash, you can query your database directly
	storedHash := "$2a$10$rFbFfz3B9wGvGn.CeUfA6..a9L5wKHVr.OBVWugS53bYJamIU54aO" // example hash

	// Test 1: Check if the password hashing algorithm is working correctly
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		log.Printf("Password verification failed: %v", err)
	} else {
		log.Println("✓ Password verification succeeded!")
	}

	// Test 2: Generate a new hash for the password and verify it
	newHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	err = bcrypt.CompareHashAndPassword(newHash, []byte(password))
	if err != nil {
		log.Printf("New hash verification failed: %v", err)
	} else {
		log.Println("✓ New hash verification succeeded!")
	}

	// Test 3: Check for case sensitivity issues
	uppercaseEmail := strings.ToUpper(email)
	if uppercaseEmail != email {
		log.Printf("Testing case sensitivity - original: %s, uppercase: %s", email, uppercaseEmail)
		// In a real implementation, you would check if your backend is case-sensitive for emails
	}

	// Test 4: Check for whitespace issues
	passwordWithSpace := " " + password + " "
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(passwordWithSpace))
	if err != nil {
		log.Printf("× Password with whitespace verification failed (expected)")
	} else {
		log.Println("! Password with whitespace verification succeeded (unexpected - check for trimming issues)")
	}

	// Test 5: Simulate the full login flow
	log.Println("\nSimulating login flow for:", email)
	log.Println("1. User submits login form with email and password")
	log.Println("2. Backend receives request with email:", email)
	log.Println("3. Backend looks up user by email in the database")
	log.Println("4. Backend compares submitted password with stored hash")

	// Simulate the hashing check
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		log.Printf("5. ✗ Password verification failed: %v", err)
		log.Println("6. Backend returns 'Invalid credentials' error")
	} else {
		log.Println("5. ✓ Password verification succeeded!")
		log.Println("6. Backend generates JWT token")
		log.Println("7. Backend returns token and user information")
	}

	// Recommendations
	fmt.Println("\nPossible issues to check:")
	fmt.Println("1. Content-Type header is set to 'application/json' in the request")
	fmt.Println("2. Request body format matches expected API format: {\"email\":\"admin@example.com\",\"password\":\"securepassword123\"}")
	fmt.Println("3. Email case sensitivity handling in the database lookup")
	fmt.Println("4. Password whitespace trimming before hashing/comparison")
	fmt.Println("5. CORS configuration is correct for your frontend")
	fmt.Println("6. API route is correctly configured at /api/v1/auth/login")
	fmt.Println("7. Network/firewall issues blocking the request")
	fmt.Println("8. Database connection issues")
}
