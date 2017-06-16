package media

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"encoding/json"
)


var dbArtist *leveldb.DB


func OpenDatabase(dbPath string) (err error) {
	dbArtist, err = leveldb.OpenFile(dbPath, &opt.Options{})
	if err != nil {
		return
	}

	return
}

func CloseDatabase() (err error) {
	if dbArtist != nil {
		if err = dbArtist.Close(); err != nil {
			return
		}
	}

	return
}

func GetArtists() (artists *Artists, err error) {
	artists = new(Artists)
	iter := dbArtist.NewIterator(nil, nil)
	for iter.Next() {
		artistBytes := iter.Value()
		artist := &Artist{}
		if err = json.Unmarshal(artistBytes, artist); err != nil {
			return
		}
		*artists = append(*artists, artist)
	}

	return
}

func StoreArtists(artists *Artists) (err error) {
	for _, artist := range(*artists) {
		var artistBytes []byte
		artistBytes, err = json.Marshal(artist)
		if err != nil {
			return
		}
		artistKey := []byte(artist.Name)
		dbArtist.Put(artistKey, artistBytes, &opt.WriteOptions{})
	}

	return
}

/*
func GetPlaylists() (playlists *Playlists, err error)
func StorePlaylists(playlists *Playlists) (err error)
func GetUsers() (users *Users, err error)
func StoreUsers(users *Users) (err error)
*/