package autoproxy

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var listurl string = "http://www.firefoxfan.com/gfwlist/gfwlist.txt"

var localfile string = "gfwlist.txt"

func downloadList() error {
	resp, err := http.Get(listurl)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("Download failed")
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(localfile, content, 0755)

	return err
}

func GetRawListContent() (string, error) {
	content, err := ioutil.ReadFile(localfile)
	if err != nil {
		log.Println("cannot read from " + localfile + ", try download it")
		downloadList()

		content, err = ioutil.ReadFile(localfile)
		if err != nil {
			return "", err
		}
	}
	return string(content), nil
}

func GetListContent() (string, error) {
	raw, err := GetRawListContent()
	if err != nil {
		return "", err
	}
	content, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func TestURL(u *url.URL) bool {
	return true
}
