package events

import (
	// import the mysql driver

	strt "RWTAPI/structures"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	sqldb "RWTAPI/sqldb"

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

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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

	vars := mux.Vars(r)
	id := vars["id"]

	var archive strt.ArchiveEntry
	var event strt.EventDetails
	sSQL := "SELECT archive.archiveID, choirevents.location, choirevents.eventDate, choirevents.title, archive.report FROM choirevents JOIN archive ON archive.eventID=choirevents.eventID WHERE choirevents.eventID = ?"
	rows, err := sqldb.DB.Query(sSQL, id)
	if err != nil {
		log.Println("Error:", err)
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
			log.Println("Error:", err)
			archive.ArchiveID = -1
			archive.Report = err.Error()
			json.NewEncoder(w).Encode(archive)
			return
		}

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
		log.Println("Error:", err)
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
			log.Println("Error:", err)
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
		log.Println("Error:", err)
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
			log.Println("Error:", err)
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
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// Initialize the database connection

	// Get the screen size and number of images to return
	screen := r.URL.Query().Get("screen")
	imagesStr := r.URL.Query().Get("archives")
	iRecords, err := strconv.Atoi(imagesStr)
	if err != nil {
		http.Error(w, "Invalid archives parameter", http.StatusBadRequest)
		return
	}

	// Determine if the device is desktop or mobile
	filepath := "images/desktop"
	if screen == "mobile" {
		filepath = "images/mobile"
	} else {
		filepath = "images/desktop"
	}

	// Create slices to hold the data
	var archives []strt.ArchiveEntry
	var events string
	var images []strt.ImageDetail
	var clips []strt.Clip
	var errorArch strt.ArchiveEntry
	var errorArray []strt.ArchiveEntry

	// Gather the archive details from the database
	sSQLDate := ""
	sSQL := "SELECT archive.archiveID, choirevents.location, choirevents.eventDate, choirevents.title, archive.report, archive.eventID FROM choirevents JOIN archive ON archive.eventID = choirevents.eventID WHERE choirevents.eventDate < CURDATE() ORDER BY choirevents.eventDate DESC LIMIT ?"
	rows, err := sqldb.DB.Query(sSQL, iRecords)
	if err != nil {
		log.Println("Error querying database:", err)
		errorArch.ArchiveID = -1
		errorArch.Report = err.Error()
		errorArray = append(errorArray, errorArch)
		json.NewEncoder(w).Encode(errorArray)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var archive strt.ArchiveEntry
		err = rows.Scan(&archive.ArchiveID, &archive.EventDetails.Location, &sSQLDate, &archive.EventDetails.Title, &archive.Report, &archive.EventDetails.EventID)
		if err != nil {
			log.Println("Error scanning row:", err)
			errorArch.ArchiveID = -1
			errorArch.Report = err.Error()
			errorArray = append(errorArray, errorArch)
			json.NewEncoder(w).Encode(errorArray)
			return
		}
		parsedDate, err := ParseMySQLDateTime(sSQLDate)
		if err != nil {
			log.Println("Error parsing date:", err)
			errorArch.ArchiveID = -1
			errorArch.Report = err.Error()
			errorArray = append(errorArray, errorArch)
			json.NewEncoder(w).Encode(errorArray)
			return
		}
		archive.EventDetails.EventDate = parsedDate.Format("2006-01-02")
		archives = append(archives, archive)
		events = events + strconv.Itoa(archive.EventDetails.EventID) + ","
	}

	// Remove the trailing comma if events is not empty
	if len(events) > 0 {
		events = events[:len(events)-1]
		events = "(" + events + ")"
	} else {
		json.NewEncoder(w).Encode(archives)
		return
	}

	// Get the images from all of the events
	sSQL = "SELECT * FROM images WHERE eventID IN " + events
	// fmt.Println(sSQL)
	rows, err = sqldb.DB.Query(sSQL)
	if err != nil {
		log.Println("Error querying images:", err)
		errorArch.ArchiveID = -1
		errorArch.Report = err.Error()
		errorArray = append(errorArray, errorArch)
		json.NewEncoder(w).Encode(errorArray)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var image strt.ImageDetail
		err = rows.Scan(&image.ImageID, &image.Filename, &image.Caption, &image.EventID, &image.Height, &image.Width)
		if err != nil {
			log.Println("Error scanning image row:", err)
			errorArch.ArchiveID = -1
			errorArch.Report = err.Error()
			errorArray = append(errorArray, errorArch)
			json.NewEncoder(w).Encode(errorArray)
			return
		}
		if filepath == "images/mobile" {
			image.Filename = strings.Replace(image.Filename, "dt", "mb", 1)
		}
		image.Filename = filepath + "/" + url.PathEscape(image.Filename)
		fmt.Println(image.Filename)
		images = append(images, image)
	}

	// Get the clips from all of the events
	sSQL = "SELECT * FROM clips WHERE eventID IN " + events
	rows, err = sqldb.DB.Query(sSQL)
	if err != nil {
		log.Println("Error querying clips:", err)
		errorArch.ArchiveID = -1
		errorArch.Report = err.Error()
		errorArray = append(errorArray, errorArch)
		json.NewEncoder(w).Encode(errorArray)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var clip strt.Clip
		err = rows.Scan(&clip.ClipID, &clip.ClipURL, &clip.EventID, &clip.Caption)
		if err != nil {
			log.Println("Error scanning clip row:", err)
			errorArch.ArchiveID = -1
			errorArch.Report = err.Error()
			errorArray = append(errorArray, errorArch)
			json.NewEncoder(w).Encode(errorArray)
			return
		}
		clips = append(clips, clip)
	}

	// Combine the data into the archive entries
	for i, archive := range archives {
		for _, image := range images {
			if image.EventID == archive.EventDetails.EventID {
				archives[i].Images = append(archives[i].Images, image)
			}
		}
		for _, clip := range clips {
			if clip.EventID == archive.EventDetails.EventID {
				archives[i].Clips = append(archives[i].Clips, clip)
			}
		}
	}

	// Return the data
	json.NewEncoder(w).Encode(archives)
}

