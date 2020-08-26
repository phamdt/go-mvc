package gomvc

import (
	"io/ioutil"
	"log"

	yaml "github.com/ghodss/yaml"
)

type GoMVCConfig struct {
	Denylist    []string
	denyListMap map[string]bool
}

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
		return config
	}
	if err := yaml.Unmarshal(c, &config); err != nil {
		log.Printf("empty config: %+v", err)
		return config
	}
	config.mapDenylist()
	return config
}

func (c *GoMVCConfig) mapDenylist() {
	c.denyListMap = map[string]bool{}
	for _, item := range c.Denylist {
		c.denyListMap[item] = true
	}
}

func (c *GoMVCConfig) IsDenylisted(path string) bool {
	ok := c.denyListMap[path]
	return ok
}
