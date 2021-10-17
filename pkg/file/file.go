package file

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

// ContentTypes returns extensions as key, and value as their type/format
func ContentTypes() map[string]string {
	return map[string]string{
		".txt": "text/plain",

		".jpg":  "image/jpeg",
		".png":  "image/png",
		".webp": "image/webp",
		".jpeg": "image/jpeg",
		".gif":  "image/gif",

		".mp3":  "audio/mpeg",
		".aac":  "audio/aac",
		".wav":  "audio/wav",
		".flac": "audio/flac",
		".weba": "audio/webm",

		".mp4":  "video/mp4",
		".mov":  "video/quicktime",
		".webm": "video/webm",
		".mkv":  "video/x-matroska",
		".wmv":  "video/x-ms-wmv",
		".avi":  "video/x-msvideo",
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

// DetectContentType gets ContentType from filename
func DetectContentType(filename string) string {
	for ext, format := range ContentTypes() {
		if strings.HasSuffix(filename, ext) {
			return format
		}
	}

	return "none"
}
