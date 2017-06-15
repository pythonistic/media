package media

import "strconv"

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

type Artist struct {
	Name   string
	Albums *Albums
	Id     string
}

type Album struct {
	Name             string
	Year             string
	Tracks           *Tracks
	CoverArtMimeType string
	CoverArtBytes    []byte
	CoverArtPath     string
	Id               string
}

type Track struct {
	Filename string
	Path     string
	Title    string
	Duration string
	Track    string
	Id       string
}

type PageContext struct {
	Artists *Artists
}