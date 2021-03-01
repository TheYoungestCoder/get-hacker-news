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
	Hits    []newsHit
	NbPages int `json:"nBPages"`
}

// HNStory is the structure for all of the hits on the front page
type HNStory struct {
	URL   string
	Title string
	Time  time.Time
}

func getPage(pageNum int) (*frontPage, error) {
	resp, err := http.Get(apiURL)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HN API: status code %d != 200", resp.StatusCode)
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	jsonString := string(bodyBytes)

	var page frontPage
	err = json.Unmarshal([]byte(jsonString), &page)

	if err != nil {
		return nil, err
	}
	return &page, nil
}

// GetHNStories will return an array of HNStory
func GetHNStories() ([]HNStory, error) {

	pageNum := 0
	firstPage, err := getPage(pageNum)

	if err != nil {
		return nil, err
	}

	stories := []HNStory{}
	for _, story := range firstPage.Hits {
		stories = append(stories, HNStory{
			Title: story.Title,
			URL:   story.URL,
			Time:  time.Unix(story.Time, 0),
		})
	}

	totalPages := firstPage.NbPages

	pageNum++
	for ; pageNum < totalPages; pageNum++ {
		page, err := getPage(pageNum)
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
	}
	return stories, nil
}
