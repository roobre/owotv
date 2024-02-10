package main

import (
	"log"
	"slices"
	"time"

	"github.com/playwright-community/playwright-go"
)

type Gotoer interface {
	Goto(url string) error
}

type playwrightGotoer struct {
	playwright.Page
}

func (pg playwrightGotoer) Goto(url string) error {
	_, err := pg.Page.Goto(url)
	return err
}

type bot struct {
	minDisplayTime time.Duration
	admins         []string
}

func (b bot) run(page Gotoer, linkChan chan linkMeta) {
	hist := history{
		admins: b.admins,
	}

	t := throttler{
		minDisplayTime: b.minDisplayTime,
		admins:         b.admins,
	}

	for linkMeta := range linkChan {
		t.throttle(linkMeta.from)

		link := linkMeta.link
		if !hist.isNew(linkMeta) {
			log.Printf("Skipping already seen link or previous poster")
			continue
		}

		err := page.Goto(link)
		if err != nil {
			log.Printf("Visiting %q: %v", link, err)
		}
	}
}

type throttler struct {
	admins         []string
	minDisplayTime time.Duration

	lastDisplayed time.Time
}

func (t *throttler) throttle(user string) {
	last := t.lastDisplayed
	t.lastDisplayed = time.Now()

	if slices.Contains(t.admins, user) {
		log.Printf("Skipping wait time for %s", user)
		return
	}

	time.Sleep(time.Until(last.Add(t.minDisplayTime)))
}

type history struct {
	admins []string

	links      map[string]bool
	lastPoster string
}

func (h *history) isNew(linkMeta linkMeta) bool {
	if slices.Contains(h.admins, linkMeta.from) {
		return true
	}

	seen := h.links[linkMeta.link]
	samePoster := h.lastPoster == linkMeta.from

	if seen || samePoster {
		return false
	}

	if h.links == nil {
		h.links = map[string]bool{}
	}

	h.links[linkMeta.link] = true
	h.lastPoster = linkMeta.from

	return true
}
