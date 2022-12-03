package main

import (
	"fmt"
	"github.com/Goscord/goscord/discord"
	"github.com/Goscord/goscord/discord/embed"
	"gorm.io/gorm"
	"log"
	"strings"
)

type SetupCommand struct{}

func (c *SetupCommand) Name() string {
	return "setup"
}

func (c *SetupCommand) Description() string {
	return "Setup and, reconfigure your advent of code leaderboard bot for this guild."
}

func (c *SetupCommand) Category() string {
	return "general"
}

const __SESSION_KEY = "sessionkey"
const __LEADBOARD_CODE = "leaderboardurl"

func (c *SetupCommand) Options() []*discord.ApplicationCommandOption {
	return []*discord.ApplicationCommandOption{
		{
			Type:        discord.ApplicationCommandOptionString,
			Name:        __SESSION_KEY,
			Description: "Your advent of code session key, this is used to fetch your leaderboard and, nothing else.",
			Required:    true,
		},
		{
			Type:        discord.ApplicationCommandOptionString,
			Name:        __LEADBOARD_CODE,
			Description: "Your advent of code leaderboard's url or, code.",
			Required:    true,
		},
	}
}

const SESSION_KEY_BLOCK_COUNT = 4
const NON_REDACTED_CHAR_COUNT = 4

func (c *SetupCommand) Execute(ctx *Context) bool {
	if ctx.interaction.Member.Permissions&discord.BitwisePermissionFlagAdministrator == 0 {
		SendPermissionsError(ctx)
		return false
	}

	sessionkey := ctx.interaction.Data.Options[0].String()
	leaderboardurl := ctx.interaction.Data.Options[1].String()

	// URL to code if needed
	leaderboardcode := leaderboardurl
	if strings.Contains(leaderboardcode, "/") {
		leaderboardcode = leaderboardurl[strings.LastIndex(leaderboardcode, "/")+1:]
	}

	if strings.Contains(leaderboardcode, "-") {
		leaderboardcode = leaderboardurl[strings.LastIndex(leaderboardcode, "-")+1:]
	}

	// Add to database
	guildid := ctx.interaction.GuildId

	err := db.Transaction(func(tx *gorm.DB) error {
		dest := GuildSettings{
			ID: guildid,
		}

		db.FirstOrCreate(&dest, guildid).Updates(map[string]interface{}{
			"SessionKey": sessionkey, "BoardCode": leaderboardcode})
		return nil
	})
	if err != nil {
		log.Print(err)
		SendDatabaseError(ctx)
		return false
	}

	// On success send a happy message!
	sessionkey_redacted := ""
	i := 0
	for ; i < len(sessionkey)-NON_REDACTED_CHAR_COUNT; i++ {
		if i%SESSION_KEY_BLOCK_COUNT == 0 {
			sessionkey_redacted += " "
		}
		sessionkey_redacted += "\\*"
	}

	for ; i < len(sessionkey); i++ {
		if i%SESSION_KEY_BLOCK_COUNT == 0 {
			sessionkey_redacted += " "
		}
		sessionkey_redacted += string(sessionkey[i])
	}

	e := embed.NewEmbedBuilder()
	message := fmt.Sprintf("**Reconfigured by:** <@%s>\n**Session key:** %s\n**Leaderboard Code:** [%s](%s)",
		ctx.interaction.Member.User.Id,
		sessionkey_redacted,
		leaderboardcode)

	e.SetTitle("Advent of Code Bot Configuration Changed")
	e.SetDescription(message)
	e.SetThumbnail(ctx.interaction.Member.User.AvatarURL())
	ThemeEmbed(e, ctx)

	// Send response
	ctx.client.Interaction.CreateResponse(ctx.interaction.Id,
		ctx.interaction.Token,
		&discord.InteractionCallbackMessage{Embeds: []*embed.Embed{e.Embed()},
			Flags: discord.MessageFlagUrgent})

	return true
}
