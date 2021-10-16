package strings

import (
	"strings"
)

// Extensions returns extensions as key, and value as their type/format
func Extensions() map[string]string {
	return map[string]string{
		".txt": "text",

		".jpg":  "image",
		".png":  "image",
		".webp": "image",
		".jpeg": "image",
		".gif":  "image",

		".mp3":  "audio",
		".aac":  "audio",
		".wav":  "audio",
		".flac": "audio",

		".mp4":  "video",
		".mov":  "video",
		".webm": "video",
		".mkv":  "video",
		".wmv":  "video",
		".avi":  "video",
	}
}

// GetExtension gets extension from filename
func GetExtension(filename string) string {
	for ext, format := range Extensions() {
		if strings.HasSuffix(filename, ext) {
			return format
		}
	}

	return "none"
}
