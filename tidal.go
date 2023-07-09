package discordtidal

import (
	"fmt"
	"time"

	"github.com/hugolgst/rich-go/client"
	"github.com/unickorn/discordtidal/discord"
	"github.com/unickorn/discordtidal/log"
	"github.com/unickorn/discordtidal/rpc"
	"github.com/unickorn/discordtidal/song"
	"github.com/unickorn/discordtidal/tidal"
)

type Status uint8

const (
	Closed Status = iota
	Opened
	Playing
	Paused
)

var (
	sleepTime       = time.Second * 5
	coverUpdateTime = 0
	t               *tidal.Track
)

// Start starts the Discord RPC update loop.
func Start() {
	discord.LoadConfig()
	defer log.Log().Sync()
	rpc.Init()
	defer rpc.Logout()
	discord.OpenDb()
	discord.Sync()

	for {
		getSong()

		// if closed, log out
		if status == Closed {
			log.Log().Debugln("Status: CLOSED")
			rpc.Logout()
			sleepTime = time.Second * 5
			time.Sleep(sleepTime)
			continue
		}

		if status == Playing {
			log.Log().Debugln("Status: PLAYING")
			// NEW SONG
			songChanged := song.Current != nil && (song.Current.Track.Title != track || !song.Current.Track.ArtistMatches(artist))
			// song exists, current unix timestamp is bigger than song start time + paused time + duration + 1
			looped := song.Current != nil && time.Now().Unix() > int64(song.Current.Track.Duration)+song.Current.StartTime+int64(song.Current.PausedTime)

			// newly started playing || song changed || looped
			if song.Current == nil || songChanged || looped {
				rpc.Login()
				log.Log().Debugln("[TRIGGER] NEW SONG/CHANGE/LOOP")
				now := time.Now()
				if looped {
					t = song.Current.Track
				} else {
					t = tidal.GetTrack(track, artist)
				}
				song.Current = &song.Song{
					StartTime: now.Unix(),
					Track:     t,
				}

				largeImage := song.Current.Track.Album.StringId()

				trackUrl := fmt.Sprintf("https://listen.tidal.com/album/%d/track/%d", song.Current.Track.Album.ID , song.Current.Track.Id)

				// if song changed, update cover and name
				if songChanged {
					log.Log().Debugln("---- [TRIGGER] SONG CHANGE")
					coverUpdateTime = 0
					if asset := discord.FetchAsset(song.Current.Track.Album); asset == nil {
						// no asset on discord -> wait for cover and set to tidal
						coverUpdateTime = 30
						largeImage = "tidal"
						sleepTime = time.Second
					}
					discord.UpdateName(song.Current.Track.Title)
					// relog to update name instantly
					rpc.Relog()
				}

				log.Log().Infoln("[TRIGGER] BUTTON URL CHANGED TO:", trackUrl)


				// set activity
				end := time.Unix(int64(song.Current.Track.Duration)+song.Current.StartTime+int64(song.Current.PausedTime), 0)
				err := client.SetActivity(client.Activity{
					Details:    "by " + song.Current.Track.FormatArtists(),
					State:      "on " + song.Current.Track.Album.Title,
					LargeImage: largeImage,
					LargeText:  song.Current.Track.Album.Title,
					Timestamps: &client.Timestamps{
						Start: &now,
						End:   &end,
					},
					Buttons: []*client.Button{
						{
							Label: "Listen on TIDAL",
							Url:   trackUrl,
						},
					},
				})
				if err != nil {
					panic(err)
				}
			}

			// tick cover update time
			if coverUpdateTime > 0 {
				coverUpdateTime--
				log.Log().Debugln("Cover update time:", coverUpdateTime)
				if coverUpdateTime == 0 {
					log.Log().Infoln("[TRIGGER] COVER UPDATE")
					coverUpdateTime = -1 // -1 is magic number for "attempt now!"
					// we don't use 0 because we don't want to update cover every second
				}
			}

			// used to be paused || needs cover update
			if song.Current.Paused || coverUpdateTime == -1 {
				// not paused obviously
				song.Current.Paused = false

				log.Log().Debugln("[TRIGGER] COVER UPDATE/UNPAUSE")

				// update cover
				largeImage := song.Current.Track.Album.StringId()
				trackUrl := fmt.Sprintf("https://listen.tidal.com/album/%d/track/%d", song.Current.Track.Album.ID , song.Current.Track.Id)

				discord.Sync()
				// it probably still won't exist because stupid discord cache hasn't updated yet
				a := discord.FetchAsset(song.Current.Track.Album)
				if a == nil {
					coverUpdateTime = 20 // reset timer
					largeImage = "tidal"
					sleepTime = time.Second
				} else {
					coverUpdateTime = 0
					sleepTime = time.Second * 5
				}

				start := time.Unix(song.Current.StartTime, 0)
				end := time.Unix(int64(uint64(song.Current.Track.Duration)+uint64(song.Current.StartTime)+song.Current.PausedTime), 0)
				err := client.SetActivity(client.Activity{
					Details:    "by " + song.Current.Track.FormatArtists(),
					State:      "on " + song.Current.Track.Album.Title,
					LargeImage: largeImage,
					LargeText:  song.Current.Track.Album.Title,
					Timestamps: &client.Timestamps{
						Start: &start,
						End:   &end,
					},
					Buttons: []*client.Button{
						{
							Label: "Listen on TIDAL",
							Url:   trackUrl,
						},
					},
				})
				if err != nil {
					panic(err)
				}
			}
		}

		// Paused a song
		if status == Paused && song.Current != nil {
			rpc.Login()
			log.Log().Debugln("Status: PAUSED")
			song.Current.PausedTime += uint64(sleepTime / time.Second)

			// update cover
			largeImage := song.Current.Track.Album.StringId()

			trackUrl := fmt.Sprintf("https://listen.tidal.com/album/%d/track/%d", song.Current.Track.Album.ID , song.Current.Track.Id)

			// it probably still won't exist because stupid discord cache hasn't updated yet
			a := discord.FetchAsset(song.Current.Track.Album)
			if a == nil {
				largeImage = "tidal"
			}

			sleepTime = time.Second * 2
			if !song.Current.Paused {
				song.Current.Paused = true
				err := client.SetActivity(client.Activity{
					Details:    "by " + song.Current.Track.FormatArtists(),
					State:      "on " + song.Current.Track.Album.Title,
					LargeImage: largeImage,
					LargeText:  song.Current.Track.Album.Title,
					Buttons: []*client.Button{
						{
							Label: "Listen on TIDAL",
							Url:   trackUrl,
						},
					},
				})
				if err != nil {
					panic(err)
				}
			}
		}

		time.Sleep(sleepTime)
	}
}
