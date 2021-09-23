package tidal

import (
	"discordtidal/log"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Album struct {
	Id         int
	Title      string
	Cover      string
	VideoCover *string
}

type Artist struct {
	Id         int
	Name       string
	ArtistType string
}

type Track struct {
	Id           int
	Title        string
	Duration     int
	TrackNumber  int
	VolumeNumber int
	AudioQuality string
	Artist       Artist
	Artists      []Artist
	Album        Album
}

// ArtistMatches returns whether the artist matches the artist of the song by comparing
// it in many formats to match songs with multiple artists.
func (t *Track) ArtistMatches(artist string) bool {
	if artist == t.FormatArtists() {
		return true
	}

	for _, a := range t.Artists {
		if !strings.Contains(artist, a.Name) {
			return false
		}
	}
	return true
}

// FormatArtists formats names of the song's artists into one string.
func (t *Track) FormatArtists() string {
	var names []string
	for _, a := range t.Artists {
		names = append(names, a.Name)
	}
	return strings.Join(names, ", ")
}

type SearchResult struct {
	Limit              int
	Offset             int
	TotalNumberOfItems int
	Items              []Track
}

// GetTrack returns a Track from Tidal with given song and artist names.
func GetTrack(songName string, artistName string) *Track {
	log.Log().Infof("searching %s by %s", songName, artistName)
	link := fmt.Sprintf("https://api.tidal.com/v1/search/tracks?countryCode=US&query=%s&limit=15", url.QueryEscape(songName+" "+artistName))

	cl := http.DefaultClient
	req, _ := http.NewRequest("GET", link, nil)
	req.Header.Set("x-tidal-token", "CzET4vdadNUFQ5JU")
	resp, err := cl.Do(req)
	if err != nil {
		panic(err)
	}
	var res SearchResult
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Log().Fatal(err)
	}

	var tTrack *Track
	//var maxQuality AudioQualityLevel
	//found := false
	for _, track := range res.Items {
		if track.Title == songName && track.FormatArtists() == artistName {
			tTrack = &track
			break
			// we don't need quality stuff for now, maybe later for some small icon indicator

			//log.Log().Infof(" | %s", track.AudioQuality)
			//ql := GetAudioQualityLevel(track.AudioQuality)
			//if found == false || ql > maxQuality {
			//	maxQuality = ql
			//	tTrack = track
			//	found = true
			//	if ql == Master {
			//		break
			//	}
			//}
		}
	}

	// just in case we couldn't find anything
	if tTrack == nil {
		tTrack = &res.Items[len(res.Items)-1]
	}

	log.Log().Infof("result: %s by %s", tTrack.Title, tTrack.FormatArtists())
	return tTrack
}
