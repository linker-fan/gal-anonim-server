package utils

import (
	"errors"
	"fmt"
	"os"
	"regexp"
)

//Validates username given as an argument
//username has to match the regex or return error otherwise
//if username is valid, function shall return nil
//@author dzania
func ValidateUsername(username string) error {
	var isValid = regexp.MustCompile(`^[a-zA-Z0-9]{3,30}$`).MatchString
	if isValid(username) {
		return nil
	} else {
		return errors.New("Username not valid")
	}
}

//Validates password given as an argument
//password has to contain at least one numeric character, one special symbol, one uppercase character and be longer than 6 characters to be valid
//if passwored is not valid, function returns an error, otherwise the function shall return nil
//@author dzania
func ValidatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("Password too short")
	}
	num := `[0-9]{1}`
	az := `[a-z]{1}`
	AZ := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`
	if b, err := regexp.MatchString(num, password); !b || err != nil {
		return errors.New("Password needs at least one number")
	}
	if b, err := regexp.MatchString(az, password); !b || err != nil {
		return errors.New("Password needs at least one small character")
	}
	if b, err := regexp.MatchString(AZ, password); !b || err != nil {
		return errors.New("Password needs at leat one uppercase character")
	}
	if b, err := regexp.MatchString(symbol, password); !b || err != nil {
		return errors.New("Password needs at least on special symbol")
	}
	return nil
}

func ValidateRoomName(name string) error {
	var isValid = regexp.MustCompile(`^[a-zA-Z0-9]{3,50}$`).MatchString
	if isValid(name) {
		return nil
	} else {
		return errors.New("Room name not valid")
	}
}

func ValidatePin(pin string) error {
	if len(pin) < 4 && len(pin) > 8 {
		return errors.New("Pin not valid")
	}

	return nil
}

func ValidatePath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}

	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func ValidateEmail(e string) error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(e) {
		return fmt.Errorf("Email not valid")
	}

	return nil

}
