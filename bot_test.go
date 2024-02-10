package main

import (
	"testing"
	"time"
)

type fakeGotoer struct {
	url string
}

func (f *fakeGotoer) Goto(url string) error {
	f.url = url
	return nil
}

func TestBot(t *testing.T) {
	t.Parallel()

	t.Run("processes first link immediately", func(t *testing.T) {
		t.Parallel()

		ch, fake := setup(time.Second)
		ch <- linkMeta{link: "foo.bar", from: "test"}

		time.Sleep(100 * time.Millisecond)

		if fake.url != "foo.bar" {
			t.Fatalf("bot did not visit the link")
		}
	})

	t.Run("respects grace time", func(t *testing.T) {
		t.Parallel()

		ch, fake := setup(time.Second)
		ch <- linkMeta{link: "foo.bar", from: "test"}
		ch <- linkMeta{link: "another", from: "another"}

		time.Sleep(500 * time.Millisecond)

		if fake.url != "foo.bar" {
			t.Fatalf("bot loaded %s too fast", fake.url)
		}
	})

	t.Run("grace time is counted since link is sent", func(t *testing.T) {
		t.Parallel()

		ch, fake := setup(time.Second)
		ch <- linkMeta{link: "foo.bar", from: "test"}
		time.Sleep(1200 * time.Millisecond)
		ch <- linkMeta{link: "another", from: "another"}
		time.Sleep(100 * time.Millisecond)

		if fake.url != "another" {
			t.Fatalf("waited too long since the first link was sent")
		}
	})

	t.Run("ignores grace time for admins", func(t *testing.T) {
		t.Parallel()

		ch, fake := setup(time.Second)
		ch <- linkMeta{link: "foo.bar", from: "test"}
		ch <- linkMeta{link: "another", from: "admin"}

		time.Sleep(100 * time.Millisecond)

		if fake.url != "another" {
			t.Fatalf("bot did not skip grace time")
		}
	})

	t.Run("ignores dupe link", func(t *testing.T) {
		t.Parallel()

		ch, fake := setup(0)
		ch <- linkMeta{link: "foo.bar", from: "test"}
		ch <- linkMeta{link: "another", from: "another"}
		ch <- linkMeta{link: "foo.bar", from: "yet another"}

		time.Sleep(time.Second)

		if fake.url != "another" {
			t.Fatalf("bot loaded dupe link %s", fake.url)
		}
	})

	t.Run("loads dupe links from admins", func(t *testing.T) {
		t.Parallel()

		ch, fake := setup(0)
		ch <- linkMeta{link: "foo.bar", from: "test"}
		ch <- linkMeta{link: "another", from: "another"}
		ch <- linkMeta{link: "foo.bar", from: "admin"}

		time.Sleep(100 * time.Millisecond)

		if fake.url != "foo.bar" {
			t.Fatalf("bot did not load dupe link %s", fake.url)
		}
	})

	t.Run("ignores same poster", func(t *testing.T) {
		t.Parallel()

		ch, fake := setup(0)
		ch <- linkMeta{link: "foo.bar", from: "test"}
		ch <- linkMeta{link: "another", from: "test"}

		time.Sleep(100 * time.Millisecond)

		if fake.url != "foo.bar" {
			t.Fatalf("bot loaded dupe %s", fake.url)
		}
	})

	t.Run("dupe poster does not add link to exclusion", func(t *testing.T) {
		t.Parallel()

		ch, fake := setup(0)
		ch <- linkMeta{link: "foo.bar", from: "test"}
		ch <- linkMeta{link: "another", from: "test"}

		time.Sleep(100 * time.Millisecond)

		if fake.url != "foo.bar" {
			t.Fatalf("bot loaded dupe %s", fake.url)
		}

		ch <- linkMeta{link: "another", from: "test2"}

		time.Sleep(100 * time.Millisecond)

		if fake.url != "another" {
			t.Fatalf("link ignored due to same poster got added to view list")
		}
	})
}

func setup(graceTime time.Duration) (chan linkMeta, *fakeGotoer) {
	linkChan := make(chan linkMeta, 100)
	fake := &fakeGotoer{}

	b := bot{
		minDisplayTime: graceTime,
		admins:         []string{"admin"},
	}
	go b.run(fake, linkChan)
	return linkChan, fake
}
