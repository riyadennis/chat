package config

import (
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/gomniauth/providers/facebook"
)

type Config struct {
	Auth Auth
}
type Auth struct {
	Security  string
	Providers []Provider
}
type Provider struct {
	Name   string
	Client string
	Secret string
	Url    string
}

func (p Provider) GetGoogleProvider() (*google.GoogleProvider) {
	if p.Name == "google" {
		return google.New(p.Client, p.Secret, p.Url)
	}
	return nil
}
func (p Provider) GetFaceBookProvider() (*facebook.FacebookProvider) {
	if p.Name == "facebook" {
		return facebook.New(p.Client, p.Secret, p.Url)
	}
	return nil
}
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
