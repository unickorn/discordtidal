package song

import "discordtidal/tidal"

type Song struct {
	StartTime  int64
	PausedTime uint64
	Paused     bool

	Track   tidal.Track
	AssetId int
}

var Current *Song
