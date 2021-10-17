package file

import (
	"testing"
)

func TestGetExtension(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "raw", args: args{"bin"}, want: "none"},
		{name: "jpg file", args: args{"bin.jpg"}, want: "image"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetExtension(tt.args.filename); got != tt.want {
				t.Errorf("GetExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectContentType(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "raw", args: args{"bin"}, want: "none"},
		{name: "jpg file", args: args{"bin.jpg"}, want: "image/jpeg"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectContentType(tt.args.filename); got != tt.want {
				t.Errorf("DetectContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}
