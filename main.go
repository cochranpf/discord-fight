package main

import (
	"discord-fight/bot"
	"discord-fight/config"
	"fmt"
)

func main() {
	if err := config.ReadConfig(); err != nil {
		fmt.Println(err.Error())
		return
	}

	//start bot
	bot.Start()

	<-make(chan struct{})
	return
}
