package main

import (
	"9fans.net/go/acme"
	"bufio"
	"flag"
	"fmt"
	//	"github.com/mmcdole/gofeed"
	"github.com/thomasduplessis/acme-rss/db"
	"github.com/thomasduplessis/acme-rss/ui"
	"log"
	"os"
	"os/user"
	"strings"
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
		fmt.Println("%v", err)
	}
	w.Name("+rss")
	w.Write("tag", []byte("Refresh"))
	go db.SyncFeeds(w, feeds)
	currentFeeds := db.GetCurrentFeeds()
	w.Write("body", []byte(strings.Join(currentFeeds, "\n")))
	ui.Listen(w)
	w.Ctl("clean")
}
