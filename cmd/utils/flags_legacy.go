















package utils

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/node"
	"gopkg.in/urfave/cli.v1"
)

var ShowDeprecated = cli.Command{
	Action:      showDeprecated,
	Name:        "show-deprecated-flags",
	Usage:       "Show flags that have been deprecated",
	ArgsUsage:   " ",
	Category:    "MISCELLANEOUS COMMANDS",
	Description: "Show flags that have been deprecated and will soon be removed",
}

var DeprecatedFlags = []cli.Flag{
	LegacyTestnetFlag,
	LegacyLightServFlag,
	LegacyLightPeersFlag,
	LegacyMinerThreadsFlag,
	LegacyMinerGasTargetFlag,
	LegacyMinerGasPriceFlag,
	LegacyMinerEtherbaseFlag,
	LegacyMinerExtraDataFlag,
}

var (
	
	LegacyMinerThreadsFlag = cli.IntFlag{
		Name:  "minerthreads",
		Usage: "Number of CPU threads to use for mining (deprecated, use --miner.threads)",
		Value: 0,
	}
	LegacyMinerGasTargetFlag = cli.Uint64Flag{
		Name:  "targetgaslimit",
		Usage: "Target gas floor for mined blocks (deprecated, use --miner.gastarget)",
		Value: eth.DefaultConfig.Miner.GasFloor,
	}
	LegacyMinerGasPriceFlag = BigFlag{
		Name:  "gasprice",
		Usage: "Minimum gas price for mining a transaction (deprecated, use --miner.gasprice)",
		Value: eth.DefaultConfig.Miner.GasPrice,
	}
	LegacyMinerEtherbaseFlag = cli.StringFlag{
		Name:  "etherbase",
		Usage: "Public address for block mining rewards (default = first account, deprecated, use --miner.etherbase)",
		Value: "0",
	}
	LegacyMinerExtraDataFlag = cli.StringFlag{
		Name:  "extradata",
		Usage: "Block extra data set by the miner (default = client version, deprecated, use --miner.extradata)",
	}

	
	LegacyLightServFlag = cli.IntFlag{
		Name:  "lightserv",
		Usage: "Maximum percentage of time allowed for serving LES requests (deprecated, use --light.serve)",
		Value: eth.DefaultConfig.LightServ,
	}
	LegacyLightPeersFlag = cli.IntFlag{
		Name:  "lightpeers",
		Usage: "Maximum number of light clients to serve, or light servers to attach to  (deprecated, use --light.maxpeers)",
		Value: eth.DefaultConfig.LightPeers,
	}

	
	LegacyTestnetFlag = cli.BoolFlag{ 
		Name:  "testnet",
		Usage: "Pre-configured test network (Deprecated: Please choose one of --goerli, --rinkeby, or --ropsten.)",
	}

	
	LegacyRPCEnabledFlag = cli.BoolFlag{
		Name:  "rpc",
		Usage: "Enable the HTTP-RPC server (deprecated, use --http)",
	}
	LegacyRPCListenAddrFlag = cli.StringFlag{
		Name:  "rpcaddr",
		Usage: "HTTP-RPC server listening interface (deprecated, use --http.addr)",
		Value: node.DefaultHTTPHost,
	}
	LegacyRPCPortFlag = cli.IntFlag{
		Name:  "rpcport",
		Usage: "HTTP-RPC server listening port (deprecated, use --http.port)",
		Value: node.DefaultHTTPPort,
	}
	LegacyRPCCORSDomainFlag = cli.StringFlag{
		Name:  "rpccorsdomain",
		Usage: "Comma separated list of domains from which to accept cross origin requests (browser enforced) (deprecated, use --http.corsdomain)",
		Value: "",
	}
	LegacyRPCVirtualHostsFlag = cli.StringFlag{
		Name:  "rpcvhosts",
		Usage: "Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard. (deprecated, use --http.vhosts)",
		Value: strings.Join(node.DefaultConfig.HTTPVirtualHosts, ","),
	}
	LegacyRPCApiFlag = cli.StringFlag{
		Name:  "rpcapi",
		Usage: "API's offered over the HTTP-RPC interface (deprecated, use --http.api)",
		Value: "",
	}
	LegacyWSListenAddrFlag = cli.StringFlag{
		Name:  "wsaddr",
		Usage: "WS-RPC server listening interface (deprecated, use --ws.addr)",
		Value: node.DefaultWSHost,
	}
	LegacyWSPortFlag = cli.IntFlag{
		Name:  "wsport",
		Usage: "WS-RPC server listening port (deprecated, use --ws.port)",
		Value: node.DefaultWSPort,
	}
	LegacyWSApiFlag = cli.StringFlag{
		Name:  "wsapi",
		Usage: "API's offered over the WS-RPC interface (deprecated, use --ws.api)",
		Value: "",
	}
	LegacyWSAllowedOriginsFlag = cli.StringFlag{
		Name:  "wsorigins",
		Usage: "Origins from which to accept websockets requests (deprecated, use --ws.origins)",
		Value: "",
	}
	LegacyGpoBlocksFlag = cli.IntFlag{
		Name:  "gpoblocks",
		Usage: "Number of recent blocks to check for gas prices (deprecated, use --gpo.blocks)",
		Value: eth.DefaultConfig.GPO.Blocks,
	}
	LegacyGpoPercentileFlag = cli.IntFlag{
		Name:  "gpopercentile",
		Usage: "Suggested gas price is the given percentile of a set of recent transaction gas prices (deprecated, use --gpo.percentile)",
		Value: eth.DefaultConfig.GPO.Percentile,
	}
	LegacyBootnodesV4Flag = cli.StringFlag{
		Name:  "bootnodesv4",
		Usage: "Comma separated enode URLs for P2P v4 discovery bootstrap (light server, full nodes) (deprecated, use --bootnodes)",
		Value: "",
	}
	LegacyBootnodesV5Flag = cli.StringFlag{
		Name:  "bootnodesv5",
		Usage: "Comma separated enode URLs for P2P v5 discovery bootstrap (light server, light nodes) (deprecated, use --bootnodes)",
		Value: "",
	}

	
	LegacyGraphQLListenAddrFlag = cli.StringFlag{
		Name:  "graphql.addr",
		Usage: "GraphQL server listening interface (deprecated, graphql can only be enabled on the HTTP-RPC server endpoint, use --graphql)",
	}
	LegacyGraphQLPortFlag = cli.IntFlag{
		Name:  "graphql.port",
		Usage: "GraphQL server listening port (deprecated, graphql can only be enabled on the HTTP-RPC server endpoint, use --graphql)",
		Value: node.DefaultHTTPPort,
	}
)


func showDeprecated(*cli.Context) {
	fmt.Println("--------------------------------------------------------------------")
	fmt.Println("The following flags are deprecated and will be removed in the future!")
	fmt.Println("--------------------------------------------------------------------")
	fmt.Println()

	for _, flag := range DeprecatedFlags {
		fmt.Println(flag.String())
	}
}
