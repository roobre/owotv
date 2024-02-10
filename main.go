package main

import (
	"context"
	"log"
	"time"
)

func main() {
	const minDisplayTime = 45 * time.Second
	admins := []string{"roobre", "moni"}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	firefox, err := browser()
	if err != nil {
		log.Fatal(err)
	}

	linkChan, err := links(ctx)
	if err != nil {
		log.Fatalf("getting stream: %v", err)
	}

	b := bot{
		minDisplayTime: minDisplayTime,
		admins:         admins,
	}

	b.run(playwrightGotoer{firefox}, linkChan)
}
