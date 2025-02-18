// backend/internal/handlers/static.go
package handlers

import (
	"net/http"
	"path/filepath"
)

// ServeStaticFiles serves static files from the specified directory
func ServeStaticFiles(w http.ResponseWriter, r *http.Request) {
	// Set the directory where the static files are located
	staticDir := "./frontend/build"

	// If the request is for the root path, serve index.html
	if r.URL.Path == "/" {
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
		return
	}

	// Use http.StripPrefix to serve files from the static directory
	http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir+"/static/"))).ServeHTTP(w, r)
}
