package repo

import (
	"encoding/json"
	"errors"
	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/repo/config"
	"path"
)

var DefaultBootstrapAddresses = []string{
	"/ip4/107.170.133.32/tcp/4001/ipfs/QmbY4yo9Eifg7DPjL7qK5JvNdiJaRAD7N76gVg4YoQsvgA", // Le Marché Serpette
	"/ip4/139.59.174.197/tcp/4001/ipfs/QmcCoBtYyduyurcLHRF14QhhA88YojJJpGFuMHoMZuU8sc", // Brixton-Village
	"/ip4/139.59.6.222/tcp/4001/ipfs/QmdzzGGc9xZq8w4z42vSHe32DZM7VXfDUFEUyfPvYNYhXE",   // Johari
	"/ip4/46.101.198.170/tcp/4001/ipfs/QmePWxsFT9wY3QuukgVDB7XZpqdKhrqJTHTXU7ECLDWJqX", // Duo Search
}

type APIConfig struct {
	Authenticated bool
	AllowedIPs    []string
	Username      string
	Password      string
	CORS          *string
	Enabled       bool
	HTTPHeaders   map[string]interface{}
	SSL           bool
	SSLCert       string
	SSLKey        string
}

type TorConfig struct {
	Password   string
	TorControl string
}

type WalletConfig struct {
	Type             string
	Binary           string
	MaxFee           int
	FeeAPI           string
	HighFeeDefault   int
	MediumFeeDefault int
	LowFeeDefault    int
	TrustedPeer      string
	RPCUser          string
	RPCPassword      string
}

var MalformedConfigError error = errors.New("Config file is malformed")

func GetAPIConfig(cfgBytes []byte) (*APIConfig, error) {
	var cfgIface interface{}
	json.Unmarshal(cfgBytes, &cfgIface)

	cfg, ok := cfgIface.(map[string]interface{})
	if !ok {
		return nil, MalformedConfigError
	}

	apiIface, ok := cfg["JSON-API"]
	if !ok {
		return nil, MalformedConfigError
	}

	api, ok := apiIface.(map[string]interface{})
	if !ok {
		return nil, MalformedConfigError
	}

	headers := make(map[string]interface{})
	h, ok := api["HTTPHeaders"]
	if h == nil || !ok {
		headers = nil
	} else {
		headers, ok = h.(map[string]interface{})
		if !ok {
			return nil, MalformedConfigError
		}
	}

	enabled, ok := api["Enabled"]
	if !ok {
		return nil, MalformedConfigError
	}
	enabledBool, ok := enabled.(bool)
	if !ok {
		return nil, MalformedConfigError
	}
	authenticated := api["Authenticated"]
	if !ok {
		return nil, MalformedConfigError
	}
	authenticatedBool, ok := authenticated.(bool)
	if !ok {
		return nil, MalformedConfigError
	}
	allowedIPs, ok := api["AllowedIPs"]
	if !ok {
		return nil, MalformedConfigError
	}
	allowedIPsIface, ok := allowedIPs.([]interface{})
	if !ok {
		return nil, MalformedConfigError
	}
	var allowedIPstrings []string
	for _, ip := range allowedIPsIface {
		ipStr, ok := ip.(string)
		if !ok {
			return nil, MalformedConfigError
		}
		allowedIPstrings = append(allowedIPstrings, ipStr)
	}

	username, ok := api["Username"]
	if !ok {
		return nil, MalformedConfigError
	}
	usernameStr, ok := username.(string)
	if !ok {
		return nil, MalformedConfigError
	}

	password, ok := api["Password"]
	if !ok {
		return nil, MalformedConfigError
	}
	passwordStr, ok := password.(string)
	if !ok {
		return nil, MalformedConfigError
	}

	c, ok := api["CORS"]
	var cors *string
	if c == nil || !ok {
		cors = nil
	} else {
		crs, ok := c.(string)
		if !ok {
			return nil, MalformedConfigError
		}
		cors = &crs
	}
	sslEnabled, ok := api["SSL"]
	if !ok {
		return nil, MalformedConfigError
	}
	sslEnabledBool, ok := sslEnabled.(bool)
	if !ok {
		return nil, MalformedConfigError
	}

	certFile, ok := api["SSLCert"]
	if !ok {
		return nil, MalformedConfigError
	}
	certFileStr, ok := certFile.(string)
	if !ok {
		return nil, MalformedConfigError
	}
	keyFile, ok := api["SSLKey"]
	if !ok {
		return nil, MalformedConfigError
	}
	keyFileStr, ok := keyFile.(string)
	if !ok {
		return nil, MalformedConfigError
	}

	apiConfig := &APIConfig{
		Authenticated: authenticatedBool,
		AllowedIPs:    allowedIPstrings,
		Username:      usernameStr,
		Password:      passwordStr,
		CORS:          cors,
		Enabled:       enabledBool,
		HTTPHeaders:   headers,
		SSL:           sslEnabledBool,
		SSLCert:       certFileStr,
		SSLKey:        keyFileStr,
	}

	return apiConfig, nil
}

