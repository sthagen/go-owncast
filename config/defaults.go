package config

import "github.com/owncast/owncast/models"

// Defaults will hold default configuration values.
type Defaults struct {
	Name            string
	Title           string
	Summary         string
	Logo            string
	Tags            []string
	PageBodyContent string

	DatabaseFilePath string
	WebServerPort    int
	RTMPServerPort   int
	StreamKey        string

	YPEnabled bool
	YPServer  string

	SegmentLengthSeconds int
	SegmentsInPlaylist   int
	StreamVariants       []models.StreamOutputVariant
}

// GetDefaults will return default configuration values.
func GetDefaults() Defaults {
	return Defaults{
		Name:    "Owncast",
		Title:   "My Owncast Server",
		Summary: "This is brief summary of whom you are or what your stream is. You can edit this description in the admin.",
		Logo:    "logo.svg",
		Tags: []string{
			"owncast",
			"streaming",
		},

		PageBodyContent: "# This is your page content that can be edited from the admin.",

		DatabaseFilePath: "data/owncast.db",

		YPEnabled: false,
		YPServer:  "https://yp.owncast.online",

		WebServerPort:  8080,
		RTMPServerPort: 1935,
		StreamKey:      "abc123",

		StreamVariants: []models.StreamOutputVariant{
			{
				IsAudioPassthrough: true,
				VideoBitrate:       1200,
				EncoderPreset:      "veryfast",
				Framerate:          24,
			},
		},
	}
}
