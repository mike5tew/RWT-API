package team

import (
	"RWTAPI/sqldb"
	strt "RWTAPI/structures"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

// ensureImageDirectories ensures that all required directories exist
func ensureImageDirectories() error {
	basePath := "/app/images/"
	dirs := []string{
		basePath,
		filepath.Join(basePath, "desktop"),
		filepath.Join(basePath, "mobile"),
		filepath.Join(basePath, "background"),
		filepath.Join(basePath, "logo"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

// GetTeam retrieves the team from the database
func GetTeam(w http.ResponseWriter, r *http.Request) {
	rows, err := sqldb.DB.Query("SELECT team.teamID, team.name, team.description, team.imageID, images.filename, images.caption FROM rwtchoir.team LEFT OUTER JOIN rwtchoir.images ON team.imageID = images.imageID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var team []strt.TeamMember
	for rows.Next() {
		var member strt.TeamMember
		var teamImageID sql.NullInt64
		var imageFilename sql.NullString
		var imageCaption sql.NullString

		if err := rows.Scan(&member.ID, &member.Name, &member.Description, &teamImageID, &imageFilename, &imageCaption); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		member.Image = strt.ImageDetail{}
		if teamImageID.Valid {
			member.Image.ImageID = int(teamImageID.Int64)
		}
		if imageFilename.Valid {
			member.Image.Filename = imageFilename.String
		}
		if imageCaption.Valid {
			member.Image.Caption = imageCaption.String
		}

		team = append(team, member)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(team); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AddMember adds a member to a team
func AddMember(w http.ResponseWriter, r *http.Request) {
	if err := ensureImageDirectories(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var member strt.TeamMember
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Adding member:", member)

	query := "INSERT INTO rwtchoir.team (name, description, imageID) VALUES (?, ?, ?)"
	result, err := sqldb.DB.Exec(query, member.Name, member.Description, member.Image.ImageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	member.ID = int(lastID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

// RemoveMember removes a member from a team
func RemoveMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var imageID sql.NullInt64
	var filename sql.NullString
	query := `SELECT t.imageID, i.filename 
			FROM rwtchoir.team t 
			LEFT JOIN rwtchoir.images i ON t.imageID = i.imageID 
			WHERE t.teamID = ?`
	err := sqldb.DB.QueryRow(query, id).Scan(&imageID, &filename)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Member not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	query = "DELETE FROM rwtchoir.team WHERE teamID = ?"
	_, err = sqldb.DB.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if imageID.Valid && filename.Valid {
		if err := DeleteMemberImage(int(imageID.Int64), filename.String); err != nil {
			log.Printf("Warning: Failed to delete associated image: %v", err)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteMemberImage deletes member image files
func DeleteMemberImage(ID int, filename string) error {
	query := "DELETE FROM rwtchoir.images WHERE imageID = ?"
	_, err := sqldb.DB.Exec(query, ID)
	if err != nil {
		return err
	}

	basePath := "/app/images/"

	// Updated to remove prefixes - just use filename directly
	dtFilePath := filepath.Join(basePath, "desktop", filename)
	if err := os.Remove(dtFilePath); err != nil {
		log.Printf("File operation error: %v, Path: %s", err, dtFilePath)
	}

	mbFilePath := filepath.Join(basePath, "mobile", filename)
	if err := os.Remove(mbFilePath); err != nil {
		log.Printf("File operation error: %v, Path: %s", err, mbFilePath)
	}

	// Also try with the old prefixes for backward compatibility
	oldDtFilePath := filepath.Join(basePath, "desktop", "dt"+filename)
	if err := os.Remove(oldDtFilePath); err != nil {
		// Ignore error as this is just for backward compatibility
		log.Printf("Old format file not found: %s", oldDtFilePath)
	}

	oldMbFilePath := filepath.Join(basePath, "mobile", "mb"+filename)
	if err := os.Remove(oldMbFilePath); err != nil {
		// Ignore error as this is just for backward compatibility
		log.Printf("Old format file not found: %s", oldMbFilePath)
	}

	return nil
}

// GetMember gets the member of a team
func GetMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var member strt.TeamMember
	query := "SELECT teamID, name, description, imageID FROM rwtchoir.team WHERE teamID = ?"
	err := sqldb.DB.QueryRow(query, id).Scan(&member.ID, &member.Name, &member.Description, &member.Image.ImageID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Member not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

// MemberPUT updates a member of a team
func MemberPUT(w http.ResponseWriter, r *http.Request) {
	if err := ensureImageDirectories(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var member strt.TeamMember
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var oldImageID sql.NullInt64
	var oldFilename sql.NullString
	checkQuery := `
		SELECT t.imageID, i.filename 
		FROM rwtchoir.team t 
		LEFT JOIN rwtchoir.images i ON t.imageID = i.imageID 
		WHERE t.teamID = ?`
	err := sqldb.DB.QueryRow(checkQuery, member.ID).Scan(&oldImageID, &oldFilename)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Member not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error fetching old image ID: %v", err), http.StatusInternalServerError)
		}
		return
	}

	if oldImageID.Valid && oldImageID.Int64 != int64(member.Image.ImageID) {
		if oldFilename.Valid {
			if err := DeleteMemberImage(int(oldImageID.Int64), oldFilename.String); err != nil {
				log.Printf("Warning: Failed to delete old image %d (%s): %v", oldImageID.Int64, oldFilename.String, err)
			}
		} else {
			log.Printf("Warning: Old image ID %d found for member %d, but filename was NULL. Cannot delete files.", oldImageID.Int64, member.ID)
			query := "DELETE FROM rwtchoir.images WHERE imageID = ?"
			_, delErr := sqldb.DB.Exec(query, oldImageID.Int64)
			if delErr != nil {
				log.Printf("Warning: Failed to delete old image DB record for ID %d: %v", oldImageID.Int64, delErr)
			}
		}
	}

	var newImageID sql.NullInt64
	if member.Image.ImageID > 0 {
		newImageID = sql.NullInt64{Int64: int64(member.Image.ImageID), Valid: true}
	} else {
		newImageID = sql.NullInt64{Valid: false}
	}

	updateQuery := "UPDATE rwtchoir.team SET name = ?, description = ?, imageID = ? WHERE teamID = ?"
	_, err = sqldb.DB.Exec(updateQuery, member.Name, member.Description, newImageID, member.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating team member: %v", err), http.StatusInternalServerError)
		return
	}

	if !newImageID.Valid {
		member.Image = strt.ImageDetail{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

// MemberPOST adds a member to a team
func MemberPOST(w http.ResponseWriter, r *http.Request) {
	if err := ensureImageDirectories(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var member strt.TeamMember
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "INSERT INTO rwtchoir.team (name, description, imageID) VALUES (?, ?, ?)"
	result, err := sqldb.DB.Exec(query, member.Name, member.Description, member.Image.ImageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	member.ID = int(lastID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}
