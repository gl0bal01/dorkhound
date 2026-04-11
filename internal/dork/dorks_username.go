package dork

import (
	"fmt"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateUsernameDorks)
}

func generateUsernameDorks(c *caseinfo.Case) []Dork {
	if len(c.Usernames) == 0 {
		return nil
	}

	highPrioritySites := map[string]bool{
		"github.com":                true,
		"gitlab.com":                true,
		"reddit.com":                true,
		"twitter.com OR site:x.com": true,
		"instagram.com":             true,
		"tiktok.com":                true,
		"youtube.com":               true,
		"steamcommunity.com":        true,
		"keybase.io":                true,
	}

	sites := []struct {
		site  string
		label string
	}{
		{"github.com", "GitHub"},
		{"gitlab.com", "GitLab"},
		{"reddit.com OR site:old.reddit.com", "Reddit"},
		{"twitter.com OR site:x.com", "Twitter/X"},
		{"instagram.com", "Instagram"},
		{"tiktok.com", "TikTok"},
		{"youtube.com", "YouTube"},
		{"twitch.tv", "Twitch"},
		{"steamcommunity.com", "Steam"},
		{"soundcloud.com", "SoundCloud"},
		{"medium.com", "Medium"},
		{"dev.to", "Dev.to"},
		{"hackerone.com", "HackerOne"},
		{"bugcrowd.com", "Bugcrowd"},
		{"keybase.io", "Keybase"},
		{"about.me", "About.me"},
		{"linktr.ee", "Linktree"},
		{"pinterest.com", "Pinterest"},
		{"flickr.com", "Flickr"},
		{"vimeo.com", "Vimeo"},
		{"behance.net", "Behance"},
		{"dribbble.com", "Dribbble"},
		{"bandcamp.com", "Bandcamp"},
		{"last.fm", "Last.fm"},
		{"strava.com", "Strava"},
		{"stackoverflow.com", "Stack Overflow"},
		{"producthunt.com", "Product Hunt"},
		{"anchor.fm", "Anchor"},
		{"mastodon.social", "Mastodon"},
		{"bsky.app", "Bluesky"},
		{"threads.net", "Threads"},
		{"discord.com", "Discord"},
	}

	var dorks []Dork
	for _, username := range c.Usernames {
		for _, s := range sites {
			priority := 2
			if highPrioritySites[s.site] {
				priority = 3
			}
			dorks = append(dorks, Dork{
				Query:    fmt.Sprintf(`"%s" site:%s`, username, s.site),
				Category: "username",
				Region:   "global",
				Priority: priority,
				Label:    fmt.Sprintf("%s username search", s.label),
			})
		}

		dorks = append(dorks,
			Dork{
				Query:    fmt.Sprintf(`"%s"`, username),
				Category: "username",
				Region:   "global",
				Priority: 1,
				Label:    "Bare username mention",
			},
			Dork{
				Query:    fmt.Sprintf(`"%s" gamertag OR "psn id" OR "steam id" OR "xbox live"`, username),
				Category: "username",
				Region:   "global",
				Priority: 2,
				Label:    "Gaming identity search",
			},
			Dork{
				Query:    fmt.Sprintf(`"%s" API_KEY OR password OR token site:pastebin.com`, username),
				Category: "username",
				Region:   "global",
				Priority: 2,
				Label:    "Credential leak check",
			},
		)
	}
	return dorks
}
