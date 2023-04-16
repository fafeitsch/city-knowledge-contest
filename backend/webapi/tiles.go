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

var tilesDirMutex sync.Mutex

type tile struct {
	timestamp   time.Time
	length      string
	contentType string
	content     []byte
}

func (r *RpcServer) serveTile(parts []string, resp http.ResponseWriter) {
	z := parts[3]
	x := parts[4]
	y := parts[5]
	resp.Header().Set("Content-Type", "image/png")
	if r.readCacheFile(z, x, y, resp) {
		return
	}
	url := strings.Replace(r.options.TileServer, "{z}", z, 1)
	url = strings.Replace(url, "{x}", x, 1)
	url = strings.Replace(url, "{y}", y, 1)
	client := http.Client{Timeout: time.Second * 10}
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
			tilesDirMutex.Lock()
			defer tilesDirMutex.Unlock()
			_ = ioutil.WriteFile(filepath.Join("tiles", z, x, y)+".png", body, 0644)
		}()
	}
	resp.Header().Set("Content-Length", strconv.Itoa(len(body)))
	resp.Header().Set("Cache-Control", "public, max-age=86400")
	_, _ = resp.Write(body)
}

func (r *RpcServer) readCacheFile(z string, x string, y string, resp http.ResponseWriter) bool {
	tilesDirMutex.Lock()
	defer tilesDirMutex.Unlock()
	cachedFile, err := os.Stat(filepath.Join("tiles", z, x, y) + ".png")
	if err != nil || time.Now().Sub(cachedFile.ModTime()) >= 24*time.Hour*180 {
		return false
	}
	tile, err := os.ReadFile(filepath.Join("tiles", z, x, y) + ".png")
	resp.Header().Set("Content-Length", strconv.Itoa(len(tile)))
	resp.Header().Set("Cache-Control", "public, max-age=86400")
	if err == nil {
		_, _ = resp.Write(tile)
	}
	return true
}
