package helper

import (
	dgo "github.com/bwmarrin/discordgo"
)

type (
	OptionCallback func(s *dgo.Session, i *dgo.InteractionCreate) *dgo.InteractionResponse
	OptionFollowup func(s *dgo.Session, i *dgo.InteractionCreate)
)

type Option struct {
	callback OptionCallback
	followup OptionFollowup
	name     string
	id       string
	style    dgo.ButtonStyle
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
	options     []*Option

	Flags dgo.MessageFlags
}

var optionMap map[string]*Option

func init() {
	optionMap = make(map[string]*Option)
}

func NewSelectMenu(prompt, id, placeholder string) *SelectMenu {
	return &SelectMenu{
		prompt:  prompt,
		id:      id,
		options: make([]*Option, 0),
		Flags:   dgo.MessageFlagsEphemeral,
	}
}

func (sm *SelectMenu) AddOption(name, id string, callback OptionCallback, followup OptionFollowup) {
	option := &Option{
		callback: callback,
		followup: followup,
		name:     name,
		id:       id,
		style:    dgo.SecondaryButton,
	}

	optionMap[option.id] = option

	sm.options = append(sm.options, option)
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
			Flags:   sm.Flags,
			Components: []dgo.MessageComponent{
				dgo.ActionsRow{
					Components: sm.OptionsAsMessageComponent(),
				},
			},
		},
	}
}

func HandleSelect(m dgo.MessageComponentInteractionData) (OptionCallback, OptionFollowup, bool) {
	if opt, ok := optionMap[m.CustomID]; ok {
		return opt.callback, opt.followup, true
	}

	return nil, nil, false
}
