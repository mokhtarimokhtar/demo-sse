package config

import (
	"encoding/json"
	"os"
)

type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Clock struct {
	Refresh uint8 `json:"refreshSecond"`
}

type Configuration struct {
	Server Server `json:"server"`
	Clock  Clock  `json:"clock"`
}

func (c *Configuration) ReadConfig(fileName string) (err error) {

	file, err := os.Open(fileName)

	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&c)
	if err != nil {
		//log.Println("Error json:", err)
		return err
	}
	return nil
}
