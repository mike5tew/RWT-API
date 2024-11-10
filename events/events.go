package events

import (
	// import the mysql driver
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	strt "packages/structures"
	"strconv"
	"strings"
	"time"

	sqldb "packages/sqldb"

	"github.com/gorilla/mux"
	"golang.org/x/exp/rand"
)

func ParseMySQLDateTime(datetimeStr string) (time.Time, error) {
	if len(datetimeStr) == 10 { // Check if the string is a date without time
		datetimeStr += " 00:00:00" // Append default time
	}
	return time.Parse("2006-01-02 15:04:05", datetimeStr)
}

// Translate the following functions into endpoints
// 1. upload
func FileDetailsPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	sqldb.Init()
	//get the json data from the request
	var uploadDetails strt.ImageDetail
	err := json.NewDecoder(r.Body).Decode(&uploadDetails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sSQL := "INSERT INTO images (filename, eventID, caption, width, height) VALUES (?, ?, ?, ?, ?)"
	// insert the data into the database
	_, err = sqldb.DB.Exec(sSQL, uploadDetails.Filename, uploadDetails.EventID, uploadDetails.Caption, uploadDetails.Width, uploadDetails.Height)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// get the imageID of the last inserted image
	sSQL = "SELECT ID FROM images ORDER BY ID DESC LIMIT 1"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var imageID int
	for rows.Next() {
		err = rows.Scan(&imageID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// return the imageID
	json.NewEncoder(w).Encode(imageID)
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	// the function recieves a multipart form with a file and a string
	// formData.append("file", file);
	// formData.append("eventID", eventID.toString());
	// formData.append("caption", caption);
	// formData.append("width", width.toString());
	// formData.append("height", height.toString());
	// formData.append("filename", filename);
	// formData.append("uploadType", uploadType.toString());
	// the function saves the file to the server and the details to the database
	// Parse the multipart form, 10 << 20 specifies a maximum upload of 10MB files
	r.ParseMultipartForm(10 << 20)
	// get the file from the form
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	// Create a file
	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	// Save the file to the folder based on the filetype from the form

	sFolderName := "images"
	switch r.FormValue("uploadType") {
	case "Background":
		// Copy the file to the destination
		sFolderName = "background"
	case "LogoImage":
		sFolderName = "logos"
	case "MobileImage":
		sFolderName = "mobile"
	}
	// To address the file to the folder we need to create the folder if it does not exist
	// check if the folder exists
	if _, err := os.Stat(sFolderName); os.IsNotExist(err) {
		// create the folder
		err = os.Mkdir(sFolderName, 0755)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// create the file
	f, err = os.OpenFile(sFolderName+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully Uploaded File\n")
	// Copy the file to the destination
	io.Copy(f, file)
	// convert the eventID to an integer
	var evID int
	var Width int
	var Height int
	evID, err = strconv.Atoi(r.FormValue("eventID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Width, err = strconv.Atoi(r.FormValue("width"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Height, err = strconv.Atoi(r.FormValue("height"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get the json data from the request
	var uploadDetails strt.ImageDetail
	uploadDetails.Filename = handler.Filename
	uploadDetails.EventID = evID
	uploadDetails.Caption = r.FormValue("caption")
	uploadDetails.Width = Width
	uploadDetails.Height = Height
	// insert the data into the database
	sqldb.Init()
	sSQL := "INSERT INTO images (filename, eventID, caption, width, height) VALUES (?, ?, ?, ?, ?)"
	_, err = sqldb.DB.Exec(sSQL, uploadDetails.Filename, uploadDetails.EventID, uploadDetails.Caption, uploadDetails.Width, uploadDetails.Height)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// get the imageID of the last inserted image
	sSQL = "SELECT ID FROM images ORDER BY ID DESC LIMIT 1"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var imageID int
	for rows.Next() {
		err = rows.Scan(&imageID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// return the imageID
	json.NewEncoder(w).Encode(imageID)
}

// func UploadBackground(w http.ResponseWriter, r *http.Request) {
// 	// Parse the multipart form, 10 << 20 specifies a maximum upload of 10MB files
// 	r.ParseMultipartForm(10 << 20)
// 	// get the file from the form
// 	file, handler, err := r.FormFile("file")
// 	if err != nil {
// 		fmt.Println("Error Retrieving the File")
// 		fmt.Println(err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer file.Close()
// 	// Create a file
// 	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
// 	if err != nil {
// 		fmt.Println(err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer f.Close()
// 	// Copy the file to the destination
// 	io.Copy(f, file)
// 	fmt.Fprintf(w, "Successfully Uploaded File\n")
// }

func ArchiveGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	vars := mux.Vars(r)
	id := vars["id"]

	sqldb.Init()
	var archive strt.ArchiveEntry
	var event strt.EventDetails
	sSQL := "SELECT archive.archiveID, choirevents.location, choirevents.eventDate, choirevents.title, archive.report FROM choirevents JOIN archive ON archive.eventID=choirevents.eventID WHERE choirevents.eventID = ?"
	rows, err := sqldb.DB.Query(sSQL, id)
	if err != nil {
		log.Fatal(err)
		archive.ArchiveID = -1
		archive.Report = err.Error()
		json.NewEncoder(w).Encode(archive)
		return
	}
	SQLDate := ""
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&archive.ArchiveID, &event.Location, &SQLDate, &event.Title, &archive.Report)
		if err != nil {
			log.Fatal(err)
			archive.ArchiveID = -1
			archive.Report = err.Error()
			json.NewEncoder(w).Encode(archive)
			return
		}
		//fmt.Println("ArchiveGET:: " + event.Location + " " + event.EventDate + " " + event.Title + " " + archive.Report)
		eventID, err := strconv.Atoi(id)
		if err != nil {
			archive.ArchiveID = -1
			archive.Report = err.Error()
			json.NewEncoder(w).Encode(archive)
			return
		}
		// convert the date to a formatted string
		parsedDate, err := ParseMySQLDateTime(SQLDate)
		if err != nil {
			archive.ArchiveID = -1
			archive.Report = err.Error()
			json.NewEncoder(w).Encode(archive)
			return
		}
		event.EventDate = parsedDate.Format("2006-01-02")

		event.EventID = eventID
		archive.EventDetails = event
		//fmt.Println("ArchiveGET:: " + archive.EventDetails.Location + " " + archive.EventDetails.EventDate + " " + archive.EventDetails.Title + " " + archive.Report)
	}
	//retrive the ImageDetails for the event
	sSQL = "SELECT * FROM images WHERE eventID = ?"
	rows, err = sqldb.DB.Query(sSQL, archive.EventDetails.EventID)
	if err != nil {
		log.Fatal(err)
		archive.ArchiveID = -1
		archive.Report = err.Error()
		json.NewEncoder(w).Encode(archive)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var image strt.ImageDetail
		err = rows.Scan(&image.ImageID, &image.Filename, &image.Caption, &image.EventID, &image.Height, &image.Width)
		if err != nil {
			log.Fatal(err)
			archive.ArchiveID = -1
			archive.Report = err.Error()
			json.NewEncoder(w).Encode(archive)
			return
		}
		archive.Images = append(archive.Images, image)
	}
	//retrive the ClipDetails for the event
	sSQL = "SELECT * FROM clips WHERE eventID = ?"
	rows, err = sqldb.DB.Query(sSQL, archive.EventDetails.EventID)
	if err != nil {
		log.Fatal(err)
		archive.ArchiveID = -1
		archive.Report = err.Error()
		json.NewEncoder(w).Encode(archive)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var clip strt.Clip
		err = rows.Scan(&clip.ClipID, &clip.ClipURL, &clip.EventID, &clip.Caption)
		if err != nil {
			log.Fatal(err)
			archive.ArchiveID = -1
			archive.Report = err.Error()
			json.NewEncoder(w).Encode(archive)
			return
		}
		archive.Clips = append(archive.Clips, clip)
	}
	// return the data
	fmt.Println("ArchiveGET")
	json.NewEncoder(w).Encode(archive)
}

// GET a given number of archive records
func ArchivesGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// the url to rerach this endpoint is /ArchivesGET/{records} or http://localhost:8080/ArchivesGET/5
	// this syntax is not returned by the swagger documentation
	// get the number of records required
	sqldb.Init()
	// split the url into an array iusing the / as the delimiter
	// the last element of the array is the number of records required
	// the array is zero based so the last element is len(array) - 1
	urlArray := strings.Split(r.URL.Path, "/")
	records := urlArray[len(urlArray)-1]

	fmt.Println("ArchivesGET " + records)
	// the records are passed as a string so we need to convert them to an integer
	iRecords, err := strconv.Atoi(records)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create a slice of archive entries
	var archives []strt.ArchiveEntry
	// create a slice of eventIDs
	var events []int
	// create a slice of image details
	var images []strt.ImageDetail
	// create a slice of clips
	var clips []strt.Clip
	// create a slice of event details
	var eventDetails []strt.EventDetails
	// Gather the archive details from the database
	sSQLDate := ""
	sSQL := "SELECT archive.archiveID, choirevents.location, choirevents.eventDate, choirevents.title, archive.report, archive.eventID FROM choirevents JOIN archive ON archive.eventID=choirevents.eventID ORDER BY choirevents.eventDate LIMIT ?"
	rows, err := sqldb.DB.Query(sSQL, iRecords)
	fmt.Println("ArchivesGET:: " + sSQL + " " + records)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var archive strt.ArchiveEntry
		var eventDetail strt.EventDetails
		err = rows.Scan(&archive.ArchiveID, &eventDetail.Location, &sSQLDate, &eventDetail.Title, &archive.Report, &eventDetail.EventID)
		// convert and format the date
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		parsedDate, err := ParseMySQLDateTime(sSQLDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		archives = append(archives, archive)
		eventDetail.EventDate = parsedDate.Format("2006-01-02")
		events = append(events, archive.EventDetails.EventID)
		eventDetails = append(eventDetails, eventDetail)
	}
	// get the images from all of the events
	sSQL = "SELECT * FROM images WHERE eventID = ?"
	for _, event := range events {
		rows, err = sqldb.DB.Query(sSQL, event)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var image strt.ImageDetail
			err = rows.Scan(&image.ImageID, &image.Filename, &image.Caption, &image.EventID, &image.Height, &image.Width)
			if err != nil {
				log.Fatal(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			images = append(images, image)
		}
	}
	// get the clips from all of the events
	sSQL = "SELECT * FROM clips WHERE eventID = ?"
	for _, event := range events {
		rows, err = sqldb.DB.Query(sSQL, event)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var clip strt.Clip
			err = rows.Scan(&clip.ClipID, &clip.ClipURL, &clip.EventID, &clip.Caption)
			if err != nil {
				log.Fatal(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			clips = append(clips, clip)
		}
	}
	// combine the data into the archive entries
	for i, archive := range archives {
		for _, eventDetail := range eventDetails {
			if archive.EventDetails.EventID == eventDetail.EventID {
				archives[i].EventDetails = eventDetail
			}
		}
		for _, image := range images {
			if archive.EventDetails.EventID == image.EventID {
				archives[i].Images = append(archives[i].Images, image)
			}
		}
		for _, clip := range clips {
			if archive.EventDetails.EventID == clip.EventID {
				archives[i].Clips = append(archives[i].Clips, clip)
			}
		}
	}
	// return the data
	json.NewEncoder(w).Encode(archives)
}

func ArchiveEntryDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	vars := mux.Vars(r)
	id := vars["id"]

	sqldb.Init()
	sSQL := "DELETE FROM archive WHERE archiveID = ?"
	_, err := sqldb.DB.Exec(sSQL, id)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var user strt.User
	sqldb.Init()
	fmt.Println("Login")
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sSQL := "SELECT * FROM users WHERE user = ? AND password = ?"
	rows, err := sqldb.DB.Query(sSQL, user.Username, user.Password)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&user.UserID, &user.Username, &user.Password, &user.Role)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if user.Username == "" {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func MessagesGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	sqldb.Init()
	var messages []strt.Message
	sSQL := "SELECT * FROM messages"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var message strt.Message
		err = rows.Scan(&message.MessageID, &message.MessageContent)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		messages = append(messages, message)
	}
	json.NewEncoder(w).Encode(messages)
}

func MessagePOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var message strt.Message
	sqldb.Init()
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sSQL := "INSERT INTO messages (messageContent) VALUES (?)"
	_, err = sqldb.DB.Exec(sSQL, message.MessageContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func MessageDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	sqldb.Init()
	fmt.Println("MessageDELETE")
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("MessageDELETE" + id)
	sSQL := "DELETE FROM messages WHERE messageID = ?"
	_, err := sqldb.DB.Exec(sSQL, id)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func MessagePUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var message strt.Message
	sqldb.Init()
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sSQL := "UPDATE messages SET messageContent = ? WHERE messageID = ?"
	_, err = sqldb.DB.Exec(sSQL, message.MessageContent, message.MessageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpcomingEventsListsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	sqldb.Init()
	var events []strt.EventDetails
	var event strt.EventDetails
	sSQL := "SELECT * FROM choirevents WHERE eventDate >= curdate() ORDER BY eventDate"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	sSQLDate := ""
	// Modify	eventID	location	eventDate	startTime	endTime	price	title	meetingPoint	invitation
	for rows.Next() {
		err = rows.Scan(&event.EventID, &event.Location, &sSQLDate, &event.StartTime, &event.EndTime, &event.Price, &event.Title, &event.MeetingPoint, &event.Invitation)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// convert and format the date
		parsedDate, err := ParseMySQLDateTime(sSQLDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		event.EventDate = parsedDate.Format("2006-01-02")
		events = append(events, event)
	}

	json.NewEncoder(w).Encode(events)
}

func UpcomingPlaylistsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	sqldb.Init()
	var events []strt.EventDetails
	var playlists []strt.PlaylistEntry
	var musicTracks []strt.MusicTrack
	var event strt.EventDetails
	var playlist strt.PlaylistEntry
	var musicTrack strt.MusicTrack
	sSQL := "SELECT choirevents.eventID, choirevents.location, choirevents.eventDate, choirevents.startTime, choirevents.endTime, choirevents.title, choirevents.meetingPoint, playlists.playlistID playlists.playorder, music.musicTrackID, music.trackName, music.lyrics, music.soprano, music.alto, music.tenor, music.allParts FROM choirevents LEFT OUTER JOIN (playlists JOIN music on playlists.musicID=music.musicTrackID) on choirevents.eventID=playlists.eventID WHERE choirevents.eventDate >= curdate() ORDER BY choirevents.eventDate, playlists.playorder"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sSQLDate := ""
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&event.EventID, &event.Location, &sSQLDate, &event.StartTime, &event.EndTime, &event.Title, &event.MeetingPoint, &playlist.PlaylistID, &playlist.Playorder, &musicTrack.ID, &musicTrack.TrackName, &musicTrack.Lyrics, &musicTrack.Soprano, &musicTrack.Alto, &musicTrack.Tenor, &musicTrack.AllParts)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// convert and format the date
		parsedDate, err := ParseMySQLDateTime(sSQLDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		event.EventDate = parsedDate.Format("2006-01-02")
		_ = append(events, event)
		_ = append(playlists, playlist)
		_ = append(musicTracks, musicTrack)
	}
	json.NewEncoder(w).Encode(events)
}

func EventArchive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	vars := mux.Vars(r)
	sqldb.Init()
	id := vars["id"]
	var archive strt.ArchiveEntry
	var event strt.EventDetails
	sSQL := "SELECT archive.archiveID, choirevents.location, choirevents.eventDate, choirevents.title, archive.report FROM choirevents JOIN archive ON archive.eventID=choirevents.eventID WHERE choirevents.eventID = ?"
	rows, err := sqldb.DB.Query(sSQL, id)
	if err != nil {
		log.Fatal(err)
		archive.ArchiveID = -1
		archive.Report = err.Error()
		json.NewEncoder(w).Encode(archive)
		return
	}
	defer rows.Close()
	sSQLDate := ""
	for rows.Next() {

		err = rows.Scan(&archive.ArchiveID, &event.Location, &sSQLDate, &event.Title, &archive.Report)
		if err != nil {
			log.Fatal(err)
			archive.ArchiveID = -1
			archive.Report = err.Error()
			json.NewEncoder(w).Encode(archive)
			return
		}
		parsedDate, err := ParseMySQLDateTime(sSQLDate)
		if err != nil {
			archive.ArchiveID = -1
			archive.Report = err.Error()
			json.NewEncoder(w).Encode(archive)
			return
		}
		event.EventID, err = strconv.Atoi(id)
		if err != nil {
			archive.ArchiveID = -1
			archive.Report = err.Error()
			json.NewEncoder(w).Encode(archive)
			return
		}
		event.EventDate = parsedDate.Format("2006-01-02")
		archive.EventDetails = event
	}
	//retrive the ImageDetails for the event
	sSQL = "SELECT * FROM images WHERE eventID = ?"
	rows, err = sqldb.DB.Query(sSQL, archive.EventDetails.EventID)
	if err != nil {
		log.Fatal(err)
		archive.ArchiveID = -1
		archive.Report = err.Error()
		json.NewEncoder(w).Encode(archive)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var image strt.ImageDetail
		err = rows.Scan(&image.ImageID, &image.Filename, &image.Caption, &image.EventID, &image.Height, &image.Width)
		if err != nil {
			log.Fatal(err)
			archive.ArchiveID = -1
			archive.Report = err.Error()
			json.NewEncoder(w).Encode(archive)
			return
		}
		archive.Images = append(archive.Images, image)
	}
	//retrive the ClipDetails for the event
	sSQL = "SELECT * FROM clips WHERE eventID = ?"
	rows, err = sqldb.DB.Query(sSQL, archive.EventDetails.EventID)
	if err != nil {
		log.Fatal(err)
		archive.ArchiveID = -1
		archive.Report = err.Error()
		json.NewEncoder(w).Encode(archive)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var clip strt.Clip
		err = rows.Scan(&clip.ClipID, &clip.ClipURL, &clip.EventID, &clip.Caption)
		if err != nil {
			log.Fatal(err)
			archive.ArchiveID = -1
			archive.Report = err.Error()
			json.NewEncoder(w).Encode(archive)
			return
		}
		archive.Clips = append(archive.Clips, clip)
	}
	// return the data
	json.NewEncoder(w).Encode(archive)
}

func EventImages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	vars := mux.Vars(r)
	sqldb.Init()
	id := vars["id"]
	var images []strt.ImageDetail
	var image strt.ImageDetail
	sSQL := "SELECT * FROM images WHERE eventID = ?"
	rows, err := sqldb.DB.Query(sSQL, id)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&image.ImageID, &image.Filename, &image.Caption, &image.EventID, &image.Height, &image.Width)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		images = append(images, image)
	}
	json.NewEncoder(w).Encode(images)
}

