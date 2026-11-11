package web_test

import (
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/OpenCommunityEdition/opence-master-server/internal/config"
	"github.com/OpenCommunityEdition/opence-master-server/internal/sign"
)

// tags is what a page's visible text is not: its markup.
var tags = regexp.MustCompile(`(?s)<[^>]*>`)

// TestGameName checks that MASTER_GAME_NAME names the game on every page a
// player or an operator reads, and that MASTER_PROJECT_URL is the source's
// link: a fork's pages say nothing of OpenCE.
func TestGameName(t *testing.T) {
	const source = "https://github.com/example/master-fork/tree/v1.2.3-fork.1"
	h := newLinkHarness(t, func(c *config.Config) {
		c.GameName, c.ProjectURL = "GOAT H1", source
	})
	c := browser()
	player, _ := sign.NewSeed()
	pages := map[string]string{}
	for _, path := range []string{"/", "/apply", "/privacy", "/link/manage"} {
		status, _, page := h.get(c, h.server.URL+path)
		if status != http.StatusOK {
			t.Fatalf("%s: %d", path, status)
		}
		pages[path] = page
	}
	ts := h.m.Now().Unix()
	_, _, pages["/link"] = h.get(c, h.server.URL+"/link?"+sign.SignPlayerQuery(player, sign.LinkLabel, h.m.KeyID, ts))
	for path, page := range pages {
		// (the whole page: its title and attributes too, not only its text)
		if strings.Contains(page, "OpenCE") {
			t.Errorf("%s names OpenCE: %s", path, tags.ReplaceAllString(page, " "))
		}
		if !strings.Contains(page, `href="`+source+`"`) {
			t.Errorf("%s has no link to the source", path)
		}
	}
	if !strings.Contains(pages["/link"], "Link Discord to the GOAT H1 game showing") ||
		!strings.Contains(pages["/link/manage"], "The GOAT H1 games you linked") ||
		!strings.Contains(pages["/apply"], "runs GOAT H1 dedicated servers") ||
		!strings.Contains(pages["/"], "A master server for GOAT H1") {
		t.Errorf("the game's name is missing: %v", pages)
	}
	if !strings.Contains(pages["/"], ">H1</span>") {
		t.Errorf("the header's badge is not the game's")
	}
}

// TestGameNameDefault checks that the defaults are upstream's.
func TestGameNameDefault(t *testing.T) {
	h := newLinkHarness(t, nil)
	_, _, page := h.get(browser(), h.server.URL+"/")
	if !strings.Contains(page, "A master server for OpenCE") || !strings.Contains(page, ">CE</span>") ||
		!strings.Contains(page, `href="`+config.DefaultProjectURL+`"`) {
		t.Fatal(page)
	}
}
