package config

type Prefix struct {
	Vault  string `cfg:"vault"`
	Consul string `cfg:"consul"`
}

type LoadApp struct {
	Prefix    Prefix    `cfg:"prefix"`
	AppName   string    `cfg:"app_name"`
	ConfigSet GetConfig `cfg:"config_set"`
}

var LoadConfig = LoadApp{
	AppName: AppName,
	ConfigSet: GetConfig{
		Consul: false,
		Vault:  false,
		File:   true,
	},
}

type GetConfig struct {
	Vault  bool
	Consul bool
	File   bool
}
