package config

var (
	DefaultAddress = ":33700"
	DefaultPort    = 33700
	DefaultConfig  = &Config{
		Root:    DefaultRoot,
		Address: DefaultAddress,
	}
)

type Config struct {
	Root    string `yaml:"root"`
	Address string `yaml:"address"`
}

func LoadRoot() string {
	return DefaultConfig.Root
}
