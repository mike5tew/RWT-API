package archive

import (
	"RWTAPI/sqldb"
	strt "RWTAPI/structures"
	tools "RWTAPI/tools"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func ArchiveGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
		parsedDate, err := tools.ParseMySQLDateTime(SQLDate)
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
	//fmt.Println("ArchiveGET")
	json.NewEncoder(w).Encode(archive)
}

// GET a given number of archive records
func ArchivesGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Initialize the database connection

	// Get the screen size and number of images to return
	imagesStr := r.URL.Query().Get("archives")

	iRecords, err := strconv.Atoi(imagesStr)
	if err != nil {
		http.Error(w, "Invalid archives parameter", http.StatusBadRequest)
		return
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
		parsedDate, err := tools.ParseMySQLDateTime(sSQLDate)
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
		// if filepath == "images/mobile" {
		// 	image.Filename = strings.Replace(image.Filename, "dt", "mb", 1)
		// }
		// fmt.Println(image.Filename)
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
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]
	// we need the eventID to delete the images and clips
	sSQL := "SELECT eventID FROM archive WHERE archiveID = ?"
	rows, err := sqldb.DB.Query(sSQL, id)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var eventID int
	for rows.Next() {
		err = rows.Scan(&eventID)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// delete the archive entry
	sSQL = "DELETE FROM archive WHERE archiveID = ?"
	_, err = sqldb.DB.Exec(sSQL, id)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// delete the images that are associated with the event
	sSQL = "DELETE FROM images WHERE eventID = ?"
	_, err = sqldb.DB.Exec(sSQL, id)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// delete the clips that are associated with the event
	sSQL = "DELETE FROM clips WHERE eventID = ?"
	_, err = sqldb.DB.Exec(sSQL, id)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func EventArchive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
		parsedDate, err := tools.ParseMySQLDateTime(sSQLDate)
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

func ArchiveEntryPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
	w.Header().Set("Content-Type", "application/json")

	var archive strt.ArchiveEntry

	err := json.NewDecoder(r.Body).Decode(&archive)

	if err != nil {
		archive.ArchiveID = -1
		archive.Report = err.Error()
		//fmt.Println("ArchiveEntryPOST 1 " + archive.Report)
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
