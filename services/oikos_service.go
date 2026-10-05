package services

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type dependenciaOikosResponse struct {
	ID     int64  `json:"Id"`
	Nombre string `json:"Nombre"`
}

type dependenciaOikosCacheEntry struct {
	nombre string
	expira time.Time
}

var dependenciaOikosCache = struct {
	sync.RWMutex
	items map[int64]dependenciaOikosCacheEntry
}{items: map[int64]dependenciaOikosCacheEntry{}}

func nombreDependenciaOikos(ctx context.Context, id int64) string {
	if id <= 0 {
		return ""
	}
	dependenciaOikosCache.RLock()
	entry, ok := dependenciaOikosCache.items[id]
	dependenciaOikosCache.RUnlock()
	if ok && time.Now().Before(entry.expira) {
		return entry.nombre
	}

	baseURL := oikosAPIURL()
	if baseURL == "" {
		return ""
	}
	endpoint := strings.TrimRight(baseURL, "/") + "/dependencia/" + url.PathEscape(strconv.FormatInt(id, 10))
	var response dependenciaOikosResponse
	if err := getJSON(ctx, endpoint, &response); err != nil {
		return ""
	}
	nombre := strings.TrimSpace(response.Nombre)
	if nombre == "" {
		return ""
	}
	dependenciaOikosCache.Lock()
	dependenciaOikosCache.items[id] = dependenciaOikosCacheEntry{nombre: nombre, expira: time.Now().Add(12 * time.Hour)}
	dependenciaOikosCache.Unlock()
	return nombre
}

func nombreProgramaFallback(id int64) string {
	return fmt.Sprintf("Programa %d", id)
}
