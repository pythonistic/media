package media

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Database is a structure to hold a reference to an open LevelDB.
type Database struct {
	*leveldb.DB
}

// Open the LevelDB file at the directory path.
func OpenDatabase(dbPath string) (database *Database, err error) {
	var db *leveldb.DB
	db, err = leveldb.OpenFile(dbPath, &opt.Options{})
	if err != nil {
		return
	}

	database = &Database{db}
	return
}

// Close the open LevelDB.
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
		err = db.Put(artistKey, artistBytes, &opt.WriteOptions{})
		if err != nil {
			return
		}
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
		err = db.Put(playlistKey, playlistBytes, &opt.WriteOptions{})
		if err != nil {
			return
		}
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
		err = db.Put(userKey, userBytes, &opt.WriteOptions{})
		if err != nil {
			return
		}
	}
	return
}

func (db *Database) GetToken(code string) (token *Token, err error) {
	limit := code + "\u0000"
	iter := db.NewIterator(&util.Range{
		Limit: []byte(limit),
		Start: []byte(code),
	}, nil)
	for iter.Next() {
		tokenBytes := iter.Value()
		token = &Token{}
		if err = json.Unmarshal(tokenBytes, token); err != nil {
			return
		}
	}

	return
}

func (db *Database) StoreToken(token *Token) (err error) {
	var tokenBytes []byte
	tokenBytes, err = json.Marshal(token)
	if err != nil {
		return
	}
	tokenKey := []byte(token.Code)
	err = db.Put(tokenKey, tokenBytes, &opt.WriteOptions{})

	return
}

func (db *Database) DeleteToken(token *Token) (err error) {
	tokenKey := []byte(token.Code)
	return db.Delete(tokenKey, &opt.WriteOptions{})
}