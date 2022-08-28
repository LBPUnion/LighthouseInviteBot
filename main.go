package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "embed"

	"github.com/Zaprit/LighthouseInviteBot/common"
	"github.com/Zaprit/LighthouseInviteBot/discordbot"
	lighthouseapi "github.com/Zaprit/LighthouseInviteBot/lighthouseAPI"
	"github.com/bwmarrin/discordgo"
	"github.com/coreos/go-systemd/daemon"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

//go:embed lighthouseBot.toml.default
var defaultConfig string

func initLog() {
	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
}

func main() {

	// note to self, don't leak discord token :P
	_, err := os.Stat("lighthouseBot.toml")
	if err != nil {
		file, err := os.Create("lighthouseBot.toml")
		if err != nil {
			logrus.WithError(err).Panicln("Failed to create lighthouseBot config")
		}
		_, er2 := file.WriteString(defaultConfig)

		if er2 != nil {
			logrus.WithError(err).Panicln("Failed to write default lighthouseBot config")
		}
		er3 := file.Close()
		if er3 != nil {
			logrus.WithError(err).Panicln("Failed to close lighthouseBot config")
		}
		os.Exit(1)
	}

	s, err := discordbot.CreateBot(common.LoadConfig().Bot.Token)
	if err != nil {
		logrus.WithError(err).Panicln("Failed to create Discord bot")
	}

	defer s.Close()

	go randomStatus(s)

	daemon.SdNotify(false, daemon.SdNotifyReady)

	logrus.Infoln("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	daemon.SdNotify(true, daemon.SdNotifyStopping)
	s.UpdateStatusComplex(discordgo.UpdateStatusData{Status: "offline"})

	logrus.Infoln("Shutting down...")

}

func randomStatus(s *discordgo.Session) {
	stat := 0
	for {
		stats, err := lighthouseapi.GetStatistics()
		if err != nil {
			logrus.Warnln("Failed to get lighthouse statistics")
			continue
		}

		var statusText string

		switch stat {
		case 0:
			statusText = strconv.Itoa(stats.RecentMatches) + " People Online"
		case 1:
			statusText = strconv.Itoa(stats.Slots) + " Levels"
		case 2:
			statusText = strconv.Itoa(stats.Users) + " Users"
		case 3:
			statusText = strconv.Itoa(stats.TeamPicks) + " Team Picks"
		case 4:
			statusText = strconv.Itoa(stats.Photos) + " Photos"
		default:
			stat = 0
			continue
		}
		idleSince := 0
		//s.UpdateStatusComplex(*newUpdateStatusData(idle, ActivityTypeGame, name, ""))
		s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Status: "online",
			Activities: []*discordgo.Activity{
				{
					Type:    discordgo.ActivityTypeGame,
					Name:    statusText,
					State:   statusText,
					Details: statusText,

					URL: "",
				},
			},
			IdleSince: &idleSince,
		})

		stat++
		if stat > 4 {
			stat = 0
		}

		time.Sleep(30 * time.Second)
	}
}
