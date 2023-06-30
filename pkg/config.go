package giga

type Config struct {
	Gigawallet GigawalletConfig
	WebAPI     WebAPIConfig
	Store      StoreConfig
	Loggers    map[string]LoggersConfig

	// Map of available networks, config.Core will be set to
	// the one specified by config.Gigawallet.Network
	Dogecoind map[string]NodeConfig
	Core      NodeConfig
}

type GigawalletConfig struct {
	// Doge Connect service domain, where is GW hosted?
	ServiceDomain string

	// Doge Connect service name, ie: Doge Payments Inc.
	ServiceName string

	// Doge Connect service icon, displayed beside name.
	ServiceIconURL string

	// A DOGENS key-hash that appears in a DOGENS DNS TXT record
	// at the ServiceDomain, will be looked up by clients to verify
	// Doge Connect messages were signed with ServiceKeySecret
	ServiceKeyHash string

	// The private key used by this GW to sign all Doge Connect
	// envelopes, consider using --service-key-secret with an
	// appropriate secret management service when deploying, rather
	// than embedding in your config file.
	ServiceKeySecret string

	// key for which Dogecoind struct to use, ie: mainnet, testnet
	Network string

	// Default number of confirmations needed to mark an invoice
	// as paid, this can be overridden per invoice using the create
	// invoice API
	ConfirmationsNeeded int
}

type NodeConfig struct {
	Host    string
	ZMQPort int
	RPCHost string
	RPCPort int
	RPCPass string
	RPCUser string
}

type WebAPIConfig struct {
	Port string
	Bind string // optional interface IP address
}

type StoreConfig struct {
	DBFile string
}

type LoggersConfig struct {
	Path  string
	Types []string
}
