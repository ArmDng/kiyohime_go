package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
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

			var url string

			// Getting the data needed from the slash commands pp

			userID := i.ApplicationCommandData().Options[0].UserValue(s).ID

			// Getting the data of the user
			user, err := s.User(userID)
			if err != nil {
				log.Printf("Unable to retrieve the user: %v", err)
			}

			// Getting the URL of the banner
			url = user.BannerURL("512")
			// Creation of the embed message

			embed := &discordgo.MessageEmbed{
				Image: &discordgo.MessageEmbedImage{
					URL: url,
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("%v, %v", user.Username, "Bannière principale"),
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

		"Janken": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

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
			},
		},

		{
			Name:        "mystère",
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

var scheduledTimes = make(map[string]bool)

func sendMessageAtDate(s *discordgo.Session, date string, message string, channel string) {
	// Check if the message has already been scheduled for the current hour
	if scheduledTimes[date] {
		log.Println("The message has already been scheduled for the current hour")
		return
	}
	// Parsing the date string
	layout := "15:04"
	parisLoc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		log.Fatalf("Error loading timezone")
	}

	t, err := time.ParseInLocation(layout, date, parisLoc)
	if err != nil {
		log.Fatalf("Error parsing date: %v", err)
	}

	// Checking if the current time is after the specified date
	now := time.Now().In(parisLoc)

	t = time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, parisLoc)

	log.Printf("The current time is %v", now)
	log.Printf("The specified time is %v", t)
	// If the specified time is in the past, add 24 hours to it
	if now.After(t) {
		log.Printf("The specified time is in the past, adding 24 hours to it")
		t = t.Add(24 * time.Hour)
	}
	// Scheduling the message to be sent at the specified date
	time.AfterFunc(t.Sub(now), func() {
		// Send the message
		s.ChannelMessageSend(channel, message)
		log.Printf("Message sent at %v", t)

		// Mark the message as scheduled for the current hour
		scheduledTimes[date] = false

		// Schedule the message for the next day
		sendMessageAtDate(s, date, message, channel)
	})
	scheduledTimes[date] = true
}
func sendAutoMessage(s *discordgo.Session) {
	sendMessageAtDate(s, "12:00", "*Regarde* ", "747540564622442569")

	sendMessageAtDate(s, "00:00", "Faîtes de beaux rêves", "747540564622442569")
}

func sendMessageAtMidnight(s *discordgo.Session) {
	s.ChannelMessageSend("747540564622442569", "Coeur sur Jeanne + Je vous souhaite une merveilleuse année 2024 !")
}

func sendMessageAtMidnight12(s *discordgo.Session) {
	s.ChannelMessageSend("747540564622442569", "Je vous aime tous !")
}

func sendMessageAt01(s *discordgo.Session) {
	s.ChannelMessageSend("747540564622442569", "Coeur sur vous !")
}

func sendMessageAt02(s *discordgo.Session) {
	s.ChannelMessageSend("747540564622442569", "Tsukihime")
}

var isTaskSchedulded1 bool

func scheduleTaskatMidnight() {
	if isTaskSchedulded1 {
		return
	}

	parisLoc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		log.Fatalf("Error loading timezone")
	}

	now := time.Now().In(parisLoc)
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day(), 00, 00, 0, 0, parisLoc)

	if now.After(nextMidnight) {
		nextMidnight = nextMidnight.Add(24 * time.Hour)
	}

	timeUntilNextTime := nextMidnight.Sub(now)

	time.AfterFunc(timeUntilNextTime, func() {
		sendMessageAtMidnight(s)
		isTaskSchedulded1 = false

		scheduleTaskatMidnight()
	})

	isTaskSchedulded1 = true
}

var isTaskSchedulded2 bool

func scheduleTaskat01() {
	if isTaskSchedulded2 {
		return
	}

	parisLoc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		log.Fatalf("Error loading timezone")
	}

	now := time.Now().In(parisLoc)
	nextTime := time.Date(now.Year(), now.Month(), now.Day(), 02, 00, 0, 0, parisLoc)

	if now.After(nextTime) {
		nextTime = nextTime.Add(24 * time.Hour)
	}

	timeUntilNextTime := nextTime.Sub(now)

	time.AfterFunc(timeUntilNextTime, func() {
		sendMessageAt01(s)
		isTaskSchedulded2 = false

		scheduleTaskat01()
	})

	isTaskSchedulded2 = true
}

var isTaskSchedulded3 bool

func scheduleTaskat02() {
	if isTaskSchedulded3 {
		return
	}

	parisLoc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		log.Fatalf("Error loading timezone")
	}

	now := time.Now().In(parisLoc)
	nextTime := time.Date(now.Year(), now.Month(), now.Day(), 01, 00, 0, 0, parisLoc)

	if now.After(nextTime) {
		nextTime = nextTime.Add(24 * time.Hour)
	}

	timeUntilMidnight := nextTime.Sub(now)

	time.AfterFunc(timeUntilMidnight, func() {
		sendMessageAt02(s)
		isTaskSchedulded3 = false

		scheduleTaskat02()
	})

	isTaskSchedulded3 = true
}

var isTaskSchedulded12 bool

func scheduleTaskat12() {
	if isTaskSchedulded12 {
		return
	}

	parisLoc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		log.Fatalf("Error loading timezone")
	}

	now := time.Now().In(parisLoc)
	nextTime := time.Date(now.Year(), now.Month(), now.Day(), 12, 00, 0, 0, parisLoc)

	if now.After(nextTime) {
		nextTime = nextTime.Add(24 * time.Hour)
	}

	timeUntilMidnight := nextTime.Sub(now)

	time.AfterFunc(timeUntilMidnight, func() {
		sendMessageAtMidnight12(s)
		isTaskSchedulded12 = false

		scheduleTaskat12()
	})

	isTaskSchedulded3 = true
}

func scheduleAllTasks() {
	scheduleTask1()
	scheduleTask2()
	scheduleTask3()
	scheduleTask12()
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
	/*
		scheduleTaskatMidnight()
		scheduleTaskat01()
		scheduleTaskat02()
		scheduleTaskat12()
	*/
	sendAutoMessage(s)

	scheduleAllTasks()
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
