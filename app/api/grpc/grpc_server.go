package grpc

import (
	"context"
	"fmt"
	wb "github.com/AlimByWater/wbroker/external/grpc/wbroker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"wbroker/app/broker"
)

// TODO add explicit logging

type (
	// Server is a grpc server wbroker
	Server struct {
		wb.UnimplementedWBrokerServer
		broker broker.Service
	}
)

// NewServer returns new instance of Server
func NewServer(b broker.Service) *Server {
	return &Server{
		broker: b,
	}
}

// Publish publishes a message to a particular topic
func (s *Server) Publish(ctx context.Context, in *wb.PublishRequest) (*wb.PublishResponse, error) {
	message := broker.Message{
		Body: in.GetBody(),
	}

	err := message.Validate()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Errorf("failed to validate: %w", err).Error())
	}

	id, err := s.broker.Publish(ctx, in.GetTopic(), message)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, fmt.Errorf("failed to publish message: %w", err).Error())
	}

	response := &wb.PublishResponse{Id: id}
	return response, nil
}

// Subscribe subscribes a client to a particular topic and starting to stream messages
func (s *Server) Subscribe(in *wb.SubscribeRequest, stream wb.WBroker_SubscribeServer) error {
	msgC, err := s.broker.Subscribe(stream.Context(), in.GetTopic())
	if err != nil {
		return status.Errorf(codes.Unavailable, fmt.Errorf("failed to subscribe: %w", err).Error())
	}

	for msg := range msgC {
		message := &wb.Message{Body: msg.Body}
		err := stream.Send(message)
		if err != nil {
			return status.Errorf(codes.Unknown, fmt.Errorf("failed to send message: %w", err).Error())
		}
	}

	return nil
}
