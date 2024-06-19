package main

import (
	"context"
	"flag"
	"fmt"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/sdk"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	defaultAddr    string = "127.0.0.1"
	defaultPort    int    = 10002
	defaultRpcAddr string = ""
)

var (
	addr    string
	port    int
	rpcAddr string

	pluginName = "node-is-alived"
)

func main() {
	flag.StringVar(&addr, "addr", defaultAddr, "address")
	flag.IntVar(&port, "port", defaultPort, "port number")
	flag.StringVar(&rpcAddr, "rpcAddr", defaultRpcAddr, "rpc address")

	flag.Parse()

	p := sdk.NewPlugin(pluginName)
	p.Register(pluginFeature)

	ctx := context.Background()
	if err := p.Start(ctx, addr, port); err != nil {
		fmt.Println("exit")
	}
}

func pluginFeature(info, opt map[string]*structpb.Value) (sdk.CallResponse, error) {
	log.Info().Str("module", pluginName).Msg("pluginFeature")

	cli, err := ethclient.Dial(rpcAddr)
	if err != nil {
		return sdk.CallResponse{
			FuncName: "pluginFeature",
			Message:  err.Error(),
			Severity: pluginpb.SEVERITY_CRITICAL,
			State:    pluginpb.STATE_FAILURE,
		}, err
	}

	ctx := context.Background()
	syncProgress, err := cli.SyncProgress(ctx)
	if err != nil {
		return sdk.CallResponse{
			FuncName: "pluginFeature",
			Message:  err.Error(),
			Severity: pluginpb.SEVERITY_CRITICAL,
			State:    pluginpb.STATE_FAILURE,
		}, err
	}

	log.Info().Str("module", pluginName).Msgf("pluginFeature: SyncProgress %v", syncProgress)

	return sdk.CallResponse{
		FuncName: "pluginFeature",
		Message:  "",
		Severity: pluginpb.SEVERITY_INFO,
		State:    pluginpb.STATE_SUCCESS,
	}, nil
}
