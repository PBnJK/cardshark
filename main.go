// Entry-point
package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pbnjk/cardshark/bot"
)

func main() {
	if err := bot.Start(os.Getenv("DISCORD_BOT_TOKEN")); err != nil {
		log.Panicf("Could not initialize bot: %v", err)
	}
}
