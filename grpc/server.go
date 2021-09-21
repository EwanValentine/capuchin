package grpc

import (
	"context"
	"encoding/json"
	"net"

	"github.com/EwanValentine/capuchin/conf"
	gw "github.com/EwanValentine/capuchin/gen/go/proto"
	query "github.com/EwanValentine/capuchin/query"
	"github.com/EwanValentine/capuchin/source"
	"google.golang.org/grpc"
)

// QueryEngine -
type QueryEngine interface {
	Exec() ([]query.Result, error)
}

// Server -
type Server struct {
	*gw.UnimplementedCapuchinQueryServiceServer
	QueryEngine QueryEngine
}

// Query -
func (s *Server) Query(ctx context.Context, req *gw.QueryRequest) (*gw.QueryResponse, error) {
	fileSource, err := source.NewFileSource().Load(req.Source)
	if err != nil {
		return nil, err
	}

	q := query.Query{
		Where:  req.Where,
		Select: req.Select,
	}
	q.Source(fileSource)

	results, err := q.Exec()
	if err != nil {
		return nil, err
	}

	d, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	return &gw.QueryResponse{
		Result: d,
	}, nil
}

// NewService -
func NewServer(conf *conf.Config) error {
	lis, err := net.Listen("tcp", conf.GRPCAddr)
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)

	gw.RegisterCapuchinQueryServiceServer(s, &Server{})

	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}
