package controller

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mike5tew/RWTProj/RWTapi/handlers/archive"
	"github.com/mike5tew/RWTProj/RWTapi/handlers/clips"
	"github.com/mike5tew/RWTProj/RWTapi/handlers/events"
	"github.com/mike5tew/RWTProj/RWTapi/handlers/images"
	"github.com/mike5tew/RWTProj/RWTapi/handlers/messages"
	"github.com/mike5tew/RWTProj/RWTapi/handlers/music"
	"github.com/mike5tew/RWTProj/RWTapi/handlers/playlists"
	"github.com/mike5tew/RWTProj/RWTapi/handlers/site"
	"github.com/mike5tew/RWTProj/RWTapi/handlers/theme"
	"github.com/mike5tew/RWTProj/RWTapi/handlers/users"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// This function provides all of the endpoints listed in events.go
func InitHandlers() {
	router := mux.NewRouter()

	// Archive routes
	router.HandleFunc("/archiveDELETE/{id}", archive.ArchiveEntryDELETE).Methods("DELETE")
	router.HandleFunc("/ArchiveEntryPOST", archive.ArchiveEntryPOST).Methods("POST")
	router.HandleFunc("/ArchiveEntryPUT", archive.ArchiveEntryPUT).Methods("PUT")
	router.HandleFunc("/ArchiveGET/{id}", archive.ArchiveGET).Methods("GET")
	router.HandleFunc("/ArchivesGET", archive.ArchivesGET).Methods("GET")

	// Clip routes
	router.HandleFunc("/ClipDELETE/{id}", clips.ClipDelete).Methods("DELETE")
	router.HandleFunc("/ClipsGET/{id}", clips.EventClips).Methods("GET")
	router.HandleFunc("/ClipPOST", clips.ClipPOST).Methods("POST")

	// Event routes
	router.HandleFunc("/EventArchiveGET/{id}", archive.EventArchive).Methods("GET")
	router.HandleFunc("/EventDELETE/{id}", events.EventDELETE).Methods("DELETE")
	router.HandleFunc("/EventGET/{id}", events.EventGET).Methods("GET")
	router.HandleFunc("/EventImagesGET", images.EventImages).Methods("GET")
	router.HandleFunc("/EventPOST", events.EventPOST).Methods("POST")
	router.HandleFunc("/EventPUT", events.EventPUT).Methods("PUT")
	router.HandleFunc("/EventsListGET", events.EventsList).Methods("GET")
	router.HandleFunc("/EventsUpcomingGET", events.EventsUpcomingGET).Methods("GET")

	// Image routes
	router.HandleFunc("/ImageBackGET", images.ImageBackGET).Methods("GET")
	router.HandleFunc("/ImageDELETE/{id}", images.ImageDELETE).Methods("DELETE")
	router.HandleFunc("/ImageFilePOST", images.ImageFilePOST).Methods("POST")
	router.HandleFunc("/ImagePOST", images.ImagesPOST).Methods("POST")
	router.HandleFunc("/ImagePUT", images.ImagePUT).Methods("PUT")

	// Message routes
	router.HandleFunc("/messageDELETE/{id}", messages.MessageDELETE).Methods("DELETE")
	router.HandleFunc("/messagePOST", messages.MessagePOST).Methods("POST")
	router.HandleFunc("/messagePUT", messages.MessagePUT).Methods("PUT")
	router.HandleFunc("/messagesGET", messages.MessagesGET).Methods("GET")

	// Music routes
	router.HandleFunc("/musicListGET", music.MusicListGET).Methods("GET")
	router.HandleFunc("/musicTrackDELETE", music.MusicTrackDELETE).Methods("DELETE")
	router.HandleFunc("/musicTrackPOST", music.MusicTrackPOST).Methods("POST")
	router.HandleFunc("/musicTrackPUT", music.MusicTrackPUT).Methods("PUT")
	router.HandleFunc("/playlistDELETE/{id}", music.PlaylistDELETE).Methods("DELETE")
	router.HandleFunc("/PlaylistGET/{id}", music.PlaylistGET).Methods("GET")
	router.HandleFunc("/PlaylistPOST", music.PlaylistPOST).Methods("POST")

	// Site routes
	router.HandleFunc("/SiteInfoGET", site.SiteInfoGET).Methods("GET")
	router.HandleFunc("/SiteInfoPUT", site.SiteInfoPUT).Methods("PUT")

	// Theme routes
	router.HandleFunc("/RandomImagesGET", theme.RandomImagesGET).Methods("GET")
	router.HandleFunc("/ThemeDetailsGET", theme.ThemeDetailsGET).Methods("GET")
	router.HandleFunc("/ThemeDetailsPUT", theme.ThemeDetailsPUT).Methods("PUT")
	router.HandleFunc("/ThemeDetailsRandom", theme.ThemeDetailsRandom).Methods("GET")

	// User routes
	router.HandleFunc("/login", users.Login).Methods("POST")

	// Additional routes
	router.HandleFunc("/upcomingPlaylistsGET", playlists.UpcomingPlaylistsGET).Methods("GET")
	router.HandleFunc("/UpcomingEventsListsGET", events.UpcomingEventsListsGET).Methods("GET")
	router.HandleFunc("/upload", images.FileDetailsPOST).Methods("POST")
	router.HandleFunc("/uploadFile", images.UploadFile).Methods("POST")
	router.HandleFunc("/instagram-oembed", events.InstagramEmbed).Methods("GET")

	// Serve static files from the public directory
	router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("/app/images"))))
	router.PathPrefix("/music/").Handler(http.StripPrefix("/music/", http.FileServer(http.Dir("/app/music"))))

	// Initialize CORS configuration
	corsWhitelist := os.Getenv("CORS_ORIGIN_WHITELIST")
	if corsWhitelist == "" {
		log.Fatal("CORS_ORIGIN_WHITELIST is not set")
	}
	whitelist := strings.Split(corsWhitelist, ",")
	c := cors.New(cors.Options{
		AllowedOrigins:   whitelist,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	})

	handler := c.Handler(contentTypeMiddleware(router))

	log.Println("File Server is running on port 8080")
	http.ListenAndServe(":8080", handler)
}
