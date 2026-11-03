package conn

import (
	"context"

	ggGrpc "github.com/ggsrc/gglib/grpc"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Config defines new client options
type Config struct {
	Addr        string
	DialOptions []grpc.DialOption
	GGOptions   []ggGrpc.ClientOption
}

func Dial(ctx context.Context, c *Config) (*grpc.ClientConn, error) {
	if c == nil {
		c = &Config{}
		envconfig.MustProcess("mtest", c)
	}
	// When c.Addr is empty, get the address from your service discovery, for
	// example: c.Addr = k8s.NewServiceAddrs().GetAddr(k8s.Xxx)
	if c.DialOptions == nil {
		c.DialOptions = append(
			c.DialOptions,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	}
	dialer := ggGrpc.NewClientWithOptions(c.GGOptions...)
	conn, err := dialer.Dial(ctx, c.Addr, c.DialOptions...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
