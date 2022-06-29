package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"wbroker/app/broker"
	grpcWBroker "wbroker/external/grpc/wbroker"
)

type (
	// Server is a grpc server wbroker
	Server struct {
		grpcWBroker.UnimplementedWBrokerServer
		broker broker.Service
	}
)

// NewServer returns new instance of Server
func NewServer(b broker.Service) *Server {
	return &Server{broker: b}
}

func (s *Server) Publish(ctx context.Context, in *grpcWBroker.PublishRequest) (*grpcWBroker.PublishResponse, error) {
	message := broker.Message{
		Body: string(in.GetBody()),
	}

	id, err := s.broker.Publish(ctx, in.GetTopic(), message)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, fmt.Errorf("failed to publish message: %w", err).Error())
	}

	response := &grpcWBroker.PublishResponse{Id: id}
	return response, nil
}

func (s *Server) Subscribe(in *grpcWBroker.SubscribeRequest, stream grpcWBroker.WBroker_SubscribeServer) error {
	msgC, err := s.broker.Subscribe(stream.Context(), in.GetTopic())
	if err != nil {
		return status.Errorf(codes.Unavailable, fmt.Errorf("failed to subscribe: %w", err).Error())
	}

	for msg := range msgC {
		message := &grpcWBroker.Message{Body: []byte(msg.Body)}
		err := stream.Send(message)
		if err != nil {
			return status.Errorf(codes.Unknown, fmt.Errorf("failed to send message: %w", err).Error())
		}
	}

	return nil
}
