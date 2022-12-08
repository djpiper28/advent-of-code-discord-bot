package main

import (
	"fmt"
	"github.com/Goscord/goscord/discord"
	"github.com/Goscord/goscord/discord/embed"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"log"
	"os"
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

	if len(entries) == 0 {
		log.Print("User has no entries")
		SendError(fmt.Sprintf("`%s` cannot be found.", name), ctx)
		return false
	}

	sort.Slice(entries, func(a int, b int) bool {
		return entries[a].Time.Unix() < entries[b].Time.Unix()
	})

	// Generate plotter data
	scorePoints := make(plotter.XYs, len(entries))
	starPoints := make(plotter.XYs, len(entries))
	for i := range scorePoints {
		scorePoints[i].Y = float64(entries[i].Score)
		starPoints[i].Y = float64(entries[i].Stars)

		scorePoints[i].X = float64(TimeToPlot(entries[i].Time))
		starPoints[i].X = float64(TimeToPlot(entries[i].Time))
	}

	// Add to graph
	p := plot.New()
	p.Title.Text = fmt.Sprintf("%s's Advent of Code Score", name)
	p.Title.TextStyle.Color = HexToRGB(TEXT_COLOUR)
	p.Legend.TextStyle.Color = HexToRGB(TEXT_COLOUR)
	p.BackgroundColor = HexToRGB(BG_COLOUR)

	p.X.Label.Text = "Time"
	p.X.Color = HexToRGB(TEXT_COLOUR)
	p.X.Label.TextStyle.Color = HexToRGB(TEXT_COLOUR)
	p.X.Tick.LineStyle.Color = HexToRGB(TEXT_COLOUR)
	p.X.Tick.Label.Color = HexToRGB(TEXT_COLOUR)

	p.Y.Label.Text = "Score"
	p.Y.Color = HexToRGB(TEXT_COLOUR)
	p.Y.Label.TextStyle.Color = HexToRGB(TEXT_COLOUR)
	p.Y.Tick.LineStyle.Color = HexToRGB(TEXT_COLOUR)
	p.Y.Tick.Label.Color = HexToRGB(TEXT_COLOUR)

	err = plotutil.AddLinePoints(p,
		"Trophies", scorePoints,
		"Stars", starPoints)
	if err != nil {
		log.Print(err)
		SendDatabaseError(ctx)
		return false
	}

	log.Printf("Plotting profile for %s", name)
	LockAndPlot(func() {
		// Save plot
		err = p.Save(PLOT_SIZE, PLOT_SIZE, PLOT_SCRATCH_FILE)
		if err != nil {
			return
		}

		image, err := os.Open(PLOT_SCRATCH_FILE)
		if err != nil {
			return
		}

		// Send plot
		_, err = ctx.client.Channel.SendMessage(ctx.interaction.ChannelId, []*os.File{image})
		log.Print(err)
	})

	if err != nil {
		log.Print(err)
		SendDatabaseError(ctx)
		return false
	}

	// Send embed
	e := embed.NewEmbedBuilder()
	score := entries[len(entries)-1].Score
	stars := entries[len(entries)-1].Stars

	message := fmt.Sprintf("Board code: `%s`\nCurrent scores: % 4d :trophy: % 3d :star:",
		gs.BoardCode,
		score,
		stars)

	e.SetTitle(fmt.Sprintf("\"%s's\" Advent of Code Profile", name))
	e.SetDescription(message)
	e.SetThumbnail(ctx.interaction.Member.User.AvatarURL())
	ThemeEmbed(e, ctx)

	// Send response
	ctx.client.Interaction.CreateResponse(ctx.interaction.Id,
		ctx.interaction.Token,
		&discord.InteractionCallbackMessage{Embeds: []*embed.Embed{e.Embed()},
			Flags: discord.MessageFlagUrgent})

	log.Printf("Plot for %s completed", name)
	return true
}
