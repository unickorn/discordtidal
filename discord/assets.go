package discord

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/df-mc/goleveldb/leveldb"
	"github.com/unickorn/discordtidal/log"
	"github.com/unickorn/discordtidal/tidal"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Asset is an image object.
type Asset struct {
	// Name is the "Asset key" of the uploaded Asset. We use album ids here.
	Name string `json:"name"`
	// Type seems to be 1 most (all) of the time, I'm not sure what it is.
	Type int `json:"type"`
	// ID is the identifier of the uploaded Asset, an 18 characters long string.
	ID string `json:"id"`
}

// AssetData is a struct holding an Asset and its data.
type AssetData struct {
	Status     AssetStatus
	Asset      Asset
	LastOpened int
}

type AssetStatus uint8

const (
	assetOK AssetStatus = iota
	assetNew
	assetDeleted
)

type temp struct {
	albumid    string
	assetid    string
	lastopened int
}

var db *leveldb.DB

// OpenDb opens the leveldb database.
func OpenDb() {
	i, err := leveldb.OpenFile("db", nil)
	if err != nil {
		panic(err)
	}
	db = i

	log.Log().Infoln("Opened database")
}

// Sync loads the assets from Discord and synchronizes the data with the database.
// The assets that are on Discord get marked or newly saved as "OK".
// If there are more than 250 assets in total, oldest 10 assets are deleted.
func Sync() {
	log.Log().Infoln("Syncing Discord cache")

	// get all the uploaded assets
	r, err := http.Get(fmt.Sprintf(endpointAssets, config.ApplicationId))
	if err != nil {
		panic(err)
	}
	var all []Asset
	err = json.NewDecoder(r.Body).Decode(&all)
	if err != nil {
		panic(err)
	}

	// set existing assets as OK on leveldb
	existing := make(map[string]byte)
	for _, a := range all {
		existing[a.ID] = 1

		// check whether on leveldb
		d, err := db.Get([]byte(a.Name), nil)
		if d == nil {
			// exists on discord but not on leveldb, so add it and set OK
			err = db.Put([]byte(a.Name), hash(a.ID, assetOK, 0), nil)
		} else {
			// exists on both, so set OK while keeping the last opened time
			_, _, lastOpened := unhash(d)
			log.Log().Debugf("[UPDATE] asset %s exists on both, setting OK!", a.Name)
			err = db.Put([]byte(a.Name), hash(a.ID, assetOK, lastOpened), nil)
		}
		if err != nil {
			panic(err)
		}
	}

	shouldDelete := len(all) > 250
	log.Log().Infof("Loaded %d assets", len(all))

	var allFromDb []temp

	// iterate over the whole db
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		asset, status, lastOpened := unhash(iter.Value())
		// delete ones that discord managed to delete
		if status == assetDeleted {
			if _, ok := existing[asset]; !ok {
				err = db.Delete(iter.Key(), nil)
				if err != nil {
					panic(err)
				}
			}
		}

		// if needed, throw the asset in a slice, so we can sort and delete the oldest ones
		// this is to make sure we don't surpass the 300 asset limit
		if shouldDelete && status == assetOK {
			allFromDb = append(allFromDb, temp{
				albumid:    string(iter.Key()),
				assetid:    asset,
				lastopened: lastOpened,
			})
		}
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		panic(err)
	}

	// sort all assets and delete the oldest ones
	if shouldDelete {
		log.Log().Infoln("Deleting older assets to stay away from discord's asset limit")
		// sort by last opened
		sort.Slice(allFromDb, func(i, j int) bool {
			return allFromDb[i].lastopened > allFromDb[j].lastopened
		})

		go func() {
			for n := 0; n < 10; n++ { // only delete 10 and sleep in between, so we don't get rate limited !
				// send delete to discord
				a := allFromDb[n]
				DoDeleteAsset(a.assetid)

				// mark db as deleted
				err = db.Put([]byte(a.albumid), hash(a.assetid, assetDeleted, a.lastopened), nil)
				if err != nil {
					panic(err)
				}
				time.Sleep(time.Second)
			}
		}()
	}
}

// FetchAsset returns an AssetData only if the album is uploaded and OK, otherwise uploads it.
func FetchAsset(album tidal.Album) *AssetData {
	ad := GetAsset(album.StringId())
	if ad != nil {
		log.Log().Debugln("[FETCH] asset exists on db")
		if ad.Status == assetOK {
			log.Log().Debugln("-- [FETCH] asset OK")
			return ad
		}
		if ad.Status == assetNew {
			log.Log().Debugln("-- [FETCH] asset NEW")
			return nil
		}
	}
	SaveAsset(album)
	return nil
}

// GetAsset attempts to read Asset data from a given album id.
func GetAsset(albumId string) *AssetData {
	h, err := db.Get([]byte(albumId), nil)
	if err != nil {
		return nil
	}
	assetId, status, lastOpened := unhash(h)
	return &AssetData{
		Status: status,
		Asset: Asset{
			Name: albumId,
			Type: 1, // !
			ID:   assetId,
		},
		LastOpened: lastOpened,
	}
}

// SaveAsset uploads an album cover to Discord and adds it to the database with the assetNew status.
func SaveAsset(album tidal.Album) {
	log.Log().Infof("Asset for album %s not found, uploading new...", album.Title)
	a := DoUploadAsset(album.StringId(), "data:image/jpg;base64,"+base64.StdEncoding.EncodeToString(album.CoverImage()))
	err := db.Put([]byte(album.StringId()), hash(a.ID, assetNew, int(time.Now().Unix())), nil)
	if err != nil {
		panic(err)
	}
}

// hash stores the asset id, status and last opened time into a byte array.
func hash(id string, status AssetStatus, lastOpened int) []byte {
	return []byte(fmt.Sprintf("%s:%d:%d", id, status, lastOpened))
}

// asset id, asset status, last opened
func unhash(h []byte) (string, AssetStatus, int) {
	str := string(h)
	split := strings.Split(str, ":")
	st, _ := strconv.ParseInt(split[1], 10, 8)
	lastOpen, _ := strconv.Atoi(split[2])
	return split[0], AssetStatus(st), lastOpen
}
