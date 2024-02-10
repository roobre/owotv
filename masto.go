package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/RasmusLindroth/go-mastodon"
)

type linkMeta struct {
	link string
	from string
}

func links(ctx context.Context) (chan linkMeta, error) {
	linksChan := make(chan linkMeta, 1024)

	masto := mastodon.NewClient(&mastodon.Config{
		Server:       "https://owo.cafe",
		ClientID:     os.Getenv("OWOTV_APP_ID"),
		ClientSecret: os.Getenv("OWOTV_APP_SECRET"),
		AccessToken:  os.Getenv("OWOTV_USER_TOKEN"),
	})

	me, err := masto.GetAccountCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("checking my account: %w", err)
	}

	log.Printf("Logged in as %s", me.Username)

	pmCh, err := masto.StreamingDirect(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting pms stream: %w", err)
	}

	log.Println("Listening for PMs")

	// visited := map[string]bool{}

	go func() {
		defer close(linksChan)

		for event := range pmCh {
			pm, isPm := event.(*mastodon.ConversationEvent)
			if !isPm {
				log.Printf("Unknown event received")
				continue
			}

			status := pm.Conversation.LastStatus
			if status.Account.ID == me.ID {
				continue
			}

			link, err := extractLink(status.Content)
			if err != nil {
				log.Printf("extracting link: %v", err)
				continue
			}

			log.Printf("%s: %s", status.Account.URL, link)

			linksChan <- linkMeta{
				link: link,
				from: status.Account.Username,
			}
		}
	}()

	return linksChan, nil
}

func extractLink(status string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(status))
	if err != nil {
		return "", fmt.Errorf("parsing status: %w", err)
	}

	link, hasHref := doc.Find("a:not(.mention)").First().Attr("href")
	if !hasHref {
		return "", errors.New("no link found")
	}

	return link, nil
}
