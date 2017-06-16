package media

import (
	"testing"
	"io/ioutil"
	"fmt"
	"path/filepath"
	"os"
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
	err := OpenDatabase(dbPath)
	if err != nil {
		t.Errorf("Failed to open the database: %v", err)
		t.FailNow()
	}

	defer func() {
		removeTempDbDir(dbPath)
	}()

	err = CloseDatabase()
	if err != nil {
		t.Errorf("Failed to close the database: %v", err)
		t.FailNow()
	}
}

func TestGetArtists(t *testing.T) {
	dbPath := getTempDbPath()
	err := OpenDatabase(dbPath)
	if err != nil {
		t.Errorf("Failed to open the database: %v", err)
		t.FailNow()
	}

	defer func() {
		err = CloseDatabase()
		removeTempDbDir(dbPath)
		if err != nil {
			t.Errorf("Failed to close the database: %v", err)
			t.FailNow()
		}

	}()

	artists := createArtists()

	if err = StoreArtists(artists); err != nil {
		t.Errorf("Failed to store artists: %v", err)
	}

	fetchedArtists, err := GetArtists()
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
