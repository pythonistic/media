package media

import "os"

const DB_ARTIST = "artists.db"
const DB_PLAYLIST = "playlists.db"
const DB_USER = "users.db"
const DB_TOKEN = "token.db"

func LoadMedia() {
	// open the database files
	artistDb, err := OpenDatabase(DB_ARTIST)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	playlistDb, err := OpenDatabase(DB_PLAYLIST)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	userDb, err := OpenDatabase(DB_USER)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}

	// try to load the artists from the database first
	artists, err := artistDb.GetArtists()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	if len(*artists) == 0 {
		artists = PrepareMedia(FSPATH_MEDIA)
		println(len(*artists))
	}

	// load the playlists
	playlists, err := playlistDb.GetPlaylists()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}

	// load the users
	users, err := userDb.GetUsers()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}

	SetArtists(artists)
	SetPlaylists(playlists)
	SetUsers(users)
}
