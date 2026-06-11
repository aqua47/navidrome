package metadatamanager

import (
	"context"

	"github.com/navidrome/navidrome/core"
	"github.com/navidrome/navidrome/model"
)

type navidromeRepo struct {
	ds      model.DataStore
	library core.Library
}

func NewRepository(ds model.DataStore, library core.Library) SongRepository {
	return &navidromeRepo{
		ds:      ds,
		library: library,
	}
}

func (r *navidromeRepo) GetSongPath(ctx context.Context, songID string) (string, error) {
	mf, err := r.ds.MediaFile(ctx).Get(songID)
	if err != nil {
		return "", err
	}
	return mf.AbsolutePath(), nil
}

func (r *navidromeRepo) RefreshSong(ctx context.Context, songID string) error {
	mf, err := r.ds.MediaFile(ctx).Get(songID)
	if err != nil {
		return err
	}
	// Tell Navidrome to re-read the file metadata and update the DB
	return r.library.ImportFile(ctx, mf.AbsolutePath(), mf.LibraryID, mf.FolderID)
}
