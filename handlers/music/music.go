package music

import (
	"RWTAPI/sqldb"
	strt "RWTAPI/structures"
	tools "RWTAPI/tools"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// GET upcoming playlists
func UpcomingPlaylistsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
		parsedDate, err := tools.ParseMySQLDateTime(event.EventDate)
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
			//fmt.Println(playlist)
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

func MusicListGET(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("musicListGET")
	w.Header().Set("Content-Type", "application/json")

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
			fmt.Println("Error in musicListGET")
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
	w.Header().Set("Content-Type", "application/json")

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

		err = rows.Scan(&musicTrack.MusicTrackID, &musicTrack.TrackName, &musicTrack.Artist, &musicTrack.Lyrics, &musicTrack.Soprano, &musicTrack.Alto, &musicTrack.Tenor, &musicTrack.AllParts, &musicTrack.Piano)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json.NewEncoder(w).Encode(musicTrack)
}

func MusicTrackPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var musicTrack strt.MusicTrack

	err := json.NewDecoder(r.Body).Decode(&musicTrack)
	if err != nil {
		fmt.Println("Error in MusicTrackPOST", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("musickTrack: ", musicTrack)
	//musicTrackID, trackName, artist, lyrics, soprano, alto, tenor, allParts, piano
	sSQL := "INSERT INTO music (trackName, artist, lyrics, soprano, alto, tenor, allParts, piano) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = sqldb.DB.Exec(sSQL, musicTrack.TrackName, musicTrack.Artist, musicTrack.Lyrics, musicTrack.Soprano, musicTrack.Alto, musicTrack.Tenor, musicTrack.AllParts, musicTrack.Piano)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error in MusicTrackPOST", err)
		return
	}
	// retieve the musicTrackID
	sSQL = "SELECT musicTrackID FROM music ORDER BY musicTrackID DESC LIMIT 1"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&musicTrack.MusicTrackID)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	//  add the id to the musictrack object and return to the client
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(musicTrack)

}

func MusicTrackDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("\n=== DELETE TRACK START: %s ===\n", id)

	// Check if record exists first
	var exists int
	err := sqldb.DB.QueryRow("SELECT COUNT(*) FROM rwtchoir.music WHERE musicTrackID = ?", id).Scan(&exists)
	if err != nil {
		log.Printf("[ERROR] Checking existence: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if exists == 0 {
		log.Printf("[ERROR] Track %s not found", id)
		http.Error(w, "Track not found", http.StatusNotFound)
		return
	}

	// Start transaction
	tx, err := sqldb.DB.Begin()
	if err != nil {
		log.Printf("[ERROR] Starting transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete playlists first
	_, err = tx.Exec("DELETE FROM rwtchoir.playlists WHERE musicID = ?", id)
	if err != nil {
		tx.Rollback()
		log.Printf("[ERROR] Deleting playlists: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the music track
	result, err := tx.Exec("DELETE FROM rwtchoir.music WHERE musicTrackID = ?", id)
	if err != nil {
		tx.Rollback()
		log.Printf("[ERROR] Deleting track: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		tx.Rollback()
		log.Printf("[ERROR] No rows affected deleting track %s", id)
		http.Error(w, "Track not found", http.StatusNotFound)
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("[ERROR] Committing transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("=== DELETE SUCCESS: Track %s deleted ===\n", id)
	w.WriteHeader(http.StatusNoContent)
}

func MusicTrackPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
	w.WriteHeader(http.StatusOK)
	//return the response
}

func PlaylistDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	id := vars["id"]
	sSQL := "DELETE FROM playlists WHERE eventID = ?"
	_, err := sqldb.DB.Exec(sSQL, id)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func PlaylistsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
	w.Header().Set("Content-Type", "application/json")

	var playlist []strt.PlaylistEntry

	err := json.NewDecoder(r.Body).Decode(&playlist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(playlist)
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

func PlaylistPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
	w.Header().Set("Content-Type", "application/json")

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
