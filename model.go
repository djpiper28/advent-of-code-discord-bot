package main

import (
	"github.com/Goscord/goscord/discord"
	"github.com/Goscord/goscord/gateway"
	"time"
)

// Bot database model
type GuildSettings struct {
	ID           uint64 `gorm:"primaryKey"`
	Name         string // A cache to make database navigation easier
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

	client.Application.RegisterCommand(client.Me().Id, "", appCmd)
	commands[cmd.Name()] = cmd
}