func ArchiveEntryDELETE(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	sSQL := "DELETE FROM archive WHERE archiveID = ?"
	_, err := sqldb.DB.Exec(sSQL, id)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func Login(w http.ResponseWriter, r *http.Request) {

	//
	var user strt.User

	fmt.Println("Login")
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sSQL := "SELECT * FROM users WHERE user = ? AND password = ?"
	rows, err := sqldb.DB.Query(sSQL, user.Username, user.Password)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&user.UserID, &user.Username, &user.Password, &user.Role)
		if err != nil {
			log.Println("Error:", err)
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

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var messages []strt.Message
	sSQL := "SELECT * FROM messages"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var message strt.Message
		err = rows.Scan(&message.MessageID, &message.MessageContent)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		messages = append(messages, message)
	}
	json.NewEncoder(w).Encode(messages)
}

func MessagePOST(w http.ResponseWriter, r *http.Request) {

	var message strt.Message
	//fmt.Println("MessagePOST")
	//fmt.Println(r.Body)

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// 	messageID	int Auto Increment
	// messageDate	timestamp NULL
	// messageFrom	varchar(60) NULL
	// messageContent	varchar(500) NULL ['']
	// eventName	varchar(100) ['']
	// eventDate	timestamp
	// eventTime	varchar(25) ['']
	// contactEmail	varchar(100) NULL
	// contactPhone	varchar(20) ['']
	// eventLocation
	sSQL := "INSERT INTO messages (messageDate, messageContent, messageFrom, eventName, eventDate, eventTime, contactEmail, contactPhone, eventLocation) VALUES (NOW(), ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = sqldb.DB.Exec(sSQL, message.MessageContent, message.MessageFrom, message.EventName, message.EventDate, message.EventTime, message.ContactEmail, message.ContactPhone, message.EventLocation)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	message.MessageContent = "message sent"
	// return a json response
	json.NewEncoder(w).Encode(message)
}

