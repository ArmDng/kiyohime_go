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
	Token          = os.Getenv("DISCORD_TOKEN")
	removeCommands = os.Getenv("RMCMD")
)

func init() { flag.Parse() }

func init() {
	var err error

	// Creation of the discord session with the bot token
	s, err = discordgo.New("Bot " + Token)
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
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "type",
					Description: "Principale ou Serveur",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Principale",
							Value: "principale",
						},
						{
							Name:  "Serveur",
							Value: "serveur",
						},
					},
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"pp": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			var (
				url   string
				title string
			)

			// Getting the data needed from the slash commands pp
			choiceOption := i.ApplicationCommandData().Options[1].StringValue()
			userID := i.ApplicationCommandData().Options[0].UserValue(s).ID

			// Getting the data of the user
			user, err := s.User(userID)
			if err != nil {
				log.Printf("Unable to retrieve the user: %v", err)
			}

			switch choiceOption {

			// If the choice was "principale"
			case "principale":
				url = user.AvatarURL("512")
				title = fmt.Sprintf("Avatar principal de %v", user.Username)

			// If the choice was "serveur"
			case "serveur":

				// Getting the data of the member (user of a server)
				member, err := s.GuildMember(getGuildID(s), userID)
				if err != nil {
					log.Printf("Unable to retrieve the member: %v", err)
					return
				}

				url = member.AvatarURL("512")
				title = fmt.Sprintf("Avatar de serveur de %v", user.Username)

			default:
				log.Printf("Kiyohime s'est perdu dans la bibliothèque de Chaldea")
				return
			}

			// Creation of the embed message

			embed := &discordgo.MessageEmbed{
				Title: title,
				Image: &discordgo.MessageEmbedImage{
					URL: url,
				},
				Color: 0x00ff00,
			}

			// Responding to the command

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

func getGuildID(s *discordgo.Session) string {
	guilds, err := s.UserGuilds(100, "", "")
	if err != nil {
		log.Fatalf("Unable to retrieve the Guild ID: %v", err)
		return ""
	}
	return guilds[0].ID
}

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
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
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	// Declare the intents
	s.Identify.Intents = discordgo.IntentsMessageContent

	defer s.Close()
	// Close the discord session

	// If term signal is received
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if removeCommands == "true" {
		log.Println("rm slash commands")

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
			if err != nil {
				log.Panicf("error deleting '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Good Night Master !")

}
