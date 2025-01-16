package events

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

func EventGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
			parsedDate, err := tools.ParseMySQLDateTime(sSQLDate)
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

	var events []strt.EventDetails
	var event strt.EventDetails
	sSQL := "SELECT eventID, location, eventDate, title FROM choirevents WHERE eventDate < curdate()"
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
		parsedDate, err := tools.ParseMySQLDateTime(sSQLDate)
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
		parsedDate, err := tools.ParseMySQLDateTime(sSQLDate)
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
	w.Header().Set("Content-Type", "application/json")

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
	w.Header().Set("Content-Type", "application/json")

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

func UpcomingEventsListsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
		parsedDate, err := tools.ParseMySQLDateTime(event.EventDate)
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

func InstagramEmbed(w http.ResponseWriter, r *http.Request) {
	// ...existing InstagramEmbed code...
}
