package games

import (
	dgo "github.com/bwmarrin/discordgo"
)

func InitBlackjack(s *dgo.Session) []*dgo.ApplicationCommand {
	return []*dgo.ApplicationCommand{
		{
			Name:        "blackjack",
			Description: "Starts a game of Blackjack",
		},
	}
}
