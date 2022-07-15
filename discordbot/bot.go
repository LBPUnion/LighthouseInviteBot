package discordbot

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Zaprit/LighthouseInviteBot/common"
	lighthouseapi "github.com/Zaprit/LighthouseInviteBot/lighthouseAPI"
	"github.com/bwmarrin/discordgo"
)

var default_permission int64 = discordgo.PermissionManageServer

var messageSendFail = "Failed to send message, most likely the specified user has DMs disabled"
var attachmentDownloadFail = "Failed to download attachment"

var commands = []*discordgo.ApplicationCommand{
	{
		Name:                     "sendinvite",
		Description:              "Invites User To Lighthouse",
		DefaultMemberPermissions: &default_permission,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "User To Invite",
				Required:    true,
			},
		},
	},
	{
		Name:                     "sendinvites",
		Description:              "Bulk invites users from a csv file",
		DefaultMemberPermissions: &default_permission,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "file",
				Description: "CSV formatted file with a list of discord user IDs",
			},
		},
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"sendinvite": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Sending...",
			},
		})
		if err != nil {
			return
		}

		option := i.ApplicationCommandData().Options[0]
		channel, err := s.UserChannelCreate(option.UserValue(s).ID)
		if err != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &messageSendFail,
			})
		}
		s.ChannelMessageSend(channel.ID, "You've been invited to join "+common.LoadConfig().Lighthouse.InstanceName)
		s.ChannelMessageSend(channel.ID, "Click here to create an account: "+lighthouseapi.GetInviteURL())
	},
	"sendinvites": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "One moment...",
			},
		})
		if err != nil {
			panic(err.Error())
		}

		url := i.ApplicationCommandData().Resolved.Attachments[i.ApplicationCommandData().Options[0].Value.(string)].URL
		resp, er2 := http.Get(url)
		if er2 != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &attachmentDownloadFail,
			})
		}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			for _, v := range strings.Split(scanner.Text(), ",") {
				channel, err := s.UserChannelCreate(v)
				if err != nil {
					s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
						Content: &messageSendFail,
					})
				}
				s.ChannelMessageSend(channel.ID, "You've been invited to join "+common.LoadConfig().Lighthouse.InstanceName)
				s.ChannelMessageSend(channel.ID, "Click here to create an account: "+lighthouseapi.GetInviteURL())
			}
		}

	},
}

func CreateBot(token string) (*discordgo.Session, error) {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return nil, err
	}

	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	// Open a websocket connection to Discord and begin listening.
	er2 := discord.Open()
	if er2 != nil {
		fmt.Println("error opening connection,", err)
		return nil, er2
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	return discord, nil
}

func DiscordNameFromID(s *discordgo.Session, id string) string {
	user, err := s.User(id)
	if err != nil {
		log.Panicln(err.Error())
	}
	return user.Username + "#" + user.Discriminator
}
