package utils

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateUsername(t *testing.T) {
	log.Println("Test ValidateUsername()")

	//valid username
	err := ValidateUsername("hyperxpizza")
	if err != nil {
		t.Fail()
	}

	//username longer than 30 characters
	err = ValidateUsername("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	if err == nil {
		t.Fail()
	}

	//username contains a special character
	err = ValidateUsername("a$!lolessa")
	if err == nil {
		t.Fail()
	}

	//empty string
	err = ValidateUsername("")
	if err == nil {
		t.Fail()
	}
}

func TestValidatePassword(t *testing.T) {
	log.Println("Test ValidatePassword()")

	//valid password
	err := ValidatePassword("testPassword1!")
	if err != nil {
		t.Fail()
	}

	//password too short
	err = ValidatePassword("aA1!")
	if err == nil {
		t.Fail()
	}

	//password without a number
	err = ValidatePassword("aA#aaaa")
	if err == nil {
		t.Fail()
	}

	//password without an uppercase character
	err = ValidatePassword("aaaaa#14")
	if err == nil {
		t.Fail()
	}

	//password without a special character
	err = ValidatePassword("aAaaa1234")
	if err == nil {
		t.Fail()
	}
}

func TestValidateRoomName(t *testing.T) {
	log.Println("Test ValidateRoomName()")

	//valid username
	err := ValidateRoomName("sometestroomname")
	if err != nil {
		t.Fail()
	}

	//username longer than 30 characters
	err = ValidateRoomName("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	if err == nil {
		t.Fail()
	}

	//username contains a special character
	err = ValidateRoomName("a$!lolessa")
	if err == nil {
		t.Fail()
	}

	//empty string
	err = ValidateRoomName("")
	if err == nil {
		t.Fail()
	}
}

func TestValidateEmail(t *testing.T) {
	t.Run("Test validate email", func(t *testing.T) {
		valid := ValidateEmail("hyperxpizza@domain.com")
		assert.NoError(t, valid)
	})

	t.Run("Test validate email special character", func(t *testing.T) {
		valid := ValidateEmail("hyper!pizza@domain.com")
		assert.Error(t, valid)
	})
}
