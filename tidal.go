package discordtidal

import (
	"discordtidal/log"
	"discordtidal/rpc"
	"discordtidal/song"
	"discordtidal/tidal"
	"github.com/hugolgst/rich-go/client"
	"time"
)

type Status uint8

const (
	Closed Status = iota
	Opened
	Playing
	Paused
)

var (
	sleepTime = time.Second * 5
)

// Start starts the Discord RPC update loop.
func Start() {
	log.Init()
	defer log.Log().Sync()
	rpc.Init()
	defer rpc.Logout()

	for {
		track, artist, status := GetSong()
		if status != Closed {
			sleepTime = time.Second
			rpc.Login()

			if status == Playing {
				// NEW SONG
				nothingPlaying := song.Current == nil
				songChanged := song.Current != nil && (song.Current.Track.Title != track || !song.Current.Track.ArtistMatches(artist))
				looped := song.Current != nil && time.Now().Unix() > int64(song.Current.Track.Duration)+song.Current.StartTime+int64(song.Current.PausedTime)+1
				if nothingPlaying || songChanged || looped {
					// Load song
					now := time.Now()

					var t tidal.Track
					if looped {
						t = song.Current.Track
					} else {
						t = *tidal.GetTrack(track, artist)
					}
					song.Current = &song.Song{
						StartTime:  now.Unix(),
						PausedTime: 0,
						Paused:     false,
						Track:      t,
					}

					if nothingPlaying || (songChanged && !looped) {
					}

					end := time.Unix(int64(uint64(song.Current.Track.Duration)+uint64(song.Current.StartTime)+song.Current.PausedTime)+1, 0)
					err := client.SetActivity(client.Activity{
						Details:    song.Current.Track.Title,
						State:      "by " + song.Current.Track.FormatArtists(),
						LargeImage: "tidal",
						LargeText:  song.Current.Track.Album.Title,
						Timestamps: &client.Timestamps{
							Start: &now,
							End:   &end,
						},
					})
					if err != nil {
						panic(err)
					}
				}

				if song.Current.Paused {
					song.Current.Paused = false

					start := time.Unix(song.Current.StartTime, 0)
					end := time.Unix(int64(uint64(song.Current.Track.Duration)+uint64(song.Current.StartTime)+song.Current.PausedTime)+1, 0)
					err := client.SetActivity(client.Activity{
						Details:    song.Current.Track.Title,
						State:      "by " + song.Current.Track.FormatArtists(),
						LargeImage: "tidal",
						LargeText:  song.Current.Track.Album.Title,
						Timestamps: &client.Timestamps{
							Start: &start,
							End:   &end,
						},
					})
					if err != nil {
						panic(err)
					}
				}
			}

			// Paused a song
			if status == Paused && song.Current != nil {
				song.Current.PausedTime++

				if !song.Current.Paused {
					song.Current.Paused = true
					err := client.SetActivity(client.Activity{
						Details:    song.Current.Track.Title,
						State:      "by " + song.Current.Track.FormatArtists(),
						LargeImage: "tidal",
						LargeText:  song.Current.Track.Album.Title,
					})
					if err != nil {
						panic(err)
					}
				}
			}
		} else {
			rpc.Logout()
		}

		time.Sleep(sleepTime)
	}
}
