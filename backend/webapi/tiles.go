package webapi

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var cacheMutex = &sync.Mutex{}

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
	resp.Header().Set("Content-Type", "image/png")
	cachedFile, err := os.Stat(filepath.Join("tiles", z, x, y) + ".png")
	if err == nil && time.Now().Sub(cachedFile.ModTime()) < 24*time.Hour*180 {
		tile, err := os.ReadFile(filepath.Join("tiles", z, x, y) + ".png")
		resp.Header().Set("Content-Length", strconv.Itoa(len(tile)))
		if err == nil {
			_, _ = resp.Write(tile)
		}
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
		http.Error(resp, "could not get tile from tile server", http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(response.Body)
	_ = os.MkdirAll(filepath.Join("tiles", z, x), os.ModePerm)
	if r.options.UseTileCache {
		go func() {
			_ = ioutil.WriteFile(filepath.Join("tiles", z, x, y)+".png", body, 0644)
		}()
	}
	resp.Header().Set("Content-Length", strconv.Itoa(len(body)))
	_, _ = resp.Write(body)
}

func (r *RpcServer) readFromCache(resp http.ResponseWriter, cached tile) {
	_, _ = resp.Write(cached.content)
}
