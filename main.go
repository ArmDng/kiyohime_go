package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
)

var (
	Token          = os.Getenv("DISCORD_TOKEN")
	removeCommands = os.Getenv("RMCMD")
)

var s *discordgo.Session

func init() {
	var err error

	// Creation of the discord session with the bot token
	s, err = discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

}

var (
	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){

		"rock": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			if i.Member.User.ID != i.Message.Interaction.User.ID {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredMessageUpdate,
				})
				return
			}
			userChoice := "rock"

			winner := manageJanken(userChoice)

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("%s", winner),
				},
			})
			if err != nil {
				log.Panicf("Unable to respond to the interaction: %v", err)
			}

		},

		"paper": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			if i.Member.User.ID != i.Message.Interaction.User.ID {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredMessageUpdate,
				})
				return
			}
			userChoice := "paper"

			winner := manageJanken(userChoice)

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("%s", winner),
				},
			})
			if err != nil {
				log.Panicf("Unable to respond to the interaction: %v", err)
			}

		},

		"scissors": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if i.Member.User.ID != i.Message.Interaction.User.ID {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredMessageUpdate,
				})
				return
			}
			userChoice := "scissors"

			winner := manageJanken(userChoice)

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("%s", winner),
				},
			})
			if err != nil {
				log.Panicf("Unable to respond to the interaction: %v", err)
			}

		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"pp": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			var (
				url        string
				title      string
				typeAvatar string
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
				typeAvatar = "Principale"

			// If the choice was "serveur"
			case "serveur":

				// Getting the data of the member (user of a server)
				member, err := s.GuildMember(i.GuildID, userID)
				if err != nil {
					log.Printf("Unable to retrieve the member: %v", err)
					return
				}

				url = member.AvatarURL("512")
				typeAvatar = "Serveur"

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
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("%v, %v", user.Username, typeAvatar),
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

		"bannière": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			var (
				url        string
				typeBanner string
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
				url = user.BannerURL("512")
				typeBanner = "Principale"

			// If the choice was "serveur"
			case "serveur":

				// Getting the data of the member (user of a server)
				member, err := s.GuildMember(i.GuildID, userID)
				if err != nil {
					log.Printf("Unable to retrieve the member: %v", err)
					return
				}

				url = member.BannerURL("512")
				typeBanner = "Serveur"

			default:
				log.Printf("Kiyohime s'est perdu dans la bibliothèque de Chaldea")
				return
			}

			// Creation of the embed message
			embed := &discordgo.MessageEmbed{
				Image: &discordgo.MessageEmbedImage{
					URL: url,
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("%v, %v", user.Username, typeBanner),
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

		"janken": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			embed := &discordgo.MessageEmbed{
				Title:       "Prépare-toi !",
				Description: "Choisis",
				Color:       0x00ff00,
			}

			components := []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.Button{
						// Label is what the user will see on the button.
						Label: "Pierre",
						// Style provides coloring of the button. There are not so many styles tho.
						Style: discordgo.SuccessButton,
						// Disabled allows bot to disable some buttons for users.
						Disabled: false,
						// CustomID is a thing telling Discord which data to send when this button will be pressed.
						CustomID: "rock",
					},

					discordgo.Button{
						Label:    "Papiers",
						Style:    discordgo.SuccessButton,
						Disabled: false,
						CustomID: "paper",
					},

					discordgo.Button{
						Label:    "Ciseaux",
						Style:    discordgo.SuccessButton,
						Disabled: false,
						CustomID: "scissors",
					},
				}},
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds:     []*discordgo.MessageEmbed{embed},
					Components: components,
				},
			})

			if err != nil {
				log.Printf("Unable to respond to the interaction: %v", err)
				return
			}
		},
	}

	// Definition of the different commands
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionAll

	commands = []*discordgo.ApplicationCommand{
		// Command to display the pfp of an user
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

		// Command for display the banner of an user

		{
			Name:        "bannière",
			Description: "Affiche la bannière d'un utilisateur",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "L'utilisateur dont vous voulez voir la bannière",
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

		{
			Name:        "janken",
			Description: "Je te défie au Pierre Feuille Ciseaux",
		},
	}
)

func manageJanken(userChoice string) string {
	var botChoice string
	var winner string
	botNumChoice := rand.Intn(3)

	switch botNumChoice {
	case 0:
		botChoice = "rock"
	case 1:
		botChoice = "paper"
	case 2:
		botChoice = "scissors"
	}

	switch userChoice {
	case "rock":
		switch botChoice {
		case "rock":
			winner = "Égalité mais je t'aurais la prochaine fois !"
		case "paper":
			winner = "J'ai gagné !"
		case "scissors":
			winner = "Nice victoire pour toi !"
		}
	case "paper":
		switch botChoice {
		case "rock":
			winner = "Égalité mais je t'aurais la prochaine fois !"
		case "paper":
			winner = "J'ai gagné !"
		case "scissors":
			winner = "Nice victoire pour toi !"
		}
	case "scissors":
		switch botChoice {
		case "rock":
			winner = "Égalité mais je t'aurais la prochaine fois !"
		case "paper":
			winner = "J'ai gagné !"
		case "scissors":
			winner = "Nice victoire pour toi !"
		}
	}

	return winner
}

func findIDGuild(s *discordgo.Session, name string) string {
	for _, guild := range s.State.Guilds {
		if guild.Name == name {
			return guild.ID
		}
	}
	return ""
}

func getRandomName(s *discordgo.Session, guildID string) string {
	members, err := s.GuildMembers(guildID, "", 1000)
	if err != nil {
		log.Fatalf("Error retrieving members: %v", err)
	}

	randomIndex := rand.Intn(len(members))
	randomMember := members[randomIndex]

	return randomMember.User.Username
}

func getLastAuthor(s *discordgo.Session, channelID string) string {
	messages, err := s.ChannelMessages(channelID, 1, "", "", "")
	if err != nil {
		log.Fatalf("Error retrieving messages: %v", err)
	}

	return messages[0].Author.Username
}

func main() {

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {

		log.Printf("Kiyohime waked up as %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}

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

	if removeCommands == "" {
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
