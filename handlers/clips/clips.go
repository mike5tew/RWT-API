package clips

import (
	"RWTAPI/sqldb"
	strt "RWTAPI/structures"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func EventClips(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]
	//fmt.Println("ClipDelete " + id)

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

func ClipPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var clip strt.Clip

	err := json.NewDecoder(r.Body).Decode(&clip)
	if err != nil {
		clip.ClipID = -1
		clip.ClipURL = "Error: " + err.Error()
		json.NewEncoder(w).Encode(clip)
		return
	}
	//fmt.Println("ClipPOST" + clip.ClipURL)
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
