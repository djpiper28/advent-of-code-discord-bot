package main

import (
	"github.com/Goscord/Bot/command"
	"github.com/Goscord/Bot/config"
	"github.com/Goscord/Bot/event"
	"github.com/Goscord/goscord"
	"github.com/Goscord/goscord/gateway"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Setup database
	databaseUrl := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Print("Migrating Database")
	db.AutoMigrate(&GuildSettings{})
	db.AutoMigrate(&LeaderboardEntry{})

	// Config stuffs
	Config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create client instance :
	client := goscord.New(&gateway.Options{
		Token:   os.Getenv("BOT_TOKEN"),
		Intents: gateway.IntentGuilds | gateway.IntentGuildMembers,
	})

	cmdMgr := command.NewCommandManager(client, Config)

	err = client.On("ready", event.OnReady(client, Config, cmdMgr))
	if err != nil {
		log.Fatal(err)
	}

	err = client.On("interactionCreate", cmdMgr.Handler(client, Config))
	if err != nil {
		log.Fatal(err)
	}

	// Login client :
	if err := client.Login(); err != nil {
		log.Fatal(err)
	}

	// Keep bot running :
	select {}
}
