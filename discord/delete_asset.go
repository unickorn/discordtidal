package discord

import (
	"fmt"
	"github.com/unickorn/discordtidal/log"
	"net/http"
)

// DeleteAsset is the interaction to delete an Asset from the Discord application.
type DeleteAsset struct {
	// id is the identifier of the Asset to be deleted.
	id string
}

// DoDeleteAsset deletes a given Asset using an id.
func DoDeleteAsset(id string) {
	d := DeleteAsset{id: id}
	r, err := Do(&d, false)
	if err != nil {
		panic(err)
	}
	log.Log().Infof("%s | delete asset: %s", r.Status, id)
}

func (d *DeleteAsset) url() string {
	return fmt.Sprintf(endpointDeleteAsset, config.ApplicationId, d.id)
}

func (d *DeleteAsset) method() string {
	return http.MethodDelete
}