func EventClips(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	vars := mux.Vars(r)
	sqldb.Init()
	id := vars["id"]
	var clips []strt.Clip
	var clip strt.Clip
	sSQL := "SELECT * FROM clips WHERE eventID = ?"
	rows, err := sqldb.DB.Query(sSQL, id)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&clip.ClipID, &clip.ClipURL, &clip.EventID, &clip.Caption)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		clips = append(clips, clip)
	}
	json.NewEncoder(w).Encode(clips)
}

func ClipDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("ClipDelete " + id)
	sqldb.Init()
	sSQL := "DELETE FROM clips WHERE clipID = ?"
	_, err := sqldb.DB.Exec(sSQL, id)
	blankClip := strt.Clip{}
	if err != nil {
		log.Fatal(err)
		blankClip.ClipID = -1
		blankClip.ClipURL = err.Error()
		json.NewEncoder(w).Encode(blankClip)
		return
	}
	blankClip.ClipID = 200
	blankClip.ClipURL = "Clip Deleted"
	json.NewEncoder(w).Encode(blankClip)
}

func RandomImagesGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	sqldb.Init()
	// the records value is a json object for a screensize
	vars := mux.Vars(r)
	var rec strt.ScreenSize
	err := json.Unmarshal([]byte(vars["scr"]), &rec)
	if err != nil {
		http.Error(w, "Invalid screen size format", http.StatusBadRequest)
		return
	}
	// based on the screensize establish if it is desktop or mobile
	filepath := "Images/Desktop"
	const (
		maxMobileWidth  = 800
		maxMobileHeight = 1280
		minPixelRatio   = 1.5
	)

	if (rec.Width <= maxMobileWidth && rec.Height <= maxMobileHeight) || float64(rec.DevicePixelRatio) >= minPixelRatio {
		filepath = "Images/Mobile"
	}

	sSQL := "SELECT filename FROM images ORDER BY RAND() LIMIT ?"
	rows, err := sqldb.DB.Query(sSQL, rec.Images)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	//declare an array of image files
	var imageName string
	var imageFiles []string
	for rows.Next() {
		err = rows.Scan(&imageName)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// if the image is for mobile then the dt prefix has to be replaced with mb
		if filepath == "Images/Mobile" {
			imageName = strings.Replace(imageName, "dt", "mb", 1)
		}
		// add the image name to the path and add it to the array
		imageName = filepath + "/" + imageName
		imageFiles = append(imageFiles, imageName)
	}
	// return the array of image files
	json.NewEncoder(w).Encode(imageFiles)
}

func ThemeDetailsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	sqldb.Init()
	var theme strt.ThemeDetails
	sqldb.Init()
	sSQL := "SELECT * FROM themedetails"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {

		err = rows.Scan(&theme.ID, &theme.BoxColour, &theme.TextColour, &theme.TextFont, &theme.BackgroundImage, &theme.TextboxColour, &theme.LogoImage, &theme.BannerColour, &theme.MenuColour, &theme.ButtonColour, &theme.ButtonHover, &theme.ButtonTextColour, &theme.MenuTextColour, &theme.TextSize)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json.NewEncoder(w).Encode(theme)
}

func ThemeDetailsPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var theme strt.ThemeDetails
	err := json.NewDecoder(r.Body).Decode(&theme)
	sqldb.Init()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sSQL := "UPDATE themedetails SET boxColour = ?, textColour = ?, textFont = ?, backgroundImage = ?, textboxColour = ?, logoimage = ?, bannerColour = ?, menuColour = ?, buttonColour = ?, buttonHover = ?, buttonTextColour = ?, menuTextColour = ?, textSize = ?"
	_, err = sqldb.DB.Exec(sSQL, theme.BoxColour, theme.TextColour, theme.TextFont, theme.BackgroundImage, theme.TextboxColour, theme.LogoImage, theme.BannerColour, theme.MenuColour, theme.ButtonColour, theme.ButtonHover, theme.ButtonTextColour, theme.MenuTextColour, theme.TextSize)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ThemeDetailsRandom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var theme strt.ThemeDetails
	sqldb.Init()
	theme.BoxColour = "#" + fmt.Sprintf("%06x", rand.Intn(0xffffff))
	theme.TextColour = "#" + fmt.Sprintf("%06x", rand.Intn(0xffffff))
	theme.TextFont = "Impact"
	theme.BackgroundImage.String = "Musical Background.png"
	theme.TextboxColour = "#" + fmt.Sprintf("%06x", rand.Intn(0xffffff))
	theme.LogoImage.String = "Choir Logo.png"
	theme.BannerColour = "#" + fmt.Sprintf("%06x", rand.Intn(0xffffff))
	theme.MenuColour = "#" + fmt.Sprintf("%06x", rand.Intn(0xffffff))
	theme.ButtonColour = "#" + fmt.Sprintf("%06x", rand.Intn(0xffffff))
	theme.ButtonHover = "#" + fmt.Sprintf("%06x", rand.Intn(0xffffff))
	theme.ButtonTextColour = "#" + fmt.Sprintf("%06x", rand.Intn(0xffffff))
	theme.MenuTextColour = "#" + fmt.Sprintf("%06x", rand.Intn(0xffffff))
	theme.TextSize = 12
	sSQL := "UPDATE themedetails SET boxColour = ?, textColour = ?, textFont = ?, backgroundImage = ?, textboxColour = ?, logoimage = ?, bannerColour = ?, menuColour = ?, buttonColour = ?, buttonHover = ?, buttonTextColour = ?, menuTextColour = ?, textSize = ?"
	_, err := sqldb.DB.Exec(sSQL, theme.BoxColour, theme.TextColour, theme.TextFont, theme.BackgroundImage, theme.TextboxColour, theme.LogoImage, theme.BannerColour, theme.MenuColour, theme.ButtonColour, theme.ButtonHover, theme.ButtonTextColour, theme.MenuTextColour, theme.TextSize)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func MusicListGET(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MusicListGET")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var musicTracks []strt.MusicTrack
	var musicTrack strt.MusicTrack
	// open the database
	sqldb.Init()
	sSQL := "SELECT * FROM music"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		fmt.Println("Error in query")
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		// musicTrackID	int Auto Increment, trackName	varchar(100)	,artist	varchar(60)	,lyrics	varchar(120)	,soprano	varchar(120)	, alto	varchar(120)	,tenor	varchar(120)	,allParts	varchar(120)	,piano
		err = rows.Scan(&musicTrack.ID, &musicTrack.TrackName, &musicTrack.Artist, &musicTrack.Lyrics, &musicTrack.Artist, &musicTrack.Soprano, &musicTrack.Alto, &musicTrack.Tenor, &musicTrack.AllParts)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error in MusicListGET")
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		musicTracks = append(musicTracks, musicTrack)
	}
	json.NewEncoder(w).Encode(musicTracks)
}

func MusicTrackGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	vars := mux.Vars(r)
	id := vars["id"]
	sqldb.Init()
	var musicTrack strt.MusicTrack
	sSQL := "SELECT * FROM music WHERE musicTrackID = ?"
	rows, err := sqldb.DB.Query(sSQL, id)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&musicTrack.ID, &musicTrack.TrackName, &musicTrack.Lyrics, &musicTrack.Soprano, &musicTrack.Alto, &musicTrack.Tenor, &musicTrack.AllParts)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json.NewEncoder(w).Encode(musicTrack)
}

func MusicTrackPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var musicTrack strt.MusicTrack
	sqldb.Init()
	err := json.NewDecoder(r.Body).Decode(&musicTrack)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sSQL := "INSERT INTO music (trackName, lyrics, soprano, alto, tenor, allParts) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = sqldb.DB.Exec(sSQL, musicTrack.TrackName, musicTrack.Lyrics, musicTrack.Soprano, musicTrack.Alto, musicTrack.Tenor, musicTrack.AllParts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func MusicTrackDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	vars := mux.Vars(r)
	id := vars["id"]
	sqldb.Init()
	sSQL := "DELETE FROM music WHERE musicTrackID = ?"
	_, err := sqldb.DB.Exec(sSQL, id)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func MusicTrackPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var musicTrack strt.MusicTrack
	sqldb.Init()
	err := json.NewDecoder(r.Body).Decode(&musicTrack)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sSQL := "UPDATE music SET trackName = ?, lyrics = ?, soprano = ?, alto = ?, tenor = ?, allParts = ? WHERE musicTrackID = ?"
	_, err = sqldb.DB.Exec(sSQL, musicTrack.TrackName, musicTrack.Lyrics, musicTrack.Soprano, musicTrack.Alto, musicTrack.Tenor, musicTrack.AllParts, musicTrack.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ClipPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var clip strt.Clip
	sqldb.Init()
	err := json.NewDecoder(r.Body).Decode(&clip)
	if err != nil {
		clip.ClipID = -1
		clip.ClipURL = "Error: " + err.Error()
		json.NewEncoder(w).Encode(clip)
		return
	}
	fmt.Println("ClipPOST" + clip.ClipURL)
	sSQL := "INSERT INTO clips (clipURL, eventID, caption) VALUES (?, ?, ?)"
	_, err = sqldb.DB.Exec(sSQL, clip.ClipURL, clip.EventID, clip.Caption)
	if err != nil {
		clip.ClipID = -1
		clip.ClipURL = "Error: " + err.Error()
		json.NewEncoder(w).Encode(clip)

		return
	}

	// retrieve the clipID
	sSQL = "SELECT clipID FROM clips ORDER BY clipID DESC LIMIT 1"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Fatal(err)
		// return the error
		clip.ClipID = -1
		clip.ClipURL = "Error: " + err.Error()
		json.NewEncoder(w).Encode(clip)

		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&clip.ClipID)
		if err != nil {
			log.Fatal(err)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})

			return
		}
	}
	// return the clip as a json object
	json.NewEncoder(w).Encode(clip)
}

func ArchiveEntryPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var archive strt.ArchiveEntry
	sqldb.Init()
	err := json.NewDecoder(r.Body).Decode(&archive)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sSQL := "UPDATE archive SET report = ? WHERE archiveID = ?"
	_, err = sqldb.DB.Exec(sSQL, archive.Report, archive.ArchiveID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ArchiveEntryPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var archive strt.ArchiveEntry
	sqldb.Init()
	err := json.NewDecoder(r.Body).Decode(&archive)

	if err != nil {
		archive.ArchiveID = -1
		archive.Report = err.Error()
		fmt.Println("ArchiveEntryPOST 1 " + archive.Report)
		json.NewEncoder(w).Encode(archive)
		return
	}
	if archive.ArchiveID > 0 {
		sSQL := "UPDATE archive SET eventID = ?, report = ? WHERE archiveID = ?"
		_, err = sqldb.DB.Exec(sSQL, archive.EventDetails.EventID, archive.Report, archive.ArchiveID)
		if err != nil {
			archive.ArchiveID = -1
			archive.Report = err.Error()
			fmt.Println("ArchiveEntryPOST 2" + archive.Report)
			json.NewEncoder(w).Encode(archive)
			return
		}

	} else {
		sSQL := "INSERT INTO archive (eventID, report) VALUES (?, ?)"
		_, err = sqldb.DB.Exec(sSQL, archive.EventDetails.EventID, archive.Report)
		if err != nil {
			archive.ArchiveID = -1
			archive.Report = err.Error()
			fmt.Println("ArchiveEntryPOST 3" + archive.Report)
			json.NewEncoder(w).Encode(archive)
			return
		}

	}
	archive.ArchiveID = 200
	archive.Report = "Archive Entry Added"
	fmt.Println("ArchiveEntryPOST 4" + archive.Report)
	json.NewEncoder(w).Encode(archive)
}

func PlaylistsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var playlists []strt.PlaylistEntry
	var playlist strt.PlaylistEntry
	sqldb.Init()
	sSQL := "SELECT * FROM playlists"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&playlist.PlaylistID, &playlist.ID, &playlist.Playorder)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		playlists = append(playlists, playlist)
	}
	json.NewEncoder(w).Encode(playlists)
}

func PlaylistPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var playlists []strt.PlaylistEntry
	sqldb.Init()
	err := json.NewDecoder(r.Body).Decode(&playlists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sSQL := "INSERT INTO playlists (musicID, playorder) VALUES (?, ?)"
	for _, playlist := range playlists {
		_, err = sqldb.DB.Exec(sSQL, playlist.ID, playlist.Playorder)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func PlaylistDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	vars := mux.Vars(r)
	sqldb.Init()
	id := vars["id"]
	sSQL := "DELETE FROM playlists WHERE playlistID = ?"
	_, err := sqldb.DB.Exec(sSQL, id)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func PlaylistPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var playlist strt.PlaylistEntry
	sqldb.Init()
	err := json.NewDecoder(r.Body).Decode(&playlist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sSQL := "UPDATE playlists SET musicID = ?, playorder = ? WHERE playlistID = ?"
	_, err = sqldb.DB.Exec(sSQL, playlist.ID, playlist.Playorder, playlist.PlaylistID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func PlaylistGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	vars := mux.Vars(r)
	id := vars["id"]
	sqldb.Init()
	var playlist strt.PlaylistEntry
	sSQL := "SELECT * FROM playlists WHERE playlistID = ?"
	rows, err := sqldb.DB.Query(sSQL, id)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&playlist.PlaylistID, &playlist.ID, &playlist.Playorder)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json.NewEncoder(w).Encode(playlist)
}

func EventGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

	vars := mux.Vars(r)
	id := vars["id"]

	sqldb.Init()

	var event strt.EventDetails
	sSQL := "SELECT * FROM choirevents WHERE eventID = ?"
	rows, err := sqldb.DB.Query(sSQL, id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error fetching event details", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var sSQLDate string
		err := rows.Scan(&event.EventID, &event.Location, &sSQLDate, &event.StartTime, &event.EndTime, &event.Price, &event.Title, &event.MeetingPoint, &event.Invitation)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error parsing event details", http.StatusInternalServerError)
			return
		}

		if sSQLDate != "" {
			parsedDate, err := ParseMySQLDateTime(sSQLDate)
			if err != nil {
				log.Println(err)
				http.Error(w, "Invalid event date format", http.StatusBadRequest)
				continue
			}
			event.EventDate = parsedDate.Format("2006-01-02")
		}
	}

	json.NewEncoder(w).Encode(event)
}

func EventsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	sqldb.Init()
	var events []strt.EventDetails
	var event strt.EventDetails
	sSQL := "SELECT eventID, location, eventDate, title FROM choirevents" // WHERE eventDate >= curdate()"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	sSQLDate := ""
	for rows.Next() {
		err = rows.Scan(&event.EventID, &event.Location, &sSQLDate, &event.Title)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		parsedDate, err := ParseMySQLDateTime(sSQLDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		event.EventDate = parsedDate.Format("2006-01-02")
		events = append(events, event)
	}
	json.NewEncoder(w).Encode(events)
}

func EventsUpcomingGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	sqldb.Init()
	var events []strt.EventDetails
	var event strt.EventDetails
	sSQL := "SELECT * FROM choirevents WHERE eventDate >= curdate()"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	sSQLDate := ""
	for rows.Next() {
		err = rows.Scan(&event.EventID, &event.Location, &sSQLDate, &event.StartTime, &event.EndTime, &event.Price, &event.Title, &event.MeetingPoint, &event.Invitation)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		parsedDate, err := ParseMySQLDateTime(sSQLDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		event.EventDate = parsedDate.Format("2006-01-02")
		events = append(events, event)
	}
	json.NewEncoder(w).Encode(events)
}

func EventPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var event strt.EventDetails
	sqldb.Init()
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sSQL := "INSERT INTO choirevents (location, eventDate, startTime, endTime, price, title, invitation, meetingPoint) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = sqldb.DB.Exec(sSQL, event.Location, event.EventDate, event.StartTime, event.EndTime, event.Price, event.Title, event.Invitation, event.MeetingPoint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func EventDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	vars := mux.Vars(r)
	id := vars["id"]
	sqldb.Init()
	sSQL := "DELETE FROM choirevents WHERE eventID = ?"
	_, err := sqldb.DB.Exec(sSQL, id)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func EventPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var event strt.EventDetails
	sqldb.Init()
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sSQL := "UPDATE choirevents SET location = ?, eventDate = ?, startTime = ?, endTime = ?, price = ?, title = ?, invitation = ?, meetingPoint = ? WHERE eventID = ?"
	_, err = sqldb.DB.Exec(sSQL, event.Location, event.EventDate, event.StartTime, event.EndTime, event.Price, event.Title, event.Invitation, event.MeetingPoint, event.EventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ImageGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var image strt.ImageDetail
	vars := mux.Vars(r)
	id := vars["id"]
	sqldb.Init()
	sSQL := "SELECT * FROM images WHERE imageID = ?"
	rows, err := sqldb.DB.Query(sSQL, id)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&image.ImageID, &image.Filename, &image.Caption, &image.EventID, &image.Height, &image.Width)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json.NewEncoder(w).Encode(image)
}

func ImagesGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var images []strt.ImageDetail
	var image strt.ImageDetail
	sqldb.Init()
	sSQL := "SELECT * FROM images"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&image.ImageID, &image.Filename, &image.Caption, &image.EventID, &image.Height, &image.Width)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		images = append(images, image)
	}
	json.NewEncoder(w).Encode(images)
}

func isValidImage(file multipart.File) bool {
	_, _, err := image.Decode(file)
	return err == nil
}

func isValidFileName(filename string) bool {
	// convert to lowercase
	filename = strings.ToLower(filename)
	return strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") || strings.HasSuffix(filename, ".png")
}

func isValidFileSize(size int64) bool {
	return size <= 50<<20
}

func isValidFileType(fileType string) bool {
	return strings.HasPrefix(fileType, "image/")
}

