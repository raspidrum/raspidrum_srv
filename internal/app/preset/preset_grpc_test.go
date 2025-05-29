package preset

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/raspidrum-srv/internal/model"
	pb "github.com/raspidrum-srv/internal/pkg/grpc"
	"github.com/stretchr/testify/assert"
)

func TestConvertPresetToProto(t *testing.T) {
	type args struct {
		preset   *model.KitPreset
		mididevs []model.MIDIDevice
	}
	tests := []struct {
		name     string
		testData string
		args     args
		want     *pb.Preset
		wantErr  bool
	}{
		{
			name:     "basic preset with single instrument channel",
			testData: "single_instrument.yaml",
			args: args{
				mididevs: []model.MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			want: &pb.Preset{
				Key:  "preset-1",
				Name: "Single Instrument",
				Channels: []*pb.Channel{
					{
						Key:  "sampler",
						Name: "Kit",
						Type: pb.ChannelType_CHANNEL_TYPE_SAMPLER,
						Volume: &pb.BaseControl{
							Key:   "s0volume",
							Name:  "Volume",
							Value: 1.0,
							Min:   makeFloat64Ptr(0),
							Max:   makeFloat64Ptr(1),
						},
					},
					{
						Key:    "ch1",
						Name:   "Kick",
						Type:   pb.ChannelType_CHANNEL_TYPE_INSTRUMENT,
						Volume: &pb.BaseControl{Key: "i0volume", Name: "Volume", Value: 0.748, Min: makeFloat64Ptr(0), Max: makeFloat64Ptr(1)},
						Pan:    &pb.BaseControl{Key: "i0pan", Name: "Pan", Value: -0.15, Min: makeFloat64Ptr(-1), Max: makeFloat64Ptr(1)},
						Instruments: []*pb.Instrument{
							{
								Key:  "0",
								Name: "Kick",
							},
						},
					},
				},
			},
		},
		{
			name:     "basic preset with single instrument channel, controls without MidiCC",
			testData: "single_instrument_controls_wo_midicc.yaml",
			args: args{
				mididevs: []model.MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			want: &pb.Preset{
				Key:  "preset-1",
				Name: "Single Instrument controls without MidiCC",
				Channels: []*pb.Channel{
					{
						Key:  "sampler",
						Name: "Kit",
						Type: pb.ChannelType_CHANNEL_TYPE_SAMPLER,
						Volume: &pb.BaseControl{
							Key:   "s0volume",
							Name:  "Volume",
							Value: 1.0,
							Min:   makeFloat64Ptr(0),
							Max:   makeFloat64Ptr(1),
						},
					},
					{
						Key:    "ch1",
						Name:   "Kick",
						Type:   pb.ChannelType_CHANNEL_TYPE_INSTRUMENT,
						Volume: &pb.BaseControl{Key: "c0volume", Name: "Volume", Value: 1.00, Min: makeFloat64Ptr(0), Max: makeFloat64Ptr(1)},
						Instruments: []*pb.Instrument{
							{
								Key:    "0",
								Name:   "Kick",
								Volume: &pb.BaseControl{Key: "i0volume", Name: "Volume", Value: 0.95, Min: makeFloat64Ptr(0), Max: makeFloat64Ptr(1)},
								Pan:    &pb.BaseControl{Key: "i0pan", Name: "Pan", Value: -0.2, Min: makeFloat64Ptr(-1), Max: makeFloat64Ptr(1)},
							},
						},
					},
				},
			},
		},
		{
			name:     "preset with multiple instruments in channel",
			testData: "two_instruments.yaml",
			args: args{
				mididevs: []model.MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			want: &pb.Preset{
				Key:  "preset-2",
				Name: "multiple instruments in channel",
				Channels: []*pb.Channel{
					{
						Key:  "sampler",
						Name: "Kit",
						Type: pb.ChannelType_CHANNEL_TYPE_SAMPLER,
						Volume: &pb.BaseControl{
							Key:   "s0volume",
							Name:  "Volume",
							Value: 1.0,
							Min:   makeFloat64Ptr(0),
							Max:   makeFloat64Ptr(1),
						},
					},
					{
						Key:    "ch1",
						Name:   "Drums",
						Type:   pb.ChannelType_CHANNEL_TYPE_INSTRUMENT,
						Volume: &pb.BaseControl{Key: "c0volume", Name: "Volume", Value: 1.00, Min: makeFloat64Ptr(0), Max: makeFloat64Ptr(1)},
						Pan:    &pb.BaseControl{Key: "c0pan", Name: "Pan", Value: 0.00, Min: makeFloat64Ptr(-1), Max: makeFloat64Ptr(1)},
						Instruments: []*pb.Instrument{
							{
								Key:    "0",
								Name:   "Kick",
								Volume: &pb.BaseControl{Key: "i0volume", Name: "Volume", Value: 0.748, Min: makeFloat64Ptr(0), Max: makeFloat64Ptr(1)},
								Pan:    &pb.BaseControl{Key: "i0pan", Name: "Pan", Value: -0.15, Min: makeFloat64Ptr(-1), Max: makeFloat64Ptr(1)},
							},
							{
								Key:    "1",
								Name:   "Tom",
								Volume: &pb.BaseControl{Key: "i1volume", Name: "Volume", Value: 0.685, Min: makeFloat64Ptr(0), Max: makeFloat64Ptr(1)},
							},
						},
					},
				},
			},
		},
		// TODO: add test for preset with single instrument with layers for test instrument tunes
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preset := loadPresetFromYAML(t, tt.testData)
			tt.args.preset = preset
			if err := preset.PrepareToLoad(tt.args.mididevs); (err != nil) != tt.wantErr {
				t.Errorf("PrepareToLoad() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := convertPresetToProto(tt.args.preset)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if diff := cmp.Diff(tt.want, got,
				cmpopts.IgnoreUnexported(pb.Preset{}),
				cmpopts.IgnoreUnexported(pb.Channel{}),
				cmpopts.IgnoreUnexported(pb.BaseControl{}),
				cmpopts.IgnoreUnexported(pb.FX{}),
				cmpopts.IgnoreUnexported(pb.FXParam{}),
				cmpopts.IgnoreUnexported(pb.Instrument{}),
				cmpopts.IgnoreUnexported(pb.Layer{})); diff != "" {
				t.Errorf("convertPresetToProto() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// Helper function to create float64 pointer
func float64Ptr(v float64) *float64 {
	return &v
}
