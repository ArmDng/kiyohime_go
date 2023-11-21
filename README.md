# kiyohime_go

Bot Discord named Kiyohime that can display avatar of a user and a member using [discordgo](https://github.com/bwmarrin/discordgo/tree/master)

## Initialization


### If you want to build it yourself

Install discordgo to your working go environment

`go get github.com/bwmarrin/discordgo`


Import the package

`import "github.com/bwmarrin/discordgo"`

Build it 

`go build`

To use it, you will need 3 variables

```
-DISCORD_TOKEN 
-GUILD_ID (If there is no guild id, the bot will register the commands globally)
-rmcmd (remove commands after shutting down or not. By default, it's true. Type "false" to not remove)
```

For example

`./kiyohime -DISCORD_TOKEN=YOUR TOKEN -GUILD_ID=YOUR GUILD_ID`

