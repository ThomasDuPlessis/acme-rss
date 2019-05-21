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

func SyncFeed(f gofeed.Feed) {
	fmt.Printf("%v %v\n", f.UpdatedParsed, len(f.Items))
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
		SyncFeed(*feed)
		w.Write("data", []byte(feed.Title+"\n"))
		file, _ := json.Marshal(feed)
		feedPath := cleanFeedName(feed.Title)
		if _, err := os.Stat(feedPath); os.IsNotExist(err) {
			_ = ioutil.WriteFile(feedPath, file, 0644)
		}
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