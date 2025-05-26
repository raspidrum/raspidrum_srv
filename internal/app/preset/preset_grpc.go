package preset

import (
	"context"
	"math"
	"strconv"

	pb "github.com/raspidrum-srv/internal/pkg/grpc"
	"github.com/spf13/afero"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/raspidrum-srv/internal/model"
	"github.com/raspidrum-srv/internal/repo"
	d "github.com/raspidrum-srv/internal/repo/db"
)

type PresetServer struct {
	pb.UnimplementedKitPresetServer
	db      *d.Sqlite
	sampler repo.SamplerRepo
	fs      afero.Fs
}

func NewPresetServer(db *d.Sqlite, sampler repo.SamplerRepo, fs afero.Fs) *PresetServer {
	return &PresetServer{
		db:      db,
		sampler: sampler,
		fs:      fs,
	}
}

func (s *PresetServer) LoadPreset(ctx context.Context, req *pb.LoadPresetRequest) (*pb.LoadPresetResponse, error) {
	preset, err := LoadPreset(req.PresetId, s.db, s.sampler, s.fs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to load preset: %v", err)
	}

	pbPreset, err := convertPresetToProto(preset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert preset: %v", err)
	}

	return &pb.LoadPresetResponse{
		Preset: pbPreset,
	}, nil
}

// convertPresetToProto converts internal KitPreset model to protobuf Preset message
func convertPresetToProto(kitPreset *model.KitPreset) (*pb.Preset, error) {
	pbPreset := &pb.Preset{
		Key:  kitPreset.Uid,
		Name: kitPreset.Name,
	}

	// Add global sampler channel first
	samplerChannel := &pb.Channel{
		// TODO: extract const "sampler"
		Key:  "sampler",
		Name: "Kit",
		Type: pb.ChannelType_CHANNEL_TYPE_SAMPLER,
	}
	pbPreset.Channels = append(pbPreset.Channels, samplerChannel)

	// Convert instrument channels
	for _, ch := range kitPreset.Channels {
		pbChannel := &pb.Channel{
			Key:  ch.Key,
			Name: ch.Name,
			Type: pb.ChannelType_CHANNEL_TYPE_INSTRUMENT,
		}

		// Get instruments in this channel
		instruments, err := kitPreset.GetChannelInstrumentsByKey(ch.Key)
		if err != nil {
			return nil, err
		}

		// Handle volume and pan based on channel configuration
		if len(instruments) == 1 {

			// Use instrument controls for single instrument with MIDI CC
			if vol, ok := instruments[0].Controls[model.CtrlVolume]; ok {
				if vol.MidiCC != 0 {
					val, min, max := normalizeVol(vol)
					pbChannel.Volume = &pb.BaseControl{Key: model.CtrlVolume, Name: vol.Name, Value: val, Min: makeFloat64Ptr(min), Max: makeFloat64Ptr(max)}

				} else {
					val, min, max := normalizeVol(ch.Controls[model.CtrlVolume])
					pbChannel.Volume = &pb.BaseControl{Key: model.CtrlVolume, Name: vol.Name, Value: val, Min: makeFloat64Ptr(min), Max: makeFloat64Ptr(max)}
				}
			}
			if pan, ok := instruments[0].Controls[model.CtrlPan]; ok {
				if pan.MidiCC != 0 {
					val, min, max := normalizePan(pan)
					pbChannel.Pan = &pb.BaseControl{Key: model.CtrlPan, Name: pan.Name, Value: val, Min: makeFloat64Ptr(min), Max: makeFloat64Ptr(max)}
				} else {
					if chpan, ok := ch.Controls[model.CtrlPan]; ok {
						val, min, max := normalizePan(chpan)
						pbChannel.Pan = &pb.BaseControl{Key: model.CtrlPan, Name: chpan.Name, Value: val, Min: makeFloat64Ptr(min), Max: makeFloat64Ptr(max)}
					}
				}
			}
		} else {
			// Use channel controls
			if vol, ok := ch.Controls[model.CtrlVolume]; ok {
				val, min, max := normalizeVol(vol)
				pbChannel.Volume = &pb.BaseControl{Key: model.CtrlVolume, Name: vol.Name, Value: val, Min: makeFloat64Ptr(min), Max: makeFloat64Ptr(max)}
			}
			if pan, ok := ch.Controls[model.CtrlPan]; ok {
				val, min, max := normalizePan(pan)
				pbChannel.Pan = &pb.BaseControl{Key: model.CtrlPan, Name: pan.Name, Value: val, Min: makeFloat64Ptr(min), Max: makeFloat64Ptr(max)}
			}
		}

		// Convert instruments
		pbChannel.Instruments = convertInstrumentToProto(instruments)

		pbPreset.Channels = append(pbPreset.Channels, pbChannel)
	}

	return pbPreset, nil
}

