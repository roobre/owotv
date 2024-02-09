package main

import (
	"context"
	"log"
	"time"
)

const minDisplayTime = 30 * time.Second

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	firefox, err := browser()
	if err != nil {
		log.Fatal(err)
	}

	linkChan, err := links(ctx)
	if err != nil {
		log.Fatalf(": %v", err)
	}

	visited := map[string]bool{}

	for link := range linkChan {
		if visited[link] {
			log.Printf("%s: already seen, get new material!", link)
			continue
		}
		visited[link] = true

		_, err = firefox.Goto(link)
		if err != nil {
			log.Printf("Visiting %q: %v", link, err)
		}

		time.Sleep(minDisplayTime)
		log.Println("Ready to load next link")
	}
}
