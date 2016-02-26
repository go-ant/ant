package setting

import (
	"github.com/robfig/config"
	"strings"
)

type appConfig struct {
	config  *config.Config
	section string
}

func loadConfig(fileName string) (*appConfig, error) {
	config, err := config.ReadDefault(fileName)
	return &appConfig{config: config}, err
}

func (c *appConfig) string(option string) (string, bool) {
	if r, err := c.config.String(c.section, option); err == nil {
		return r, true
	}
	return "", false
}

func (c *appConfig) int(option string) (int, bool) {
	if r, err := c.config.Int(c.section, option); err == nil {
		return r, true
	}
	return 0, false
}

func (c *appConfig) Section(section string) *appConfig {
	c.section = section
	return c
}

func (c *appConfig) StringDefault(option, dfault string) string {
	if r, found := c.string(option); found {
		return r
	}
	return dfault
}

func (c *appConfig) IntDefault(option string, dfault int) int {
	if r, found := c.int(option); found {
		return r
	}
	return dfault
}

func (c *appConfig) ArrayDefault(option string, dfault []string) []string {
	opt, found := c.string(option)
	if found {
		return strings.Split(opt, "|")
	}
	return dfault
}
