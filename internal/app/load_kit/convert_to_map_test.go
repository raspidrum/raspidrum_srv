package loadkit

import (
	"path"
	"testing"
)

func TestTransformKitFormat(t *testing.T) {
	kitPath := path.Join(getProjectPath(), "../_presets/", "SMDrums")

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "convert",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := TransformKitFormat(kitPath); (err != nil) != tt.wantErr {
				t.Errorf("TransformKitFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
