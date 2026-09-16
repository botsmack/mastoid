package pkg

import (
	"github.com/mattn/go-mastodon"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStoreCredentialsWritesSecureConfig(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	viper.Reset()
	t.Cleanup(viper.Reset)

	credentials := &Credentials{
		Server:      "https://mastodon.example",
		GrantToken:  "grant-token",
		AccessToken: "access-token",
		Application: &mastodon.Application{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			AuthURI:      "https://mastodon.example/oauth/authorize",
			RedirectURI:  "urn:ietf:wg:oauth:2.0:oob",
		},
	}

	if err := StoreCredentials(credentials); err != nil {
		t.Fatalf("StoreCredentials returned error: %v", err)
	}

	configFile := filepath.Join(tempHome, ".mastoid", "config.yml")

	info, err := os.Stat(configFile)
	if err != nil {
		t.Fatalf("could not stat config file: %v", err)
	}

	if got := info.Mode().Perm(); got != 0600 {
		t.Fatalf("expected config permissions 0600, got %04o", got)
	}

	configDir := filepath.Dir(configFile)

	dirInfo, err := os.Stat(configDir)
	if err != nil {
		t.Fatalf("could not stat config directory: %v", err)
	}

	if got := dirInfo.Mode().Perm(); got != 0700 {
		t.Fatalf("expected config directory permissions 0700, got %04o", got)
	}
}

func TestStoreCredentialsUsesLoadedConfigFile(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "custom.yml")

	if err := os.WriteFile(configFile, []byte("server: https://old.example\n"), 0600); err != nil {
		t.Fatalf("could not create config file: %v", err)
	}

	viper.Reset()
	t.Cleanup(viper.Reset)

	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("could not load config file: %v", err)
	}

	credentials := &Credentials{
		Server:      "https://mastodon.example",
		GrantToken:  "grant-token",
		AccessToken: "access-token",
		Application: &mastodon.Application{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			AuthURI:      "https://mastodon.example/oauth/authorize",
			RedirectURI:  "urn:ietf:wg:oauth:2.0:oob",
		},
	}

	if err := StoreCredentials(credentials); err != nil {
		t.Fatalf("StoreCredentials returned error: %v", err)
	}

	if got := viper.ConfigFileUsed(); got != configFile {
		t.Fatalf("expected config file %q, got %q", configFile, got)
	}

	info, err := os.Stat(configFile)
	if err != nil {
		t.Fatalf("could not stat config file: %v", err)
	}

	if got := info.Mode().Perm(); got != 0600 {
		t.Fatalf("expected config permissions 0600, got %04o", got)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("could not read config file: %v", err)
	}

	content := string(data)

	for _, expected := range []string{
		"client-id",
		"client-secret",
		"access-token",
		"https://mastodon.example",
	} {
		if !strings.Contains(content, expected) {
			t.Errorf("config file does not contain expected value %q", expected)
		}
	}
}
