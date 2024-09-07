// Entry-point
package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pbnjk/hh/bot"
)

func main() {
	err := bot.New(os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		bot.Quit()
		log.Panicf("Could not initialize bot: %v", err)
	}

	bot.Run()
}
