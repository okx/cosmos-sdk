package config

import (
	"bytes"
	"text/template"

	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
)

const defaultConfigTemplate = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

##### main base config options #####

# The minimum gas prices a validator is willing to accept for processing a
# transaction. A transaction's fees must meet the minimum of any denomination
# specified in this config (e.g. 0.25token1;0.0001token2).
minimum-gas-prices = "{{ .BaseConfig.MinGasPrices }}"

# HaltHeight contains a non-zero block height at which a node will gracefully
# halt and shutdown that can be used to assist upgrades and testing.
#
# Note: Commitment of state will be attempted on the corresponding block.
halt-height = {{ .BaseConfig.HaltHeight }}

# HaltTime contains a non-zero minimum block time (in Unix seconds) at which
# a node will gracefully halt and shutdown that can be used to assist upgrades
# and testing.
#
# Note: Commitment of state will be attempted on the corresponding block.
halt-time = {{ .BaseConfig.HaltTime }}
##### backend configuration options #####
[backend]
enable_backend = "{{ .BackendConfig.EnableBackend }}"
enable_mkt_compute = "{{ .BackendConfig.EnableMktCompute }}"
log_sql = "{{ .BackendConfig.LogSQL }}"
clean_ups_kept_days = "{{ .BackendConfig.CleanUpsKeptDays }}"
clean_ups_time = "{{ .BackendConfig.CleanUpsTime }}"
[backend.orm_engine]
engine_type = "{{ .BackendConfig.OrmEngine.EngineType }}"
connect_str = "{{ js .BackendConfig.OrmEngine.ConnectStr }}"
[stream]
engine = "{{ .StreamConfig.Engine }}"
klines_query_connect = "{{ .StreamConfig.KlineQueryConnect }}"
worker_id = "{{ .StreamConfig.WorkerId }}"
redis_scheduler = "{{ .StreamConfig.RedisScheduler }}"
redis_lock = "{{ .StreamConfig.RedisLock }}"
local_lock_dir = "{{ js .StreamConfig.LocalLockDir }}"
market_service_enable = "{{ .StreamConfig.MarketServiceEnable }}"
cache_queue_capacity = "{{ .StreamConfig.CacheQueueCapacity }}"
market_pulsar_topic = "{{ .StreamConfig.MarketPulsarTopic }}"
market_pulsar_partition = "{{ .StreamConfig.MarketPulsarPartition }}"
market_quotations_eureka_name = "{{ .StreamConfig.MarketQuotationsEurekaName }}"
eureka_server_url = "{{ .StreamConfig.EurekaServerUrl }}"
rest_application_name = "{{ .StreamConfig.RestApplicationName }}"
nacos_server_url = "{{ .StreamConfig.NacosServerUrl }}"
nacos_namespace_id = "{{ .StreamConfig.NacosNamespaceId }}"
pushservice_pulsar_public_topic = "{{ .StreamConfig.PushservicePulsarPublicTopic }}"
pushservice_pulsar_private_topic = "{{ .StreamConfig.PushservicePulsarPrivateTopic }}"
pushservice_pulsar_depth_topic = "{{ .StreamConfig.PushservicePulsarDepthTopic }}"
redis_require_pass = "{{ .StreamConfig.RedisRequirePass }}"
`

var configTemplate *template.Template

func init() {
	var err error
	tmpl := template.New("appConfigFileTemplate")
	if configTemplate, err = tmpl.Parse(defaultConfigTemplate); err != nil {
		panic(err)
	}
}

// ParseConfig retrieves the default environment configuration for the
// application.
func ParseConfig() (*Config, error) {
	conf := DefaultConfig()
	err := viper.Unmarshal(conf)
	return conf, err
}

// WriteConfigFile renders config using the template and writes it to
// configFilePath.
func WriteConfigFile(configFilePath string, config *Config) {
	var buffer bytes.Buffer

	if err := configTemplate.Execute(&buffer, config); err != nil {
		panic(err)
	}

	cmn.MustWriteFile(configFilePath, buffer.Bytes(), 0644)
}
