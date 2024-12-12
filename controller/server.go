package controller

import (
	"RWTAPI/events"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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
	router.HandleFunc("/archiveDELETE/{id}", events.ArchiveEntryDELETE).Methods("DELETE")
	router.HandleFunc("/ArchiveEntryPOST", events.ArchiveEntryPOST).Methods("POST")
	router.HandleFunc("/ArchiveEntryPUT", events.ArchiveEntryPUT).Methods("PUT")
	router.HandleFunc("/ArchiveGET/{id}", events.ArchiveGET).Methods("GET")
	router.HandleFunc("/ArchivesGET", events.ArchivesGET).Methods("GET")
	router.HandleFunc("/ClipDELETE/{id}", events.ClipDelete).Methods("DELETE")
	router.HandleFunc("/ClipsGET/{id}", events.EventClips).Methods("GET")
	router.HandleFunc("/ClipPOST", events.ClipPOST).Methods("POST")
	router.HandleFunc("/EventArchiveGET/{id}", events.EventArchive).Methods("GET")
	router.HandleFunc("/EventDELETE/{id}", events.EventDELETE).Methods("DELETE")
	router.HandleFunc("/EventGET/{id}", events.EventGET).Methods("GET")
	router.HandleFunc("/EventImagesGET", events.EventImages).Methods("GET")
	router.HandleFunc("/EventPOST", events.EventPOST).Methods("POST")
	router.HandleFunc("/EventPUT", events.EventPUT).Methods("PUT")
	router.HandleFunc("/EventsListGET", events.EventsList).Methods("GET")
	router.HandleFunc("/EventsUpcomingGET", events.EventsUpcomingGET).Methods("GET")
	router.HandleFunc("/ImageDELETE/{id}", events.ImageDELETE).Methods("DELETE")
	router.HandleFunc("/ImageGET/{id}", events.ImageGET).Methods("GET")
	router.HandleFunc("/ImageFilePOST", events.ImageFilePOST).Methods("POST")
	router.HandleFunc("/ImagePOST", events.ImagesPOST).Methods("POST")
	router.HandleFunc("/ImagePUT", events.ImagePUT).Methods("PUT")
	router.HandleFunc("/login", events.Login).Methods("POST")
	router.HandleFunc("/messageDELETE/{id}", events.MessageDELETE).Methods("DELETE")
	router.HandleFunc("/messagePOST", events.MessagePOST).Methods("POST")
	router.HandleFunc("/messagePUT", events.MessagePUT).Methods("PUT")
	router.HandleFunc("/messagesGET", events.MessagesGET).Methods("GET")
	router.HandleFunc("/MusicListGET", events.MusicListGET).Methods("GET")
	router.HandleFunc("/MusicTrackDELETE", events.MusicTrackDELETE).Methods("DELETE")
	router.HandleFunc("/MusicTrackPOST", events.MusicTrackPOST).Methods("POST")
	router.HandleFunc("/MusicTrackPUT", events.MusicTrackPUT).Methods("PUT")
	router.HandleFunc("/PlaylistGET/{id}", events.PlaylistGET).Methods("GET")
	router.HandleFunc("/PlaylistPOST", events.PlaylistPOST).Methods("POST")
	router.HandleFunc("/SiteInfoGET", events.SiteInfoGET).Methods("GET")
	router.HandleFunc("/SiteInfoPUT", events.SiteInfoPUT).Methods("PUT")
	router.HandleFunc("/RandomImagesGET", events.RandomImagesGET).Methods("GET")
	router.HandleFunc("/ThemeDetailsGET", events.ThemeDetailsGET).Methods("GET")
	router.HandleFunc("/ThemeDetailsPUT", events.ThemeDetailsPUT).Methods("PUT")
	router.HandleFunc("/ThemeDetailsRandom", events.ThemeDetailsRandom).Methods("GET")
	router.HandleFunc("/upcomingPlaylistsGET", events.UpcomingPlaylistsGET).Methods("GET")
	router.HandleFunc("/UpcomingEventsListsGET", events.UpcomingEventsListsGET).Methods("GET")
	router.HandleFunc("/upload", events.FileDetailsPOST).Methods("POST")
	router.HandleFunc("/uploadFile", events.UploadFile).Methods("POST")
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
	fmt.Println("CORS Origin Whitelist:", whitelist)

	c := cors.New(cors.Options{
		AllowedOrigins:   whitelist,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	})

	handler := c.Handler(contentTypeMiddleware(router))

	// handler := c.Handler(router)
	fmt.Println("File Server is running on port 8080")
	http.ListenAndServe(":8080", handler)
}
