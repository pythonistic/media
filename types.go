package media

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

func (artists Artists) Contains(artist *Artist) (*Artist, bool) {
	for _, myArtist := range artists {
		if myArtist.Name == artist.Name {
			return myArtist, true
		}
	}

	return nil, false
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

func (albums Albums) Contains(album *Album) (*Album, bool) {
	for _, myAlbum := range albums {
		if myAlbum.Name == album.Name {
			return myAlbum, true
		}
	}

	return nil, false
}

type Tracks []*Track

func (tracks Tracks) Len() int {
	return len(tracks)
}

func (tracks Tracks) Swap(i, j int) {
	tracks[i], tracks[j] = tracks[j], tracks[i]
}

func (tracks Tracks) Less(i, j int) bool {
	return tracks[i].Track < tracks[j].Track
}

type Artist struct {
	Name string
	Albums Albums
}

type Album struct {
	Name string
	Year string
	Tracks Tracks
	CoverArtMimeType string
	CoverArtBytes []byte
}

type Track struct {
	Filename        string
	Path            string
	Title           string
	Duration        string
	Track           string
}

type PageContext struct {
	Artists []*Artist
}