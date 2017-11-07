package config

import (
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/gomniauth/providers/google"
)

type Config struct {
	Auth Auth
}
type Auth struct {
	Security string
	Providers []Provider
}
type Provider struct {
	Name string
	Client string
	Secret string
	Url string
}
func (p Provider) NewProvider(){
	if p.Name == "google"{
		google.New(p.Client, p.Secret, p.Url)
	}
}
func ParseConfig(pathToFile string) *Config {
	cFile, err := os.Open(pathToFile)
	if err != nil {
		logrus.Fatal("unable to open the file")
	}
	configCon, err := ioutil.ReadAll(cFile)
	if err != nil {
		logrus.Fatal("Unable to read the config file")
	}
	config := &Config{}
	err = yaml.Unmarshal(configCon, config)
	if err != nil {
		logrus.Fatalf("Invalid YAML file %s", err.Error())
	}
	return config
}

