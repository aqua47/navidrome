package tests

import (
	"context"
	"io"
)

type MockMusicFileManager struct {
	DeleteSongFunc    func(ctx context.Context, songID string) error
	UpdateArtworkFunc func(ctx context.Context, songID string, data io.Reader, mimeType string) error
	UpdateTagsFunc    func(ctx context.Context, songID string, tags map[string]string) error
}

func NewMockMusicFileManager() *MockMusicFileManager {
	return &MockMusicFileManager{}
}

func (m *MockMusicFileManager) DeleteSong(ctx context.Context, songID string) error {
	if m.DeleteSongFunc != nil {
		return m.DeleteSongFunc(ctx, songID)
	}
	return nil
}

func (m *MockMusicFileManager) UpdateArtwork(ctx context.Context, songID string, data io.Reader, mimeType string) error {
	if m.UpdateArtworkFunc != nil {
		return m.UpdateArtworkFunc(ctx, songID, data, mimeType)
	}
	return nil
}

func (m *MockMusicFileManager) UpdateTags(ctx context.Context, songID string, tags map[string]string) error {
	if m.UpdateTagsFunc != nil {
		return m.UpdateTagsFunc(ctx, songID, tags)
	}
	return nil
}
