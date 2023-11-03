package main

import (
	"github.com/Goscord/goscord/discord"
	"github.com/Goscord/goscord/discord/embed"
)

type HelpCommand struct{}

func (c *HelpCommand) Name() string {
	return "help"
}

func (c *HelpCommand) Description() string {
	return "Shows help about the AOC bot and, how to set it up and, use it."
}

func (c *HelpCommand) Category() string {
	return "general"
}

func (c *HelpCommand) Options() []*discord.ApplicationCommandOption {
	return []*discord.ApplicationCommandOption{}
}

func (c *HelpCommand) Execute(ctx *Context) bool {
	commandRequests++
	e := embed.NewEmbedBuilder()
	e.SetTitle("Advent of Code Bot Help")
	e.SetDescription(`**Setup**
In order to setup the bot you need
 1. A session id, this can be got by going reading it from your browser's local storage in inspect element, it is a cookie saved as session.
 2. You also need to copy paste the leaderboard code of your leaderboard
 3. Use /setup, you need Admnistrator permissions to do this

**Privacy**
Data is stored in a database and, cannot be accessed externally, all data is not shared.`)
	e.SetThumbnail(ctx.interaction.Member.User.AvatarURL())
	ThemeEmbed(e, ctx)

	// Send response
	ctx.client.Interaction.CreateResponse(ctx.interaction.Id,
		ctx.interaction.Token,
		&discord.InteractionCallbackMessage{Embeds: []*embed.Embed{e.Embed()},
			Flags: discord.MessageFlagEphemeral})

	return true
}
