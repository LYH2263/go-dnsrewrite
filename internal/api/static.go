package api

import (
	"net/http"
	"os"
	"path/filepath"
)

func staticHandler(webDir string) http.Handler {
	if webDir == "" {
		webDir = "web"
	}
	fs := http.FileServer(http.Dir(webDir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := filepath.Join(webDir, filepath.Clean(r.URL.Path))
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
			return
		}
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			fs.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
	})
}
