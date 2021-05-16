package discord

import (
	"discordtidal/log"
	"fmt"
	"net/http"
)

// UpdateInfo is the interaction to update general information about a Discord application.
type UpdateInfo struct {
	Name                    string `json:"name"`
	Description             string `json:"description"`
	Icon                    string `json:"icon"`
	InteractionsEndpointUrl string `json:"interactions_endpoint_url"`
	TermsOfServiceUrl       string `json:"terms_of_service_url"`
	PrivacyPolicyUrl        string `json:"privacy_policy_url"`
	EndpointInterface
}

// UpdateName updates the name of the Discord application.
func UpdateName(name string) {
	c := UpdateInfo{
		Icon: "f96bae9313fdfc20cc27429322d3d7e8",
		Name: name,
	}
	run, err := Do(&c, true)
	if err != nil {
		panic(err)
	}
	log.Log().Infof("%s | update name: %s", run.Status, name)
}

func (c *UpdateInfo) url() string {
	return fmt.Sprintf(editApplicationData, config.ApplicationId)
}

func (c *UpdateInfo) method() string {
	return http.MethodPatch
}
