package gateway

import (
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/go-redis/redis"

	pb "github.com/grpc-streamer/proto"
)

const (
	//callStatusChannel はcall statusのpub sub keyです
	callStatusChannel = "callStatusChannel_"
	callStatusKey     = "callStatusKey_"
)

//CallStatusRedisGateway は，topicとgenreで使用するbaseのファンミーティングリストです
type CallStatusRedisGateway struct {
	redis *redis.Client
}

//NewBaseRankinggRedisGateway is
func NewCallStatusRedisGateway(redis *redis.Client) *CallStatusRedisGateway {
	return &CallStatusRedisGateway{
		redis: redis,
	}
}

//Subscribe is redis string to delete
func (sv *CallStatusRedisGateway) Subscribe(influencerUUID string) (*pb.GetCallStatusResponse, error) {
	subscriber := sv.redis.Subscribe(callStatusChannel + influencerUUID)

	message, err := subscriber.ReceiveMessage()
	if err != nil {
		return nil, toGatewayError(err)
	}

	var callStatusResponse pb.GetCallStatusResponse

	err = proto.Unmarshal([]byte(message.Payload), &callStatusResponse)
	if err != nil {
		return nil, toGatewayError(err)
	}

	return &callStatusResponse, nil
}

//Subscribe is redis string to delete
func (sv *CallStatusRedisGateway) Get(influencerUUID string) (*pb.GetCallStatusResponse, error) {
	bytes, err := sv.redis.Get(callStatusKey + influencerUUID).Bytes()
	if err != nil {
		return nil, toGatewayError(err)
	}

	var callStatusResponse pb.GetCallStatusResponse

	err = proto.Unmarshal(bytes, &callStatusResponse)
	if err != nil {
		return nil, toGatewayError(err)
	}

	return &callStatusResponse, nil
}

func (sv *CallStatusRedisGateway) Publish(influencerUUID string, pbCallStatus *pb.GetCallStatusResponse) error {
	callStatus, err := proto.Marshal(pbCallStatus)
	if err != nil {
		return toGatewayError(err)
	}
	if err := sv.redis.Publish(callStatusChannel+influencerUUID, callStatus).Err(); err != nil {
		return toGatewayError(err)
	}

	if err := sv.redis.Set(callStatusKey+influencerUUID, callStatus, time.Hour*1).Err(); err != nil {
		return toGatewayError(err)
	}

	return nil
}
