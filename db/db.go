package db

import (
	"9fans.net/go/acme"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mmcdole/gofeed"
	"io/ioutil"
	"os"
	"log"
	"strings"
	"time"
)

var (
	Dir = flag.String("dir", "", "Directory which stores conf file and content cache")
)

func cleanFeedName(feed_name string) string {
	name := strings.Replace(feed_name, " ", "_", -1)
	name = strings.Replace(name, "'", "", -1)
	return *Dir + strings.Replace(name, ":", "", -1) + ".json" 
}

func unCleanFeedName(feed_name string) string {
	name := strings.Replace(feed_name, "_", " ", -1)
	return strings.Replace(name, ".json", " ", -1)
}

func WriteFeed(feed *gofeed.Feed, overwrite bool) {
	data, _ := json.Marshal(feed)
	feedPath := cleanFeedName(feed.Title)
	if _, err := os.Stat(feedPath); os.IsNotExist(err) {
		if err := ioutil.WriteFile(feedPath, data, 0644); err != nil {
			log.Printf("failed to write feed %v to %v", feed.Title, feedPath)
		}
	}
	if overwrite {
		if err := ioutil.WriteFile(feedPath, data, 0644); err != nil {
			log.Printf("failed to write feed %v to %v", feed.Title, feedPath)
		}
	}
}

// Returns the feed as it is on disk
func ReadFeed(feed *gofeed.Feed) (*gofeed.Feed, error) {
	// Read file in to sync with.
	feedPath := cleanFeedName(feed.Title)
	oldFile, _ := ioutil.ReadFile(feedPath)
	oldFeed := &gofeed.Feed{}

	if err := json.Unmarshal([]byte(oldFile), oldFeed); err != nil {
		log.Printf("Could not read in json format from %v: %v", feedPath, err)
		return nil, err
	}
	return oldFeed, nil
}

func getLatestFeedDate(i *gofeed.Item) *time.Time {
	if i.UpdatedParsed != nil && i.UpdatedParsed.After(*i.PublishedParsed) {
		return i.UpdatedParsed
	}
	return i.PublishedParsed
}

func SyncFeeds(w *acme.Win, feeds []string) {
	fp := gofeed.NewParser()
	for _, f := range feeds {
		if f == "" {
			continue
		}
		feed, err := fp.ParseURL(f)
		if err != nil {
			fmt.Println(err)
		}
		w.Write("data", []byte(feed.Title+"\n"))
		WriteFeed(feed, false)

		oldFeed, err := ReadFeed(feed)
		if err != nil {
			return
		}
		lastItemDate := feed.Items[len(feed.Items)-1].PublishedParsed
		var cut int
		for i, item := range oldFeed.Items {
			if lastItemDate.Before(*getLatestFeedDate(item)) {
				cut = i
				break
			}
		}
		feed.Items = append(feed.Items, oldFeed.Items[cut:]...)
		WriteFeed(feed, true)
	}
}


func GetCurrentFeeds() []string {
    files, err := ioutil.ReadDir(*Dir)
    if err != nil {
        log.Fatal(err)
    }
	var currentFeeds []string
    for _, f := range files {
		actualName := unCleanFeedName(f.Name())
		currentFeeds = append(currentFeeds, actualName)
	}
	return currentFeeds
}