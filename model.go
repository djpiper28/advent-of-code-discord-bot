package main

import (
	"github.com/Goscord/goscord/discord"
	"github.com/Goscord/goscord/discord/embed"
	"github.com/Goscord/goscord/gateway"
	"log"
	"time"
)

// Bot database model
type GuildSettings struct {
	ID         string `gorm:"primaryKey"`
	SessionKey string
	BoardCode  string
	Year       string
}

type LeaderboardEntry struct {
	PK        string `gorm:"primaryKey"`
	Name      string `gorm:"index"`
	ID        int    `gorm:"index"`
	Score     int
	Stars     int
	Event     string    `gorm:"index"`
	Time      time.Time `gorm:"index"`
	BoardCode string    `gorm:"index"`
}

/*
type ApiCompletionDayLevel struct {
	StarIndex int `json:"star_index"`
	StarTime  int `json:"get_star_ts"`
}
*/

// Api structs
type ApiMember struct {
	Score int    `json:"local_score"`
	Name  string `json:"name"`
	Stars int    `json:"stars"`
	ID    int    `json:"id"`

	/*
		// Unused at the moment
		CompletionDayLevel map[string]map[string]ApiCompletionDayLevel `json:"completion_day_level,omitempty"`
		GlobalScore        int                                         `json:"global_score,omitempty"`
		LastStarTime       int                                         `json:"last_star_ts"`
	*/
}

type ApiLeaderboard struct {
	Event   string               `json:"event"`
	Members map[string]ApiMember `json:"members"`
	OwnerId int                  `json:"owner_id"`
}

type Context struct {
	client      *gateway.Session
	interaction *discord.Interaction
}

type Command interface {
	Name() string
	Description() string
	Category() string
	Options() []*discord.ApplicationCommandOption
	Execute(ctx *Context) bool
}

func Register(cmd Command, client *gateway.Session, commands map[string]Command) {
	appCmd := &discord.ApplicationCommand{
		Name:        cmd.Name(),
		Type:        discord.ApplicationCommandChat,
		Description: cmd.Description(),
		Options:     cmd.Options(),
	}

	_, err := client.Application.RegisterCommand(client.Me().Id, "", appCmd)
	if err != nil {
		log.Print(err)
	}
	commands[cmd.Name()] = cmd
}

func ThemeEmbed(e *embed.Builder, ctx *Context) {
	e.SetFooter(ctx.client.Me().Username, ctx.client.Me().AvatarURL())
	e.SetColor(embed.Green)
}

func SendDatabaseError(ctx *Context) {
	e := embed.NewEmbedBuilder()

	e.SetTitle("An Error Occurred During Your Command")
	e.SetDescription("A database error occured.")
	ThemeEmbed(e, ctx)

	// Send response
	ctx.client.Interaction.CreateResponse(ctx.interaction.Id,
		ctx.interaction.Token,
		&discord.InteractionCallbackMessage{Embeds: []*embed.Embed{e.Embed()},
			Flags: discord.MessageFlagEphemeral})
}

func SendPermissionsError(ctx *Context) {
	e := embed.NewEmbedBuilder()

	e.SetTitle("This Command Requires Administrator Permissions To Run")
	ThemeEmbed(e, ctx)

	// Send response
	ctx.client.Interaction.CreateResponse(ctx.interaction.Id,
		ctx.interaction.Token,
		&discord.InteractionCallbackMessage{Embeds: []*embed.Embed{e.Embed()},
			Flags: discord.MessageFlagEphemeral})
}

func SendError(message string, ctx *Context) {
	e := embed.NewEmbedBuilder()

	e.SetTitle("An Error Occurred During Your Command")
	e.SetDescription(message)
	ThemeEmbed(e, ctx)

	// Send response
	ctx.client.Interaction.CreateResponse(ctx.interaction.Id,
		ctx.interaction.Token,
		&discord.InteractionCallbackMessage{Embeds: []*embed.Embed{e.Embed()},
			Flags: discord.MessageFlagEphemeral})
}