func GetWalletConfig(cfgBytes []byte) (*WalletConfig, error) {
	var cfgIface interface{}
	json.Unmarshal(cfgBytes, &cfgIface)
	cfg, ok := cfgIface.(map[string]interface{})
	if !ok {
		return nil, MalformedConfigError
	}

	walletIface, ok := cfg["Wallet"]
	if !ok {
		return nil, MalformedConfigError
	}
	wallet, ok := walletIface.(map[string]interface{})
	if !ok {
		return nil, MalformedConfigError
	}
	feeAPI, ok := wallet["FeeAPI"]
	if !ok {
		return nil, MalformedConfigError
	}
	feeAPIstr, ok := feeAPI.(string)
	if !ok {
		return nil, MalformedConfigError
	}
	trustedPeer, ok := wallet["TrustedPeer"]
	if !ok {
		return nil, MalformedConfigError
	}
	trustedPeerStr, ok := trustedPeer.(string)
	if !ok {
		return nil, MalformedConfigError
	}
	low, ok := wallet["LowFeeDefault"]
	if !ok {
		return nil, MalformedConfigError
	}
	lowFloat, ok := low.(float64)
	if !ok {
		return nil, MalformedConfigError
	}
	medium, ok := wallet["MediumFeeDefault"]
	if !ok {
		return nil, MalformedConfigError
	}
	mediumFloat, ok := medium.(float64)
	if !ok {
		return nil, MalformedConfigError
	}
	high, ok := wallet["HighFeeDefault"]
	if !ok {
		return nil, MalformedConfigError
	}
	highFloat, ok := high.(float64)
	if !ok {
		return nil, MalformedConfigError
	}
	maxFee, ok := wallet["MaxFee"]
	if !ok {
		return nil, MalformedConfigError
	}
	maxFeeFloat, ok := maxFee.(float64)
	if !ok {
		return nil, MalformedConfigError
	}
	walletType, ok := wallet["Type"]
	if !ok {
		return nil, MalformedConfigError
	}
	walletTypeStr, ok := walletType.(string)
	if !ok {
		return nil, MalformedConfigError
	}
	binary, ok := wallet["Binary"]
	if !ok {
		return nil, MalformedConfigError
	}
	binaryStr, ok := binary.(string)
	if !ok {
		return nil, MalformedConfigError
	}
	rpcUser, ok := wallet["RPCUser"]
	if !ok {
		return nil, MalformedConfigError
	}
	rpcUserStr, ok := rpcUser.(string)
	if !ok {
		return nil, MalformedConfigError
	}
	rpcPassword, ok := wallet["RPCPassword"]
	if !ok {
		return nil, MalformedConfigError
	}
	rpcPasswordStr, ok := rpcPassword.(string)
	if !ok {
		return nil, MalformedConfigError
	}
	wCfg := &WalletConfig{
		Type:             walletTypeStr,
		Binary:           binaryStr,
		MaxFee:           int(maxFeeFloat),
		FeeAPI:           feeAPIstr,
		HighFeeDefault:   int(highFloat),
		MediumFeeDefault: int(mediumFloat),
		LowFeeDefault:    int(lowFloat),
		TrustedPeer:      trustedPeerStr,
		RPCUser:          rpcUserStr,
		RPCPassword:      rpcPasswordStr,
	}
	return wCfg, nil
}

