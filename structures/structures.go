package structures

import (
	"database/sql"
	// "time"
)

// We need to recretae the structures in Go in order to be able to use them in the frontend

type ImageDetail struct {
	ImageID      int
	Filename     string
	Caption      string
	EventDetails EventDetails
	Rows         int
	Cols         int
	Height       int
	Width        int
	EventID      int
	Imagetype    int
}

type DatURLResponse struct {
	ReturnedFile sql.NullString
	FileDetails  ImageDetail
}

type URLdetails struct {
	URL string
}

type EventDetails struct {
	EventID      int
	Location     string
	EventDate    string
	StartTime    string
	EndTime      string
	Invitation   string
	MeetingPoint string
	Price        string
	Title        string
	Playlist     []PlaylistEntry
}

type ThemeDetails struct {
	ID               int
	BoxColour        string
	TextColour       string
	TextFont         string
	BackgroundImage  sql.NullString
	TextboxColour    string
	LogoImage        sql.NullString
	BannerColour     string
	MenuColour       string
	ButtonColour     string
	ButtonHover      string
	ButtonTextColour string
	MenuTextColour   string
	TextSize         int
}

type PlaylistEntry struct {
	ID         int
	PlaylistID int
	EventID    int
	MusicTrack MusicTrack
	Playorder  int
}

type Clip struct {
	ID      int
	ClipURL string
	EventID int
	Caption string
}

type ArchiveEntry struct {
	ArchiveID    int
	NextFile     string
	Imagecaption string
	NextURL      string
	Clipcaption  string
	EventDetails EventDetails
	Report       string
	Images       []ImageDetail
	Clips        []Clip
}

type Message struct {
	MessageID      int
	MessageDate    string
	MessageFrom    string
	MessageContent string
}

type User struct {
	Username string
	Password string
	Role     string
}

type ImageFiles struct {
	MainImage   sql.NullString
	MobileImage sql.NullString
	EventID     int
}

type SiteInfo struct {
	ID            int
	HomeTitle     string
	HomeText      string
	AboutTitle    string
	AboutText     string
	ArchiveTitle  string
	ArchiveText   string
	NoticesTitle  string
	NoticesText   string
	BookingTitle  string
	BookingText   string
	MembersTitle  string
	MembersText   string
	AppealTitle   string
	AppealText    string
	SettingsTitle string
	SettingsText  string
}

type MusicTrack struct {
	ID        int
	TrackName string
	Lyrics    string
	Artist    string
	Soprano   string
	Alto      string
	Tenor     string
	AllParts  string
	Piano     string
}

// export interface ImageDetail {
//     imageID: number;
//     filename: string;
//     caption: string;
//     eventDetails: EventDetails;
//     rows: number;
//     cols: number;
//     height: number;
//     width: number;
//     eventID: number;
//     imagetype: number;
// }

// export function EmptyImageDetail(): ImageDetail {
//     return {
//         imageID: 0,
//         filename: "",
//         caption: "",
//         eventDetails: EmptyEventDetails(),
//         rows: 1,
//         cols: 1,
//         height: 0,
//         width: 0,
//         eventID: 0,
//         imagetype: 0
//     };
// }

// export interface DatURLResponse {
//     returnedFile: File;
//     fileDetails: ImageDetail;
// }
// export function EmptyDatURLResponse(): DatURLResponse {
//     return {
//         returnedFile: new File([], ""),
//         fileDetails: EmptyImageDetail()
//     };
// }

// export interface URLdetails {
//     url: string;
// }
// export function EmptyURLdetails(): URLdetails {
//     return {
//         url: "http://localhost:3001/"
//     };
// }
// //eventID, location, eventDate, startTime, endTime, price, title
// export interface EventDetails {
//     eventID: number;
//     location: string;
//     eventDate: Date;
//     startTime: string;
//     endTime: string;
//     invitation: string;
//     meetingPoint: string;
//     price: string;
//     title: string;
//     playlist: PlaylistEntry[];
// }
// // themeDetails, boxColour, textColour, textFont, backgroundImage, textboxColour, logoimage, bannerColour, menuColour, buttonColour, buttonHover, buttonTextColour, menuTextColour
// export interface ThemeDetails {
//     boxColour: string;
//     textColour: string;
//     textFont: string;
//     backgroundImage: string;
//     textboxColour: string;
//     logoImage: string;
//     bannerColour: string;
//     menuColour: string;
//     buttonColour: string;
//     buttonHover: string;
//     buttonTextColour: string;
//     menuTextColour: string;
//     textSize: number;
// }
// export function EmptyThemeDetails(): ThemeDetails {
//     return {
//         boxColour: "",
//         textColour: "",
//         textFont: "",
//         backgroundImage: "",
//         textboxColour: "",
//         logoImage: "",
//         bannerColour: "",
//         menuColour: "",
//         buttonColour: "",
//         buttonHover: "",
//         buttonTextColour: "",
//         menuTextColour: "",
//         textSize: 0
//     };
// }

// export function DefaultThemeDetails(): ThemeDetails {
//     return {
//         boxColour: "white",
//         textColour: "black",
//         textFont: "Arial",
//         backgroundImage: "Musical Background.png",
//         textboxColour: "white",
//         logoImage: "Choir Logo.png",
//         bannerColour: "blue",
//         menuColour: "blue",
//         buttonColour: "blue",
//         buttonHover: "darkblue",
//         buttonTextColour: "white",
//         menuTextColour: "white",
//         textSize: 24
//     };
// }

// export function EmptyEventDetails(): EventDetails {
//     return {
//         eventID: 0,
//         location: "",
//         eventDate: new Date(),
//         startTime: "",
//         endTime: "",
//         invitation: "",
//         meetingPoint: "",
//         price: "",
//         title: "",
//         playlist: [],
//     };
// }

