package auth

import (
	"fmt"
	"log"
	"testing"
)

func TestGenerateJWTToken(t *testing.T) {
	log.Println("Test GenerateJWTToken()")
	tokenString, expTime, err := GenerateJWTToken("someUsername", 1, false)
	if err != nil {
		t.Fail()
	}

	fmt.Printf("TokenString: %s\n", tokenString)
	fmt.Printf("ExpTime: %s\n", expTime.String())

}

func TestIsTokenValid(t *testing.T) {
	log.Println("Test IsTokenValid()")
	tokenString, _, err := GenerateJWTToken("someUsername", 1, false)
	if err != nil {
		t.Fail()
	}

	username, id, isAdmin, err := IsTokenValid(tokenString)
	if err != nil {
		t.Fail()
	}

	if username != "someUsername" {
		t.Fail()
	}

	if id != 1 {
		t.Fail()
	}

	if isAdmin {
		t.Fail()
	}

	fmt.Printf("Username: %s\n", username)
	fmt.Printf("ID: %d\n", id)
	fmt.Printf("IsAdmin: %v\n", isAdmin)

}
