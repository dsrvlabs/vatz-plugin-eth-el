package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/sdk"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	defaultAddr                string = "127.0.0.1"
	defaultPort                int    = 10003
	defaultRpcAddr             string = ""
	defaultBlockIntervalSecond int    = 12

	minimumSamples = 10
)

type syncStamp struct {
	Time        time.Time
	BlockNumber uint64
}

var (
	addr                string
	port                int
	rpcAddr             string
	blockIntervalSecond int

	pluginName = "block-sync"

	syncHistory = make([]syncStamp, 0)
)

func main() {
	flag.StringVar(&addr, "addr", defaultAddr, "address")
	flag.IntVar(&port, "port", defaultPort, "port number")
	flag.StringVar(&rpcAddr, "rpcAddr", defaultRpcAddr, "rpc address")
	flag.IntVar(&blockIntervalSecond, "blocktime", defaultBlockIntervalSecond, "block time in second")

	flag.Parse()

	log.Info().Str("module", pluginName).Msgf("addr %s", addr)
	log.Info().Str("module", pluginName).Msgf("port %d", port)
	log.Info().Str("module", pluginName).Msgf("rpcAddr %s", rpcAddr)
	log.Info().Str("module", pluginName).Msgf("blockInterval %d", blockIntervalSecond)

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
	blockNumber, err := cli.BlockNumber(ctx)
	if err != nil {
		return sdk.CallResponse{
			FuncName: "pluginFeature",
			Message:  err.Error(),
			Severity: pluginpb.SEVERITY_CRITICAL,
			State:    pluginpb.STATE_FAILURE,
		}, err
	}

	newStamp := syncStamp{
		Time:        time.Now(),
		BlockNumber: blockNumber,
	}

	syncHistory = append(syncHistory, newStamp)

	if len(syncHistory) < minimumSamples {
		return sdk.CallResponse{
			FuncName: "pluginFeature",
			Message:  "calculating...",
			Severity: pluginpb.SEVERITY_CRITICAL,
			State:    pluginpb.STATE_PENDING,
		}, nil
	} else {
		syncHistory = syncHistory[len(syncHistory)-10:]
	}

	timeDiff := syncHistory[len(syncHistory)-1].Time.Sub(syncHistory[0].Time)
	blockDiff := syncHistory[len(syncHistory)-1].BlockNumber - syncHistory[0].BlockNumber

	curRate := float64(blockDiff) / timeDiff.Seconds()
	normalRate := 1.0 / float64(blockIntervalSecond)

	log.Info().Str("module", pluginName).Msgf("pluginFeature: rate %f / %f", curRate, normalRate)

	const stallThreshold = 0.1

	if curRate < normalRate*stallThreshold {
		return sdk.CallResponse{
			FuncName: "pluginFeature",
			Message:  "block sync is stalled",
			Severity: pluginpb.SEVERITY_CRITICAL,
			State:    pluginpb.STATE_FAILURE,
		}, nil
	}

	return sdk.CallResponse{
		FuncName: "pluginFeature",
		Message:  "syncing",
		Severity: pluginpb.SEVERITY_INFO,
		State:    pluginpb.STATE_SUCCESS,
	}, nil
}
