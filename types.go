package media

import "strconv"

// Artists is the parent media structure, consisting of []*Artist -> []*Album -> []*Track
type Artists []*Artist

func (artists Artists) Len() int {
	return len(artists)
}

func (artists Artists) Swap(i, j int) {
	artists[i], artists[j] = artists[j], artists[i]
}

func (artists Artists) Less(i, j int) bool {
	return artists[i].Name < artists[j].Name
}

func (artists Artists) Get(i int) *Artist {
	return artists[i]
}

func (artists Artists) Contains(artist *Artist) (*Artist, bool) {
	for _, myArtist := range artists {
		if myArtist.Name == artist.Name {
			return myArtist, true
		}
	}

	return nil, false
}

func (artists Artists) GetByName(name string) (*Artist, bool) {
	return artists.Contains(&Artist{Name: name, })
}

// Albums is the type containing all the *Album structs and tracks for a specific artist.
type Albums []*Album

func (albums Albums) Len() int {
	return len(albums)
}

func (albums Albums) Swap(i, j int) {
	albums[i], albums[j] = albums[j], albums[i]
}

func (albums Albums) Less(i, j int) bool {
	return albums[i].Name < albums[j].Name
}

func (albums Albums) Get(i int) *Album {
	return albums[i]
}

func (albums Albums) Contains(album *Album) (*Album, bool) {
	for _, myAlbum := range albums {
		if myAlbum.Name == album.Name {
			return myAlbum, true
		}
	}

	return nil, false
}

func (albums Albums) GetByName(name string) (*Album, bool) {
	return albums.Contains(&Album{Name: name, })
}

// Tracks contains all the *Track structs for an Album.
type Tracks []*Track

func (tracks Tracks) Len() int {
	return len(tracks)
}

func (tracks Tracks) Swap(i, j int) {
	tracks[i], tracks[j] = tracks[j], tracks[i]
}

func (tracks Tracks) Less(i, j int) bool {
	first, err := strconv.Atoi(tracks[i].Track)
	if err == nil {
		second, err := strconv.Atoi(tracks[j].Track)
		if err == nil {
			return first < second
		}
	}
	return tracks[i].Track < tracks[j].Track
}

func (tracks Tracks) GetByTrackTitle(track string, title string) (*Track, bool) {
	for _, myTrack := range tracks {
		if myTrack.Track == track && myTrack.Title == title {
			return myTrack, true
		}
	}

	return nil, false
}

// Playlists is the type containing all the Playlist structs
type Playlists []*Playlist

func (playlists Playlists) Len() int {
	return len(playlists)
}

func (playlists Playlists) Swap(i, j int) {
	playlists[i], playlists[j] = playlists[j], playlists[i]
}

func (playlists Playlists) Less(i, j int) bool {
	return playlists[i].Name < playlists[j].Name
}

func (playlists Playlists) Get(i int) *Playlist {
	return playlists[i]
}

func (playlists Playlists) Contains(playlist *Playlist) (*Playlist, bool) {
	for _, myPlaylist := range playlists {
		if myPlaylist.Name == playlist.Name && myPlaylist.User == playlist.User{
			return myPlaylist, true
		}
	}

	return nil, false
}

func (playlists Playlists) GetByName(user *User, name string) (*Playlist, bool) {
	return playlists.Contains(&Playlist{Name: name, User: user})
}

// Users is the type containing all the User structs
type Users []*User

func (users Users) Len() int {
	return len(users)
}

func (users Users) Swap(i, j int) {
	users[i], users[j] = users[j], users[i]
}

func (users Users) Less(i, j int) bool {
	return users[i].Name < users[j].Name
}

func (users Users) Get(i int) *User {
	return users[i]
}

func (users Users) Contains(user *User) (*User, bool) {
	for _, myUser := range users {
		if myUser.Email == user.Email {
			return myUser, true
		}
	}

	return nil, false
}

func (users Users) GetByEmail(email string) (*User, bool) {
	return users.Contains(&User{Email: email})
}

// Artist is a single artist
type Artist struct {
	Name   string
	Albums *Albums
	Id     string
}

// Album is a single album, belonging to an Artist
type Album struct {
	Name             string
	Year             string
	Tracks           *Tracks
	CoverArtMimeType string
	CoverArtBytes    []byte
	CoverArtPath     string
	Id               string
}

// Track is a single track, belonging to an Album
type Track struct {
	Artist   string
	Album    string
	Filename string
	Path     string
	Title    string
	Duration string
	Track    string
	Id       string
}

// Playlist is an ordered collection of Tracks
type Playlist struct {
	User   *User
	Name   string
	Tracks []*Track
}

// User is the user account record
type User struct {
	Name string
	Email string
	// consider using a map of cookies, IPs, and registration dates for passwordless logins
}

// PageContext contains the structs needed to render the media page.
type PageContext struct {
	Artists   *Artists
	Playlists []*Playlist		// a slice of *Playlist with only playlists belonging to user
	User      *User
}