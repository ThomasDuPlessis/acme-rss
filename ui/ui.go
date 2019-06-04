package ui

import (
	"9fans.net/go/acme"
	"bufio"
	"errors"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/thomasduplessis/acme-rss/db"
	"io"
	"os"
	"os/user"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	Feeds = map[string]gofeed.Feed{}
	htmlre = regexp.MustCompile("<br( /)*>")
)


func SetFeeds(feeds []gofeed.Feed) {
	for _, f := range feeds {
		Feeds[strings.Trim(f.Title, " ")] = f
	}
}

type WinReader struct {
	w *acme.Win
}

func (win WinReader) Read(b []byte) (int, error) {
	return win.w.Read("body", b)
}

func readLine(w *acme.Win, charaddr int) (string, error) {
	w.Seek("body", 0, 0)
	br := bufio.NewReader(&WinReader{w: w})
	offset := 0
	for {
		line, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", err
		}
		offset += utf8.RuneCountInString(line)
		if offset > charaddr {

			return line, nil
		}
		if err == io.EOF {
			return "", errors.New(fmt.Sprintf("%v is past total charcount: %v", charaddr, offset))
		}
	}
	panic("unreachable")
}


func ListenFeedPage(w *acme.Win, feed *gofeed.Feed) {
	for e := range w.EventChan() {
		switch e.C2 {
		case 'L': // button 3 in body
			//ignore expansions
			if e.OrigQ0 != e.OrigQ1 {
				continue
			}
			line, err := readLine(w, e.OrigQ0)
			if err != nil {
				fmt.Printf("%v", err)
				continue
			}
			line = strings.Trim(line, " \n")
			// Check if this is home page or feed page
			nw, err := acme.New()
			if err != nil {
				fmt.Printf("could not create win: %v", err)
				continue
			}
			found := false
			for _, item := range feed.Items {
				if item.Title == line {
					content := htmlre.ReplaceAllString(item.Description +"\n" +  item.Content, "\n")
					nw.Write("body", []byte(content))
					found = true
					break
				}
			}
			if !found {
				var titles []string
				for _, item := range feed.Items {
					titles = append(titles, item.Title)
				}
				fmt.Printf("%v not found in: %v", line, strings.Join(titles, "\n"))
				nw.Ctl("delete")
				nw.CloseFiles()
			}
			nw.Name("feedItem")
			nw.Write("tag", []byte(line))
			nw.Ctl("clean")
		case 'x': // exec
			switch string(e.Text) {
			case "Del":
				w.Ctl("delete")
				w.CloseFiles()
				return
			}
		}
	}
}

func Listen(w *acme.Win) {
	for e := range w.EventChan() {
		switch e.C2 {
		case 'L': // button 3 in body
			//ignore expansions
			if e.OrigQ0 != e.OrigQ1 {
				continue
			}
			line, err := readLine(w, e.OrigQ0)
			if err != nil {
				fmt.Printf("%v", err)
				continue
			}
			line = strings.Trim(line, " \n")
			// Check if this is home page or feed page
			nw, err := acme.New()
			if err != nil {
				fmt.Printf("could not create win: %v", err)
				continue
			}
			nw.Name("feed")
			nw.Write("tag", []byte(line))
			var titles []string
			
			if feed, ok := Feeds[line]; ok {
				for _, item := range feed.Items {
					titles = append(titles, item.Title)
				}
				nw.Write("body",[]byte(strings.Join(titles, "\n")))
				go ListenFeedPage(nw, &feed)
			} else {
				var names []string
				for name := range Feeds {
					names = append(names, name)
				}
				fmt.Printf("Could not find %v in %v", line, strings.Join(names, ", "))
			}
			nw.Ctl("clean")
		case 'x': // exec
			switch string(e.Text) {
			case "Del":
				w.Ctl("delete")
				w.CloseFiles()
				return
			case "Refresh":
				go Refresh(w)
				return
			}
		}
	}
}

func SetDir(usr *user.User) {
	if len(*db.Dir) < 1 {
		*db.Dir = usr.HomeDir + "/feeds/"
	}
}

func getFeeds(usr *user.User) []string {
	file, err := os.Open(usr.HomeDir + "/.feeds")
	if err != nil {
		fmt.Printf("%v", err)
		return []string{}
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}


func Refresh(w *acme.Win) {
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	feeds := getFeeds(usr)
	db.SyncFeeds(feeds)
	SetFeeds(db.ReadInFeedsOnDisk())
	currentFeeds := db.GetCurrentFeeds()
	w.Write("body", []byte(strings.Join(currentFeeds, "\n")))
}