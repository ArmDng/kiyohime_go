package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

func init() {

	flag.StringVar(&Token, "DISCORD_TOKEN", "", "Bot Token")
	flag.Parse()
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

	// Creation of the discord session with the bot token

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating discord session", err)
		return
	}

	// call the func

	dg.AddHandler(messageCreate)

	// Declare the intents
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and start listening
	err = dg.Open()

	if err != nil {
		fmt.Println("error opening connection", err)
		return
	}

	fmt.Println("Kiyohime is fighting. Ctrl-C to make her sleep")

	// If term signal is received
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Close the discord_session
	dg.Close()
}
