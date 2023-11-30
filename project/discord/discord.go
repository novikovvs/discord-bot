package discord

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"np-discord-bot/structs"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	GuildID        = flag.String("guild", os.Getenv("GUILD_ID"), "Test guild ID")
	BotToken       = flag.String("token", os.Getenv("BOT_TOKEN"), "Bot access token")
	AppID          = flag.String("app", os.Getenv("APP_ID"), "Application ID")
	Cleanup        = flag.Bool("cleanup", true, "Cleanup of commands")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commandsInline after shutdowning or not")
	characters     []structs.Character
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	jsonFile, err := os.Open(os.Getenv("CHAMPIONS_LIST_PATH"))

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &characters)
}

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	commandsInline = []*discordgo.ApplicationCommand{
		{
			Name: "rickroll-em",
			Type: discordgo.UserApplicationCommand,
		},
		{
			Name: "roll-pick",
			// All commandsInline and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Подскажет какого перса пикнуть в лижке:)))",
		},
		//{
		//	Name:                     "permission-overview",
		//	Description:              "Command for demonstration of default command permissions",
		//	DefaultMemberPermissions: &defaultMemberPermissions,
		//	DMPermission:             &dmPermission,
		//},
		//{
		//	Name:        "basic-command-with-files",
		//	Description: "Basic command with files",
		//},
		//{
		//	Name:        "localized-command",
		//	Description: "Localized command. Description and name may vary depending on the Language setting",
		//	NameLocalizations: &map[discordgo.Locale]string{
		//		discordgo.ChineseCN: "本地化的命令",
		//	},
		//	DescriptionLocalizations: &map[discordgo.Locale]string{
		//		discordgo.ChineseCN: "这是一个本地化的命令",
		//	},
		//	Options: []*discordgo.ApplicationCommandOption{
		//		{
		//			Name:        "localized-option",
		//			Description: "Localized option. Description and name may vary depending on the Language setting",
		//			NameLocalizations: map[discordgo.Locale]string{
		//				discordgo.ChineseCN: "一个本地化的选项",
		//			},
		//			DescriptionLocalizations: map[discordgo.Locale]string{
		//				discordgo.ChineseCN: "这是一个本地化的选项",
		//			},
		//			Type: discordgo.ApplicationCommandOptionInteger,
		//			Choices: []*discordgo.ApplicationCommandOptionChoice{
		//				{
		//					Name: "First",
		//					NameLocalizations: map[discordgo.Locale]string{
		//						discordgo.ChineseCN: "一的",
		//					},
		//					Value: 1,
		//				},
		//				{
		//					Name: "Second",
		//					NameLocalizations: map[discordgo.Locale]string{
		//						discordgo.ChineseCN: "二的",
		//					},
		//					Value: 2,
		//				},
		//			},
		//		},
		//	},
		//},
		//{
		//	Name:        "options",
		//	Description: "Command for demonstrating options",
		//	Options: []*discordgo.ApplicationCommandOption{
		//
		//		{
		//			Type:        discordgo.ApplicationCommandOptionString,
		//			Name:        "string-option",
		//			Description: "String option",
		//			Required:    true,
		//		},
		//		{
		//			Type:        discordgo.ApplicationCommandOptionInteger,
		//			Name:        "integer-option",
		//			Description: "Integer option",
		//			MinValue:    &integerOptionMinValue,
		//			MaxValue:    10,
		//			Required:    true,
		//		},
		//		{
		//			Type:        discordgo.ApplicationCommandOptionNumber,
		//			Name:        "number-option",
		//			Description: "Float option",
		//			MaxValue:    10.1,
		//			Required:    true,
		//		},
		//		{
		//			Type:        discordgo.ApplicationCommandOptionBoolean,
		//			Name:        "bool-option",
		//			Description: "Boolean option",
		//			Required:    true,
		//		},
		//
		//		// Required options must be listed first since optional parameters
		//		// always come after when they're used.
		//		// The same concept applies to Discord's Slash-commandsInline API
		//
		//		{
		//			Type:        discordgo.ApplicationCommandOptionChannel,
		//			Name:        "channel-option",
		//			Description: "Channel option",
		//			// Channel type mask
		//			ChannelTypes: []discordgo.ChannelType{
		//				discordgo.ChannelTypeGuildText,
		//				discordgo.ChannelTypeGuildVoice,
		//			},
		//			Required: false,
		//		},
		//		{
		//			Type:        discordgo.ApplicationCommandOptionUser,
		//			Name:        "user-option",
		//			Description: "User option",
		//			Required:    false,
		//		},
		//		{
		//			Type:        discordgo.ApplicationCommandOptionRole,
		//			Name:        "role-option",
		//			Description: "Role option",
		//			Required:    false,
		//		},
		//	},
		//},
		//{
		//	Name:        "subcommands",
		//	Description: "Subcommands and command groups example",
		//	Options: []*discordgo.ApplicationCommandOption{
		//		// When a command has subcommands/subcommand groups
		//		// It must not have top-level options, they aren't accesible in the UI
		//		// in this case (at least not yet), so if a command has
		//		// subcommands/subcommand any groups registering top-level options
		//		// will cause the registration of the command to fail
		//
		//		{
		//			Name:        "subcommand-group",
		//			Description: "Subcommands group",
		//			Options: []*discordgo.ApplicationCommandOption{
		//				// Also, subcommand groups aren't capable of
		//				// containing options, by the name of them, you can see
		//				// they can only contain subcommands
		//				{
		//					Name:        "nested-subcommand",
		//					Description: "Nested subcommand",
		//					Type:        discordgo.ApplicationCommandOptionSubCommand,
		//				},
		//			},
		//			Type: discordgo.ApplicationCommandOptionSubCommandGroup,
		//		},
		//		// Also, you can create both subcommand groups and subcommands
		//		// in the command at the same time. But, there's some limits to
		//		// nesting, count of subcommands (top level and nested) and options.
		//		// Read the intro of slash-commandsInline docs on Discord dev portal
		//		// to get more information
		//		{
		//			Name:        "subcommand",
		//			Description: "Top-level subcommand",
		//			Type:        discordgo.ApplicationCommandOptionSubCommand,
		//		},
		//	},
		//},
		//{
		//	Name:        "responses",
		//	Description: "Interaction responses testing initiative",
		//	Options: []*discordgo.ApplicationCommandOption{
		//		{
		//			Name:        "resp-type",
		//			Description: "Response type",
		//			Type:        discordgo.ApplicationCommandOptionInteger,
		//			Choices: []*discordgo.ApplicationCommandOptionChoice{
		//				{
		//					Name:  "Channel message with source",
		//					Value: 4,
		//				},
		//				{
		//					Name:  "Deferred response With Source",
		//					Value: 5,
		//				},
		//			},
		//			Required: true,
		//		},
		//	},
		//},
		//{
		//	Name:        "followups",
		//	Description: "Followup messages",
		//},
	}

	commandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"basic-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there! Congratulations, you just executed your first slash command",
				},
			})
		},
		"basic-command-with-files": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there! Congratulations, you just executed your first slash command with a file in the response",
					Files: []*discordgo.File{
						{
							ContentType: "text/plain",
							Name:        "test.txt",
							Reader:      strings.NewReader("Hello Discord!!"),
						},
					},
				},
			})
		},
		"localized-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			responses := map[discordgo.Locale]string{
				discordgo.ChineseCN: "你好！ 这是一个本地化的命令",
			}
			response := "Hi! This is a localized message"
			if r, ok := responses[i.Locale]; ok {
				response = r
			}
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})
			if err != nil {
				log.Println(err)
			}
		},
		"options": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			// This example stores the provided arguments in an []interface{}
			// which will be used to format the bot's response
			margs := make([]interface{}, 0, len(options))
			msgformat := "You learned how to use command options! " +
				"Take a look at the value(s) you entered:\n"

			// Get the value from the option map.
			// When the option exists, ok = true
			if option, ok := optionMap["string-option"]; ok {
				// Option values must be type asserted from interface{}.
				// Discordgo provides utility functions to make this simple.
				margs = append(margs, option.StringValue())
				msgformat += "> string-option: %s\n"
			}

			if opt, ok := optionMap["integer-option"]; ok {
				margs = append(margs, opt.IntValue())
				msgformat += "> integer-option: %d\n"
			}

			if opt, ok := optionMap["number-option"]; ok {
				margs = append(margs, opt.FloatValue())
				msgformat += "> number-option: %f\n"
			}

			if opt, ok := optionMap["bool-option"]; ok {
				margs = append(margs, opt.BoolValue())
				msgformat += "> bool-option: %v\n"
			}

			if opt, ok := optionMap["channel-option"]; ok {
				margs = append(margs, opt.ChannelValue(nil).ID)
				msgformat += "> channel-option: <#%s>\n"
			}

			if opt, ok := optionMap["user-option"]; ok {
				margs = append(margs, opt.UserValue(nil).ID)
				msgformat += "> user-option: <@%s>\n"
			}

			if opt, ok := optionMap["role-option"]; ok {
				margs = append(margs, opt.RoleValue(nil, "").ID)
				msgformat += "> role-option: <@&%s>\n"
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, they will be discussed in "responses"
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
		"permission-overview": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			perms, err := s.ApplicationCommandPermissions(s.State.User.ID, i.GuildID, i.ApplicationCommandData().ID)

			var restError *discordgo.RESTError
			if errors.As(err, &restError) && restError.Message != nil && restError.Message.Code == discordgo.ErrCodeUnknownApplicationCommandPermissions {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: ":x: No permission overwrites",
					},
				})
				return
			} else if err != nil {
				log.Println(err)
			}

			if err != nil {
				log.Println(err)
			}
			format := "- %s %s\n"

			channels := ""
			users := ""
			roles := ""

			for _, o := range perms.Permissions {
				emoji := "❌"
				if o.Permission {
					emoji = "☑"
				}

				switch o.Type {
				case discordgo.ApplicationCommandPermissionTypeUser:
					users += fmt.Sprintf(format, emoji, "<@!"+o.ID+">")
				case discordgo.ApplicationCommandPermissionTypeChannel:
					allChannels, _ := discordgo.GuildAllChannelsID(i.GuildID)

					if o.ID == allChannels {
						channels += fmt.Sprintf(format, emoji, "All channels")
					} else {
						channels += fmt.Sprintf(format, emoji, "<#"+o.ID+">")
					}
				case discordgo.ApplicationCommandPermissionTypeRole:
					if o.ID == i.GuildID {
						roles += fmt.Sprintf(format, emoji, "@everyone")
					} else {
						roles += fmt.Sprintf(format, emoji, "<@&"+o.ID+">")
					}
				}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Permissions overview",
							Description: "Overview of permissions for this command",
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:  "Users",
									Value: users,
								},
								{
									Name:  "Channels",
									Value: channels,
								},
								{
									Name:  "Roles",
									Value: roles,
								},
							},
						},
					},
					AllowedMentions: &discordgo.MessageAllowedMentions{},
				},
			})
		},
		"subcommands": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			content := ""

			// As you can see, names of subcommands (nested, top-level)
			// and subcommand groups are provided through the arguments.
			switch options[0].Name {
			case "subcommand":
				content = "The top-level subcommand is executed. Now try to execute the nested one."
			case "subcommand-group":
				options = options[0].Options
				switch options[0].Name {
				case "nested-subcommand":
					content = "Nice, now you know how to execute nested commandsInline too"
				default:
					content = "Oops, something went wrong.\n" +
						"Hol' up, you aren't supposed to see this message."
				}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
		},
		"responses": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Responses to a command are very important.
			// First of all, because you need to react to the interaction
			// by sending the response in 3 seconds after receiving, otherwise
			// interaction will be considered invalid and you can no longer
			// use the interaction token and ID for responding to the user's request

			content := ""
			// As you can see, the response type names used here are pretty self-explanatory,
			// but for those who want more information see the official documentation
			switch i.ApplicationCommandData().Options[0].IntValue() {
			case int64(discordgo.InteractionResponseChannelMessageWithSource):
				content =
					"You just responded to an interaction, sent a message and showed the original one. " +
						"Congratulations!"
				content +=
					"\nAlso... you can edit your response, wait 5 seconds and this message will be changed"
			default:
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
				})
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
				}
				return
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong",
				})
				return
			}
			time.AfterFunc(time.Second*5, func() {
				content := content + "\n\nWell, now you know how to create and edit responses. " +
					"But you still don't know how to delete them... so... wait 10 seconds and this " +
					"message will be deleted."
				_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &content,
				})
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
				time.Sleep(time.Second * 10)
				s.InteractionResponseDelete(i.Interaction)
			})
		},
		"followups": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Followup messages are basically regular messages (you can create as many of them as you wish)
			// but work as they are created by webhooks and their functionality
			// is for handling additional messages after sending a response.

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					// Note: this isn't documented, but you can use that if you want to.
					// This flag just allows you to create messages visible only for the caller of the command
					// (user who triggered the command)
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "Surprise!",
				},
			})
			msg, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "Followup message has been created, after 5 seconds it will be edited",
			})
			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong",
				})
				return
			}
			time.Sleep(time.Second * 5)

			content := "Now the original message is gone and after 10 seconds this message will ~~self-destruct~~ be deleted."
			s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
				Content: &content,
			})

			time.Sleep(time.Second * 10)

			s.FollowupMessageDelete(i.Interaction, msg.ID)

			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "For those, who didn't skip anything and followed tutorial along fairly, " +
					"take a unicorn :unicorn: as reward!\n" +
					"Also, as bonus... look at the original interaction response :D",
			})
		},
		"rickroll-em": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Operation rickroll has begun",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				log.Println(err)
			}

			ch, err := s.UserChannelCreate(
				i.ApplicationCommandData().TargetID,
			)
			if err != nil {
				_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: fmt.Sprintf("Mission failed. Cannot send a message to this user: %q", err.Error()),
					Flags:   discordgo.MessageFlagsEphemeral,
				})
				if err != nil {
					log.Println(err)
				}
			}
			_, err = s.ChannelMessageSend(
				ch.ID,
				fmt.Sprintf("%s sent you this: https://youtu.be/dQw4w9WgXcQ", i.Member.Mention()),
			)
			if err != nil {
				log.Println(err)
			}
		},
		"roll-pick": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)
			index := r1.Intn(len(characters))

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: characters[index].Icon,
				},
			})
			if err != nil {
				log.Println(err)
			}
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func Start() {

	comm, _ := s.ApplicationCommands(*AppID, *GuildID)
	for _, IDs := range comm {
		s.ApplicationCommandDelete(*AppID, *GuildID, IDs.ID)
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	cmdIDs := make(map[string]string, len(commandsInline))

	for _, cmd := range commandsInline {

		rcmd, err := s.ApplicationCommandCreate(*AppID, *GuildID, cmd)
		if err != nil {
			log.Fatalf("Cannot create slash command %q: %v", cmd.Name, err)
		}

		cmdIDs[rcmd.ID] = rcmd.Name

	}

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")

	if !*Cleanup {
		return
	}

	for id, name := range cmdIDs {
		err := s.ApplicationCommandDelete(*AppID, *GuildID, id)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commandsInline...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commandsInline))
	for i, v := range commandsInline {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop = make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commandsInline...")
		// // We need to fetch the commandsInline, since deleting requires the command ID.
		// // We are doing this from the returned commandsInline on line 375, because using
		// // this will delete all the commandsInline, which might not be desirable, so we
		// // are deleting only the commandsInline that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commandsInline: %v", err)
		// }

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "/roll" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}