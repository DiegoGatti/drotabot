package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"

	opendotaapi "github.com/DiegoGatti/drotabot/OpenDotaApi"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "help",
			Description: "Gets list of commands and short descriptions",
		},
		{
			Name:        "search",
			Description: "Searches Dota 2 player by name",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "query",
					Description: "Query",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "How may i help?\n ```Use comands:\n\n/search\n/match\n/player```",
				},
			})
		},
		"search": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			var msgformat string

			if option, ok := optionMap["query"]; ok {
				msgformat = option.StringValue()
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: CmdSearch(msgformat),
				},
			})
		},
	}
)

func CmdSearch(q string) string {
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

	return msg
}
