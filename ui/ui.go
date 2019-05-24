package ui

import "9fans.net/go/acme"

func Listen(w *acme.Win) {
	for e := range w.EventChan() {
		switch e.C2 {
			case 'L': // button 3 in body
				w.Ctl("clean")
				//ignore expansions
				if e.OrigQ0 != e.OrigQ1 {
					continue
				}
				charAddr := fmt.Sprintf("#%v", e.OrigQ0)
				// Check if this is home page or feed page
				
		}
	}
}
