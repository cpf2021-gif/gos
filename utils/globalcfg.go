package utils

import (
	"github.com/cpf2021-gif/gos/tiface"

	"github.com/spf13/viper"
)

var GlobalConfig *GlobalCfg

func init() {
	GlobalConfig = &GlobalCfg{
		ServerCfg: ServerCfg{
			TcpServer: nil,
			Host:      "0.0.0.0",
			TcpPort:   8999,
			Name:      "[gos] Server v0.1",
		},
		Gos: GosCfg{
			Version:       "0.1",
			MaxConn:       1024,
			MaxPacketSize: 4096,
		},
	}
}

func LoadConfig(path string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(GlobalConfig)
	if err != nil {
		panic(err)
	}
}

type GlobalCfg struct {
	ServerCfg ServerCfg `mapstructure:"server" json:"server" yaml:"server"`
	Gos       GosCfg    `mapstructure:"gos" json:"gos" yaml:"gos"`
}

type ServerCfg struct {
	TcpServer tiface.IServer
	Host      string `mapstructure:"host" json:"host" yaml:"host"`
	TcpPort   int    `mapstructure:"tcp_port" json:"tcp_port" yaml:"tcp_port"`
	Name      string `mapstructure:"name" json:"name" yaml:"name"`
}

type GosCfg struct {
	Version       string `mapstructure:"version" json:"version" yaml:"version"`
	MaxConn       int    `mapstructure:"max_conn" json:"max_conn" yaml:"max_conn"`
	MaxPacketSize int    `mapstructure:"max_packet_size" json:"max_packet_size" yaml:"max_packet_size"`
}
