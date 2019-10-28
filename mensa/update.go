package mensa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/inconshreveable/go-update"
)

func doUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		// error handling
	}
	return err
}

//GetURLLatestRelease gets the url of the latest release in github
func GetURLLatestRelease() string {

	resp, err := http.Get("https://api.github.com/repos/pfaaj/go-mensa-kl/releases/latest")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(body, &objmap)

	var tagName string

	err = json.Unmarshal(*objmap["tag_name"], &tagName)

	if err != nil {
		return ""
	}

	return fmt.Sprintf("https://github.com/pfaaj/go-mensa-kl/releases/download/%s/go-mensa", tagName)

}

//IsNewRelease returns true if there is a new binary release at github
func IsNewRelease(release string) bool {

	return LatestRelease != release
}
