package site

import (
	"RWTAPI/sqldb"
	strt "RWTAPI/structures"
	"encoding/json"
	"log"
	"net/http"
)

func SiteInfoGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var siteInfo strt.SiteInfo
	sSQL := "SELECT * FROM siteinfo"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		//id, HomeTitle, HomeText, AboutTitle, AboutText, ArchiveTitle, ArchiveText, NoticesTitle, NoticesText, BookingTitle, BookingText, MembersTitle, MembersText, AppealTitle, AppealText, SettingsTitle, SettingsText
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
	w.Header().Set("Content-Type", "application/json")

	var siteInfo strt.SiteInfo
	err := json.NewDecoder(r.Body).Decode(&siteInfo)

	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sSQL := "UPDATE siteinfo SET HomeTitle = ?, HomeText = ?, AboutTitle = ?, AboutText = ?, ArchiveTitle = ?, ArchiveText = ?, NoticesTitle = ?, NoticesText = ?, BookingTitle = ?, BookingText = ?, MembersTitle = ?, MembersText = ?, AppealTitle = ?, AppealText = ?, SettingsTitle = ?, SettingsText = ? WHERE id = ?"
	_, err = sqldb.DB.Exec(sSQL, siteInfo.HomeTitle, siteInfo.HomeText, siteInfo.AboutTitle, siteInfo.AboutText, siteInfo.ArchiveTitle, siteInfo.ArchiveText, siteInfo.NoticesTitle, siteInfo.NoticesText, siteInfo.BookingTitle, siteInfo.BookingText, siteInfo.MembersTitle, siteInfo.MembersText, siteInfo.AppealTitle, siteInfo.AppealText, siteInfo.SettingsTitle, siteInfo.SettingsText, siteInfo.ID)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(siteInfo)

}
