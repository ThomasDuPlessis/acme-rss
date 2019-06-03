package ui

import (
	"9fans.net/go/acme"
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

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


func ListenFeedPage(w *acme.Win) {

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
			// Check if this is home page or feed page
			nw, err := acme.New()
			if err != nil {
				fmt.Printf("could not create win: %v", err)
				continue
			}
			nw.Name("feed")
			nw.Write("tag", []byte(line))
			nw.Ctl("clean")
			go Listen(nw)
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
