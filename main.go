package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token          = flag.String("DISCORD_TOKEN", "", "Bot Token")
	guildID        = flag.String("GUILD_ID", "", "GUILD ID, If not passed - bot registers commands globally")
	removeCommands = flag.Bool("rmcmd", true, "type 'false' to not remove command after extinction")
)

func init() { flag.Parse() }

func init() {
	var err error

	// Creation of the discord session with the bot token
	s, err = discordgo.New("Bot " + *Token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

}

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionSendMessages

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "pp",
			Description: "Affiche la pp d'un utilisateur",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "L'utilisateur dont vous voulez voir la pp",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"pp": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			userID := i.ApplicationCommandData().Options[0].UserValue(s).ID

			user, err := s.User(userID)
			if err != nil {
				log.Fatalf("Unable to get the user")
			}

			// Creation of the embed message

			embed := &discordgo.MessageEmbed{
				Title: fmt.Sprintf("Avatar de %v", user.Username),
				Image: &discordgo.MessageEmbedImage{
					URL: user.AvatarURL("512"),
				},
				Color: 0x00ff00,
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})
		},
	}
)

var s *discordgo.Session

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore the function if the author of the message is the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "Morgan" {

		author := discordgo.MessageEmbedAuthor{
			Name: "O.V.",
		}

		embed := discordgo.MessageEmbed{
			Author: &author,
			Title:  "Microbe",
		}

		s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	}

	if m.Content == "Queen" {
		s.ChannelMessageSend(m.ChannelID, "Morgan !")
	}

}

func main() {

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Kiyohime waked up as %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Open a websocket connection to Discord
	err := s.Open()
	if err != nil {
		log.Fatalf("error connecting discord session: %v", err)
	}

	fmt.Println("Kiyohime is fighting. Ctrl-C to make her sleep")

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *guildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	s.AddHandler(messageCreate)

	// Declare the intents
	s.Identify.Intents = discordgo.IntentsAll

	// Close the discord session
	defer s.Close()

	// If term signal is received
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if *removeCommands {
		log.Println("rm")

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *guildID, v.ID)
			if err != nil {
				log.Panicf("error deleting '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Good Night Master !")

}
