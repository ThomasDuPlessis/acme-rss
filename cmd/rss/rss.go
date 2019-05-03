package main

import (
	"9fans.net/go/acme"
	"bufio"
	"flag"
	"fmt"
	"github.com/mmcdole/gofeed"
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

func setDir(usr *user.User, feed_name string) {
	if len(*dir) < 1 {
		fmt.Println(usr.HomeDir + "/feeds/")
		*dir = usr.HomeDir + "/feeds/"
	}
	feedPath := *dir + feed_name
	if _, err := os.Stat(feedPath); os.IsNotExist(err) {
		fmt.Println("creating: " + feedPath)
		os.MkdirAll(*dir+"/"+cleanFeedName(feed_name), os.ModePerm)
	}
}

func main() {
	args := flag.Args()
	if (len(args) > 0) {
		fmt.Println(args)
	}
	usr, err := user.Current()
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}
	fp := gofeed.NewParser()
	feeds := getFeeds(usr)

	w, err := acme.New()
	if err != nil {
		fmt.Println("error creating acme window")
	}
	tag := "rss"
	w.Name(tag) 
	w.Write("tag", []byte(tag))


	for _, f := range feeds {
		if f == "" {
			continue
		}
		feed, err := fp.ParseURL(f)
		if err != nil {
			fmt.Println(err)
		}
		w.Write("data", []byte(feed.Title + "\n"))
		setDir(usr, feed.Title)

	}

}
