package config

import (
	"linker-fan/gal-anonim-server/server/utils"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`
	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"postgres"`
	Jwt struct {
		TokenSecret string `yaml:"tokenSecret"`
		ExpTime     int64  `yaml:"expTime"`
		Issuer      string `yaml:"issuer"`
	} `yaml:"jwt"`
	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
	FileStorage struct {
		Endpoint        string `yaml:"endpoint"`
		AccessKeyID     string `yaml:"accessKeyID"`
		SecretAccessKey string `yaml:"secretAccessKey"`
		Secure          bool   `yaml:"secure"`
	} `yaml:"filestorage"`
}

func NewConfig(path string) (*Config, error) {
	config := &Config{}

	err := utils.ValidatePath(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
