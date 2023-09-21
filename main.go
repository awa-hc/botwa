package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Estructura para almacenar información sobre los canales de voz
type VoiceChannelInfo struct {
	UserID    string
	ChannelID string
	GuildID   string
}

var voiceChannels = make(map[string]*VoiceChannelInfo)

func main() {
	// Crear una nueva sesión de Discord
	token := "Bot MTExNjAwOTE2ODk0NDUxNzEyMA.Gtt-lJ.K8W4QZaOyHFydVRVCFOrNF9BG_qpormuig-s9o"
	sess, err := discordgo.New(token)
	if err != nil {
		log.Fatal(err)
	}

	// Agregar manejadores de eventos
	sess.AddHandler(onVoiceStateUpdate)
	sess.AddHandler(onMessageCreate)

	// Definir los intents que necesita el bot
	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Abrir la conexión a Discord
	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	fmt.Println("Bot is running!")

	// Esperar una señal para cerrar el bot
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func onVoiceStateUpdate(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	if vs.ChannelID != "" {
		voiceChannels[vs.UserID] = &VoiceChannelInfo{
			UserID:    vs.UserID,
			ChannelID: vs.ChannelID,
			GuildID:   vs.GuildID,
		}
	} else {
		delete(voiceChannels, vs.UserID)
	}
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if voiceInfo, ok := voiceChannels[m.Author.ID]; ok {
		channelID := voiceInfo.ChannelID
		GuildID := voiceInfo.GuildID
		s.ChannelMessageSend(channelID, m.Content)
		fmt.Printf("El usuario %s está en el canal de voz %s en el servidor %s \n", m.Author.Username, channelID, GuildID)
	} else {
		fmt.Printf("El usuario %s no está en ningún canal de voz\n", m.Author.Username)
	}

	if m.Content == "awa.join" {
		guildID := m.GuildID
		channelID := m.ChannelID
		connectToVC(s, guildID, channelID)
	}

}

func connectToVC(s *discordgo.Session, guildID, channelID string) {
	_, err := s.ChannelVoiceJoin(guildID, channelID, false, false)
	if err != nil {
		fmt.Println("Error al unir al bot al canal de voz:", err)
		return
	}

	fmt.Printf("Bot se ha unido al canal de voz %s en el servidor %s.\n", channelID, guildID)

}
