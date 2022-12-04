package main

import (
	"fmt"
	"github.com/Goscord/goscord"
	"github.com/Goscord/goscord/discord"
	"github.com/Goscord/goscord/gateway"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"runtime"
	"time"
)

var db *gorm.DB

func main() {
	fmt.Printf(" -> Environment information: \"%s\"\n", runtime.Version())
	fmt.Println("Please send above data in any bug reports or support queries.")
	log.SetFlags(log.Llongfile | log.Ldate | log.Ltime | log.Lmicroseconds)

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Setup database
	databaseUrl := os.Getenv("DATABASE_URL")
	db, err = gorm.Open(postgres.Open(databaseUrl), &gorm.Config{}) // *gorm.DB
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

	// Setup commands map
	commands := make(map[string]Command)

	// Add all commands here:
	commandsList := []Command{
		new(SetupCommand),
		new(LeaderboardCommand),
	}

	// Create client instance
	client := goscord.New(&gateway.Options{
		Token:   os.Getenv("BOT_TOKEN"),
		Intents: gateway.IntentGuilds | gateway.IntentGuildMembers,
	})

	// Setup events
	err = client.On("ready", func() {
		log.Print("Clearing old slash commands")
		cmds, err := client.Application.GetCommands(client.Me().Id, "")
		if err != nil {
			log.Print(err)
		} else {
			for i := range cmds {
				err = client.Application.DeleteCommand(client.Me().Id, "", cmds[i].Id)
				if err != nil {
					log.Print(err)
				}
			}
		}

		log.Print("Registering slash commands")
		for i := range commandsList {
			Register(commandsList[i], client, commands)
		}

		err = client.SetActivity(&discord.Activity{Name: "/setup or /aocrank", Type: discord.ActivityListening})
		if err != nil {
			log.Print(err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	err = client.On("interactionCreate", func(interaction *discord.Interaction) {
		if interaction.Member == nil {
			return
		}

		if interaction.Member.User.Bot {
			return
		}

		cmd := commands[interaction.Data.Name]

		if cmd != nil {
			_ = cmd.Execute(&Context{client: client, interaction: interaction})
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	// Login client :
	if err := client.Login(); err != nil {
		log.Fatal(err)
	}

	// Keep bot running :
	log.Print("Bot started")
	go UpdateThread()
	select {}
}
