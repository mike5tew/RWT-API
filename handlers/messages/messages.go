package messages

import (
	"RWTAPI/sqldb"
	strt "RWTAPI/structures"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

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
		//messageID, messageDate, messageFrom, messageContent, eventName, eventDate, eventTime, contactEmail, contactPhone, eventLocation
		err = rows.Scan(&message.MessageID, &message.MessageDate, &message.MessageFrom, &message.MessageContent, &message.EventName, &message.EventDate, &message.EventTime, &message.ContactEmail, &message.ContactPhone, &message.EventLocation)
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
	//extract the messageID from the URL which is just a number at the end of the URL
	url := r.URL.Path
	// split the URL into parts using the / character
	parts := strings.Split(url, "/")
	// get the last part of the URL which is the messageID
	id := parts[len(parts)-1]

	sSQL := "DELETE FROM messages WHERE messageID = ?"
	_, err := sqldb.DB.Exec(sSQL, id)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// return a json response
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
