package media

import (
	"testing"
	"io/ioutil"
	"fmt"
	"path/filepath"
	"os"
	"time"
)

const UT_DB_NAME = "media_ut.db"

func getTempDbPath() string {
	tmpDir, err := ioutil.TempDir(".", "UT_")
	if err != nil {
		fmt.Printf("ERROR: failed to create temp dir: %v", err)
	}

	return filepath.Join(tmpDir, UT_DB_NAME)
}

func removeTempDbDir(path string) {
	err := os.RemoveAll(filepath.Dir(path))
	if err != nil {
		fmt.Printf("ERROR: failed to remove temp dir: %v", err)
	}
}

func TestOpenDatabase(t *testing.T) {
	dbPath := getTempDbPath()
	db, err := OpenDatabase(dbPath)
	if err != nil {
		t.Errorf("Failed to open the database: %v", err)
		t.FailNow()
	}

	defer func() {
		removeTempDbDir(dbPath)
	}()

	err = db.CloseDatabase()
	if err != nil {
		t.Errorf("Failed to close the database: %v", err)
		t.FailNow()
	}
}

func TestGetArtists(t *testing.T) {
	dbPath := getTempDbPath()
	db, err := OpenDatabase(dbPath)
	if err != nil {
		t.Errorf("Failed to open the database: %v", err)
		t.FailNow()
	}

	defer func() {
		err = db.CloseDatabase()
		removeTempDbDir(dbPath)
		if err != nil {
			t.Errorf("Failed to close the database: %v", err)
			t.FailNow()
		}

	}()

	artists := createArtists()

	if err = db.StoreArtists(artists); err != nil {
		t.Errorf("Failed to store artists: %v", err)
	}

	fetchedArtists, err := db.GetArtists()
	if err != nil {
		t.Errorf("Failed to get artists: %v", err)
	}

	// expecting two artists
	if len(*fetchedArtists) != 2 {
		t.Errorf("Expected 2 artists, got %d", len(*fetchedArtists))
	}

	// expecting Giraffe
	if fetchedArtists.Get(0).Name != "Giraffe" {
		t.Errorf("Expected Giraffe, got %s", fetchedArtists.Get(0).Name)
	}
}

func TestGetPlaylists(t *testing.T) {
	dbPath := getTempDbPath()
	db, err := OpenDatabase(dbPath)
	if err != nil {
		t.Errorf("Failed to open the database: %v", err)
		t.FailNow()
	}

	defer func() {
		err = db.CloseDatabase()
		removeTempDbDir(dbPath)
		if err != nil {
			t.Errorf("Failed to close the database: %v", err)
			t.FailNow()
		}

	}()

	playlists := createPlaylists()

	if err = db.StorePlaylists(playlists); err != nil {
		t.Errorf("Failed to store playlists: %v", err)
	}

	fetchedPlaylists, err := db.GetPlaylists()
	if err != nil {
		t.Errorf("Failed to get playlists: %v", err)
	}

	// expecting two playlists
	if len(*fetchedPlaylists) != 2 {
		t.Errorf("Expected 2 playlists, got %d", len(*fetchedPlaylists))
	}

	// expecting Playlist 1
	if fetchedPlaylists.Get(0).Name != "Playlist 1" {
		t.Errorf("Expected Playlist 1, got %s", fetchedPlaylists.Get(0).Name)
	}

	// expecting Playlist 2
	if fetchedPlaylists.Get(1).Name != "Playlist 2" {
		t.Errorf("Expected Playlist 2, got %s", fetchedPlaylists.Get(1).Name)
	}

	if fetchedPlaylists.Get(1).Tracks[0].Title != "Beta" {
		t.Errorf("Expected playlist 2, track 1 to be Beta, got %s", fetchedPlaylists.Get(1).Tracks[0].Title)
	}
}

func TestGetUsers(t *testing.T) {
	dbPath := getTempDbPath()
	db, err := OpenDatabase(dbPath)
	if err != nil {
		t.Errorf("Failed to open the database: %v", err)
		t.FailNow()
	}

	defer func() {
		err = db.CloseDatabase()
		removeTempDbDir(dbPath)
		if err != nil {
			t.Errorf("Failed to close the database: %v", err)
			t.FailNow()
		}

	}()

	users := createUsers()

	if err = db.StoreUsers(users); err != nil {
		t.Errorf("Failed to store users: %v", err)
	}

	fetchedUsers, err := db.GetUsers()
	if err != nil {
		t.Errorf("Failed to load users: %v", err)
	}

	if len(*fetchedUsers) != 2 {
		t.Errorf("Unexpected number of users, wanted 2, got %d", len(*fetchedUsers))
	}

	if users.Get(0).Email != "alice@unittest.com" {
		t.Errorf("Incorrect user email, expected alice@unittest.com, got %s", users.Get(0).Email)
	}
}

func TestGetToken(t *testing.T) {
	dbPath := getTempDbPath()
	db, err := OpenDatabase(dbPath)
	if err != nil {
		t.Errorf("Failed to open the database: %v", err)
		t.FailNow()
	}

	defer func() {
		err = db.CloseDatabase()
		removeTempDbDir(dbPath)
		if err != nil {
			t.Errorf("Failed to close the database: %v", err)
			t.FailNow()
		}

	}()

	// create and store some tokens
	token1 := &Token{
		Code: "alpha",
		Email: "12345@67890",
		Expiration: time.Now(),
	}

	token2 := &Token{
		Code: "beta",
		Email: "09876@54321",
		Expiration: time.Now(),
	}

	err = db.StoreToken(token1)
	if err != nil {
		t.Errorf("Failed to store token 1: %v", err)
		t.FailNow()
	}
	err = db.StoreToken(token2)
	if err != nil {
		t.Errorf("Failed to store token 2: %v", err)
		t.FailNow()
	}

	// get a token that doesn't exist
	badToken, err := db.GetToken("gamma")
	if err != nil {
		t.Errorf("Unepxected error for token gamma: %v", err)
		t.FailNow()
	}
	if badToken != nil {
		t.Errorf("BadToken should be nil, got %v", badToken)
		t.FailNow()
	}

	// get token 1
	fetchedToken, err := db.GetToken("alpha")
	if err != nil {
		t.Errorf("Excepted error for token alpha: %v", err)
		t.FailNow()
	}
	if fetchedToken == nil {
		t.Error("Unexpected nil token for alpha")
		t.FailNow()
	}
	if fetchedToken.Code != token1.Code {
		t.Errorf("Code mismatch, expected %s got %s", token1.Code, fetchedToken.Code)
	}
	if fetchedToken.Email != token1.Email {
		t.Errorf("Email mismatch, expected %s got %s", token1.Email, fetchedToken.Email)
	}
	if fetchedToken.Expiration != token1.Expiration {
		t.Errorf("Expiration mismatch, expected %v got %v", token1.Expiration, fetchedToken.Expiration)
	}

}