// ImageFilePOST is used to upload an image file to the server
// The images will either be for the background, desktop or mobile and will be stored in the images directory
func ImageFilePOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	fmt.Println("Content-Type:", r.Header.Get("Content-Type"))

	// Parse the multipart form TO 50MB
	err := r.ParseMultipartForm(50 << 20)
	if err != nil {
		fmt.Println("Error parsing multipart form:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve the file from the form
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error retrieving the file:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Validate file type
	fileType := handler.Header.Get("Content-Type")
	if !isValidFileType(fileType) {
		fmt.Println("Invalid file type:", fileType)
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	// Validate file size
	if !isValidFileSize(handler.Size) {
		fmt.Println("File size exceeds limit:", handler.Size)
		http.Error(w, "File size exceeds limit", http.StatusBadRequest)
		return
	}

	// Validate file name
	if !isValidFileName(handler.Filename) {
		fmt.Println("Invalid file name:", handler.Filename)
		http.Error(w, "Invalid file name", http.StatusBadRequest)
		return
	}

	// Validate image content (if applicable)
	if strings.HasPrefix(fileType, "data:image/") {
		if !isValidImage(file) {
			fmt.Println("Invalid image content")
			http.Error(w, "Invalid image content", http.StatusBadRequest)
			return
		}
	}

	// Create a temporary file
	tempFile, err := os.CreateTemp("temp-images", "upload-*.png")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	// Read the file bytes
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file bytes:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the file bytes to the temporary file
	_, err = tempFile.Write(fileBytes)
	if err != nil {
		fmt.Println("Error writing to temporary file:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Determine the destination directory based on the filename prefix
	var destinationDir string
	switch {
	case strings.HasPrefix(handler.Filename, "bg"):
		destinationDir = "images/background/"
	case strings.HasPrefix(handler.Filename, "dt"):
		destinationDir = "images/desktop/"
	case strings.HasPrefix(handler.Filename, "mb"):
		destinationDir = "images/mobile/"
	default:
		destinationDir = "images/"
	}

	// Move the file to the destination directory
	err = os.Rename(tempFile.Name(), destinationDir+handler.Filename)
	if err != nil {
		fmt.Println("Error moving file:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialize the image details
	var image strt.ImageDetail
	image.Filename = handler.Filename

	// Get additional file details from the form
	caption := r.FormValue("caption")
	eventID := r.FormValue("eventID")
	height := r.FormValue("height")
	width := r.FormValue("width")

	// Convert eventID, height, and width to integers
	image.EventID, err = strconv.Atoi(eventID)
	if err != nil {
		fmt.Println("Error converting eventID:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	image.Height, err = strconv.Atoi(height)
	if err != nil {
		fmt.Println("Error converting height:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	image.Width, err = strconv.Atoi(width)
	if err != nil {
		fmt.Println("Error converting width:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	image.Caption = caption
	// if the image has an mb prefix then it doesnt need to go on the database
	// so return a json encode response
	if strings.HasPrefix(handler.Filename, "mb") {
		json.NewEncoder(w).Encode(image)
		return
	}
	// Insert the file details into the database
	sqldb.Init()
	sSQL := "INSERT INTO images (filename, caption, eventID, height, width) VALUES (?, ?, ?, ?, ?)"
	_, err = sqldb.DB.Exec(sSQL, image.Filename, image.Caption, image.EventID, image.Height, image.Width)
	if err != nil {
		fmt.Println("Error inserting into database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Println("Record uploaded successfully")
	// Retrieve the imageID of the newly inserted record
	sSQL = "SELECT imageID FROM images ORDER BY imageID DESC LIMIT 1"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		fmt.Println("Error querying database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type imageID struct {
		ImageID int
	}
	Im := imageID{}
	defer rows.Close()
	fmt.Println("Rows:", rows)
	if rows.Next() {
		err = rows.Scan(&Im.ImageID)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Respond with the image details
	// translate the image object to a string
	sResp := fmt.Sprintf("%v", Im)
	fmt.Println(sResp)
	// fmt.Println(json.NewEncoder(w).Encode(Im))
	json.NewEncoder(w).Encode(Im)
}

func ImagesPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var images []strt.ImageDetail
	sqldb.Init()
	err := json.NewDecoder(r.Body).Decode(&images)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sSQL := "INSERT INTO images (filename, caption, eventID, height, width) VALUES (?, ?, ?, ?, ?)"
	for _, image := range images {
		_, err = sqldb.DB.Exec(sSQL, image.Filename, image.Caption, image.EventID, image.Height, image.Width)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sSQL = "SELECT imageID FROM images LIMIT 1 ORDER BY imageID DESC"
		rows, err := sqldb.DB.Query(sSQL)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&image.ImageID)
			if err != nil {
				log.Fatal(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	//	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:")
	//  return the image object
	json.NewEncoder(w).Encode(images)

}

func ImageDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	vars := mux.Vars(r)
	id := vars["id"]
	sqldb.Init()
	// we need to get the filename so we can delete the files
	sSQL := "SELECT filename FROM images WHERE imageID = ?"
	rows, err := sqldb.DB.Query(sSQL, id)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var filename string
	for rows.Next() {
		err = rows.Scan(&filename)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// delete the file. If the prefix is bg, dt or mb then we need to delete from the appropriate directory
	if strings.HasPrefix(filename, "bg") {
		err = os.Remove("images/background/" + filename)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if strings.HasPrefix(filename, "dt") {
		err = os.Remove("Images/Desktop/" + filename)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// then remove the prefix and replace it with mb
		mbfile := strings.Replace(filename, "dt", "mb", 1)
		err = os.Remove("Images/Mobile/" + mbfile)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if strings.HasPrefix(filename, "lg") {
		err = os.Remove("Images/Logo/" + filename)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// now delete the record from the database

	sSQL = "DELETE FROM images WHERE imageID = ?"
	_, err = sqldb.DB.Exec(sSQL, id)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Image deleted successfully"})
}

func ImagePUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var image strt.ImageDetail
	sqldb.Init()
	err := json.NewDecoder(r.Body).Decode(&image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sSQL := "UPDATE images SET filename = ?, caption = ?, eventID = ?, height = ?, width = ? WHERE imageID = ?"
	_, err = sqldb.DB.Exec(sSQL, image.Filename, image.Caption, image.EventID, image.Height, image.Width, image.ImageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func SiteInfoGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var siteInfo strt.SiteInfo
	sqldb.Init()
	sSQL := "SELECT * FROM siteinfo"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Error in SiteInfoGET")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&siteInfo.ID, &siteInfo.HomeTitle, &siteInfo.HomeText, &siteInfo.AboutTitle, &siteInfo.AboutText, &siteInfo.ArchiveTitle, &siteInfo.ArchiveText, &siteInfo.NoticesTitle, &siteInfo.NoticesText, &siteInfo.BookingTitle, &siteInfo.BookingText, &siteInfo.MembersTitle, &siteInfo.MembersText, &siteInfo.AppealTitle, &siteInfo.AppealText, &siteInfo.SettingsTitle, &siteInfo.SettingsText)
		if err != nil {

			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json.NewEncoder(w).Encode(siteInfo)
}

func SiteInfoPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var siteInfo strt.SiteInfo
	sqldb.Init()
	err := json.NewDecoder(r.Body).Decode(&siteInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sSQL := "UPDATE siteinfo SET HomeTitle = ?, HomeText = ?, AboutTitle = ?, AboutText = ?, ArchiveTitle = ?, ArchiveText = ?, NoticesTitle = ?, NoticesText = ?, BookingTitle = ?, BookingText = ?, MembersTitle = ?, MembersText = ?, AppealTitle = ?, AppealText = ?, SettingsTitle = ?, SettingsText = ?"
	_, err = sqldb.DB.Exec(sSQL, siteInfo.HomeTitle, siteInfo.HomeText, siteInfo.AboutTitle, siteInfo.AboutText, siteInfo.ArchiveTitle, siteInfo.ArchiveText, siteInfo.NoticesTitle, siteInfo.NoticesText, siteInfo.BookingTitle, siteInfo.BookingText, siteInfo.MembersTitle, siteInfo.MembersText, siteInfo.AppealTitle, siteInfo.AppealText, siteInfo.SettingsTitle, siteInfo.SettingsText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// to build the project run the following command
// go build -o bin/choirapi main.go
// to run the project run the following command
// ./bin/choirapi
