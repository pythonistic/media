package media

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"encoding/json"
)


type Database struct {
	*leveldb.DB
}

func OpenDatabase(dbPath string) (database *Database, err error) {
	var db *leveldb.DB
	db, err = leveldb.OpenFile(dbPath, &opt.Options{})
	if err != nil {
		return
	}

	database = &Database{db}
	return
}

func (db *Database) CloseDatabase() (err error) {
	if err = db.Close(); err != nil {
		return
	}

	return
}

func (db *Database) GetArtists() (artists *Artists, err error) {
	artists = new(Artists)
	iter := db.NewIterator(nil, nil)
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

func (db *Database) StoreArtists(artists *Artists) (err error) {
	for _, artist := range(*artists) {
		var artistBytes []byte
		artistBytes, err = json.Marshal(artist)
		if err != nil {
			return
		}
		artistKey := []byte(artist.Name)
		db.Put(artistKey, artistBytes, &opt.WriteOptions{})
	}

	return
}

func (db *Database) GetPlaylists() (playlists *Playlists, err error) {
	playlists = new(Playlists)
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		playlistBytes := iter.Value()
		playlist := &Playlist{}
		if err = json.Unmarshal(playlistBytes, playlist); err != nil {
			return
		}
		*playlists = append(*playlists, playlist)
	}

	return
}

func (db *Database) StorePlaylists(playlists *Playlists) (err error) {
	for _, playlist := range(*playlists) {
		var playlistBytes []byte
		playlistBytes, err = json.Marshal(playlist)
		if err != nil {
			return
		}
		playlistKey := []byte(playlist.Name)
		db.Put(playlistKey, playlistBytes, &opt.WriteOptions{})
	}

	return
}

func (db *Database) GetUsers() (users *Users, err error) {
	users = new(Users)
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		userBytes := iter.Value()
		user := &User{}
		if err = json.Unmarshal(userBytes, user); err != nil {
			return
		}
		*users = append(*users, user)
	}

	return
}

func (db *Database) StoreUsers(users *Users) (err error) {
	for _, user := range(*users) {
		var userBytes []byte
		userBytes, err = json.Marshal(user)
		if err != nil {
			return
		}
		userKey := []byte(user.Email)
		db.Put(userKey, userBytes, &opt.WriteOptions{})
	}
	return
}
