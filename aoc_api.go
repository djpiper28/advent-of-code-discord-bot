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

func getMostRecentEntriesNoTimeLimit(gs GuildSettings) ([]LeaderboardEntry, error) {
	db := db.Model(&LeaderboardEntry{})

	var ret []LeaderboardEntry
	db = db.Raw(`SELECT DISTINCT ON (board_code, id) name, stars, score, time, pk, id, board_code 
    FROM leaderboard_entries
    WHERE board_code = ?
    ORDER BY board_code, id, time DESC;`, gs.BoardCode).Scan(&ret)
	if db.Error != nil {
		return nil, db.Error
	}

	return ret, nil
}

func getMostRecentEntries(gs GuildSettings) ([]LeaderboardEntry, error) {
	db := db.Model(&LeaderboardEntry{})
	fifteenMinsAgo := time.Now().Add(-15 * time.Minute)

	var ret []LeaderboardEntry
	db = db.Raw(`SELECT DISTINCT ON (board_code, id) name, stars, score, time, pk, id, board_code 
    FROM leaderboard_entries
    WHERE board_code = ? AND time >= ? 
    ORDER BY board_code, id, time DESC;`, gs.BoardCode, fifteenMinsAgo).Scan(&ret)
	if db.Error != nil {
		return nil, db.Error
	}

	return ret, nil
}

func updateLeaderBoard(gs GuildSettings) ([]LeaderboardEntry, error) {
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
	retmap := make(map[int]LeaderboardEntry)
	for _, val := range rawret.Members {
		retmap[val.ID] = LeaderboardEntry{
			Time:      time.Now(),
			Stars:     val.Stars,
			Score:     val.Score,
			ID:        val.ID,
			Name:      val.Name,
			Event:     rawret.Event,
			PK:        uuid.New().String(),
			BoardCode: gs.BoardCode,
		}
		ret = append(ret, retmap[val.ID])
	}

	// Get old entries, such that only changes are created. Update timestamps otherwise
	entriesraw, err := getMostRecentEntriesNoTimeLimit(gs)
	if err == nil {
		entries := make(map[int]LeaderboardEntry)
		for _, entry := range entriesraw {
			entries[entry.ID] = entry
		}

		for id, entry := range entries {
			retentry, cont := retmap[id]

			// Guard against removal
			if !cont {
				log.Print("Entry removed from leaderboard")
			} else {
				if entry.Score == retentry.Score && entry.Stars == retentry.Stars {
					err = db.Exec(`UPDATE leaderboard_entries
    SET time = ?, name = ?
    WHERE pk = ?;`,
						time.Now(),
						retentry.Name,
						entry.PK).Error
					log.Printf("Updated time for %s", entry.Name)
					if err != nil {
						log.Print("Cannot update cache with compression ", err)
						break
					}
				} else {
					err = db.Exec(`UPDATE leaderboard_entries
    SET time = ?, stars = ?, score = ?, name = ?
    WHERE pk = ?;`,
						time.Now(),
						retentry.Stars,
						retentry.Score,
						retentry.Name,
						entry.PK).Error
					log.Printf("Updated time for %s", entry.Name)
					if err != nil {
						log.Print("Cannot update cache with compression ", err)
						break
					}
				}

				delete(retmap, id)
			}
		}
	} else {
		log.Print(err)
	}

	newentries := make([]LeaderboardEntry, 0)
	for _, entry := range retmap {
		entry.PK = uuid.New().String()
		newentries = append(newentries, entry)

		entry.PK = uuid.New().String()
		newentries = append(newentries, entry)
	}

	if len(newentries) != 0 {
		log.Printf("Found %d new entries", len(newentries))
		err = db.Create(newentries).Error
	}

	return ret, err
}

// / This will update a leaderboard and, return the (inserted) leadboard entry
func GetLeaderboard(gs GuildSettings) ([]LeaderboardEntry, error) {
	// Get cache
	cache, err := getMostRecentEntries(gs)
	if err != nil {
		return []LeaderboardEntry{}, err
	}

	// Cache was found
	if len(cache) != 0 {
		return cache, nil
	}

	// No cache was found
	log.Print("No cache was found")
	return updateLeaderBoard(gs)
}

func GetProfile(name string, gs GuildSettings) ([]LeaderboardEntry, error) {
	db := db.Model(&LeaderboardEntry{})

	var ret []LeaderboardEntry
	db = db.Where("board_code = ?", gs.BoardCode).Where("name = ?", name).Find(&ret)
	return ret, db.Error
}

func UpdateThread() {
	// Fetch all unique boards, then update them
	for true {
		if time.Now().Month() == time.December {
			func() {
				defer func() {
					err := recover()
					if err != nil {
						log.Print(err)
					}
				}()

				// Get all guilds
				var guilds []GuildSettings
				db := db.Model(&GuildSettings{})

				db = db.Find(&guilds)
				if db.Error != nil {
					log.Print(db.Error)
				}

				year := fmt.Sprintf("%d", time.Now().Year())

				// Try and update each board. This uses the board settings for each board until one works.
				guildsuniq := make(map[string]GuildSettings)
				for _, gs := range guilds {
					if gs.Year == year {
						// If the board has been successfully queried then do not query again
						_, cont := guildsuniq[gs.BoardCode]
						if !cont {
							ent, err := getMostRecentEntries(gs)

							if err != nil {
								log.Print(err)
							} else if len(ent) == 0 {
								log.Print("Polling indicates update is needed for: ", gs.BoardCode)
								_, err := updateLeaderBoard(gs)

								if err != nil {
									log.Print(err)
								} else {
									// Tag as queried
									guildsuniq[gs.BoardCode] = gs
								}
							}
						}
					}
				}
			}()
		}

		time.Sleep(time.Second)
	}
}
