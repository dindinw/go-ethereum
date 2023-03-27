package params

type MeerChainConfig struct {
	ChainID string // chainId identifies the current chain and is used for replay protection
}

var (
	QngMainnetChainConfig = &MeerChainConfig{
		ChainID: "813",
	}
	QngTestnetChainConfig = &MeerChainConfig{
		ChainID: "8131",
	}
	AmanaChainConfig = &MeerChainConfig{
		ChainID: "8132",
	}
	AmanaTestnetChainConfig = &MeerChainConfig{
		ChainID: "81321",
	}
	FlanaChainConfig = &MeerChainConfig{
		ChainID: "8133",
	}
	FlanaTestnetChainConfig = &MeerChainConfig{
		ChainID: "81331",
	}
	MizanaChainConfig = &MeerChainConfig{
		ChainID: "8134",
	}
	MizanaTestnetChainConfig = &MeerChainConfig{
		ChainID: "81341",
	}
)

func init() {
	NetworkNames[QngMainnetChainConfig.ChainID] = "qng"
	NetworkNames[QngTestnetChainConfig.ChainID] = "qng-test"
	NetworkNames[AmanaChainConfig.ChainID] = "amana"
	NetworkNames[AmanaTestnetChainConfig.ChainID] = "amana-test"
	NetworkNames[FlanaChainConfig.ChainID] = "flana"
	NetworkNames[FlanaTestnetChainConfig.ChainID] = "flana-test"
	NetworkNames[MizanaChainConfig.ChainID] = "mizana"
	NetworkNames[MizanaTestnetChainConfig.ChainID] = "mizana-test"
}
