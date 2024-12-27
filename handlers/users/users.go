package users

import (
	"RWTAPI/sqldb"
	strt "RWTAPI/structures"
	"encoding/json"
	"log"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {

	//
	var user strt.User

	//fmt.Println("Login")
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
