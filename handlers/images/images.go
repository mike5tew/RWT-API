package images

import (
	"RWTAPI/sqldb"
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

	"github.com/gorilla/mux"
)

func FileDetailsPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")

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

func EventImages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

func RandomImagesGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	fmt.Println("RandomImagesGET")
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
		// Adjust the filename for the frontend
		fLocal := filepath + "/" + image.Filename
		image.Filename = filepath + "/" + url.PathEscape(image.Filename)
		fmt.Println("filename ", image.Filename)
		// Confirm the file exists
		if _, err := os.Stat(fLocal); os.IsNotExist(err) {
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

func ImageBackGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var images []strt.ImageDetail
	// we are retieving the background and logo images
	sSQL := "SELECT * FROM images WHERE eventID < 1"
	rows, err := sqldb.DB.Query(sSQL)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var image strt.ImageDetail
		err = rows.Scan(&image.ImageID, &image.Filename, &image.Caption, &image.EventID, &image.Height, &image.Width)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		urlversion := url.PathEscape(image.Filename)

		// now we need to retieve the image file
		// check the filename prefix
		// if the filename starts with bg then it is a background image
		if image.EventID == 0 {
			// check if the file exists
			localfile := "/app/images/background/" + image.Filename
			if _, err := os.Stat(localfile); os.IsNotExist(err) {
				// File does not exist
				fmt.Println("File does not exist" + localfile)
				image.Filename = ""
			} else {
				// File exists, adjust the path for the frontend
				image.ImageURL = "images/background/" + urlversion
			}
			images = append(images, image)

		} else {
			// check if the file exists
			localfile := "/app/images/logo/" + image.Filename
			if _, err := os.Stat(localfile); os.IsNotExist(err) {
				// File does not exist
				fmt.Println("File does not exist" + localfile)
				image.Filename = ""
			} else {
				// File exists, adjust the path for the frontend
				image.ImageURL = "images/logo/" + urlversion
			}
			images = append(images, image)
		}
	}
	json.NewEncoder(w).Encode(images)
}

func ImagesGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
	//fmt.Println("Content-Type:", r.Header.Get("Content-Type"))
	w.Header().Set("Content-Type", "application/json")

	// Parse the multipart form to 50MB
	err := r.ParseMultipartForm(50 << 20)
	if err != nil {
		//fmt.Println("Error parsing multipart form:", err)
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
	case strings.HasPrefix(handler.Filename, "dt"):
		destinationDir = "/app/images/desktop/"
	case strings.HasPrefix(handler.Filename, "mb"):
		destinationDir = "/app/images/mobile/"
	default:
		destinationDir = "/app/images/"
	}

	// Initialize the image details
	var image strt.ImageDetail
	image.Filename = handler.Filename

	// Get additional file details from the form
	caption := r.FormValue("caption")
	eventID := r.FormValue("eventID")
	// is event id a number or q string
	var evid int
	evid, err = strconv.Atoi(eventID)
	if err != nil {
		fmt.Println("Error converting eventID:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	height := r.FormValue("height")
	width := r.FormValue("width")
	//fmt.Println("Destination directory:", destinationDir)
	//fmt.Println("Event ID:", eventID)
	if destinationDir == "/app/images/" {
		// Check if the file is a background image
		if evid == 0 {
			fmt.Println("Background image")
			destinationDir = "/app/images/background/"
		}
		if evid == -1 {
			fmt.Println("Logo image")
			destinationDir = "/app/images/logo/"
		}
	}

	//create the folder if it doesn't exist
	if _, err := os.Stat(destinationDir); os.IsNotExist(err) {
		// create the folder
		err = os.Mkdir(destinationDir, 0755)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
	fmt.Println("Image details:", image)
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

	var ImgRet strt.ImageDetail
	ImgRet.Filename = image.Filename
	ImgRet.Caption = image.Caption
	ImgRet.EventID = image.EventID
	ImgRet.ImageURL = strings.Replace(destinationDir, "/app", "", 1) + url.PathEscape(image.Filename)
	fmt.Println("ImageURL:", ImgRet.ImageURL)
	if rows.Next() {
		err = rows.Scan(&ImgRet.ImageID)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// Set the Content-Type header
	w.Header().Set("Content-Type", "application/json")
	// Return the image details as JSON
	err = json.NewEncoder(w).Encode(ImgRet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ImagesPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

	//	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:")
	//  return the image object
	json.NewEncoder(w).Encode(images)

}

func ImageDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//extract the id from the URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get image details before deletion
	var image strt.ImageDetail
	if err := sqldb.DB.QueryRow("SELECT filename FROM images WHERE imageID = ?", id).Scan(&image.Filename); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete from database
	_, err := sqldb.DB.Exec("DELETE FROM images WHERE imageID = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete files from disk
	basePath := "/app/images/"
	// Delete main image
	if err := os.Remove(basePath + image.Filename); err != nil {
		log.Printf("Error deleting main image: %v", err)
	}
	// Delete mobile version
	if err := os.Remove(basePath + "mb" + image.Filename); err != nil {
		log.Printf("Error deleting mobile image: %v", err)
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Image deleted successfully"})
}

func ImagePUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
