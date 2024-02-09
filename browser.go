package main

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
)

func browser() (playwright.Page, error) {
	err := playwright.Install(&playwright.RunOptions{
		Browsers: []string{"firefox"},
	})
	if err != nil {
		return nil, fmt.Errorf("installing firefox: %w", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("running playwright: %w", err)
	}

	firefox, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
		FirefoxUserPrefs: map[string]interface{}{
			"layout.css.devPixelsPerPx": "1.2",
			"media.autoplay.default":    0, // Allowed.
		},
	})
	if err != nil {
		return nil, fmt.Errorf("launching firefox: %w", err)
	}

	page, err := firefox.NewPage(playwright.BrowserNewPageOptions{
		AcceptDownloads: playwright.Bool(false),
		NoViewport:      playwright.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("creating page: %w", err)
	}

	_, err = page.Goto("https://addons.mozilla.org/en-US/firefox/addon/ublock-origin/")
	if err != nil {
		return nil, fmt.Errorf("opening first page: %w", err)
	}

	return page, err
}
