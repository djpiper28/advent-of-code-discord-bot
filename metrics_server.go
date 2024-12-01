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
	model := db.Model(&GuildSettings{}).Select("count (*)").Count(&guildCount)
	if model.Error != nil {
		log.Println("Cannot get guild count metric", model.Error)
		w.Write(metricsErrorMessage)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var leaderboardEntries int64
	model = db.Model(&LeaderboardEntry{}).Select("count(distinct(pk))").Count(&leaderboardEntries)
	if model.Error != nil {
		log.Println("Cannot get leaderboard entries metric")
		w.Write(metricsErrorMessage)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	message := []byte(fmt.Sprintf(`aoc_guild_count %d
aoc_leaderboard_entries %d
aoc_command_requests %d
aoc_command_errors %d
aoc_cache_hits %d
aoc_cache_misses %d`,
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
