package metadatamanager

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/bogem/id3v2/v2"
)

type SongRepository interface {
	GetSongPath(ctx context.Context, songID string) (string, error)
	RefreshSong(ctx context.Context, songID string) error
}

type MetadataService interface {
	UpdateTags(ctx context.Context, songID string, tags map[string]string) error
}

type mp3Service struct {
	repo SongRepository
}

func NewService(repo SongRepository) MetadataService {
	return &mp3Service{repo: repo}
}

func (s *mp3Service) UpdateTags(ctx context.Context, songID string, tags map[string]string) error {
	path, err := s.repo.GetSongPath(ctx, songID)
	if err != nil {
		return fmt.Errorf("could not retrieve song path: %w", err)
	}

	// Verify file access
	if info, err := os.Stat(path); err != nil || info.IsDir() {
		return fmt.Errorf("file is inaccessible or is a directory: %s", path)
	}

	// Ensure we are only processing MP3 files as id3v2 library is specific to ID3 tags
	if !strings.HasSuffix(strings.ToLower(path), ".mp3") {
		return fmt.Errorf("metadata editing is currently only supported for MP3 files")
	}

	fmt.Printf("Updating tags for: %s\n", path)

	// Open the MP3 file
	tag, err := id3v2.Open(path, id3v2.Options{Parse: true})
	if err != nil {
		return fmt.Errorf("error opening MP3 file: %w", err)
	}
	defer tag.Close()

	for key, value := range tags {
		switch key {
		case "title":
			tag.SetTitle(value)
		case "artist":
			tag.SetArtist(value)
		case "album":
			tag.SetAlbum(value)
		case "albumArtist":
			tag.SetAlbumArtist(value)
		case "genre":
			tag.SetGenre(value)
		case "comment":
			commentFrame := id3v2.CommentFrame{
				Encoding:    id3v2.EncodingUTF8,
				Language:    "eng", // Default language, could be made configurable
				Description: "",
				Text:        value,
			}
			tag.AddCommentFrame(commentFrame)
		case "year":
			// TDRC frame is used for the recording date/year in ID3v2.4
			tag.AddFrame("TDRC", tag.TextFrame(value))
		case "trackNumber":
			// SetTrack handles strings like "01" or "01/12"
			tag.SetTrack(value)
		case "disc":
			// SetDisc handles strings like "1" or "1/2"
			tag.SetDisc(value)
		default:
			fmt.Printf("Warning: Tag '%s' not supported for MP3 metadata. Ignoring.\n", key)
		}
	}

	// Save the changes
	if err = tag.Save(); err != nil {
		return fmt.Errorf("error saving MP3 tags: %w", err)
	}

	// Trigger a rescan of this song so Navidrome updates its database
	return s.repo.RefreshSong(ctx, songID)
}