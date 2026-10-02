package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	Token     string
	GuildID   string
	ChannelId string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Authorizartion Token")
	flag.StringVar(&GuildID, "guild", "", "Id of the testing guild")
	flag.StringVar(&ChannelId, "channel", "", "ID of the testing channel")
	flag.Parse()
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")

	}

	if Token == "" {
		Token = os.Getenv("DISCORD_TOKEN")
	}

	if GuildID == "" {
		GuildID = os.Getenv("GUILD_ID")
	}

	if ChannelId == "" {
		ChannelId = os.Getenv("CHANEL_ID")
	}

	dg, err := discordgo.New("Bot " + Token)

	if err != nil {
		fmt.Println("error creating a discord bot", err)
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents |= discordgo.IntentsGuildMessages
	dg.Identify.Intents |= discordgo.IntentMessageContent
	dg.Identify.Intents |= discordgo.IntentAutoModerationExecution

	err = dg.Open()

	if err != nil {
		fmt.Println("error opening connection for the bot", err)
		return
	}

	fmt.Println("bot is now running .... press CTRL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "kida fer" {
		s.ChannelMessageSend(m.ChannelID, "vdiya ji tusi sunao")
	}

	if m.Content == "lachi" {
		s.ChannelMessageSend(m.ChannelID, "yes?")
	}

	if m.Content == "time" {
		currentTime := time.Now()
		s.ChannelMessageSend(m.ChannelID, currentTime.String())
	}

	if m.Content == "should i sleep?" {
		s.ChannelMessageSend(m.ChannelID, "nah the city need you")
	}
}
