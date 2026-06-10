package metadatamanager

import (
	"fmt"
	"os"
)

type SongRepository interface {
	GetSongPath(songID string) (string, error)
	RefreshSong(songID string) error
}

type MetadataService interface {
	UpdateTags(songID string, tags map[string]string) error
}

type mp3Service struct {
	repo SongRepository
}

func NewService(repo SongRepository) MetadataService {
	return &mp3Service{repo: repo}
}

func (s *mp3Service) UpdateTags(songID string, tags map[string]string) error {
	path, err := s.repo.GetSongPath(songID)
	if err != nil {
		return fmt.Errorf("could not retrieve song path: %w", err)
	}

	// Verify file access
	if info, err := os.Stat(path); err != nil || info.IsDir() {
		return fmt.Errorf("file is inaccessible or is a directory: %s", path)
	}

	fmt.Printf("Updating tags for: %s\n", path)

	// Next step: Use a library like "github.com/bogem/id3v2" 
	// to open the file at 'path' and apply the provided 'tags'.

	return s.repo.RefreshSong(songID)
}