package media

func createPlaylists() *Playlists {
	artists := createArtists()

	album0 := *artists.Get(0).Albums.Get(0)
	album1 := *artists.Get(0).Albums.Get(1)

	album0.Tracks = new(Tracks)
	*album0.Tracks = append(*album0.Tracks, &Track{
		Album: album0.Name,
		Track: "1",
		Artist: artists.Get(0).Name,
		Title: "Alpha",
	}, &Track{
		Album: album0.Name,
		Track: "2",
		Artist: artists.Get(0).Name,
		Title: "Beta",
	})
	*album1.Tracks = append(*album1.Tracks, &Track{
		Album: album1.Name,
		Track: "1",
		Artist: artists.Get(0).Name,
		Title: "Gamma",
	}, &Track{
		Album: album1.Name,
		Track: "2",
		Artist: artists.Get(0).Name,
		Title: "Delta",
	})

	alpha, _ := album0.Tracks.GetByTrackTitle("1", "Alpha")
	beta, _ := album0.Tracks.GetByTrackTitle("2", "Beta")
	gamma, _ := album1.Tracks.GetByTrackTitle("1", "Gamma")
	delta, _ := album1.Tracks.GetByTrackTitle("2", "Delta")

	playlist1 := &Playlist{
		Name: "Playlist 1",
		Tracks: []*Track{
			alpha,
			delta,
		},
	}

	playlist2 := &Playlist {
		Name: "Playlist 2",
		Tracks: []*Track{
			beta,
			gamma,
		},
	}

	playlists := new(Playlists)
	*playlists = append(*playlists, playlist1, playlist2)

	return playlists
}

func createUsers() *Users {
	users := new(Users)

	user1 := &User{
		Name: "Alice",
		Email: "alice@unittest.com",
	}

	user2 := &User{
		Name: "Bob",
		Email: "bob@unittest.com",
	}

	*users = append(*users, user1, user2)

	return users
}