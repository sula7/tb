package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/clbanning/mxj/v2"

	"tb/storage"
)

func GetTagContent(url, tag string) ([]storage.NewsItem, error) {
	respBytes, err := makeHTTPReq(url)
	if err != nil {
		return nil, err
	}

	var newsItems []storage.NewsItem

	re, err := regexp.Compile(`<channel>[\s\S]*?<\/channel>`)
	if err != nil {
		return nil, err
	}

	respChannel := re.FindString(string(respBytes))

	xmlMap := mxj.New()
	xmlMap, err = mxj.NewMapXml([]byte(respChannel))

	channel := xmlMap["channel"]
	for _, item := range channel.(map[string]interface{})["item"].([]interface{}) {
		tagContent := item.(map[string]interface{})[tag]

		switch tagContent.(type) {
		case map[string]interface{}:
			tagContent = tagContent.(map[string]interface{})["#text"] //this is "mxj" lib specified for correct <guid> process
		}

		re, err := regexp.Compile(`<.*?>`)
		if err != nil {
			return nil, err
		}

		newsItem := storage.NewsItem{
			URL:        url,
			TagContent: strings.TrimSpace(re.ReplaceAllString(fmt.Sprint(tagContent), "")),
		}

		newsItems = append(newsItems, newsItem)
	}

	return newsItems, err
}

func makeHTTPReq(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
