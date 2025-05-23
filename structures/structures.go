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
	ImageID   int    `json:"imageID"`
	ImageURL  string `json:"imageURL"`
	Filename  string `json:"filename"`
	Caption   string `json:"caption"`
	Rows      int    `json:"rows"`
	Cols      int    `json:"cols"`
	Height    int    `json:"height"`
	Width     int    `json:"width"`
	EventID   int    `json:"eventID"`
	Imagetype string `json:"imagetype"`
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
	DateString   string          `json:"DateString"`
	StartTime    string          `json:"StartTime"`
	EndTime      string          `json:"EndTime"`
	Invitation   string          `json:"Invitation"`
	MeetingPoint string          `json:"MeetingPoint"`
	Price        string          `json:"Price"`
	Title        string          `json:"Title"`
	Playlist     []PlaylistEntry `json:"Playlist"`
}

type ThemeDetails struct {
	ID               int    `json:"ID"`
	BoxColour        string `json:"BoxColour"`
	TextColour       string `json:"TextColour"`
	TextFont         string `json:"TextFont"`
	BackgroundImage  string `json:"BackgroundImage"`
	TextboxColour    string `json:"TextboxColour"`
	LogoImage        string `json:"LogoImage"`
	BannerColour     string `json:"BannerColour"`
	MenuColour       string `json:"MenuColour"`
	ButtonColour     string `json:"ButtonColour"`
	ButtonHover      string `json:"ButtonHover"`
	ButtonTextColour string `json:"ButtonTextColour"`
	MenuTextColour   string `json:"MenuTextColour"`
	TextSize         int    `json:"TextSize"`
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

// messageID	int Auto Increment
// messageDate	timestamp NULL
// messageFrom	varchar(60) NULL
// messageContent	varchar(500) NULL ['']
// eventName	varchar(100) ['']
// eventDate	timestamp
// eventTime	varchar(25) ['']
// contactEmail	varchar(100) NULL
// contactPhone	varchar(20) ['']
// eventLocation

type Message struct {
	MessageID      int    `json:"MessageID"`
	MessageDate    string `json:"MessageDate"`
	MessageFrom    string `json:"MessageFrom"`
	MessageContent string `json:"MessageContent"`
	EventName      string `json:"EventName"`
	EventDate      string `json:"EventDate"`
	EventTime      string `json:"EventTime"`
	ContactEmail   string `json:"ContactEmail"`
	ContactPhone   string `json:"ContactPhone"`
	EventLocation  string `json:"EventLocation"`
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
	MusicTrackID int    `json:"MusicTrackID"`
	TrackName    string `json:"TrackName"`
	Lyrics       string `json:"Lyrics"`
	Artist       string `json:"Artist"`
	Soprano      string `json:"Soprano"`
	Alto         string `json:"Alto"`
	Tenor        string `json:"Tenor"`
	Bass         string `json:"Bass"`
	AllParts     string `json:"AllParts"`
	Piano        string `json:"Piano"`
	ExtraTitle   string `json:"ExtraTitle"`
	ExtraLink    string `json:"ExtraLink"`
}

// id, name, description, image
type TeamMember struct {
	ID          int         `json:"ID"`
	Name        string      `json:"Name"`
	Description string      `json:"Description"`
	Image       ImageDetail `json:"Image"`
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
