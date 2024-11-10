package controller

import (
	"fmt"
	"net/http"
	"packages/events"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// This function provides all of the endpoints listed in events.go
func InitHandlers() {
	router := mux.NewRouter()
	router.HandleFunc("/archiveDELETE/{id}", events.ArchiveEntryDELETE).Methods("DELETE")
	router.HandleFunc("/ArchiveEntryPOST", events.ArchiveEntryPOST).Methods("POST")
	router.HandleFunc("/ArchiveEntryPUT", events.ArchiveEntryPUT).Methods("PUT")
	router.HandleFunc("/ArchiveGET/{id}", events.ArchiveGET).Methods("GET")
	router.HandleFunc("/ArchivesGET/{records}", events.ArchivesGET).Methods("GET")
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
	router.HandleFunc("/RandomImagesGET/{scr}", events.RandomImagesGET).Methods("GET")
	router.HandleFunc("/ThemeDetailsGET", events.ThemeDetailsGET).Methods("GET")
	router.HandleFunc("/ThemeDetailsPUT", events.ThemeDetailsPUT).Methods("PUT")
	router.HandleFunc("/ThemeDetailsRandom", events.ThemeDetailsRandom).Methods("GET")
	router.HandleFunc("/upcomingPlaylistsGET", events.UpcomingPlaylistsGET).Methods("GET")
	router.HandleFunc("/UpcomingEventsListsGET", events.UpcomingEventsListsGET).Methods("GET")
	router.HandleFunc("/upload", events.FileDetailsPOST).Methods("POST")
	router.HandleFunc("/uploadFile", events.UploadFile).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	handler := c.Handler(router)
	fmt.Println("Server is running on port 8086")
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.ListenAndServe(":8086", handler)
	// in this instance the url of the giz.jpg file in the public directory would be?
	// http://localhost:8086/dt8giz.jpg
	// this is showing as not found because... the file is not in the public directory
	// the file is in the public/images directory
	//dt8giz.jpg is in the public directory but still not visible because the file is not in the public directory.  The public directory should be in the same directory as the server.go file
}
