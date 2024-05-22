package core

/*

   Titouan

   - Server core
   - router
   - load assets
*/
import (
	"Forum/core/handlers"
	"fmt"
	"net/http"
)

// WebServer represents the core of the web server.
type WebServer struct {
	Core      *http.ServeMux
	Port      int
	AssetsDir string
}

// Router sets up the routes for the web server.
func (s WebServer) Router() {
	s.Core.HandleFunc("/login", handlers.LoginHandler)
	s.Core.HandleFunc("/register", handlers.RegisterHandler)
	s.Core.HandleFunc("/", handlers.ForumHandler)
}

// Launch starts the web server.
func (s WebServer) Launch() {
	http.ListenAndServe(fmt.Sprintf(":%d", s.Port), s.Core)
}

// LoadAssets serves static assets.
func (s WebServer) LoadAssets() {
	fs := http.FileServer(http.Dir(s.AssetsDir))
	s.Core.Handle("/"+s.AssetsDir+"/", http.StripPrefix("/"+s.AssetsDir+"/", fs))
	fmt.Printf("Server is running on port %d...\n", s.Port)
}
