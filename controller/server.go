package controller

import (
	"fmt"
	"net/http"
	"packages/events"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// The functions are:
// // 1. upload
// // 2. uploadLogo
// // 3. uploadBackground
// // 4. archives
// // 5. messages
// // 6. login
// // 7. upcomingPlaylists
// // 8. themeDetails
// // 9. themeDetailsPUT
// // 10. themeDetailsRandom
// // 11. musicList
// // 12. uploadClips
// // 13. updateArchiveEntry
// // 14. insertArchiveEntry
// // 15. archivePOST
// // 16. archiveDELETE
// // 17. archivePUT
// // 18. playlistGET
// // 19. playlistPOST
// // 20. EventDets
// // 21. images
// // 22. eventComplete
// // 23. randomImages
// // 24. ArchiveFromEvent
// // 25. clipsFromEvent
// // 26. messagesDELETE
// // 27. loginAddUser
// // 28. loginDeleteUser

// This function provides all of the endpoints listed in events.go
func InitHandlers() {
	router := mux.NewRouter()
	router.HandleFunc("/upload", events.FileDetailsPOST).Methods("POST")
	router.HandleFunc("/uploadFile", events.UploadFile).Methods("POST")
	router.HandleFunc("/ArchivesGET/{records}", events.ArchivesGET).Methods("GET")
	router.HandleFunc("/archivePOST", events.ArchiveEntryPOST).Methods("POST")
	router.HandleFunc("/archiveDELETE/{id}", events.ArchiveEntryDELETE).Methods("DELETE")
	router.HandleFunc("/messagesGET", events.MessagesGET).Methods("GET")
	router.HandleFunc("/messagePOST", events.MessagePOST).Methods("GET")
	router.HandleFunc("/messageDELETE/{id}", events.MessageDELETE).Methods("DELETE")
	router.HandleFunc("/messagePUT", events.MessagePUT).Methods("PUT")
	router.HandleFunc("/upcomingPlaylistsGET", events.UpcomingPlaylistsGET).Methods("GET")
	router.HandleFunc("/EventArchiveGET/{id}", events.EventArchive).Methods("GET")
	router.HandleFunc("/EventImagesGET", events.EventImages).Methods("GET")
	router.HandleFunc("/ClipsGET/{id}", events.EventClips).Methods("GET")
	router.HandleFunc("/ClipDELETE/{id}", events.ClipDelete).Methods("DELETE")
	router.HandleFunc("/EventGET/{id}", events.EventGET).Methods("GET")
	router.HandleFunc("/EventsList", events.EventsList).Methods("GET")
	router.HandleFunc("/EventPOST", events.EventPOST).Methods("POST")
	router.HandleFunc("/EventPUT", events.EventPUT).Methods("PUT")
	router.HandleFunc("/EventDELETE/{id}", events.EventDELETE).Methods("DELETE")
	router.HandleFunc("/EventsUpcomingGET", events.EventsUpcomingGET).Methods("GET")
	router.HandleFunc("/RandomImagesGET/{numbreq}", events.RandomImagesGET).Methods("GET")
	router.HandleFunc("/ThemeDetailsGET", events.ThemeDetailsGET).Methods("GET")
	router.HandleFunc("/ThemeDetailsPUT", events.ThemeDetailsPUT).Methods("PUT")
	router.HandleFunc("/ThemeDetailsRandom", events.ThemeDetailsRandom).Methods("GET")
	router.HandleFunc("/MusicListGET", events.MusicListGET).Methods("GET")
	router.HandleFunc("/MusicTrackPOST", events.MusicTrackPOST).Methods("POST")
	router.HandleFunc("/MusicTrackDELETE", events.MusicTrackDELETE).Methods("DELETE")
	router.HandleFunc("/MusicTrackPUT", events.MusicTrackPUT).Methods("PUT")
	router.HandleFunc("/UploadClips", events.UploadClips).Methods("POST")
	router.HandleFunc("/ArchiveEntryPUT", events.ArchiveEntryPUT).Methods("PUT")
	router.HandleFunc("/ArchiveEntryPOST", events.ArchiveEntryPOST).Methods("POST")
	router.HandleFunc("/ArchiveEntryDELETE/{id}", events.ArchiveEntryDELETE).Methods("DELETE")
	router.HandleFunc("/PlaylistGET/{id}", events.PlaylistGET).Methods("GET")
	router.HandleFunc("/PlaylistPOST", events.PlaylistPOST).Methods("POST")
	router.HandleFunc("/EventGET/{id}", events.EventGET).Methods("GET")
	router.HandleFunc("/EventPOST", events.EventPOST).Methods("POST")
	router.HandleFunc("/EventPUT", events.EventPUT).Methods("PUT")
	router.HandleFunc("/EventDELETE/{id}", events.EventDELETE).Methods("DELETE")
	router.HandleFunc("/ImageGET/{id}", events.ImageGET).Methods("GET")
	router.HandleFunc("/ImagePOST", events.ImagesPOST).Methods("POST")
	router.HandleFunc("/ImageDELETE/{id}", events.ImageDELETE).Methods("DELETE")
	router.HandleFunc("/ImagePUT", events.ImagePUT).Methods("PUT")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	handler := c.Handler(router)
	http.Handle("/", router)
	fmt.Println("Server is running on port 8086")
	http.ListenAndServe(":8086", handler)
}
