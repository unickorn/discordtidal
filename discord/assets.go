package discord

import (
	"discordtidal/log"
	"discordtidal/tidal"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/df-mc/goleveldb/leveldb"
	"net/http"
	"sort"
	"strconv"
	"time"
)

// Asset is an image object
type Asset struct {
	// Name is the "Asset key" of the uploaded Asset. We use album ids here.
	Name string `json:"name"`
	// Type seems to be 1 most (all) of the time, I'm not sure what it is.
	Type int `json:"type"`
	// Id is the identifier of the uploaded Asset, an 18 characters long string.
	Id string `json:"id"`
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

	log.Log().Info("opened database")
}

// Sync loads the assets from Discord and synchronizes the data with the database.
// The assets that are on Discord get marked or newly saved as "OK".
// If there are more than 250 assets in total, oldest 10 assets are deleted.
func Sync() {
	log.Log().Info("syncing discord")

	t := 0
	// get all of the uploaded ones, set them as OK
	r, err := http.Get(fmt.Sprintf("https://discord.com/api/v9/oauth2/applications/%s/assets", config.ApplicationId))
	if err != nil {
		panic(err)
	}
	var all []Asset
	err = json.NewDecoder(r.Body).Decode(&all)
	if err != nil {
		panic(err)
	}

	existing := make(map[string]byte)
	for _, a := range all {
		t++
		existing[a.Id] = 1
		d, err := db.Get([]byte(a.Name), nil)
		if d == nil {
			err = db.Put([]byte(a.Name), hash(a.Id, assetOK, 0), nil)
		} else {
			_, _, lastOpened := unhash(d)
			err = db.Put([]byte(a.Name), hash(a.Id, assetOK, lastOpened), nil)
		}
		if err != nil {
			panic(err)
		}
	}

	shouldDelete := t > 250
	log.Log().Infof("loaded %d assets", t)

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

		// if it's needed, throw the Asset in a slice so we can sort and delete oldest ones
		// this is to make sure we don't surpass the 300 Asset limit
		if shouldDelete && status == assetOK {
			allFromDb = append(allFromDb, temp{
				albumid:    string(iter.Key()),
				assetid:    asset,
				lastopened: lastOpened,
			})
		}
	}
	iter.Release()

	// sort all assets and delete the oldest ones
	if shouldDelete {
		log.Log().Info("deleting stuff")
		sort.Slice(allFromDb, func(i, j int) bool {
			return allFromDb[i].lastopened > allFromDb[j].lastopened
		})

		go func() {
			for n := 0; n < 10; n++ { // only delete 10 so we don't get rate limited !!!!
				// send delete to discord!!!!
				a := allFromDb[n]
				DoDeleteAsset(a.assetid)

				// delete from db too !!!!!
				err = db.Put([]byte(a.albumid), hash(a.assetid, assetDeleted, a.lastopened), nil)
				if err != nil {
					panic(err)
				}

				time.Sleep(time.Second * 2)
			}
		}()
	}

	err = iter.Error()
	if err != nil {
		panic(err)
	}
}

// FetchAsset returns an AssetData only if the album is uploaded and OK, otherwise uploads it.
func FetchAsset(album tidal.Album) *AssetData {
	ad := GetAsset(album.StringId())
	if ad != nil {
		if ad.Status == assetOK {
			return ad
		}
		if ad.Status == assetNew {
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
			Id:   assetId,
		},
		LastOpened: lastOpened,
	}
}

// SaveAsset uploads an album cover to Discord and adds it to the database with the Status "new".
func SaveAsset(album tidal.Album) {
	log.Log().Infof("asset for album %s not found, uploading new...", album.Title)
	a := DoUploadAsset(album.StringId(), "data:image/jpg;base64,"+base64.StdEncoding.EncodeToString(album.GetCover()))
	err := db.Put([]byte(album.StringId()), hash(a.Id, assetNew, int(time.Now().Unix())), nil)
	if err != nil {
		panic(err)
	}
}

// this is bad but it's my first time using leveldb ok
// plus it doesn't really matter we don't need perfect performance
func hash(id string, status AssetStatus, lastOpened int) []byte {
	return []byte(id + strconv.Itoa(int(status)) + strconv.Itoa(lastOpened))
}

// asset id, asset status, last opened
func unhash(h []byte) (string, AssetStatus, int) {
	str := string(h)
	s, _ := strconv.Atoi(str[18:19])
	l, _ := strconv.Atoi(str[19 : len(h)-1])
	return str[0:18], AssetStatus(s), l
}
