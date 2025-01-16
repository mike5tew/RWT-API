package theme

import (
	"RWTAPI/sqldb"
	strt "RWTAPI/structures"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/exp/rand"
)

func ThemeDetailsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
	//check the background image file is there
	urlversion := url.PathEscape(theme.BackgroundImage)
	localfile2 := "/app/images/background/" + theme.BackgroundImage
	if _, err := os.Stat(localfile2); os.IsNotExist(err) {
		// File does not exist
		fmt.Println("File2 does not exist" + localfile2)
		theme.BackgroundImage = ""
	} else {
		// File exists, adjust the path for the frontend
		theme.BackgroundImage = "images/background/" + urlversion
	}
	json.NewEncoder(w).Encode(theme)
}

func ThemeDetailsPUT(w http.ResponseWriter, r *http.Request) {
	var theme strt.ThemeDetails
	err := json.NewDecoder(r.Body).Decode(&theme)
	if err != nil {
		log.Printf("Error decoding theme: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate TextSize
	if theme.TextSize <= 0 {
		theme.TextSize = 12 // Set default if invalid
	}

	sSQL := "UPDATE themedetails SET boxColour = ?, textColour = ?, textFont = ?, backgroundImage = ?, textboxColour = ?, logoimage = ?, bannerColour = ?, menuColour = ?, buttonColour = ?, buttonHover = ?, buttonTextColour = ?, menuTextColour = ?, textSize = ?"
	_, err = sqldb.DB.Exec(sSQL,
		theme.BoxColour,
		theme.TextColour,
		theme.TextFont,
		theme.BackgroundImage,
		theme.TextboxColour,
		theme.LogoImage,
		theme.BannerColour,
		theme.MenuColour,
		theme.ButtonColour,
		theme.ButtonHover,
		theme.ButtonTextColour,
		theme.MenuTextColour,
		theme.TextSize)

	if err != nil {
		log.Printf("Error updating theme: %v", err)
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
	theme.BackgroundImage = "Musical Background.png"
	theme.TextboxColour = "#" + fmt.Sprintf("%06x", rand.Intn(0xffffff))
	theme.LogoImage = "Choir Logo.png"
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
