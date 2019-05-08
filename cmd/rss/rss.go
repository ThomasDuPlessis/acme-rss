package main

import (
	"9fans.net/go/acme"
	"bufio"
	"flag"
	"fmt"
	"github.com/mmcdole/gofeed"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

var (
	dir = flag.String("dir", "", "Directory which stores conf file and content cache")
	w   *acme.Win

//	feeds = []string{"http://esr.ibiblio.org/?feed=rss2"}
)

func cleanFeedName(feed_name string) string {
	return strings.Replace(feed_name, " ", "_", -1)
}

func unCleanFeedName(feed_name string) string {
	return strings.Replace(feed_name, "_", " ", -1)
}

func setDir(usr *user.User, ) {
	if len(*dir) < 1 {
		*dir = usr.HomeDir + "/feeds/"
	}
}

func getFeeds(usr *user.User) []string {
	file, err := os.Open(usr.HomeDir + "/.feeds")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func createDir(feed_name string) {

	feedPath := *dir + cleanFeedName(feed_name)
	if _, err := os.Stat(feedPath); os.IsNotExist(err) {
		fmt.Println("creating: " + feedPath)
		os.MkdirAll(*dir+"/"+cleanFeedName(feed_name), os.ModePerm)
	}
}

func syncFeeds(w *acme.Win, feeds []string) {
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
		createDir(feed.Title)
	}
}

func main() {
	args := flag.Args()
	if len(args) > 0 {
		fmt.Println(args)
	}
	usr, err := user.Current()
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}
	setDir(usr)
	feeds := getFeeds(usr)

	w, err := acme.New()
	if err != nil {
		fmt.Println("error creating acme window")
	}
	w.Write("tag", []byte("rss"))
	go syncFeeds(w, feeds)
    files, err := ioutil.ReadDir(*dir)
    if err != nil {
        log.Fatal(err)
    }
	var currentFeeds []string
    for _, f := range files {
		actualName := unCleanFeedName(f.Name())
		w.Write("data", []byte(actualName + "\n"))
		currentFeeds = append(currentFeeds, actualName)
    }

}
