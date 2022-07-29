package steamcmd

import (
	"testing"
)

func Test_downloadSteamcmd(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "works",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := downloadSteamcmd(); (err != nil) != tt.wantErr {
				t.Errorf("downloadSteamcmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
