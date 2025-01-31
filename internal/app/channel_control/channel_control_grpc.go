package channelcontrol

import (
	"io"

	pb "github.com/raspidrum-srv/internal/pkg/grpc"
)

// TODO: добавить сюда UseCase, к-й будет инжектиться в конструкторе и вызываться из main.go при старте приложения
type ChannelControlServer struct {
	pb.UnimplementedChannelControlServer
}

func NewChannelControlServer() *ChannelControlServer {
	return &ChannelControlServer{}
}

func (s *ChannelControlServer) SetValue(stream pb.ChannelControl_SetValueServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		// TODO: здесь должна быть логика: маппинг на физический контрол (семплера или jack), реализуемая через usecase. Т.е. здесь вызов useCase
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
