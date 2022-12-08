package main

import (
	"fmt"
	"github.com/Goscord/goscord/discord"
	"github.com/Goscord/goscord/discord/embed"
	"log"
	"sort"
)

type LeaderboardCommand struct{}

func (c *LeaderboardCommand) Name() string {
	return "aocrank"
}

func (c *LeaderboardCommand) Description() string {
	return "View the advent of code leaderboard and, marvel at your high rank!"
}

func (c *LeaderboardCommand) Category() string {
	return "general"
}

func (c *LeaderboardCommand) Options() []*discord.ApplicationCommandOption {
	return []*discord.ApplicationCommandOption{}
}

func (c *LeaderboardCommand) Execute(ctx *Context) bool {
	var gs GuildSettings
	r := db.First(&gs, ctx.interaction.GuildId)
	if r == nil {
		log.Print("Cannot find any matching guilds")
		SendDatabaseError(ctx)
		return false
	}

	entries, err := GetLeaderboard(gs)
	if err != nil {
		log.Print(err)
		SendDatabaseError(ctx)
		return false
	}

	sort.Slice(entries, func(a int, b int) bool {
		return entries[a].Score > entries[b].Score
	})

	e := embed.NewEmbedBuilder()
	message := fmt.Sprintf("Board code: `%s`\n", gs.BoardCode)

	for i, entry := range entries {
		icon := "trophy"
		if i == 0 {
			icon = "first_place"
		} else if i == 1 {
			icon = "second_place"
		} else if i == 2 {
			icon = "third_place"
		}

		message += fmt.Sprintf("`% 4d` :%s: `% 3d` :star: **%s**\n",
			entry.Score,
			icon,
			entry.Stars,
			entry.Name)
	}

	e.SetTitle("Advent of Code Leaderboard")
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
