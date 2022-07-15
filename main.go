package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "embed"

	"github.com/Zaprit/LighthouseInviteBot/common"
	"github.com/Zaprit/LighthouseInviteBot/discordbot"
)

//go:embed lighthouseBot.toml.default
var defaultConfig string

func main() {
	// note to self, don't leak discord token :P
	_, err := os.Stat("lighthouseBot.toml")
	if err != nil {
		file, err := os.Create("lighthouseBot.toml")
		if err != nil {
			panic(err.Error())
		}
		_, er2 := file.WriteString(defaultConfig)

		if er2 != nil {
			panic(er2.Error())
		}
		er3 := file.Close()
		if er3 != nil {
			panic(er3.Error())
		}
		os.Exit(1)
	}

	discordbot.CreateBot(common.LoadConfig().Bot.Token)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("\nShutting down...")

}
