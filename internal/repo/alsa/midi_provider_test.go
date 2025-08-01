//go:build linux
// +build linux

package alsa

import (
	"regexp"
	"testing"
)

func TestAlsaMidiProvider_getDevIdFromUDevEvent(t *testing.T) {
	type fields struct {
		devPathRegexp *regexp.Regexp
	}
	type args struct {
		devId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantId string
		want   bool
	}{
		{
			name: "controlC3",
			fields: fields{
				devPathRegexp: regexp.MustCompile(`sound/card(\d+)/midiC(\d+)D(\d+)`),
			},
			args: args{
				devId: "/devices/platform/scb/fd500000.pcie/pci0000:00/0000:00:00.0/0000:01:00.0/usb1/1-1/1-1.1/1-1.1:1.0/sound/card3/controlC3",
			},
			wantId: "",
			want:   false,
		},
		{
			name: "midiC3D0",
			fields: fields{
				devPathRegexp: regexp.MustCompile(`sound/card(\d+)/midiC(\d+)D(\d+)`),
			},
			args: args{
				devId: "/devices/platform/scb/fd500000.pcie/pci0000:00/0000:00:00.0/0000:01:00.0/usb1/1-1/1-1.1/1-1.1:1.0/sound/card3/midiC3D0",
			},
			wantId: "3:0",
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			did, got := getDevIdFromUDevEvent(tt.fields.devPathRegexp, tt.args.devId)
			if did != tt.wantId {
				t.Errorf("AlsaMidiProvider.getDevIdFromUDevEvent() got = %v, want %v", did, tt.wantId)
			}
			if got != tt.want {
				t.Errorf("AlsaMidiProvider.getDevIdFromUDevEvent() got = %v, want %v", got, tt.want)
			}
		})
	}
}
