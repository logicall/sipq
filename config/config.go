package config

import (
	"io/ioutil"
	"os"

	"github.com/henryscala/sipq/util"

	"github.com/naoina/toml"
)

//configurations of the executable
type ExeConfig struct {
	Server []Server
}

type Server struct {
	Type string //udp or tcp
	Ip   string //v4 or v6
	Port int    //number between 1 to 65535
}

var TheExeConfig *ExeConfig

var DefaultConfig string = `
	#example of the default config 
	[[server]]
		type="udp"
		ip="127.0.0.1"
		port=50600
	[[server]]
		type="tcp"
		ip="127.0.0.1"
		port=50600				
	`

func readExeConfigBytes(buf []byte) (*ExeConfig, error) {
	var cfg ExeConfig
	err := toml.Unmarshal(buf, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func ReadExeConfigFile(fileName string) (*ExeConfig, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return readExeConfigBytes(buf)
}

func init() {
	var err error

	TheExeConfig, err = readExeConfigBytes([]byte(DefaultConfig))
	util.ErrorPanic(err)

}
