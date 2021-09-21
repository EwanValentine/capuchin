package grpc

import (
	"context"
	"encoding/json"
	"log"
	"net"

	"github.com/EwanValentine/capuchin/conf"
	gw "github.com/EwanValentine/capuchin/gen/go/proto"
	query "github.com/EwanValentine/capuchin/query"
	"github.com/EwanValentine/capuchin/source"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
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
	log.Println("got here...")
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

	se := []*structpb.Struct{}
	for _, res := range results {
		e := &structpb.Struct{}
		b, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}

		if err := protojson.Unmarshal(b, e); err != nil {
			return nil, err
		}

		se = append(se, e)
	}

	return &gw.QueryResponse{
		Results: se,
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

	return s.Serve(lis)
}
