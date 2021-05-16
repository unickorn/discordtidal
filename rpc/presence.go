package rpc

import (
	"discordtidal/discord"
	"discordtidal/log"
	"github.com/hugolgst/rich-go/client"
	"os"
	"os/signal"
	"syscall"
)

var (
	loggedIn = false
)

func Init() {
	Login()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		Logout()
		os.Exit(1)
	}()
}

func Login() {
	if !loggedIn {
		err := client.Login(discord.GetConfig().ApplicationId)
		if err != nil {
			log.Log().Fatal(err)
		}
		loggedIn = true

		log.Log().Info("logged into discord")
	}
}

func Logout() {
	if loggedIn {
		client.Logout()
		loggedIn = false

		log.Log().Info("logged out of discord")
	}
}

func Relog() {
	log.Log().Info("reloading presence")

	client.Logout()
	err := client.Login(discord.GetConfig().ApplicationId)
	if err != nil {
		panic(err)
	}
}
