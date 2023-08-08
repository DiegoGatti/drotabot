package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	opendotaapi "github.com/DiegoGatti/drotabot/OpenDotaApi"
	"github.com/bwmarrin/discordgo"
)

const OpenDotaAPIURL = "https://api.opendota.com/api/"

// Variables used for command line parameters
var (
	GuildID        string
	BotToken       string
	RemoveCommands bool
	dgs            *discordgo.Session
)

// Parse command-line flags.
func init() {
	flag.StringVar(&GuildID, "gid", "", "GuildID")
	flag.StringVar(&BotToken, "t", "", "Bot Token")
	flag.Parse()
}

// Create a new Discord session using the provided bot token.
func init() {
	var err error

	dgs, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("error creating Discord session, %v", err)
	}
}

func init() {
	dgs.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {

	// Register the messageCreate func as a callback for MessageCreate events.
	dgs.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dgs.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err := dgs.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dgs.Close()
}

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer
)

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch {
	case strings.Contains(m.Content, "!help"):
		CmdHelp(s, m.ChannelID)
	case strings.Contains(m.Content, "!search"):
		CmdSearch(s, m.ChannelID, strings.TrimSpace(strings.Replace(m.Content, "!search", "", 1)))
	default:
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s is not recognized as a valid command. Try !help", m.Content))
	}

	if strings.Contains(m.Content, "!search") {

	}
}

func CmdHelp(s *discordgo.Session, ChId string) {
	_, err := s.ChannelMessageSend(ChId, "How may i help?\n ```Use comands:\n\n!search\n!match\n!player```")
	if err != nil {
		fmt.Println(err)
	}
}

func CmdSearch(s *discordgo.Session, ChId string, q string) {
	c, err := opendotaapi.NewClientWithResponses("https://api.opendota.com/api")
	if err != nil {
		panic(err)
	}

	var params = opendotaapi.GetSearchParams{
		Q: q,
	}

	resp, err := c.GetSearchWithResponse(context.Background(), &params)
	if err != nil {
		panic(err)
	}

	var respDecoded []opendotaapi.SearchResponse
	json.Unmarshal(resp.Body, &respDecoded)

	msg := fmt.Sprintf("found %d results.", len(respDecoded))

	s.ChannelMessageSend(ChId, msg)
}
