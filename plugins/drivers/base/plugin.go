package base

import (
	"context"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/nomad/plugins/drivers/base/proto"
	"google.golang.org/grpc"
)

func LaunchDriver(d DriverPlugin) error {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{},
		Plugins: map[string]plugin.Plugin{
			DriverGoPlugin: &PluginDriver{impl: d},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
	return nil
}

type ServeConfig struct {
	EventRecorderRPCAddr string
	HandshakeConfig      *plugin.HandshakeConfig
	GRPCServer           func([]grpc.ServerOption) *grpc.Server
}

func Serve(cfg ServeConfig) error {

	return nil
}

type DriverFactory func(interface{}) DriverPlugin

// PluginDriver wraps a DriverPlugin and implements go-plugins GRPCPlugin
// interface to expose the the interface over gRPC
type PluginDriver struct {
	plugin.NetRPCUnsupportedPlugin
	impl DriverPlugin
}

func NewDriverPlugin(d DriverPlugin) plugin.GRPCPlugin {
	return &PluginDriver{impl: d}
}

func (p *PluginDriver) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterDriverServer(s, &driverPluginServer{
		impl:   p.impl,
		broker: broker,
	})
	return nil
}

func (p *PluginDriver) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &driverPluginClient{client: proto.NewDriverClient(c)}, nil
}