func MessageDELETE(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	fmt.Println("MessageDELETE")
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("MessageDELETE" + id)
	sSQL := "DELETE FROM messages WHERE messageID = ?"
	_, err := sqldb.DB.Exec(sSQL, id)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func MessagePUT(w http.ResponseWriter, r *http.Request) {

	var message strt.Message

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

	var events []strt.EventDetails
	var event strt.EventDetails
	sSQL := "SELECT * FROM choirevents WHERE eventDate >= curdate() ORDER BY eventDate"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	// Modify	eventID	location	eventDate	startTime	endTime	price	title	meetingPoint	invitation
	for rows.Next() {
		err = rows.Scan(&event.EventID, &event.Location, &event.EventDate, &event.StartTime, &event.EndTime, &event.Price, &event.Title, &event.MeetingPoint, &event.Invitation)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// convert and format the date
		parsedDate, err := ParseMySQLDateTime(event.EventDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// format the date into a javascript date object
		event.DateString = parsedDate.Format("2006-01-02")
		events = append(events, event)
	}
	json.NewEncoder(w).Encode(events)
}

// GET upcoming playlists
func UpcomingPlaylistsGET(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Initialize the database connection

	var events []strt.EventDetails

	sSQL := `SELECT * FROM choirevents WHERE eventDate >= CURDATE() ORDER BY eventDate`

	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Println("Error querying database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var event strt.EventDetails
		err = rows.Scan(&event.EventID, &event.Location, &event.EventDate, &event.StartTime, &event.EndTime, &event.Price, &event.Title, &event.MeetingPoint, &event.Invitation)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert the date to a formatted string
		parsedDate, err := ParseMySQLDateTime(event.EventDate)
		if err != nil {
			log.Println("Error parsing date:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		event.DateString = parsedDate.Format("2006-01-02")
		events = append(events, event)
	}

	// Close the first rows before starting the second query
	rows.Close()
	// loop through events to get the playlists
	for i := range events {
		sSQL = `SELECT playlistID, eventID, musicID, playorder, trackName, artist, lyrics, soprano, alto, tenor, allParts, piano FROM playlists JOIN music ON playlists.musicID = music.musicTrackID WHERE eventID = ? ORDER BY playorder`
		//fmt.Println(sSQL, event.EventID)
		playlistRows, err := sqldb.DB.Query(sSQL, events[i].EventID)
		if err != nil {
			log.Println("Error querying database:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer playlistRows.Close()

		for playlistRows.Next() {
			var playlist strt.PlaylistEntry
			err = playlistRows.Scan(&playlist.ID, &playlist.EventID, &playlist.MusicTrack.MusicTrackID, &playlist.Playorder, &playlist.MusicTrack.TrackName, &playlist.MusicTrack.Artist, &playlist.MusicTrack.Lyrics, &playlist.MusicTrack.Soprano, &playlist.MusicTrack.Alto, &playlist.MusicTrack.Tenor, &playlist.MusicTrack.AllParts, &playlist.MusicTrack.Piano)
			if err != nil {
				log.Println("Error scanning row:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println(playlist)
			events[i].Playlist = append(events[i].Playlist, playlist)
		}
		playlistRows.Close()
	}

	// Return the combined data
	if err := json.NewEncoder(w).Encode(events); err != nil {
		log.Println("Error encoding JSON:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func EventArchive(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	vars := mux.Vars(r)

	id := vars["id"]
	var archive strt.ArchiveEntry
	var event strt.EventDetails
	sSQL := "SELECT archive.archiveID, choirevents.location, choirevents.eventDate, choirevents.title, archive.report FROM choirevents JOIN archive ON archive.eventID=choirevents.eventID WHERE choirevents.eventID = ?"
	rows, err := sqldb.DB.Query(sSQL, id)
	if err != nil {
		log.Println("Error:", err)
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
			log.Println("Error:", err)
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
		log.Println("Error:", err)
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
			log.Println("Error:", err)
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
		log.Println("Error:", err)
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
			log.Println("Error:", err)
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

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	vars := mux.Vars(r)

	id := vars["id"]
	var images []strt.ImageDetail
	var image strt.ImageDetail
	sSQL := "SELECT * FROM images WHERE eventID = ?"
	rows, err := sqldb.DB.Query(sSQL, id)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&image.ImageID, &image.Filename, &image.Caption, &image.EventID, &image.Height, &image.Width)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		images = append(images, image)
	}
	json.NewEncoder(w).Encode(images)
}

func EventClips(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	vars := mux.Vars(r)

	id := vars["id"]
	var clips []strt.Clip
	var clip strt.Clip
	sSQL := "SELECT * FROM clips WHERE eventID = ?"
	rows, err := sqldb.DB.Query(sSQL, id)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&clip.ClipID, &clip.ClipURL, &clip.EventID, &clip.Caption)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		clips = append(clips, clip)
	}
	json.NewEncoder(w).Encode(clips)
}

func ClipDelete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("ClipDelete " + id)

	sSQL := "DELETE FROM clips WHERE clipID = ?"
	_, err := sqldb.DB.Exec(sSQL, id)
	blankClip := strt.Clip{}
	if err != nil {
		log.Println("Error:", err)
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
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// Get the screen size and number of images to return
	screen := r.URL.Query().Get("screen")
	imagesStr := r.URL.Query().Get("images")
	imageReq, err := strconv.Atoi(imagesStr)
	if err != nil {
		http.Error(w, "Invalid width parameter", http.StatusBadRequest)
		return
	}
	// create an empty archive entry
	var Arch strt.ArchiveEntry
	// Determine if the device is desktop or mobile
	filepath := "/app/images/desktop"
	prefix := "dt%"
	if screen == "mobile" {
		filepath = "/app/images/mobile"
		prefix = "mb%"
	}

	// Check if the database connection is initialized
	if sqldb.DB == nil {
		log.Println("Database connection is not initialized")
		// initialize the database connection
	}

	// Query the database for random images
	sSQL := "SELECT imageID, filename, caption, eventID, height, width FROM images WHERE filename LIKE ? ORDER BY RAND() LIMIT ?"
	rows, err := sqldb.DB.Query(sSQL, prefix, imageReq)
	if err != nil {
		log.Println("Error querying database:", err)
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var image strt.ImageDetail
		err = rows.Scan(&image.ImageID, &image.Filename, &image.Caption, &image.EventID, &image.Height, &image.Width)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}

		// Adjust filename for mobile
		if filepath == "/app/images/mobile" {
			image.Filename = strings.Replace(image.Filename, "dt", "mb", 1)
		}
		image.Filename = filepath + "/" + url.PathEscape(image.Filename)

		// Confirm the file exists
		if _, err := os.Stat(image.Filename); os.IsNotExist(err) {
			log.Println("File does not exist:", image.Filename)
			image.Filename = "https://via.placeholder.com/80"
		}
		// remove the /app prefix
		image.Filename = image.Filename[4:]
		// Add the image to the array
		Arch.Images = append(Arch.Images, image)
	}

	// collect the same number of clips for this event if they exist
	sSQL = "SELECT clipID, clipURL, eventID, caption FROM clips ORDER BY RAND() LIMIT ?"
	rows, err = sqldb.DB.Query(sSQL, imageReq)
	if err != nil {
		log.Println("Error querying database:", err)
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	// Add the clips to the array
	for rows.Next() {
		var clip strt.Clip
		err = rows.Scan(&clip.ClipID, &clip.ClipURL, &clip.EventID, &clip.Caption)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}
		Arch.Clips = append(Arch.Clips, clip)
	}

	// // Fetch Instagram oEmbed data for each clip
	// for i := range Arch.Clips {
	// 	if strings.Contains(Arch.Clips[i].ClipURL, "instagram.com") {
	// 		// this is the embed text for the clip
	// 		Arch.Clips[i].ClipURL = "<blockquote class=\"instagram-media\" data-instgrm-captioned data-instgrm-permalink=\"" + Arch.Clips[i].ClipURL + "\" data-instgrm-version=\"12\"></blockquote>"
	// apiURL := "https://api.instagram.com/oembed?url=" + Arch.Clips[i].ClipURL
	// resp, err := http.Get(apiURL)
	// if err != nil {
	// 	log.Println("Error fetching Instagram oEmbed data:", err)
	// 	http.Error(w, "Error fetching Instagram oEmbed data", http.StatusInternalServerError)
	// 	return
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	http.Error(w, "Error fetching Instagram oEmbed data", resp.StatusCode)
	// 	return
	// }

	// contentType := resp.Header.Get("Content-Type")
	// if !strings.Contains(contentType, "application/json") {
	// 	log.Println("Invalid content type:", contentType, resp)
	// 	http.Error(w, "Invalid content type", http.StatusInternalServerError)
	// 	return
	// }

	// var data map[string]interface{}
	// if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
	// 	log.Println("Error decoding Instagram oEmbed response:", err)
	// 	http.Error(w, "Error decoding Instagram oEmbed response", http.StatusInternalServerError)
	// 	return
	// }

	// Add the Instagram data to the clip
	// 		if html, ok := data["html"].(string); ok {
	// 			Arch.Clips[i].ClipURL = html
	// 		}
	// 	}
	// }

	// Return the archive entry with images and clips
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(Arch); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func ThemeDetailsGET(w http.ResponseWriter, r *http.Request) {

	var theme strt.ThemeDetails

	sSQL := "SELECT * FROM themedetails"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {

		err = rows.Scan(&theme.ID, &theme.BoxColour, &theme.TextColour, &theme.TextFont, &theme.BackgroundImage, &theme.TextboxColour, &theme.LogoImage, &theme.BannerColour, &theme.MenuColour, &theme.ButtonColour, &theme.ButtonHover, &theme.ButtonTextColour, &theme.MenuTextColour, &theme.TextSize)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json.NewEncoder(w).Encode(theme)
}

func ThemeDetailsPUT(w http.ResponseWriter, r *http.Request) {

	var theme strt.ThemeDetails
	err := json.NewDecoder(r.Body).Decode(&theme)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sSQL := "UPDATE themedetails SET boxColour = ?, textColour = ?, textFont = ?, backgroundImage = ?, textboxColour = ?, logoimage = ?, bannerColour = ?, menuColour = ?, buttonColour = ?, buttonHover = ?, buttonTextColour = ?, menuTextColour = ?, textSize = ?"
	_, err = sqldb.DB.Exec(sSQL, theme.BoxColour, theme.TextColour, theme.TextFont, theme.BackgroundImage, theme.TextboxColour, theme.LogoImage, theme.BannerColour, theme.MenuColour, theme.ButtonColour, theme.ButtonHover, theme.ButtonTextColour, theme.MenuTextColour, theme.TextSize)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ThemeDetailsRandom(w http.ResponseWriter, r *http.Request) {

	var theme strt.ThemeDetails

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
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func MusicListGET(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("MusicListGET")

	var musicTracks []strt.MusicTrack
	var musicTrack strt.MusicTrack
	// open the database

	sSQL := "SELECT * FROM music"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		fmt.Println("Error in query")
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		// musicTrackID	int Auto Increment, trackName	varchar(100)	,artist	varchar(60)	,lyrics	varchar(120)	,soprano	varchar(120)	, alto	varchar(120)	,tenor	varchar(120)	,allParts	varchar(120)	,piano
		err = rows.Scan(&musicTrack.MusicTrackID, &musicTrack.TrackName, &musicTrack.Artist, &musicTrack.Lyrics, &musicTrack.Soprano, &musicTrack.Alto, &musicTrack.Tenor, &musicTrack.AllParts, &musicTrack.Piano)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error in MusicListGET")
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		musicTracks = append(musicTracks, musicTrack)
		//mt.Println(musicTracks)
	}
	json.NewEncoder(w).Encode(musicTracks)
}

func MusicTrackGET(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	var musicTrack strt.MusicTrack
	sSQL := "SELECT * FROM music WHERE musicTrackID = ?"
	rows, err := sqldb.DB.Query(sSQL, id)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&musicTrack.MusicTrackID, &musicTrack.TrackName, &musicTrack.Lyrics, &musicTrack.Soprano, &musicTrack.Alto, &musicTrack.Tenor, &musicTrack.AllParts)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json.NewEncoder(w).Encode(musicTrack)
}

func MusicTrackPOST(w http.ResponseWriter, r *http.Request) {

	var musicTrack strt.MusicTrack

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

	vars := mux.Vars(r)
	id := vars["id"]

	sSQL := "DELETE FROM music WHERE musicTrackID = ?"
	_, err := sqldb.DB.Exec(sSQL, id)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func MusicTrackPUT(w http.ResponseWriter, r *http.Request) {

	var musicTrack strt.MusicTrack

	err := json.NewDecoder(r.Body).Decode(&musicTrack)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sSQL := "UPDATE music SET trackName = ?, lyrics = ?, soprano = ?, alto = ?, tenor = ?, allParts = ? WHERE musicTrackID = ?"
	_, err = sqldb.DB.Exec(sSQL, musicTrack.TrackName, musicTrack.Lyrics, musicTrack.Soprano, musicTrack.Alto, musicTrack.Tenor, musicTrack.AllParts, musicTrack.MusicTrackID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ClipPOST(w http.ResponseWriter, r *http.Request) {

	var clip strt.Clip

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
		log.Println("Error:", err)
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
			log.Println("Error:", err)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})

			return
		}
	}
	// return the clip as a json object
	json.NewEncoder(w).Encode(clip)
}

func ArchiveEntryPUT(w http.ResponseWriter, r *http.Request) {

	var archive strt.ArchiveEntry

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

	var archive strt.ArchiveEntry

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

	var playlists []strt.PlaylistEntry
	var playlist strt.PlaylistEntry

	sSQL := "SELECT * FROM playlists"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&playlist.PlaylistID, &playlist.ID, &playlist.Playorder)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		playlists = append(playlists, playlist)
	}
	json.NewEncoder(w).Encode(playlists)
}

func PlaylistPOST(w http.ResponseWriter, r *http.Request) {

	var playlist []strt.PlaylistEntry

	err := json.NewDecoder(r.Body).Decode(&playlist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//dewlete the existing playlist for the event
	sSQL := "DELETE FROM playlists WHERE eventID = ?"
	_, err = sqldb.DB.Exec(sSQL, playlist[0].EventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// insert the new playlist
	for i := range playlist {
		sSQL := "INSERT INTO playlists (eventID, musicID, playorder) VALUES (?, ?, ?)"
		_, err = sqldb.DB.Exec(sSQL, playlist[i].EventID, playlist[i].MusicTrack.MusicTrackID, playlist[i].Playorder)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	var response strt.PlaylistEntry
	response.ID = 200
	response.MusicTrack.TrackName = "Playlist Added"
	json.NewEncoder(w).Encode(response)
}

func PlaylistDELETE(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id := vars["id"]
	sSQL := "DELETE FROM playlists WHERE playlistID = ?"
	_, err := sqldb.DB.Exec(sSQL, id)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func PlaylistPUT(w http.ResponseWriter, r *http.Request) {

	var playlist strt.PlaylistEntry

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

	vars := mux.Vars(r)
	id := vars["id"]
	//fmt.Println("PlaylistGET " + id)

	var playlists []strt.PlaylistEntry

	var sSQL = `SELECT playlistID, eventID, musicID, playorder, trackName, artist, lyrics, soprano, alto, tenor, allParts, piano 
                FROM playlists 
                JOIN music ON playlists.musicID = music.musicTrackID 
                WHERE eventID = ? 
                ORDER BY playorder`

	rows, err := sqldb.DB.Query(sSQL, id)
	if err != nil {
		log.Println("Error querying database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var playlist strt.PlaylistEntry
		err = rows.Scan(&playlist.PlaylistID, &playlist.EventID, &playlist.MusicTrack.MusicTrackID, &playlist.Playorder, &playlist.MusicTrack.TrackName, &playlist.MusicTrack.Artist, &playlist.MusicTrack.Lyrics, &playlist.MusicTrack.Soprano, &playlist.MusicTrack.Alto, &playlist.MusicTrack.Tenor, &playlist.MusicTrack.AllParts, &playlist.MusicTrack.Piano)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		playlists = append(playlists, playlist)
	}

	// Return the playlists, even if empty
	json.NewEncoder(w).Encode(playlists)
}

func EventGET(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

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

	var events []strt.EventDetails
	var event strt.EventDetails
	sSQL := "SELECT eventID, location, eventDate, title FROM choirevents" // WHERE eventDate >= curdate()"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	sSQLDate := ""
	for rows.Next() {
		err = rows.Scan(&event.EventID, &event.Location, &sSQLDate, &event.Title)
		if err != nil {
			log.Println("Error:", err)
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

	var events []strt.EventDetails
	var event strt.EventDetails
	sSQL := "SELECT * FROM choirevents WHERE eventDate >= curdate()"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	sSQLDate := ""
	for rows.Next() {
		err = rows.Scan(&event.EventID, &event.Location, &sSQLDate, &event.StartTime, &event.EndTime, &event.Price, &event.Title, &event.MeetingPoint, &event.Invitation)
		if err != nil {
			log.Println("Error:", err)
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

	var event strt.EventDetails

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//fmt.Println(event)
	sSQL := "INSERT INTO choirevents (location, eventDate, startTime, endTime, price, title, invitation, meetingPoint) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = sqldb.DB.Exec(sSQL, event.Location, event.DateString, event.StartTime, event.EndTime, event.Price, event.Title, event.Invitation, event.MeetingPoint)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	event.Title = "Event Added"
	json.NewEncoder(w).Encode(event)
}

func EventDELETE(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	sSQL := "DELETE FROM choirevents WHERE eventID = ?"
	_, err := sqldb.DB.Exec(sSQL, id)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func EventPUT(w http.ResponseWriter, r *http.Request) {

	var event strt.EventDetails

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

	var image strt.ImageDetail
	vars := mux.Vars(r)
	id := vars["id"]

	sSQL := "SELECT * FROM images WHERE imageID = ?"
	rows, err := sqldb.DB.Query(sSQL, id)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&image.ImageID, &image.Filename, &image.Caption, &image.EventID, &image.Height, &image.Width)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json.NewEncoder(w).Encode(image)
}

func ImagesGET(w http.ResponseWriter, r *http.Request) {

	var images []strt.ImageDetail
	var image strt.ImageDetail

	sSQL := "SELECT * FROM images"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&image.ImageID, &image.Filename, &image.Caption, &image.EventID, &image.Height, &image.Width)
		if err != nil {
			log.Println("Error:", err)
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
	fmt.Println("Content-Type:", r.Header.Get("Content-Type"))

	// Parse the multipart form to 50MB
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

	// Ensure the temporary directory exists
	tempDir := "/root/temp-images"
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		err = os.MkdirAll(tempDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating temporary directory:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	fmt.Println("Filename:", handler.Filename)
	fmt.Println("Size:", handler.Size)
	fmt.Println("Header:", handler.Header)

	// Create a temporary file
	tempFile, err := os.CreateTemp(tempDir, "temp-*.png")
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
		destinationDir = "/app/images/background/"
	case strings.HasPrefix(handler.Filename, "dt"):
		destinationDir = "/app/images/desktop/"
	case strings.HasPrefix(handler.Filename, "mb"):
		destinationDir = "/app/images/mobile/"
	default:
		destinationDir = "/app/images/"
	}

	// Move the file to the destination directory
	//err = os.Rename(tempFile.Name(), destinationDir+handler.Filename)
	err = os.WriteFile(filepath.Join(destinationDir, handler.Filename), fileBytes, 0666)
	// if err != nil {
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

	// If the image has an mb prefix then it doesn't need to go in the database
	// so return a JSON encoded response
	if strings.HasPrefix(handler.Filename, "mb") {
		json.NewEncoder(w).Encode(image)
		return
	}

	// Insert the file details into the database
	sSQL := "INSERT INTO images (filename, caption, eventID, height, width) VALUES (?, ?, ?, ?, ?)"
	_, err = sqldb.DB.Exec(sSQL, image.Filename, image.Caption, image.EventID, image.Height, image.Width)
	if err != nil {
		fmt.Println("Error inserting into database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve the imageID of the newly inserted record
	sSQL = "SELECT imageID FROM images ORDER BY imageID DESC LIMIT 1"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		fmt.Println("Error querying database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type imageID struct {
		ImageID int
	}
	Im := imageID{}
	if rows.Next() {
		err = rows.Scan(&Im.ImageID)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Respond with the image details
	json.NewEncoder(w).Encode(Im)
}

func ImagesPOST(w http.ResponseWriter, r *http.Request) {

	var images []strt.ImageDetail

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
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&image.ImageID)
			if err != nil {
				log.Println("Error:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	w.WriteHeader(http.StatusOK)

	//	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:")
	//  return the image object
	json.NewEncoder(w).Encode(images)

}

func ImageDELETE(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	// we need to get the filename so we can delete the files
	sSQL := "SELECT filename FROM images WHERE imageID = ?"
	rows, err := sqldb.DB.Query(sSQL, id)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var filename string
	for rows.Next() {
		err = rows.Scan(&filename)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// delete the file. If the prefix is bg, dt or mb then we need to delete from the appropriate directory
	if strings.HasPrefix(filename, "bg") {
		err = os.Remove("images/background/" + filename)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if strings.HasPrefix(filename, "dt") {
		err = os.Remove("Images/Desktop/" + filename)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// then remove the prefix and replace it with mb
		mbfile := strings.Replace(filename, "dt", "mb", 1)
		err = os.Remove("Images/mobile/" + mbfile)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if strings.HasPrefix(filename, "lg") {
		err = os.Remove("Images/Logo/" + filename)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// now delete the record from the database

	sSQL = "DELETE FROM images WHERE imageID = ?"
	_, err = sqldb.DB.Exec(sSQL, id)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Image deleted successfully"})
}

func ImagePUT(w http.ResponseWriter, r *http.Request) {

	var image strt.ImageDetail

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

	var siteInfo strt.SiteInfo

	sSQL := "SELECT * FROM siteinfo"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Println("Error:", err)
		fmt.Println("Error in SiteInfoGET")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&siteInfo.ID, &siteInfo.HomeTitle, &siteInfo.HomeText, &siteInfo.AboutTitle, &siteInfo.AboutText, &siteInfo.ArchiveTitle, &siteInfo.ArchiveText, &siteInfo.NoticesTitle, &siteInfo.NoticesText, &siteInfo.BookingTitle, &siteInfo.BookingText, &siteInfo.MembersTitle, &siteInfo.MembersText, &siteInfo.AppealTitle, &siteInfo.AppealText, &siteInfo.SettingsTitle, &siteInfo.SettingsText)
		if err != nil {

			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json.NewEncoder(w).Encode(siteInfo)
}

func SiteInfoPUT(w http.ResponseWriter, r *http.Request) {

	var siteInfo strt.SiteInfo

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
	//return the object
	json.NewEncoder(w).Encode(siteInfo)
}

func InstagramEmbed(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	instagramURL := query.Get("url")
	if instagramURL == "" {
		http.Error(w, "Missing URL parameter", http.StatusBadRequest)
		return
	}

	apiURL := "https://api.instagram.com/oembed?url=" + url.QueryEscape(instagramURL)
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Println("Error fetching Instagram oEmbed data:", err)
		http.Error(w, "Error fetching Instagram oEmbed data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error fetching Instagram oEmbed data", resp.StatusCode)
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println("Error decoding Instagram oEmbed response:", err)
		http.Error(w, "Error decoding Instagram oEmbed response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
