package song

import "github.com/unickorn/discordtidal/tidal"

// Song is a struct that holds information about the current song.
type Song struct {
	StartTime  int64
	PausedTime uint64
	Paused     bool

	Track   *tidal.Track
	AssetId int
}

// Current is the current song.
var Current *Song
