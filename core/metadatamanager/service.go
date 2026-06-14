package metadatamanager

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/bogem/id3v2/v2"
	"github.com/navidrome/navidrome/log"
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

	log.Info(ctx, "Updating MP3 tags", "path", path, "songID", songID)

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
			tag.AddTextFrame("TPE2", id3v2.EncodingUTF8, value)
		case "genre":
			tag.SetGenre(value)
		case "comment":
			tag.AddCommentFrame(id3v2.CommentFrame{
				Encoding: id3v2.EncodingUTF8,
				Language: "eng",
				Text:     value,
			})
		case "year":
			tag.SetYear(value)
		case "trackNumber":
			tag.AddTextFrame("TRCK", id3v2.EncodingUTF8, value)
		case "disc":
			tag.AddTextFrame("TPOS", id3v2.EncodingUTF8, value)
		default:
			log.Warn(ctx, "Tag not supported for MP3 metadata. Ignoring.", "tag", key)
		}
	}

	// Save the changes
	if err = tag.Save(); err != nil {
		return fmt.Errorf("error saving MP3 tags: %w", err)
	}

	// Trigger a rescan of this song so Navidrome updates its database
	return s.repo.RefreshSong(ctx, songID)
}
