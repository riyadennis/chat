package config

import (
	"github.com/stretchr/gomniauth/common"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/google"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// Config struct holds configuration read from the file
type Config struct {
	Auth *Auth
	TemplatePath string
}

// Auth struct holds authentication information
type Auth struct {
	Security  string
	Providers  []*Provider
}

// Provider struct has the details of the providers like facebook, google etc
type Provider struct {
	Name   string
	Client string
	Secret string
	URL    string
}

func (p Provider) GetProvider()common.Provider{
	switch p.Name {
	case "google":
		// creates google client
		return google.New(p.Client, p.Secret, p.URL)
	case "facebook":
		// creates facebook client
		return facebook.New(p.Client, p.Secret, p.URL)
	default:
		return nil
	}
}

// ParseConfig gets the config file and create a struct
func ParseConfig(pathToFile string) (*Config, error) {
	cFile, err := os.Open(pathToFile)
	if err != nil {
		return nil, err
	}
	configCon, err := ioutil.ReadAll(cFile)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(configCon, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
