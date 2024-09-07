// TODO
package bot

import (
	"fmt"
	"os"
	"os/signal"

	dgo "github.com/bwmarrin/discordgo"

	"github.com/pbnjk/cardshark/bot/games"
	"github.com/pbnjk/cardshark/helper"
)

var (
	session        *dgo.Session
	registeredCmds []*dgo.ApplicationCommand

	Commands           []*dgo.ApplicationCommand
	ComponentsHandlers helper.HandlerMap
	CommandHandlers    helper.HandlerMap
)

func init() {
	Commands = make([]*dgo.ApplicationCommand, 0)
}

// Syncs commands across runs, deleting old commands and creating or editing
// new ones
func syncCommands() error {
	existingCommands, err := session.ApplicationCommands(session.State.User.ID, "")
	if err != nil {
		return fmt.Errorf("Failed to fetch commands for guild %s: %v", "", err)
	}

	desiredMap := make(map[string]*dgo.ApplicationCommand)
	for _, cmd := range Commands {
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
	for _, cmd := range Commands {
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

// Starts the games
func createGames() {
	Commands = append(Commands, games.InitBlackjack(session)...)
}

// Creates the necessary handlers
func createHandlers() {
	session.AddHandler(func(s *dgo.Session, i *dgo.Ready) {
		fmt.Printf("Bot is up! Logged as: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	})

	session.AddHandler(func(s *dgo.Session, i *dgo.InteractionCreate) {
		switch i.Type {
		case dgo.InteractionApplicationCommand:
			if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case dgo.InteractionMessageComponent:
			if h, ok := ComponentsHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
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

// Creates and starts new bot
func Start(token string) error {
	if err := createSession(token); err != nil {
		return err
	}

	createGames()
	createHandlers()

	if err := openSession(); err != nil {
		return err
	}

	defer Quit()

	fmt.Printf("The %+v\n", Commands)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	return nil
}

func Quit() {
	fmt.Println("Quitting...")
	session.Close()
}
