package main

import (
	"github.com/Goscord/Bot/command"
	"github.com/Goscord/Bot/config"
	"github.com/Goscord/Bot/event"
	"github.com/Goscord/goscord"
	"github.com/Goscord/goscord/gateway"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	client *gateway.Session
	Config *config.Config
	cmdMgr *command.CommandManager
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	databaseUrl := os.Getenv("DATABASE_URL")

	// Load envionment variables :
	godotenv.Load()

	Config, err = config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create client instance :
	client = goscord.New(&gateway.Options{
		Token:   os.Getenv("BOT_TOKEN"),
		Intents: gateway.IntentGuilds | gateway.IntentGuildMessages | gateway.IntentGuildMembers,
	})
}
