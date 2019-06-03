package core

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"x_yield/plugins"
)

type Core struct {
	Config  *Config
	Plugins []*plugins.Plugin
}

func NewCore() *Core {
	c := &Core{}
	return c
}

// LoadConfig loads the given config file and applies it to c
func (c *Core) LoadConfig(path string) error {
	var err error
	data, err := loadConfig(path)
	if err != nil {
		return fmt.Errorf("Error loading %s, %s", path, err)
	}
	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return nil
}

func loadConfig(config string) ([]byte, error) {
	u, err := url.Parse(config)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "https", "http":
		return fetchConfig(u)
	default:
		// If it isn't a https scheme, try it as a file.
	}
	return ioutil.ReadFile(config)
}

func fetchConfig(u *url.URL) ([]byte, error) {
	v := os.Getenv("X_YIELD_TOKEN")

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Token "+v)
	req.Header.Add("Accept", "application/yaml")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve remote config: %s", resp.Status)
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