func convertInstrumentToProto(instruments []*model.PresetInstrument) []*pb.Instrument {
	res := make([]*pb.Instrument, 0)
	for _, instr := range instruments {
		pbInstrument := &pb.Instrument{
			Key:  strconv.FormatInt(instr.Id, 10),
			Name: instr.Name,
		}

		// Convert other controls to tunes
		for key, ctrl := range instr.Controls {
			if ctrl.Type == model.CtrlVolume || ctrl.Type == model.CtrlPan {
				// Add volume and pan for instrument only if channel has multiple instruments
				if len(instruments) > 1 {
					if cv, ok := instr.Controls[ctrl.Type]; ok {
						if ctrl.Type == model.CtrlVolume {
							val, min, max := normalizeVol(cv)
							pbInstrument.Volume = &pb.BaseControl{Key: ctrl.Type, Name: cv.Name, Value: val, Min: makeFloat64Ptr(min), Max: makeFloat64Ptr(max)}
						} else {
							val, min, max := normalizePan(cv)
							pbInstrument.Pan = &pb.BaseControl{Key: ctrl.Type, Name: cv.Name, Value: val, Min: makeFloat64Ptr(min), Max: makeFloat64Ptr(max)}
						}
					}
				}
			} else {
				// TODO: make unique control key across all instruments and store in KitPreset
				tune := &pb.FX{
					Key:  key,
					Name: ctrl.Name,
					// TODO: sort by control key
					Order: int32(len(pbInstrument.Tunes)),
					Params: []*pb.FXParam{{
						Key:   key,
						Name:  ctrl.Name,
						Type:  pb.FXParamType_FX_PARAM_TYPE_RANGE,
						Value: float64(ctrl.Value),
						// TODO: detect class of contol: midiCC or sendFX. For midiCC set min and max as const 0 and 127
						Min: makeFloat64Ptr(0),
						Max: makeFloat64Ptr(127),
					}},
				}
				pbInstrument.Tunes = append(pbInstrument.Tunes, tune)
			}
		}

		// Convert layers
		for key, layer := range instr.Layers {
			pbLayer := &pb.Layer{
				Key:  key,
				Name: layer.Name,
			}

			if vol, ok := layer.Controls[model.CtrlVolume]; ok {
				val, min, max := normalizeVol(vol)
				pbLayer.Volume = &pb.BaseControl{Key: model.CtrlVolume, Name: vol.Name, Value: val, Min: makeFloat64Ptr(min), Max: makeFloat64Ptr(max)}
			}
			if pan, ok := layer.Controls[model.CtrlPan]; ok {
				val, min, max := normalizePan(pan)
				pbLayer.Pan = &pb.BaseControl{Key: model.CtrlPan, Name: pan.Name, Value: val, Min: makeFloat64Ptr(min), Max: makeFloat64Ptr(max)}
			}

			pbInstrument.Layers = append(pbInstrument.Layers, pbLayer)
		}

		res = append(res, pbInstrument)
	}
	return res
}

func normalizeVol(ctrl model.PresetControl) (val float64, min float64, max float64) {
	if ctrl.MidiCC != 0 {
		// val from 0 to 1 with 3 decimal places
		return roundFloat(float64(ctrl.Value/127), 3), 0, 1
	}
	return roundFloat(float64(ctrl.Value), 3), 0, 1
}

func normalizePan(ctrl model.PresetControl) (val float64, min float64, max float64) {
	if ctrl.MidiCC != 0 {
		return roundFloat(float64((ctrl.Value*2/127)-1), 3), -1, 1
	}
	return roundFloat(float64(ctrl.Value), 3), -1, 1
}

func makeFloat64Ptr(v float64) *float64 {
	return &v
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
