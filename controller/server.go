package controller

import (
	"log"
	"net/http"
	"os"
	"strings"

	"RWTAPI/handlers/archive"
	"RWTAPI/handlers/clips"
	"RWTAPI/handlers/events"
	"RWTAPI/handlers/images"
	"RWTAPI/handlers/messages"
	"RWTAPI/handlers/music"
	"RWTAPI/handlers/site"
	"RWTAPI/handlers/theme"
	"RWTAPI/handlers/users"

	"github.com/gorilla/mux"
)

func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		// Enhanced CORS headers
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Set content type and other headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Transfer-Encoding", "identity") // Disable chunked encoding
		w.Header().Set("Connection", "close")           // Ensure connection closes properly
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(nw, r)

		// Ensure proper response completion
		if nw.status == http.StatusOK && nw.written == 0 {
			w.Write([]byte("{}"))
		}
	})
}

// Add a custom response writer to track status and body size
type responseWriter struct {
	http.ResponseWriter
	status  int
	written int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

// This function provides all of the endpoints listed in events.go
func InitHandlers() {
	router := mux.NewRouter()

	// Archive routes
	router.HandleFunc("/archiveDELETE/{id}", archive.ArchiveEntryDELETE).Methods("DELETE")
	router.HandleFunc("/ArchiveEntryPOST", archive.ArchiveEntryPOST).Methods("POST")
	router.HandleFunc("/ArchiveEntryPUT", archive.ArchiveEntryPUT).Methods("PUT")
	router.HandleFunc("/ArchiveGET/{id}", archive.ArchiveGET).Methods("GET")
	router.HandleFunc("/ArchivesGET", archive.ArchivesGET).Queries("screen", "{screen}", "archives", "{archives}").Methods("GET") // Clip routes
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
	router.HandleFunc("/RandomImagesGET", images.RandomImagesGET).Queries("screen", "{screen}", "images", "{images}").Methods("GET")
	router.HandleFunc("/ThemeDetailsGET", theme.ThemeDetailsGET).Methods("GET")
	router.HandleFunc("/ThemeDetailsPUT", theme.ThemeDetailsPUT).Methods("PUT")
	router.HandleFunc("/ThemeDetailsRandom", theme.ThemeDetailsRandom).Methods("GET")

	// User routes
	router.HandleFunc("/login", users.Login).Methods("POST")

	// Additional routes
	router.HandleFunc("/upcomingPlaylistsGET", music.UpcomingPlaylistsGET).Methods("GET")
	router.HandleFunc("/UpcomingEventsListsGET", events.UpcomingEventsListsGET).Methods("GET")
	router.HandleFunc("/upload", images.FileDetailsPOST).Methods("POST")
	router.HandleFunc("/uploadFile", images.UploadFile).Methods("POST")
	router.HandleFunc("/instagram-oembed", events.InstagramEmbed).Methods("GET")

	// Serve static files from the public directory
	router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("/app/images"))))
	router.PathPrefix("/music/").Handler(http.StripPrefix("/music/", http.FileServer(http.Dir("/app/music"))))
	// Add font files handler
	router.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts/", http.FileServer(http.Dir("/app/fonts"))))

	// Initialize CORS configuration
	corsWhitelist := os.Getenv("CORS_ORIGIN_WHITELIST")
	if corsWhitelist == "" {
		log.Fatal("CORS_ORIGIN_WHITELIST is not set")
	}
	whitelist := strings.Split(corsWhitelist, ",")

	handler := contentTypeMiddleware(router)

	log.Println("File Server is running on port 8080")
	log.Println("CORS Origin Whitelist: ", whitelist)
	http.ListenAndServe(":8080", handler)
}
