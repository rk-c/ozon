package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type dbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

type ozonClient struct {
	Id  string `yaml:"id"`
	Key string `yaml:"key"`
}

type googleClient struct {
	File   string `yaml:"api_file"`
	Folder string `yaml:"folder"`
}

type config struct {
	Db              dbConfig     `yaml:"db"`
	Ozon            ozonClient   `yaml:"ozon_client"`
	OzonOrdersDepth string       `yaml:"ozon_orders_depth"`
	Google          googleClient `yaml:"google_client"`
}

type ReadConfigFileError struct {
	filename string
	err      error
}

func NewReadConfigFileError(filename string, err error) ReadConfigFileError {
	return ReadConfigFileError{
		filename: filename,
		err:      err,
	}
}

func (e ReadConfigFileError) Error() string {
	return fmt.Sprintf("couldn`t read config file %s: %s", e.filename, e.err.Error())
}

type UnmarshalConfigError struct {
	err error
}

func NewUnmarshalConfigError(err error) UnmarshalConfigError {
	return UnmarshalConfigError{
		err: err,
	}
}

func (e UnmarshalConfigError) Error() string {
	return fmt.Sprintf("couldn`t unmarshal config: %s", e.err.Error())
}

func configure(fileName string) (*config, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, NewReadConfigFileError(fileName, err)
	}

	var cnf config

	if errUnmarshal := yaml.Unmarshal(data, &cnf); errUnmarshal != nil {
		return nil, NewUnmarshalConfigError(errUnmarshal)
	}

	return &cnf, nil
}
