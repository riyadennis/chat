package config

import (
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProviderGetFaceBookProvider(t *testing.T) {
	provider := Provider{
		Name:   "facebook",
		Client: "TESTCLIENT",
		Secret: "THISISASCERET",
		URL:    "http://localhost:8080/auth/callback/facebook/",
	}
	facebookProvider := provider.GetFaceBookProvider()
	assert.IsType(t, facebookProvider, &facebook.FacebookProvider{})
	assert.Equal(t, facebookProvider.Name(), "facebook")
}
func TestProviderGetGoogleProvider(t *testing.T) {
	provider := Provider{
		Name:   "google",
		Client: "TESTCLIENT",
		Secret: "THISISASCERET",
		URL:    "http://localhost:8080/auth/callback/google/",
	}
	googleProvider := provider.GetGoogleProvider()
	assert.IsType(t, googleProvider, &google.GoogleProvider{})
	assert.Equal(t, googleProvider.Name(), "google")
}
func TestParseConfigNoFile(t *testing.T) {
	config, err := ParseConfig("invalid_config.yaml")
	assert.Error(t, err)
	assert.IsType(t, config, &Config{})
}
func TestParseConfigValidFile(t *testing.T) {
	config, err := ParseConfig("config.test.yaml")
	assert.NoError(t, err)
	assert.IsType(t, config, &Config{})
	assert.Equal(t, config.Auth.Security, "ChangeMe")
}
