package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	Token     string //To store value of Token from config.json.
	BotPrefix string //To store value of BotPrefix from config.json.

	config *configStruct //To store value extracted from config.json.
)

type configStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
}

func ReadConfig() error {
	fmt.Println("Reading config file...")
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(string(file))
	if err := json.Unmarshal(file, &config); err != nil {
		fmt.Println(err.Error())
		return err
	}

	Token = config.Token
	BotPrefix = config.BotPrefix

	//If there isn't any error we will return nil.
	return nil

}
