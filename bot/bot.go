// TODO
package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	dgo "github.com/bwmarrin/discordgo"

	"github.com/pbnjk/hh/bot/helper"
)

var (
	session        *dgo.Session
	registeredCmds []*dgo.ApplicationCommand

	chooseGameMenu *dgo.InteractionResponse
)

func constructChooseGameMenu() *dgo.InteractionResponse {
	sm := helper.NewSelectMenu("Selecione o jogo que quer jogar", "basic", "Type stuff...?")

	sm.AddOption("21", "o1", func(s *dgo.Session, i *dgo.InteractionCreate) *dgo.InteractionResponse {
		return &dgo.InteractionResponse{
			Type: dgo.InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "You chose the first one",
				Flags:   dgo.MessageFlagsEphemeral,
			},
		}
	}, nil)

	sm.AddOption("Option 2", "o2", func(s *dgo.Session, i *dgo.InteractionCreate) *dgo.InteractionResponse {
		return &dgo.InteractionResponse{
			Type: dgo.InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "You chose the second one",
				Flags:   dgo.MessageFlagsEphemeral,
			},
		}
	}, nil)

	sm.AddOption("Option 3", "o3", func(s *dgo.Session, i *dgo.InteractionCreate) *dgo.InteractionResponse {
		return &dgo.InteractionResponse{
			Type: dgo.InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "You chose the third one",
				Flags:   dgo.MessageFlagsEphemeral,
			},
		}
	}, nil)

	return sm.AsInteractionResponse()
}

func init() {
	chooseGameMenu = constructChooseGameMenu()
}

// Slash commands that this bot is able to respond to
var commands = []*dgo.ApplicationCommand{
	{
		Name:        "basic",
		Description: "A basic command",
	},
}

func handleBasicCmd(_ *dgo.Session, i *dgo.InteractionCreate) {
	err := session.InteractionRespond(i.Interaction, chooseGameMenu)
	if err != nil {
		log.Panicf("Panicked with %v!!\n", err)
	}

	fmt.Println("Done!")
}

func handleResponseCmd(s *dgo.Session, i *dgo.InteractionCreate) {
	if cb, fw, ok := helper.HandleSelect(i.MessageComponentData()); ok {
		err := session.InteractionRespond(i.Interaction, cb(s, i))
		if err != nil {
			log.Panicf("Panicked with %v!!\n", err)
		}

		if fw != nil {
			fw(s, i)
		}
	}
}

// Syncs commands across runs, deleting old commands and creating or editing
// new ones
func syncCommands() error {
	existingCommands, err := session.ApplicationCommands(session.State.User.ID, "")
	if err != nil {
		return fmt.Errorf("Failed to fetch commands for guild %s: %v", "", err)
	}

	desiredMap := make(map[string]*dgo.ApplicationCommand)
	for _, cmd := range commands {
		desiredMap[cmd.Name] = cmd
	}

	existingMap := make(map[string]*dgo.ApplicationCommand)
	for _, cmd := range existingCommands {
		existingMap[cmd.Name] = cmd
	}

	// Delete commands not in the desired list
	for _, cmd := range existingCommands {
		if _, found := desiredMap[cmd.Name]; !found {
			err := session.ApplicationCommandDelete(session.State.User.ID, "", cmd.ID)
			if err != nil {
				return fmt.Errorf("Failed to delete command %s: %v\n", cmd.Name, err)
			} else {
				fmt.Printf("Successfully deleted command %s\n", cmd.Name)
			}
		}
	}

	// Create or update existing commands
	for _, cmd := range commands {
		if existingCmd, found := existingMap[cmd.Name]; found {
			// Edit existing command
			_, err := session.ApplicationCommandEdit(session.State.User.ID, "", existingCmd.ID, cmd)
			if err != nil {
				return fmt.Errorf("Failed to edit command %s: %v\n", cmd.Name, err)
			} else {
				fmt.Printf("Successfully edited command %s\n", cmd.Name)
			}
		} else {
			// Create new command
			_, err := session.ApplicationCommandCreate(session.State.User.ID, "", cmd)
			if err != nil {
				return fmt.Errorf("Failed to create command %s: %v", cmd.Name, err)
			} else {
				fmt.Printf("Successfully created command %s", cmd.Name)
			}
		}
	}

	return nil
}

// Creates the discordgo session from a token
func createSession(token string) error {
	var err error

	session, err = dgo.New("Bot " + token)
	if err != nil {
		return err
	}

	return nil
}

// Creates the necessary handlers
func createHandlers() {
	session.AddHandler(func(s *dgo.Session, i *dgo.Ready) {
		fmt.Printf("Bot is up! Logged as: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	})

	session.AddHandler(func(s *dgo.Session, i *dgo.InteractionCreate) {
		switch i.Type {
		case dgo.InteractionApplicationCommand:
			switch i.ApplicationCommandData().Name {
			case "basic":
				handleBasicCmd(s, i)
			default:
				log.Panicf("Cannot handle '%v' interaction", i.ApplicationCommandData().Name)
			}
		case dgo.InteractionMessageComponent:
			handleResponseCmd(s, i)
		}
	})
}

// Starts the bot session
func openSession() error {
	if err := session.Open(); err != nil {
		return err
	}

	if err := syncCommands(); err != nil {
		return err
	}

	return nil
}

func New(token string) error {
	if err := createSession(token); err != nil {
		return err
	}

	createHandlers()

	if err := openSession(); err != nil {
		return err
	}

	return nil
}

func Run() {
	defer Quit()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func Quit() {
	fmt.Println("Quitting...")
	session.Close()
}
