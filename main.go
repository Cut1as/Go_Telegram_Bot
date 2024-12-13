package main

import (
	"fmt"
	"os"

	"github.com/mymmrac/telego"
)

func main() {

	botToken := os.Getenv("токен")

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//updates, _ := bot.UpdatesViaLongPolling(nil)

	defer bot.StopLongPolling()

}
