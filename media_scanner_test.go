package media

import "testing"

func createArtists() *Artists {
	zAlbums := make([]*Album,0)
	zAlbums = append(zAlbums, &Album{
		Name: "Zaphod",
		Tracks: new(Tracks),
	}, &Album{
		Name: "Ford",
		Tracks: new(Tracks),
	})
	aAlbums := make([]*Album,0)
	aAlbums = append(aAlbums, &Album{
		Name: "Trillian",
		Tracks: new(Tracks),
	}, &Album{
		Name: "Arthur",
		Tracks: new(Tracks),
	})

	z := Albums(zAlbums)
	a := Albums(aAlbums)
	artists := make([]*Artist, 0)
	artists = append(artists, &Artist{
		Name: "Zebra",
		Albums: &z,
	}, &Artist{
		Name: "Giraffe",
		Albums: &a,
	})

	ar := Artists(artists)
	return &ar
}

func TestSortMedia(t *testing.T) {
	artists := createArtists()
	sortMedia(artists)

	if artists.Get(0).Name != "Giraffe" {
		t.Errorf("First artist was %s, expected %s", artists.Get(0).Name, "Giraffe")
	}
	if artists.Get(0).Albums.Get(0).Name != "Arthur" {
		t.Errorf("First album was %s, expected %s", artists.Get(0).Albums.Get(0).Name, "Arthur")
	}

}

func TestCleanupTrack(t *testing.T) {
	cases := map[string]string{
		"1/10": "1",
		"999": "999",
		"2 of 10": "2",
		"apple": "apple",
		"": "2",
		"  ": "2",
	}

	for in, expected := range cases {
		actual := cleanUpTrack(1, in)
		if actual != expected {
			t.Errorf("cleanupTrack() failed; expected %s, got %s", expected, actual)
		}
	}
}
