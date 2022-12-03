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
	ID           string `gorm:"primaryKey"`
	SessionKey   string
	BoardCode    string
	LastPollTime time.Time `gorm:"index"`
}

type LeaderboardEntry struct {
	Name  string `gorm:"primaryKey"`
	Id    uint   `gorm:"primaryKey"`
	Score uint
	Stars uint
	Event string `gorm:"index"`
}

// Api structs
type ApiMember struct {
	Score int    `json:"local_score"`
	Name  string `json:"name"`
	Stars int    `json:"stars"`
	ID    int    `json:"id"`
}

type ApiLeaderboard struct {
	Event   string      `json:"event"`
	Members []ApiMember `json:"members"`
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
	ThemeEmbed(e, ctx)

	// Send response
	ctx.client.Interaction.CreateResponse(ctx.interaction.Id,
		ctx.interaction.Token,
		&discord.InteractionCallbackMessage{Embeds: []*embed.Embed{e.Embed()},
			Flags: discord.MessageFlagUrgent})
}
