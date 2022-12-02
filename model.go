package main

import (
	"time"
)

// Bot database model
type GuildSettings struct {
	ID           uint64
	Name         string // A cache to make database navigation easier
	SessionKey   string
	BoardCode    string
	LastPollTime time.Time
}

type LeaderboardEntry struct {
	Name  string
	Id    uint
	Score uint
	Stars uint
	Event string
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
