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

To use it, you will need 2 variables

```
-DISCORD_TOKEN=YOUR_TOKEN
-RMCMD=true (if you want to rm slash commands after extinction. If not, don't put it)

```

For example

` DISCORD_TOKEN=YOUR TOKEN ./kiyohime  `

### Dockerfile

To build your image, do :

`docker build -t name_of_your_image:multistage .`

Don't forget to put your discord token to run it