package gomvc

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// NewGoMVCConfig is a constructor for a gomvc configuration read into memory
func NewGoMVCConfig(configDir string) GoMVCConfig {
	var config GoMVCConfig
	if configDir == "" {
		// use defaults
		log.Println("no config provided, using defaults")
		return config
	}
	c, err := ioutil.ReadFile(configDir)
	if err != nil {
		log.Printf("error reading config: %v", err)
	} else {
		err = yaml.Unmarshal(c, &config)
		if err != nil {
			log.Printf("empty config: %+v", err)
			return config
		}
		config.mapBlacklist()
	}
	return config
}

type GoMVCConfig struct {
	Blacklist    []string
	blacklistMap map[string]bool
}

func (c *GoMVCConfig) mapBlacklist() {
	c.blacklistMap = map[string]bool{}
	for _, item := range c.Blacklist {
		c.blacklistMap[item] = true
	}
}

func (c *GoMVCConfig) IsBlacklisted(path string) bool {
	ok, _ := c.blacklistMap[path]
	return ok
}
