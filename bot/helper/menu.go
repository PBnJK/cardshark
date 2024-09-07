package helper

import (
	dgo "github.com/bwmarrin/discordgo"
)

type Option struct {
	name  string
	id    string
	style dgo.ButtonStyle
}

func (o Option) AsButton() dgo.Button {
	return dgo.Button{
		Label:    o.name,
		CustomID: o.id,
		Style:    o.style,
		Disabled: false,
	}
}

type SelectMenu struct {
	prompt      string
	id          string
	placeholder string
	options     []Option
}

func NewSelectMenu(prompt, id, placeholder string) *SelectMenu {
	return &SelectMenu{
		prompt:  prompt,
		id:      id,
		options: make([]Option, 0),
	}
}

func (sm *SelectMenu) AddOption(name, id string) {
	var style dgo.ButtonStyle
	if len(sm.options) == 0 {
		style = dgo.PrimaryButton
	} else {
		style = dgo.SecondaryButton
	}

	sm.options = append(sm.options, Option{
		name:  name,
		id:    id,
		style: style,
	})
}

func (sm *SelectMenu) OptionsAsMessageComponent() []dgo.MessageComponent {
	options := make([]dgo.MessageComponent, len(sm.options))
	for i, o := range sm.options {
		options[i] = o.AsButton()
	}

	return options
}

func (sm *SelectMenu) AsInteractionResponse() *dgo.InteractionResponse {
	return &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: sm.prompt,
			Flags:   dgo.MessageFlagsEphemeral, // TODO: Make this option customizable
			Components: []dgo.MessageComponent{
				dgo.ActionsRow{
					Components: sm.OptionsAsMessageComponent(),
				},
			},
		},
	}
}
