package config

import (
	"fmt"
	"testing"
)

func TestNewCofig(t *testing.T) {
	config, err := NewConfig("./../config.yml")
	if err != nil {
		t.Fail()
	}

	fmt.Println(config)
}
