package main

import (
	"9fans.net/go/acme"
	"flag"
	"fmt"
	"github.com/mmcdole/gofeed"
	"os"
	"os/user"
	"strings"
)

var (
	dir   = flag.String("dir", "", "Directory which stores conf file and content cache")
	w     *acme.Win
	feeds = []string{"http://esr.ibiblio.org/?feed=rss2"}
)

func cleanFeedName(feed_name string) string {
	return strings.Replace(feed_name, " ", "_", -1)
}

func setDir(feed_name string) {
	if len(*dir) < 1 {
		usr, err := user.Current()
		if err != nil {
			fmt.Println("ERROR: ", err)
			return
		}
		fmt.Println(usr.HomeDir + "/feeds/")
		*dir = usr.HomeDir + "/feeds/"
	}
	fmt.Println("creating: " + *dir + "/" + feed_name)
	os.MkdirAll(*dir+"/"+cleanFeedName(feed_name), os.ModePerm)
}

func main() {
	args := flag.Args()
	fmt.Println(args)
	fp := gofeed.NewParser()
	for _, f := range feeds {
		feed, err := fp.ParseURL(f)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(feed.Title)
		setDir(feed.Title)

	}
	w, err := acme.New()
	if err != nil {
		fmt.Println("error creating acme window")
	}
	fmt.Println(w, err)
}