// //playlistID, eventID, musicID, playorder
// export interface PlaylistEntry {
//     id: number;
//     playlistID: number;
//     eventID: number;
//     musicTrack: MusicTrack;
//     playorder: number;
// }

// export function EmptyPlaylistEntry(): PlaylistEntry {
//     return {
//         id: 0,
//         playlistID: 0,
//         eventID: 0,
//         musicTrack: EmptyMusicTrack(),
//         playorder: 0
//     };
// }

// //clipID, clipURL, eventID, caption
// export interface Clip {
//     id: number;
//     clipURL: string;
//     eventID: number;
//     caption: string;
// }
// export function EmptyClip(): Clip {
//     return {
//         id: 0,
//         clipURL: "",
//         eventID: 0,
//         caption: ""
//     };
// }
// //archiveID, eventID, report
// export interface ArchiveEntry {
//     archiveID: number;
//     nextFile: string;
//     imagecaption: string;
//     nextURL: string;
//     clipcaption: string;
//     eventDetails: EventDetails;
//     report: string;
//     images: ImageDetail[];
//     clips: Clip[];
// }
// export function EmptyArchiveEntry(): ArchiveEntry {
//     return {
//         archiveID: 0,
//         eventDetails: EmptyEventDetails(),
//         nextFile: "",
//         imagecaption: "",
//         nextURL: "",
//         clipcaption: "",
//         report: "",
//         images: [],
//         clips: []
//     };
// }
// // musicTrackID, trackName, lyrics, soprano, alto, tenor, allParts
// export interface MusicTrack {
//     musicTrackID: number;
//     trackName: string;
//     artist: string;
//     lyrics: string;
//     soprano: string;
//     alto: string;
//     tenor: string;
//     allParts: string;
//     piano: string;
// }

// export function IsMobile() {
//     const regex = /Mobi|Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i;
//     if(regex.test(navigator.userAgent)) {
//         return "/images/mobile/";
//      } else {
//         return "/images/";
//     }
//   }

// export  const DatetoString =(date: Date) => {
//     var nMonth: number = date.getMonth() + 1; // getMonth() returns a value between 0 and 11 so add 1 to get the correct month
//     var sMonth: string = nMonth.toString();
//     if (nMonth < 10) {
//       sMonth = '0' + sMonth;
//     }
//     var nDay: number = date.getDate();
//     var sDay: string = nDay.toString();
//     if (nDay < 10) {
//       sDay = '0' + sDay;
//     }
//     var sDate = date.getFullYear() + "-" + sMonth + "-" + sDay;
//     return sDate;
//   }

//   export function StringtoDate(DateObj: string) {
//     // format of string is 2024-06-19T23:00:00.000Z
//     var sDate = DateObj.split("-");
//     var nMonth = parseInt(sDate[1])
//     var sDay = sDate[2].split("T");
//     var nDay = parseInt(sDay[0]);
//     return nDay + "/" + nMonth
//   }

//   export interface Message {
//     messageID: number;
//     messageDate: string
//     messageFrom: string;
//     messageContent: string;
//     }
// export function EmptyMessage(): Message {
//     return {
//         messageID: 0,
//         messageDate: "",
//         messageFrom: "",
//         messageContent: ""
//     };
// }
// export interface User {
//     username: string;
//     password: string;
//     role: string;
//   }
// export function EmptyUser(): User {
//     return {
//         username: "",
//         password: "",
//         role: ""
//     };
// }

//   declare var jpgImage: {
//     prototype: File;
//     new(fileBits: [BlobPart], fileName: string, options?: FilePropertyBag): File;
// };

// export interface ImageFiles {
//     mainImage: File;
//     mobileImage: File;
//     eventID: number;
// }

// //id, HomeTitle, HomeText, AboutTitle, AboutText, ArchiveTitle, ArchiveText, NoticesTitle, NoticesText, BookingTitle, BookingText, MembersTitle, MembersText, AppealTitle, AppealText, SettingsTitle, SettingsText

// export interface SiteInfo {
//     id: number;
//     HomeTitle: string;
//     HomeText: string;
//     AboutTitle: string;
//     AboutText: string;
//     ArchiveTitle: string;
//     ArchiveText: string;
//     NoticesTitle: string;
//     NoticesText: string;
//     BookingTitle: string;
//     BookingText: string;
//     MembersTitle: string;
//     MembersText: string;
//     AppealTitle: string;
//     AppealText: string;
//     SettingsTitle: string;
//     SettingsText: string;
// }

// export function EmptySiteInfo(): siteInfo {
//     return {
//         id: 0,
//         HomeTitle: "",
//         HomeText: "",
//         AboutTitle: "",
//         AboutText: "",
//         ArchiveTitle: "",
//         ArchiveText: "",
//         NoticesTitle: "",
//         NoticesText: "",
//         BookingTitle: "",
//         BookingText: "",
//         MembersTitle: "",
//         MembersText: "",
//         AppealTitle: "",
//         AppealText: "",
//         SettingsTitle: "",
//         SettingsText: ""
//     };
// }

// export interface MusicTrack {
//     // using the same schema as the music.json file
//     id : number,
//     trackName : string,
//     lyrics : string,
//     soprano : string,
//     alto : string,
//     tenor : string,
//     allParts : string,
//     piano : string
// }

// export function EmptyMusicTrack(): MusicTrack {
//     return {
//         id: 0,
//         trackName: "",
//         lyrics: "",
//         soprano: "",
//         alto: "",
//         tenor: "",
//         allParts: "",
//         piano: ""
//     };
// }
