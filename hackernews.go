package hackernews

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const apiURL = "https://hn.algolia.com/api/v1/search?tags=front_page"

type newsHit struct {
	URL   string
	Title string
	Time  int64 `json:"created_at_i"`
}

type frontPage struct {
	Hits []newsHit
}

// HNStory is the structure for all of the hits on the front page
type HNStory struct {
	URL   string
	Title string
	Time  time.Time
}

// GetHNStories will return an array of HNStory
func GetHNStories() ([]HNStory, error) {
	resp, err := http.Get(apiURL)

	var jsonString string
	var page frontPage
	var stories []HNStory

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HN API: status code %d != 200", resp.StatusCode)
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	jsonString = string(bodyBytes)
	err = json.Unmarshal([]byte(jsonString), &page)

	if err != nil {
		return nil, err
	}

	for _, story := range page.Hits {
		stories = append(stories, HNStory{
			Title: story.Title,
			URL:   story.URL,
			Time:  time.Unix(story.Time, 0),
		})
	}

	return stories, nil
}
