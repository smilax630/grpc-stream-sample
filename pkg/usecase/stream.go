package usecase

import (
	"github.com/grpc-streamer/pkg/gateway"
	pb "github.com/grpc-streamer/proto"
)

// StreamUseCase is
type StreamUseCase struct {
	callStatusGateway gateway.CallStatusRedisGateway
}

// NewStreamUsecase return StreamUseCase
func NewStreamUsecase(callStatusGateway gateway.CallStatusRedisGateway) *StreamUseCase {
	return &StreamUseCase{callStatusGateway: callStatusGateway}
}

//GetCallStatus is ...
func (i *StreamUseCase) GetCallStatus(in *pb.GetCallStatusRequest, stream pb.StreamService_GetCallStatusServer) error {
	// 初回でも値を返すため
	callStatus, err := i.callStatusGateway.Get(in.GetInfluencerUuid())
	if err != nil {
		return err
	}
	err = stream.Send(callStatus)
	if err != nil {
		return err
	}

	for {
		callStatus, err := i.callStatusGateway.Subscribe(in.GetInfluencerUuid())
		if err != nil {
			return err
		}
		err = stream.Send(callStatus)
		if err != nil {
			return err
		}
	}
}
