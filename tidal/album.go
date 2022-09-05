package tidal

import (
	"fmt"
	"github.com/unickorn/discordtidal/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Album represents a Tidal album.
type Album struct {
	ID         int
	Title      string
	Cover      string
	VideoCover *string
}

// StringId returns the album's ID as a string.
func (a *Album) StringId() string {
	return strconv.Itoa(a.ID)
}

// CoverImage returns the album's cover image.
func (a *Album) CoverImage() []byte {
	url := fmt.Sprintf("https://resources.tidal.com/images/%s/640x640.jpg", strings.Replace(a.Cover, "-", "/", -1))
	log.Log().Infoln(url)
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