func GetTorConfig(cfgBytes []byte) (*TorConfig, error) {
	var cfgIface interface{}
	json.Unmarshal(cfgBytes, &cfgIface)

	cfg, ok := cfgIface.(map[string]interface{})
	if !ok {
		return nil, MalformedConfigError
	}

	tcIface, ok := cfg["Tor-config"]
	if !ok {
		return nil, MalformedConfigError
	}
	tc, ok := tcIface.(map[string]interface{})

	pw, ok := tc["Password"]
	if !ok {
		return nil, MalformedConfigError
	}
	pwStr, ok := pw.(string)
	if !ok {
		return nil, MalformedConfigError
	}
	controlUrl, ok := tc["TorControl"]
	if !ok {
		return nil, MalformedConfigError
	}
	controlUrlStr, ok := controlUrl.(string)
	if !ok {
		return nil, MalformedConfigError
	}

	return &TorConfig{TorControl: controlUrlStr, Password: pwStr}, nil
}

func GetDropboxApiToken(cfgBytes []byte) (string, error) {
	var cfgIface interface{}
	json.Unmarshal(cfgBytes, &cfgIface)

	cfg, ok := cfgIface.(map[string]interface{})
	if !ok {
		return "", MalformedConfigError
	}

	token, ok := cfg["Dropbox-api-token"]
	if !ok {
		return "", MalformedConfigError
	}
	tokenStr, ok := token.(string)
	if !ok {
		return "", MalformedConfigError
	}

	return tokenStr, nil
}

func GetCrosspostGateway(cfgBytes []byte) ([]string, error) {
	var cfgIface interface{}
	json.Unmarshal(cfgBytes, &cfgIface)
	var urls []string

	cfg, ok := cfgIface.(map[string]interface{})
	if !ok {
		return urls, MalformedConfigError
	}

	gwys, ok := cfg["Crosspost-gateways"]
	if !ok {
		return urls, MalformedConfigError
	}
	gatewayList, ok := gwys.([]interface{})
	if !ok {
		return urls, MalformedConfigError
	}

	for _, gw := range gatewayList {
		gwStr, ok := gw.(string)
		if !ok {
			return urls, MalformedConfigError
		}
		urls = append(urls, gwStr)
	}

	return urls, nil
}

func GetResolverUrl(cfgBytes []byte) (string, error) {
	var cfgIface interface{}
	json.Unmarshal(cfgBytes, &cfgIface)

	cfg, ok := cfgIface.(map[string]interface{})
	if !ok {
		return "", MalformedConfigError
	}

	r, ok := cfg["Resolver"]
	if !ok {
		return "", MalformedConfigError
	}
	resolverStr, ok := r.(string)
	if !ok {
		return "", MalformedConfigError
	}

	return resolverStr, nil
}

func extendConfigFile(r repo.Repo, key string, value interface{}) error {
	if err := r.SetConfigKey(key, value); err != nil {
		return err
	}
	return nil
}

func InitConfig(repoRoot string) (*config.Config, error) {
	bootstrapPeers, err := config.ParseBootstrapPeers(DefaultBootstrapAddresses)
	if err != nil {
		return nil, err
	}

	datastore := datastoreConfig(repoRoot)

	conf := &config.Config{

		// Setup the node's default addresses.
		// NOTE: two swarm listen addrs, one TCP, one UTP.
		Addresses: config.Addresses{
			Swarm: []string{
				"/ip4/0.0.0.0/tcp/4001",
				"/ip6/::/tcp/4001",
				"/ip4/0.0.0.0/tcp/9005/ws",
				"/ip6/::/tcp/9005/ws",
			},
			API:     "",
			Gateway: "/ip4/127.0.0.1/tcp/4002",
		},

		Datastore: datastore,
		Bootstrap: config.BootstrapPeerStrings(bootstrapPeers),
		Discovery: config.Discovery{config.MDNS{
			Enabled:  false,
			Interval: 10,
		}},

		// Setup the node mount points
		Mounts: config.Mounts{
			IPFS: "/ipfs",
			IPNS: "/ipns",
		},

		Ipns: config.Ipns{
			ResolveCacheSize: 128,
			RecordLifetime:   "7d",
			RepublishPeriod:  "24h",
			QuerySize:        5,
		},

		Gateway: config.Gateway{
			RootRedirect: "",
			Writable:     false,
			PathPrefixes: []string{},
		},
	}

	return conf, nil
}

func datastoreConfig(repoRoot string) config.Datastore {
	dspath := path.Join(repoRoot, "datastore")
	return config.Datastore{
		Path:               dspath,
		Type:               "leveldb",
		StorageMax:         "10GB",
		StorageGCWatermark: 90, // 90%
		GCPeriod:           "1h",
	}
}
