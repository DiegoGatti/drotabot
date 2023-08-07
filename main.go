package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	opendotaapi "github.com/DiegoGatti/drotabot/OpenDotaApi"
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

const KuteGoAPIURL = "https://kutego-api-xxxxx-ew.a.run.app"
const OpenDotaAPIURL = "https://api.opendota.com/api/"

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
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
	dg.Close()
}

type Gopher struct {
	Name string `json: "name"`
}

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
		CmdSearch(s, m.ChannelID)
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

func CmdSearch(s *discordgo.Session, ChId string) {
	q := strings.TrimSpace(strings.Replace(m.Content, "!search", "", 1))

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
