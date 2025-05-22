package preset

import (
	"testing"

	"github.com/raspidrum-srv/internal/model"
	pb "github.com/raspidrum-srv/internal/pkg/grpc"
	"github.com/stretchr/testify/assert"
)

func TestConvertInstrumentToProto(t *testing.T) {
	tests := []struct {
		name        string
		instruments []*model.PresetInstrument
		want        []*pb.Instrument
	}{
		{
			name: "single instrument without controls",
			instruments: []*model.PresetInstrument{
				{
					Id:       1,
					Name:     "Kick",
					Controls: map[string]model.PresetControl{},
					Layers:   map[string]model.PresetLayer{},
				},
			},
			want: []*pb.Instrument{
				{
					Key:  "1",
					Name: "Kick",
				},
			},
		},
		{
			name: "multiple instruments with volume and pan",
			instruments: []*model.PresetInstrument{
				{
					Id:   1,
					Name: "Kick",
					Controls: map[string]model.PresetControl{
						"volume": {Name: "Volume", Type: "volume", Value: 100},
						"pan":    {Name: "Pan", Type: "pan", Value: 64},
					},
					Layers: map[string]model.PresetLayer{},
				},
				{
					Id:   2,
					Name: "Snare",
					Controls: map[string]model.PresetControl{
						"volume": {Name: "Volume", Type: "volume", Value: 90},
						"pan":    {Name: "Pan", Type: "pan", Value: 32},
					},
					Layers: map[string]model.PresetLayer{},
				},
			},
			want: []*pb.Instrument{
				{
					Key:    "1",
					Name:   "Kick",
					Volume: float64Ptr(100),
					Pan:    float64Ptr(64),
				},
				{
					Key:    "2",
					Name:   "Snare",
					Volume: float64Ptr(90),
					Pan:    float64Ptr(32),
				},
			},
		},
		{
			name: "instrument with layers and controls",
			instruments: []*model.PresetInstrument{
				{
					Id:   1,
					Name: "Kick",
					Controls: map[string]model.PresetControl{
						"pitch": {Name: "Pitch", Type: "pitch", Value: 50},
					},
					Layers: map[string]model.PresetLayer{
						"layer1": {
							Name: "Main",
							Controls: map[string]model.PresetControl{
								"volume": {Name: "Volume", Type: "volume", Value: 100},
								"pan":    {Name: "Pan", Type: "pan", Value: 64},
							},
						},
					},
				},
			},
			want: []*pb.Instrument{
				{
					Key:  "1",
					Name: "Kick",
					Tunes: []*pb.FX{
						{
							Key:   "pitch",
							Name:  "Pitch",
							Order: 0,
							Params: []*pb.FXParam{
								{
									Key:   "pitch",
									Name:  "Pitch",
									Type:  pb.FXParamType_FX_PARAM_TYPE_RANGE,
									Value: 50,
									Min:   float64Ptr(0),
									Max:   float64Ptr(127),
								},
							},
						},
					},
					Layers: []*pb.Layer{
						{
							Key:    "layer1",
							Name:   "Main",
							Volume: float64Ptr(100),
							Pan:    float64Ptr(64),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertInstrumentToProto(tt.instruments)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvertPresetToProto(t *testing.T) {
	tests := []struct {
		name    string
		preset  *model.KitPreset
		want    *pb.Preset
		wantErr bool
	}{
		{
			name: "basic preset with single instrument channel",
			preset: &model.KitPreset{
				Uid:  "preset-1",
				Name: "Test Preset",
				Channels: []model.PresetChannel{
					{
						Key:  "ch1",
						Name: "Channel 1",
						Controls: map[string]model.PresetControl{
							"volume": {Name: "Volume", Type: "volume", Value: 100},
							"pan":    {Name: "Pan", Type: "pan", Value: 64},
						},
					},
				},
				Instruments: []model.PresetInstrument{
					{
						Id:         1,
						Name:       "Kick",
						ChannelKey: "ch1",
						Controls: map[string]model.PresetControl{
							"volume": {Name: "Volume", Type: "volume", Value: 90, MidiCC: 7},
							"pan":    {Name: "Pan", Type: "pan", Value: 54, MidiCC: 10},
						},
					},
				},
			},
			want: &pb.Preset{
				Key:  "preset-1",
				Name: "Test Preset",
				Channels: []*pb.Channel{
					{
						Key:  "sampler",
						Name: "Kit",
						Type: pb.ChannelType_CHANNEL_TYPE_SAMPLER,
					},
					{
						Key:    "ch1",
						Name:   "Channel 1",
						Type:   pb.ChannelType_CHANNEL_TYPE_INSTRUMENT,
						Volume: 90,
						Pan:    54,
						Instruments: []*pb.Instrument{
							{
								Key:  "1",
								Name: "Kick",
							},
						},
					},
				},
			},
		},
		{
			name: "basic preset with single instrument channel, controls without MidiCC",
			preset: &model.KitPreset{
				Uid:  "preset-1",
				Name: "Test Preset",
				Channels: []model.PresetChannel{
					{
						Key:  "ch1",
						Name: "Channel 1",
						Controls: map[string]model.PresetControl{
							"volume": {Name: "Volume", Type: "volume", Value: 100},
							"pan":    {Name: "Pan", Type: "pan", Value: 64},
						},
					},
				},
				Instruments: []model.PresetInstrument{
					{
						Id:         1,
						Name:       "Kick",
						ChannelKey: "ch1",
						Controls: map[string]model.PresetControl{
							"volume": {Name: "Volume", Type: "volume", Value: 90},
							"pan":    {Name: "Pan", Type: "pan", Value: 54},
						},
					},
				},
			},
			want: &pb.Preset{
				Key:  "preset-1",
				Name: "Test Preset",
				Channels: []*pb.Channel{
					{
						Key:  "sampler",
						Name: "Kit",
						Type: pb.ChannelType_CHANNEL_TYPE_SAMPLER,
					},
					{
						Key:    "ch1",
						Name:   "Channel 1",
						Type:   pb.ChannelType_CHANNEL_TYPE_INSTRUMENT,
						Volume: 100,
						Pan:    64,
						Instruments: []*pb.Instrument{
							{
								Key:  "1",
								Name: "Kick",
							},
						},
					},
				},
			},
		},
		{
			name: "basic preset with single instrument channel, controls without MidiCC, channel hasn't pan",
			preset: &model.KitPreset{
				Uid:  "preset-1",
				Name: "Test Preset",
				Channels: []model.PresetChannel{
					{
						Key:  "ch1",
						Name: "Channel 1",
						Controls: map[string]model.PresetControl{
							"volume": {Name: "Volume", Type: "volume", Value: 100},
						},
					},
				},
				Instruments: []model.PresetInstrument{
					{
						Id:         1,
						Name:       "Kick",
						ChannelKey: "ch1",
						Controls: map[string]model.PresetControl{
							"volume": {Name: "Volume", Type: "volume", Value: 90},
							"pan":    {Name: "Pan", Type: "pan", Value: 54},
						},
					},
				},
			},
			want: &pb.Preset{
				Key:  "preset-1",
				Name: "Test Preset",
				Channels: []*pb.Channel{
					{
						Key:  "sampler",
						Name: "Kit",
						Type: pb.ChannelType_CHANNEL_TYPE_SAMPLER,
					},
					{
						Key:    "ch1",
						Name:   "Channel 1",
						Type:   pb.ChannelType_CHANNEL_TYPE_INSTRUMENT,
						Volume: 100,
						Pan:    0,
						Instruments: []*pb.Instrument{
							{
								Key:  "1",
								Name: "Kick",
							},
						},
					},
				},
			},
		},
		{
			name: "preset with multiple instruments in channel",
			preset: &model.KitPreset{
				Uid:  "preset-2",
				Name: "Multi Preset",
				Channels: []model.PresetChannel{
					{
						Key:  "ch1",
						Name: "Drums",
						Controls: map[string]model.PresetControl{
							"volume": {Name: "Volume", Type: "volume", Value: 100},
							"pan":    {Name: "Pan", Type: "pan", Value: 64},
						},
					},
				},
				Instruments: []model.PresetInstrument{
					{
						Id:         1,
						Name:       "Kick",
						ChannelKey: "ch1",
						Controls: map[string]model.PresetControl{
							"volume": {Name: "Volume", Type: "volume", Value: 90},
							"pan":    {Name: "Pan", Type: "pan", Value: 32},
						},
					},
					{
						Id:         2,
						Name:       "Snare",
						ChannelKey: "ch1",
						Controls: map[string]model.PresetControl{
							"volume": {Name: "Volume", Type: "volume", Value: 85},
							"pan":    {Name: "Pan", Type: "pan", Value: 96},
						},
					},
				},
			},
			want: &pb.Preset{
				Key:  "preset-2",
				Name: "Multi Preset",
				Channels: []*pb.Channel{
					{
						Key:  "sampler",
						Name: "Kit",
						Type: pb.ChannelType_CHANNEL_TYPE_SAMPLER,
					},
					{
						Key:    "ch1",
						Name:   "Drums",
						Type:   pb.ChannelType_CHANNEL_TYPE_INSTRUMENT,
						Volume: 100,
						Pan:    64,
						Instruments: []*pb.Instrument{
							{
								Key:    "1",
								Name:   "Kick",
								Volume: float64Ptr(90),
								Pan:    float64Ptr(32),
							},
							{
								Key:    "2",
								Name:   "Snare",
								Volume: float64Ptr(85),
								Pan:    float64Ptr(96),
							},
						},
					},
				},
			},
		},
		{
			name: "preset with invalid channel key",
			preset: &model.KitPreset{
				Uid:  "preset-3",
				Name: "Invalid Preset",
				Channels: []model.PresetChannel{
					{
						Key:  "ch1",
						Name: "Channel 1",
					},
				},
				Instruments: []model.PresetInstrument{
					{
						Id:         1,
						Name:       "Kick",
						ChannelKey: "invalid-channel",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertPresetToProto(tt.preset)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// Helper function to create float64 pointer
func float64Ptr(v float64) *float64 {
	return &v
}
