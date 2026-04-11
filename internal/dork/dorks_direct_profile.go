package dork

import (
	"fmt"
	"net/url"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateDirectProfileDorks)
}

func generateDirectProfileDorks(c *caseinfo.Case) []Dork {
	if len(c.Usernames) == 0 {
		return nil
	}

	type platform struct {
		name string
		tmpl string // %s = url.PathEscape(username)
	}

	platforms := []platform{
		{"Telegram", "https://t.me/%s"},
		{"Keybase", "https://keybase.io/%s"},
		{"Twitter", "https://twitter.com/%s"},
		{"X", "https://x.com/%s"},
		{"Instagram", "https://www.instagram.com/%s/"},
		{"TikTok", "https://www.tiktok.com/@%s"},
		{"Reddit", "https://www.reddit.com/user/%s"},
		{"YouTube", "https://www.youtube.com/@%s"},
		{"Medium", "https://medium.com/@%s"},
		{"Dev.to", "https://dev.to/%s"},
		{"Dribbble", "https://dribbble.com/%s"},
		{"SoundCloud", "https://soundcloud.com/%s"},
		{"Last.fm", "https://www.last.fm/user/%s"},
		{"Steam", "https://steamcommunity.com/id/%s"},
		{"Twitch", "https://www.twitch.tv/%s"},
		{"Flickr", "https://www.flickr.com/people/%s"},
		{"About.me", "https://about.me/%s"},
		{"Linktree", "https://linktr.ee/%s"},
		{"Bluesky", "https://bsky.app/profile/%s.bsky.social"},
		{"Mastodon", "https://mastodon.social/@%s"},
	}

	var dorks []Dork
	for _, username := range c.Usernames {
		esc := url.PathEscape(username)
		for _, p := range platforms {
			dorks = append(dorks, Dork{
				Query:    fmt.Sprintf(p.tmpl, esc),
				Category: "direct-profile",
				Region:   "global",
				Priority: 3,
				Label:    fmt.Sprintf("Direct profile: %s/%s", p.name, username),
			})
		}
	}
	return dorks
}
