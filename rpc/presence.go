package rpc

import (
	"discordtidal/log"
	"github.com/hugolgst/rich-go/client"
	"os"
	"os/signal"
	"syscall"
)

var (
	applicationId = "1295735155817709648"
	loggedIn      = false
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
		err := client.Login(applicationId)
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
