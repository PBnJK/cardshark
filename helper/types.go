package helper

import (
	dgo "github.com/bwmarrin/discordgo"
)

type (
	InteractionFunc func(s *dgo.Session, i *dgo.InteractionCreate)
	HandlerMap      map[string]InteractionFunc
)
