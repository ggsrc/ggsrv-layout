package {{ cookiecutter.project_name.strip().replace("-", "_") }}

import (
	"context"

	"google.golang.org/grpc"

	"{{ cookiecutter.go_module_path.strip('/') }}/api/go/client/conn"
	pb "{{ cookiecutter.go_module_path.strip('/') }}/api/go/{{ cookiecutter.project_name.strip() }}/v1"
)

// Client defines a wrapper for grpc client
type Client interface {
	pb.GServiceClient
}

type gClient struct {
	pb.GServiceClient
	conn *grpc.ClientConn
}

func NewClient(ctx context.Context, conf *conn.Config) (Client, error) {
	c, err := conn.Dial(ctx, conf)
	if err != nil {
		return nil, err
	}
	gCli := pb.NewGServiceClient(c)
	return &gClient{gCli, c}, nil
}
