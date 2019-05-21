package main

import (
	"9fans.net/go/acme"
	"bufio"
	"flag"
	"fmt"
	//	"github.com/mmcdole/gofeed"
	"github.com/thomasduplessis/acme-rss/db"
	"log"
	"os"
	"os/user"
)

var (
	w   *acme.Win
)

func setDir(usr *user.User) {
	if len(*db.Dir) < 1 {
		*db.Dir = usr.HomeDir + "/feeds/"
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
	db.SyncFeeds(w, feeds)
	currentFeeds := db.GetCurrentFeeds()
    for _, feedname := range currentFeeds {
		w.Write("data", []byte(feedname + "\n"))
    }
}
