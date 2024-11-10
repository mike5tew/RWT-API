package structures

import (
	"database/sql"
	"mime/multipart"
)

// the golang equivalent of a FormData object would be a struct with the fields you need
type FormData struct {
	File *multipart.FileHeader
}

// type tempResponse struct {
// 	fData      FormData
// 	uploadType string
// }

type ImageDetail struct {
	ImageID   int    `json:"ImageID"`
	ImageURL  string `json:"ImageURL"`
	Filename  string `json:"Filename"`
	Caption   string `json:"Caption"`
	Rows      int    `json:"Rows"`
	Cols      int    `json:"Cols"`
	Height    int    `json:"Height"`
	Width     int    `json:"Width"`
	EventID   int    `json:"EventID"`
	Imagetype int    `json:"Imagetype"`
}

type DatURLResponse struct {
	ReturnedFile sql.NullString `json:"ReturnedFile"`
	FileDetails  ImageDetail    `json:"FileDetails"`
}

type URLdetails struct {
	URL string `json:"URL"`
}

type EventDetails struct {
	EventID      int             `json:"EventID"`
	Location     string          `json:"Location"`
	EventDate    string          `json:"EventDate"`
	StartTime    string          `json:"StartTime"`
	EndTime      string          `json:"EndTime"`
	Invitation   string          `json:"Invitation"`
	MeetingPoint string          `json:"MeetingPoint"`
	Price        string          `json:"Price"`
	Title        string          `json:"Title"`
	Playlist     []PlaylistEntry `json:"Playlist"`
}

type ThemeDetails struct {
	ID               int            `json:"ID"`
	BoxColour        string         `json:"BoxColour"`
	TextColour       string         `json:"TextColour"`
	TextFont         string         `json:"TextFont"`
	BackgroundImage  sql.NullString `json:"BackgroundImage"`
	TextboxColour    string         `json:"TextboxColour"`
	LogoImage        sql.NullString `json:"LogoImage"`
	BannerColour     string         `json:"BannerColour"`
	MenuColour       string         `json:"MenuColour"`
	ButtonColour     string         `json:"ButtonColour"`
	ButtonHover      string         `json:"ButtonHover"`
	ButtonTextColour string         `json:"ButtonTextColour"`
	MenuTextColour   string         `json:"MenuTextColour"`
	TextSize         int            `json:"TextSize"`
}

type PlaylistEntry struct {
	ID         int        `json:"ID"`
	PlaylistID int        `json:"PlaylistID"`
	EventID    int        `json:"EventID"`
	MusicTrack MusicTrack `json:"MusicTrack"`
	Playorder  int        `json:"Playorder"`
}

type Clip struct {
	ClipID  int    `json:"ClipID"`
	ClipURL string `json:"ClipURL"`
	EventID int    `json:"EventID"`
	Caption string `json:"Caption"`
}

type ArchiveEntry struct {
	ArchiveID    int           `json:"ArchiveID"`
	NextFile     string        `json:"NextFile"`
	Imagecaption string        `json:"Imagecaption"`
	NextURL      string        `json:"NextURL"`
	Clipcaption  string        `json:"Clipcaption"`
	EventDetails EventDetails  `json:"EventDetails"`
	Report       string        `json:"Report"`
	Images       []ImageDetail `json:"Images"`
	Clips        []Clip        `json:"Clips"`
}

type Message struct {
	MessageID      int    `json:"MessageID"`
	MessageDate    string `json:"MessageDate"`
	MessageFrom    string `json:"MessageFrom"`
	MessageContent string `json:"MessageContent"`
}

type User struct {
	UserID   int    `json:"UserID"`
	Username string `json:"Username"`
	Password string `json:"Password"`
	Role     string `json:"Role"`
}

type ImageFiles struct {
	MainImage   sql.NullString `json:"MainImage"`
	MobileImage sql.NullString `json:"MobileImage"`
	EventID     int            `json:"EventID"`
}

type SiteInfo struct {
	ID            int    `json:"ID"`
	HomeTitle     string `json:"HomeTitle"`
	HomeText      string `json:"HomeText"`
	AboutTitle    string `json:"AboutTitle"`
	AboutText     string `json:"AboutText"`
	ArchiveTitle  string `json:"ArchiveTitle"`
	ArchiveText   string `json:"ArchiveText"`
	NoticesTitle  string `json:"NoticesTitle"`
	NoticesText   string `json:"NoticesText"`
	BookingTitle  string `json:"BookingTitle"`
	BookingText   string `json:"BookingText"`
	MembersTitle  string `json:"MembersTitle"`
	MembersText   string `json:"MembersText"`
	AppealTitle   string `json:"AppealTitle"`
	AppealText    string `json:"AppealText"`
	SettingsTitle string `json:"SettingsTitle"`
	SettingsText  string `json:"SettingsText"`
}

type MusicTrack struct {
	ID        int    `json:"ID"`
	TrackName string `json:"TrackName"`
	Lyrics    string `json:"Lyrics"`
	Artist    string `json:"Artist"`
	Soprano   string `json:"Soprano"`
	Alto      string `json:"Alto"`
	Tenor     string `json:"Tenor"`
	AllParts  string `json:"AllParts"`
	Piano     string `json:"Piano"`
}

//	export interface ScreenSize {
//	    width: number;
//	    height: number;
//	    devicePixelRatio: number;
//	    images: number
//	  }
type ScreenSize struct {
	Width            int `json:"Width"`
	Height           int `json:"Height"`
	DevicePixelRatio int `json:"DevicePixelRatio"`
	Images           int `json:"Images"`
}
