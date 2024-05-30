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
	s.Core.HandleFunc("/profil", handlers.ProfilHandler)
	s.Core.HandleFunc("/settings", handlers.SettingsHandler)
	s.Core.HandleFunc("/settings-password", handlers.SettingsPasswordChange)
	s.Core.HandleFunc("/settings-pseudo", handlers.SettingsPseudoChange)
	s.Core.HandleFunc("/settings-deleteaccount", handlers.SettingsDeleteAccount)
	s.Core.HandleFunc("/add-post", handlers.AddPostHandler)
	s.Core.HandleFunc("/viewpost", handlers.ViewPostHandler)
	s.Core.HandleFunc("/likepost", handlers.LikePostHandler)
	s.Core.HandleFunc("/dislikepost", handlers.DislikePostHandler)
	s.Core.HandleFunc("/likecomment", handlers.LikeCommentHandler)
	s.Core.HandleFunc("/dislikecomment", handlers.DislikeCommentHandler)
	s.Core.HandleFunc("/panel", handlers.PanelAdminHandler)
	s.Core.HandleFunc("/panel-ban", handlers.PanelBanHandler)
	s.Core.HandleFunc("/panel-unban", handlers.PanelUnbanHandler)
	s.Core.HandleFunc("/panel-grade", handlers.PanelGradeHandler)
	s.Core.HandleFunc("/panel-addcat", handlers.PanelAddCatHandler)
	s.Core.HandleFunc("/panel-action-ban", handlers.PanelActionBanHandler)
	s.Core.HandleFunc("/panel-action-unban", handlers.PanelActionUnbanHandler)
	s.Core.HandleFunc("/panel-changegrade", handlers.PanelActionChangeGradeHandler)
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
