package hackernews

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const apiURL = "https://hn.algolia.com/api/v1/search?tags=front_page"

type newsHit struct {
	URL   string
	Title string
	Time  int
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
func GetHNStories() []HNStory {
	resp, httpErr := http.Get(apiURL)
	var jsonString string
	var page frontPage
	var stories []HNStory

	if httpErr != nil || resp.StatusCode != http.StatusOK {
		panic(httpErr)
	}

	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	jsonString = string(bodyBytes)
	parseErr := json.Unmarshal([]byte(jsonString), &page)

	if parseErr != nil {
		panic(parseErr)
	}

	for _, story := range page.Hits {
		stories = append(stories, HNStory{
			Title: story.Title,
			URL:   story.URL,
			Time:  time.Unix(int64(story.Time), 0),
		})
	}

	return stories
}
