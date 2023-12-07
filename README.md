# advent-of-code-discord-bot

Instead of doing advent of code I did this. idk why.
[Discord Invite Link](https://discord.com/api/oauth2/authorize?client_id=1047611666604503061&permissions=8&scope=bot%20applications.commands)

## Setting Up The `.env` File

```sh 
METRICS_SERVER=localhost:6563
DATABASE_URL=postgres:// ...
BOT_TOKEN=my discord bot token
ENABLE_POLLING=true
PROXY="socks5://localhost:9050"
TOR_CONTROLLER="localhost:9051"
```

## Postgres?

You can setup Postgres as a docker image, or setup Postgres locally. Both should work just fine.
