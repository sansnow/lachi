package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
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

	enabled := true
	rule, err := dg.AutoModerationRuleCreate(GuildID, &discordgo.AutoModerationRule{
		Name:        "Auto Moderation Testing",
		EventType:   discordgo.AutoModerationEventMessageSend,
		TriggerType: discordgo.AutoModerationEventTriggerKeyword,
		TriggerMetadata: &discordgo.AutoModerationTriggerMetadata{
			KeywordFilter: []string{"*cat*"},
			RegexPatterns: []string{"(c|b)at"},
		},

		Enabled: &enabled,
		Actions: []discordgo.AutoModerationAction{
			{Type: discordgo.AutoModerationRuleActionBlockMessage},
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully created the rule")

	defer dg.AutoModerationRuleDelete(GuildID, rule.ID)

	dg.AddHandlerOnce(func(s *discordgo.Session, e *discordgo.AutoModerationActionExecution) {
		_, err = dg.AutoModerationRuleEdit(GuildID, rule.ID, &discordgo.AutoModerationRule{
			TriggerMetadata: &discordgo.AutoModerationTriggerMetadata{
				KeywordFilter: []string{"cat"},
			},
			Actions: []discordgo.AutoModerationAction{
				{Type: discordgo.AutoModerationRuleActionTimeout, Metadata: &discordgo.AutoModerationActionMetadata{Duration: 60}},
				{Type: discordgo.AutoModerationRuleActionSendAlertMessage, Metadata: &discordgo.AutoModerationActionMetadata{
					ChannelID: e.ChannelID,
				}},
			},
		})

		fmt.Println("auto mod rule triggered")
		if err != nil {
			dg.AutoModerationRuleDelete(GuildID, rule.ID)
			panic(err)
		}

		s.ChannelMessageSend(e.ChannelID, "meow lol",)

		var counter int
		var counterMutex sync.Mutex
		dg.AddHandler(func(s *discordgo.Session, e *discordgo.AutoModerationActionExecution) {
			action := "unknown"
			switch e.Action.Type {
			case discordgo.AutoModerationRuleActionBlockMessage:
				action = "block message"
			case discordgo.AutoModerationRuleActionSendAlertMessage:
				action = "send alert message into <#" + e.Action.Metadata.ChannelID + ">"
			case discordgo.AutoModerationRuleActionTimeout:
				action = "timeout"
			}

			counterMutex.Lock()
			counter++
			switch counter {
			case 1:
				counterMutex.Unlock()
				s.ChannelMessageSend(e.ChannelID, "Nothing has changed, right? "+
					"Well, since separate gateway events are fired per each action (current is "+action+"), ",)
			}
				dg.Close()
				dg.AutoModerationRuleDelete(GuildID, rule.ID)
				os.Exit(0)
		})
	})

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
