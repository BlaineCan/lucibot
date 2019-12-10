package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
	luciFound bool
	luciID string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	luciFound = false
	luciID = "165335653515526146"

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(onReady)
	dg.AddHandler(guildCreate)
	dg.AddHandler(presenceUpdate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bonfire LIT.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.ID == luciID {
		if rand.Intn(15-1)+1 != 1 {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		}
	}
}

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	botStatus := discordgo.UpdateStatusData{new(int), new(discordgo.Game), false, "invisible"}
	s.UpdateStatusComplex(botStatus)
}

func guildCreate(s *discordgo.Session, g *discordgo.GuildCreate) {
	for _, p := range g.Guild.Presences {
		if !luciFound && p.User.ID == luciID {
			luciFound = true
			botStatus := discordgo.UpdateStatusData{p.Since, p.Game, false, string(p.Status)}
			s.UpdateStatusComplex(botStatus)
		}
	}
}

func presenceUpdate(s *discordgo.Session, p *discordgo.PresenceUpdate) {
	if p.User.ID == luciID {
		luciStatus := string(p.Status)
		if luciStatus == "offline" {
			luciStatus = "invisible"
		}
		botStatus := discordgo.UpdateStatusData{p.Since, p.Game, false, luciStatus}
		s.UpdateStatusComplex(botStatus)
	}
}
