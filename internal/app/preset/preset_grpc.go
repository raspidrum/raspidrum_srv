package preset

import (
	"context"
	"io"
	"math"
	"strconv"

	pb "github.com/raspidrum-srv/internal/pkg/grpc"
	"github.com/spf13/afero"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/raspidrum-srv/internal/model"
	"github.com/raspidrum-srv/internal/repo"
	d "github.com/raspidrum-srv/internal/repo/db"
)

type PresetServer struct {
	pb.UnimplementedKitPresetServer
	pb.UnimplementedChannelControlServer
	db           *d.Sqlite
	sampler      repo.SamplerRepo
	ctrlHandler  *SamplerControlHandler
	fs           afero.Fs
	loadedPreset *model.KitPreset
}

func NewPresetServer(db *d.Sqlite, sampler repo.SamplerRepo, fs afero.Fs) *PresetServer {
	return &PresetServer{
		db:      db,
		sampler: sampler,
		fs:      fs,
	}
}

func (s *PresetServer) LoadPreset(ctx context.Context, req *pb.GetPresetRequest) (*pb.PresetResponse, error) {
	preset, chnls, err := LoadPreset(req.PresetId, s.db, s.sampler, s.fs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to load preset: %v", err)
	}

	s.loadedPreset = preset
	s.ctrlHandler = NewSamplerControlHandler(s.sampler, chnls)

	pbPreset, err := convertPresetToProto(preset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert preset: %v", err)
	}

	return &pb.PresetResponse{
		Preset: pbPreset,
	}, nil
}

func (s *PresetServer) GetPreset(ctx context.Context, req *pb.GetPresetRequest) (*pb.PresetResponse, error) {
	if s.loadedPreset == nil || s.loadedPreset.Id != req.PresetId {
		return s.LoadPreset(ctx, req)
	}

	pbPreset, err := convertPresetToProto(s.loadedPreset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert preset: %v", err)
	}

	return &pb.PresetResponse{
		Preset: pbPreset,
	}, nil
}

func (s *PresetServer) SetValue(stream grpc.BidiStreamingServer[pb.ControlValue, pb.ControlValue]) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		// set value
		err = s.loadedPreset.SetControlValue(in.Key, float32(in.Value), s.ctrlHandler)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to set control value: %v", err)
		}
		// send response
		out := &pb.ControlValue{
			Key:   in.Key,
			Seq:   in.Seq,
			Value: in.Value,
		}
		if err := stream.Send(out); err != nil {
			return err
		}
	}
}

// convertPresetToProto converts internal KitPreset model to protobuf Preset message
func convertPresetToProto(kitPreset *model.KitPreset) (*pb.Preset, error) {
	pbPreset := &pb.Preset{
		Id:   kitPreset.Id,
		Key:  kitPreset.Uid,
		Name: kitPreset.Name,
	}

	// Add global sampler channel first
	samplerChannel := &pb.Channel{
		// TODO: extract const "sampler"
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

		for ctrl := range ch.GetControls() {
			val, min, max := ctrl.GetNormalizedValue()
			if ctrl.Type == model.CtrlVolume {
				pbChannel.Volume = &pb.BaseControl{Key: string(ctrl.Key), Name: ctrl.Name, Value: roundFloat(float64(val), 3), Min: makeFloat64Ptr(min), Max: makeFloat64Ptr(max)}
			}
			if ctrl.Type == model.CtrlPan {
				pbChannel.Pan = &pb.BaseControl{Key: string(ctrl.Key), Name: ctrl.Name, Value: roundFloat(float64(val), 3), Min: makeFloat64Ptr(min), Max: makeFloat64Ptr(max)}
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

		for ctrl := range instr.GetControls() {
			val, min, max := ctrl.GetNormalizedValue()
			switch ctrl.Type {
			case model.CtrlVolume:
				pbInstrument.Volume = &pb.BaseControl{Key: ctrl.Key, Name: ctrl.Name, Value: roundFloat(float64(val), 3), Min: makeFloat64Ptr(min), Max: makeFloat64Ptr(max)}
			case model.CtrlPan:
				pbInstrument.Pan = &pb.BaseControl{Key: ctrl.Key, Name: ctrl.Name, Value: roundFloat(float64(val), 3), Min: makeFloat64Ptr(min), Max: makeFloat64Ptr(max)}
			default:
				// Convert other controls to tunes
				tune := &pb.FX{
					Key:  ctrl.Key,
					Name: ctrl.Name,
					// TODO: sort by control key
					Order: int32(len(pbInstrument.Tunes)),
					Params: []*pb.FXParam{{
						Key:   ctrl.Key,
						Name:  ctrl.Name,
						Type:  pb.FXParamType_FX_PARAM_TYPE_RANGE,
						Value: roundFloat(float64(val), 3),
						Min:   makeFloat64Ptr(min),
						Max:   makeFloat64Ptr(max),
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
			for ctrl := range layer.GetControls() {
				val, min, max := ctrl.GetNormalizedValue()
				switch ctrl.Type {
				case model.CtrlVolume:
					pbLayer.Volume = &pb.BaseControl{Key: ctrl.Key, Name: ctrl.Name, Value: roundFloat(float64(val), 3), Min: makeFloat64Ptr(min), Max: makeFloat64Ptr(max)}
				case model.CtrlPan:
					pbLayer.Pan = &pb.BaseControl{Key: ctrl.Key, Name: ctrl.Name, Value: roundFloat(float64(val), 3), Min: makeFloat64Ptr(min), Max: makeFloat64Ptr(max)}
				}
			}

			pbInstrument.Layers = append(pbInstrument.Layers, pbLayer)
		}

		res = append(res, pbInstrument)
	}
	return res
}

func makeFloat64Ptr(v float32) *float64 {
	res := float64(v)
	return &res
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
