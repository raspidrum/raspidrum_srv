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
					Controls: map[string]*model.PresetControl{},
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
					Controls: map[string]*model.PresetControl{
						"volume": {Name: "volume", Type: "volume", Value: 100},
						"pan":    {Name: "pan", Type: "pan", Value: 64},
					},
					Layers: map[string]model.PresetLayer{},
				},
				{
					Id:   2,
					Name: "Snare",
					Controls: map[string]*model.PresetControl{
						"volume": {Name: "volume", Type: "volume", Value: 90},
						"pan":    {Name: "pan", Type: "pan", Value: 32},
					},
					Layers: map[string]model.PresetLayer{},
				},
			},
			want: []*pb.Instrument{
				{
					Key:    "1",
					Name:   "Kick",
					Volume: &pb.BaseControl{Key: "volume", Name: "volume", Value: 100, Min: makeFloat64Ptr(0), Max: makeFloat64Ptr(1)},
					Pan:    &pb.BaseControl{Key: "pan", Name: "pan", Value: 64, Min: makeFloat64Ptr(-1), Max: makeFloat64Ptr(1)},
				},
				{
					Key:    "2",
					Name:   "Snare",
					Volume: &pb.BaseControl{Key: "volume", Name: "volume", Value: 90, Min: makeFloat64Ptr(0), Max: makeFloat64Ptr(1)},
					Pan:    &pb.BaseControl{Key: "pan", Name: "pan", Value: 32, Min: makeFloat64Ptr(-1), Max: makeFloat64Ptr(1)},
				},
			},
		},
		{
			name: "instrument with layers and controls",
			instruments: []*model.PresetInstrument{
				{
					Id:   1,
					Name: "Kick",
					Controls: map[string]*model.PresetControl{
						"pitch": {Name: "Pitch", Type: "pitch", Value: 50},
					},
					Layers: map[string]model.PresetLayer{
						"layer1": {
							Name: "Main",
							Controls: map[string]*model.PresetControl{
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
							Volume: &pb.BaseControl{Key: "volume", Name: "Volume", Value: 100, Min: makeFloat64Ptr(0), Max: makeFloat64Ptr(1)},
							Pan:    &pb.BaseControl{Key: "pan", Name: "Pan", Value: 64, Min: makeFloat64Ptr(-1), Max: makeFloat64Ptr(1)},
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
						Controls: map[string]*model.PresetControl{
							"volume": {Name: "Volume", Type: "volume", Value: 1.00},
							"pan":    {Name: "Pan", Type: "pan", Value: 0.64},
						},
					},
				},
				Instruments: []model.PresetInstrument{
					{
						Id:         1,
						Name:       "Kick",
						ChannelKey: "ch1",
						Controls: map[string]*model.PresetControl{
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
						Volume: &pb.BaseControl{Key: "volume", Name: "Volume", Value: 0.709, Min: makeFloat64Ptr(0), Max: makeFloat64Ptr(1)},
						Pan:    &pb.BaseControl{Key: "pan", Name: "Pan", Value: -0.15, Min: makeFloat64Ptr(-1), Max: makeFloat64Ptr(1)},
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
						Controls: map[string]*model.PresetControl{
							"volume": {Name: "Volume", Type: "volume", Value: 1.00},
							"pan":    {Name: "Pan", Type: "pan", Value: 0.64},
						},
					},
				},
				Instruments: []model.PresetInstrument{
					{
						Id:         1,
						Name:       "Kick",
						ChannelKey: "ch1",
						Controls: map[string]*model.PresetControl{
							"volume": {Name: "Volume", Type: "volume", Value: 0.90},
							"pan":    {Name: "Pan", Type: "pan", Value: 0.54},
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
						Volume: &pb.BaseControl{Key: "volume", Name: "Volume", Value: 1.00, Min: makeFloat64Ptr(0), Max: makeFloat64Ptr(1)},
						Pan:    &pb.BaseControl{Key: "pan", Name: "Pan", Value: 0.64, Min: makeFloat64Ptr(-1), Max: makeFloat64Ptr(1)},
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
						Controls: map[string]*model.PresetControl{
							"volume": {Name: "Volume", Type: "volume", Value: 1.00},
						},
					},
				},
				Instruments: []model.PresetInstrument{
					{
						Id:         1,
						Name:       "Kick",
						ChannelKey: "ch1",
						Controls: map[string]*model.PresetControl{
							"volume": {Name: "Volume", Type: "volume", Value: 0.90},
							"pan":    {Name: "Pan", Type: "pan", Value: 0.54},
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
						Volume: &pb.BaseControl{Key: "volume", Name: "Volume", Value: 1.00, Min: makeFloat64Ptr(0), Max: makeFloat64Ptr(1)},
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
						Controls: map[string]*model.PresetControl{
							"volume": {Name: "Volume", Type: "volume", Value: 1.00},
							"pan":    {Name: "Pan", Type: "pan", Value: 0.64},
						},
					},
				},
				Instruments: []model.PresetInstrument{
					{
						Id:         1,
						Name:       "Kick",
						ChannelKey: "ch1",
						Controls: map[string]*model.PresetControl{
							"volume": {Name: "Volume", Type: "volume", Value: 0.90},
							"pan":    {Name: "Pan", Type: "pan", Value: 0.32},
						},
					},
					{
						Id:         2,
						Name:       "Snare",
						ChannelKey: "ch1",
						Controls: map[string]*model.PresetControl{
							"volume": {Name: "Volume", Type: "volume", Value: 0.85},
							"pan":    {Name: "Pan", Type: "pan", Value: 0.96},
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
						Volume: &pb.BaseControl{Key: "volume", Name: "Volume", Value: 1.00, Min: makeFloat64Ptr(0), Max: makeFloat64Ptr(1)},
						Pan:    &pb.BaseControl{Key: "pan", Name: "Pan", Value: 0.64, Min: makeFloat64Ptr(-1), Max: makeFloat64Ptr(1)},
						Instruments: []*pb.Instrument{
							{
								Key:    "1",
								Name:   "Kick",
								Volume: &pb.BaseControl{Key: "volume", Name: "Volume", Value: 0.90, Min: makeFloat64Ptr(0), Max: makeFloat64Ptr(1)},
								Pan:    &pb.BaseControl{Key: "pan", Name: "Pan", Value: 0.32, Min: makeFloat64Ptr(-1), Max: makeFloat64Ptr(1)},
							},
							{
								Key:    "2",
								Name:   "Snare",
								Volume: &pb.BaseControl{Key: "volume", Name: "Volume", Value: 0.85, Min: makeFloat64Ptr(0), Max: makeFloat64Ptr(1)},
								Pan:    &pb.BaseControl{Key: "pan", Name: "Pan", Value: 0.96, Min: makeFloat64Ptr(-1), Max: makeFloat64Ptr(1)},
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
