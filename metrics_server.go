package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var metricsErrorMessage = []byte("Cannot get metrics")
var commandRequests int
var commandErrors int
var cacheHits int
var cacheMisses int

func ServeMetrics(w http.ResponseWriter, r *http.Request) {
	var guildCount int64
	db := db.Model(&GuildSettings{}).Count(&guildCount)
	if db.Error != nil {
		log.Println("Cannot get guild count metric", db.Error)
		w.Write(metricsErrorMessage)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var leaderboardEntries int64
	db = db.Model(&LeaderboardEntry{}).Count(&leaderboardEntries)
	if db.Error != nil {
		log.Println("Cannot get leaderboard entries metric")
		w.Write(metricsErrorMessage)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	message := []byte(fmt.Sprintf(`aoc_guild_count %d
aoc_leaderboard_entries %d
aoc_command_requests %d
aoc_command_errors %d
aoc_cache_hits
aoc_cache_misses`,
		guildCount,
		leaderboardEntries,
		commandRequests,
		commandErrors,
		cacheHits,
		cacheMisses))
	w.Write(message)
}

func StartMetricsServer() {
	log.Println("Start metrics server")
	http.HandleFunc("/metrics", ServeMetrics)
	err := http.ListenAndServe(os.Getenv("METRICS_SERVER"), nil)

	log.Println("Cannot start the metrics server", err)
}
