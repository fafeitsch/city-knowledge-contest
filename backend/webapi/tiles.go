package webapi

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var cacheMutex = &sync.Mutex{}
var tileCache = map[string]tile{}

func ClearTileCache() {
	timer := time.NewTicker(15 * time.Second)
	for {
		select {
		case <-timer.C:
			cacheMutex.Lock()
			now := time.Now()
			counter := 0
			for key, tile := range tileCache {
				if now.Sub(tile.timestamp) > 60*time.Minute {
					delete(tileCache, key)
					counter = counter + 1
				}
			}
			cacheMutex.Unlock()
			log.Printf("removed %d tiles from cache", counter)
		}
	}
}

type tile struct {
	timestamp   time.Time
	length      string
	contentType string
	content     []byte
}

func (r *RpcServer) serveTile(parts []string, resp http.ResponseWriter) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	z := parts[2]
	x := parts[3]
	y := parts[4]
	cached, ok := tileCache[fmt.Sprintf("%s/%s/%s", z, x, y)]
	if ok {
		r.readFromCache(resp, cached)
		return
	}
	url := strings.Replace(r.options.TileServer, "{z}", z, 1)
	url = strings.Replace(url, "{x}", x, 1)
	url = strings.Replace(url, "{y}", y, 1)
	client := http.Client{Timeout: 0}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "github.com/fafeitsch/city-knowledge-contest")
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("%v", err)
		http.Error(resp, "could not get tile from tile server", http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(response.Body)
	cachedTile := tile{
		timestamp:   time.Now(),
		length:      response.Header.Get("Content-Length"),
		contentType: response.Header.Get("Content-Type"),
		content:     body,
	}
	tileCache[fmt.Sprintf("%s/%s/%s", z, x, y)] = cachedTile
	r.readFromCache(resp, cachedTile)
}

func (r *RpcServer) readFromCache(resp http.ResponseWriter, cached tile) {
	resp.Header().Set("Content-Type", cached.contentType)
	resp.Header().Set("Content-Length", cached.length)
	_, _ = resp.Write(cached.content)
}
