package discord

import (
	"encoding/json"
	"fmt"
	"github.com/unickorn/discordtidal/log"
	"net/http"
)

// UploadAsset is the interaction to upload new assets to a Discord application.
type UploadAsset struct {
	// Image is the image data.
	Image string `json:"image"` // data:image/png;base64,...
	// Name is the name of the Asset we're uploading.
	Name string `json:"name"`
	// Type seems to be 1 most of the time, I'm not sure what that is.
	Type int `json:"type"`
}

// DoUploadAsset uploads a new Asset and returns the response.
func DoUploadAsset(name string, image string) Asset {
	u := UploadAsset{
		Name:  name,
		Image: image,
		Type:  1,
	}
	r, err := Do(&u, true)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	var response Asset
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		panic(err)
	}
	log.Log().Infof("%s | upload asset (id: %s, name: %s)", r.Status, response.ID, response.Name)
	return response
}

func (u *UploadAsset) url() string {
	return fmt.Sprintf(endpointAssets, config.ApplicationId)
}

func (u *UploadAsset) method() string {
	return http.MethodPost
}
