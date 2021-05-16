package tidal

import (
	"discordtidal/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Album struct {
	Id         int
	Title      string
	Cover      string
	VideoCover *string
}

func (a *Album) StringId() string {
	return strconv.Itoa(a.Id)
}

func (a *Album) GetCover() []byte {
	url := "https://resources.tidal.com/images/" + strings.Replace(a.Cover, "-", "/", -1) + "/640x640.jpg"
	log.Log().Info(url)

	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	return b
}
