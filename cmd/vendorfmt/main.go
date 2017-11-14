package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/magiconair/vendorfmt"
)

func main() {
	log.SetFlags(0)

	files := os.Args[1:]

	if len(files) == 0 {
		files = []string{"vendor/vendor.json"}
	}

	for _, f := range files {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			log.Print("Error reading ", f)
			continue
		}
		b, err = vendorfmt.Format(b)
		if err != nil {
			log.Print("Error formatting %s: %s", f, err)
			continue
		} else {
			// Add the ending newline because json.Marshall[Indent]()
			// strips it.
			b = append(b, '\n')
		}
		if err := ioutil.WriteFile(f, b, 0644); err != nil {
			log.Print("Error writing %s: %s", f, err)
			continue
		}
		log.Print("Formatted ", f)
	}
}
