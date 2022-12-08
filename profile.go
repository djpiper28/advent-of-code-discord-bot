package main

import (
	"fmt"
	"github.com/Goscord/goscord/discord"
	"github.com/Goscord/goscord/discord/embed"
	"log"
	"sort"
)

type ProfileCommand struct{}

func (c *ProfileCommand) Name() string {
	return "aocprofile"
}

func (c *ProfileCommand) Description() string {
	return "View the profile of a user and, see their rank over time for!"
}

func (c *ProfileCommand) Category() string {
	return "general"
}

const __NAME = "name"

func (c *ProfileCommand) Options() []*discord.ApplicationCommandOption {
	return []*discord.ApplicationCommandOption{
		{
			Type:        discord.ApplicationCommandOptionString,
			Name:        __NAME,
			Description: "The name of the person that you wish to lookup",
			Required:    true,
		},
	}
}

func (c *ProfileCommand) Execute(ctx *Context) bool {
	name := ctx.interaction.Data.Options[0].String()

	var gs GuildSettings
	r := db.First(&gs, ctx.interaction.GuildId)
	if r == nil {
		log.Print("Cannot find any matching guilds")
		SendDatabaseError(ctx)
		return false
	}

	entries, err := GetProfile(name, gs)
	if err != nil {
		log.Print(err)
		SendDatabaseError(ctx)
		return false
	}

	sort.Slice(entries, func(a int, b int) bool {
		return entries[a].Time.Unix() > entries[b].Time.Unix()
	})

	e := embed.NewEmbedBuilder()
	message := fmt.Sprintf("Board code: `%s`\n", gs.BoardCode)
	if len(entries) == 0 {
		message = fmt.Sprintf("The user `%s` cannot be found in your guild.", name)
	}

	e.SetTitle(fmt.Sprintf("Advent of Code Profile: ", name))
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
