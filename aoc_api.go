package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// / This cookie jar is from https://stackoverflow.com/questions/12756782/go-http-post-and-use-cookies
type Jar struct {
	lk      sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewJar() *Jar {
	jar := new(Jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

// SetCookies handles the receipt of the cookies in a reply for the
// given URL.  It may or may not choose to save the cookies, depending
// on the jar's policy and implementation.
func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.lk.Lock()
	jar.cookies[u.Host] = cookies
	jar.lk.Unlock()
}

// Cookies returns the cookies to send in a request for the given URL.
// It is up to the implementation to honor the standard cookie use
// restrictions such as in RFC 6265.
func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}

// / This will update a leaderboard and, return the (inserted) leadboard entry
func GetLeaderboard(gs GuildSettings) ([]LeaderboardEntry, error) {
	// TODO: Return cached value if it is less than 15 minutes old. Respecting the AOC api guides is nice!

	url_s := fmt.Sprintf("https://adventofcode.com/%s/leaderboard/private/view/%s.json",
		gs.Year,
		gs.BoardCode)
	url, err := url.Parse(url_s)
	if err != nil {
		return []LeaderboardEntry{}, err
	}

	log.Print("Fetching leaderboard ", url_s)

	cookie := http.Cookie{
		Name:     "session",
		Value:    gs.SessionKey,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	jar := NewJar()
	jar.SetCookies(url, []*http.Cookie{&cookie})

	client := http.Client{Jar: jar}
	resp, err := client.Get(url_s)

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return []LeaderboardEntry{}, err
	}

	var rawret ApiLeaderboard
	err = json.Unmarshal(bytes, &rawret)
	if err != nil {
		return []LeaderboardEntry{}, err
	}

	// Map api leaderboard to entry
	ret := make([]LeaderboardEntry, 0)
	for _, val := range rawret.Members {
		ret = append(ret, LeaderboardEntry{
			Time:      time.Now(),
			Stars:     val.Stars,
			Score:     val.Score,
			ID:        val.ID,
			Name:      val.Name,
			Event:     rawret.Event,
			PK:        uuid.New().String(),
			BoardCode: gs.BoardCode,
		})
	}

	// Insert the new data
	err = db.Create(ret).Error
	return ret, err
}

func UpdateThread() {

}
