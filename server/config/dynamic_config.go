package config

type IDynamicConfig interface {
	GetMempoolRecheck() bool
}

var DynamicConfig IDynamicConfig

func SetDynamicConfig(c IDynamicConfig) {
	DynamicConfig = c
}
