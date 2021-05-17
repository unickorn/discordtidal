package discordtidal

import (
	"discordtidal/discord"
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
	sleepTime        = time.Second * 5
	needsCoverUpdate = false
	coverUpdateTime  = 0
)

// Start starts the Discord RPC update loop.
func Start() {
	log.Init()
	defer log.Log().Sync()
	discord.LoadConfig()
	rpc.Init()
	defer rpc.Logout()
	discord.OpenDb()
	discord.Sync()

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

					albumId := song.Current.Track.Album.StringId()
					if nothingPlaying || (songChanged && !looped) {
						coverUpdateTime = 0
						needsCoverUpdate = false
						if asset := discord.FetchAsset(song.Current.Track.Album); asset == nil {
							coverUpdateTime = 20
							albumId = "tidal"
						}
						discord.UpdateName(song.Current.Track.Title)
						rpc.Relog()
					}

					end := time.Unix(int64(uint64(song.Current.Track.Duration)+uint64(song.Current.StartTime)+song.Current.PausedTime)+1, 0)
					err := client.SetActivity(client.Activity{
						Details:    "by " + song.Current.Track.FormatArtists(),
						State:      "on " + song.Current.Track.Album.Title,
						LargeImage: albumId,
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

				if coverUpdateTime > 0 {
					coverUpdateTime--
					if coverUpdateTime <= 0 {
						needsCoverUpdate = true
					}
				}

				if song.Current.Paused || needsCoverUpdate {
					song.Current.Paused = false

					albumId := song.Current.Track.Album.StringId()
					if needsCoverUpdate {
						discord.Sync()
						// it probably still won't exist because stupid discord cache hasn't updated yet
						a := discord.FetchAsset(song.Current.Track.Album)
						if a == nil {
							coverUpdateTime = 20
							albumId = "tidal"
						}
						needsCoverUpdate = false
					} else {
						if a := discord.FetchAsset(song.Current.Track.Album); a == nil {
							albumId = "tidal"
						}
					}
					start := time.Unix(song.Current.StartTime, 0)
					end := time.Unix(int64(uint64(song.Current.Track.Duration)+uint64(song.Current.StartTime)+song.Current.PausedTime)+1, 0)
					err := client.SetActivity(client.Activity{
						Details:    "by " + song.Current.Track.FormatArtists(),
						State:      "on " + song.Current.Track.Album.Title,
						LargeImage: albumId,
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

			// Just opened Tidal
			//if status == Opened {
			// to be honest, we shouldn't really do anything on this one
			// it's a mess to update the discord again to the TIDAL name and relog, spotify integration
			// also doesn't show "opened spotify" either

			//err := client.SetActivity(client.Activity{
			//	Details:    "TIDAL",
			//	State:      "Opened",
			//	LargeImage: "tidal",
			//	LargeText:  "TIDAL",
			//})
			//if err != nil {
			//	panic(err)
			//}
			//}

			// Paused a song
			if status == Paused && song.Current != nil {
				song.Current.PausedTime++

				if !song.Current.Paused {
					song.Current.Paused = true
					err := client.SetActivity(client.Activity{
						Details:    "by " + song.Current.Track.FormatArtists(),
						State:      "on " + song.Current.Track.Album.Title,
						LargeImage: song.Current.Track.Album.StringId(),
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
