package gostructwalker

import "fmt"

type Config struct {
	TagKey string
}

func (cfg Config) validate() error {
	if cfg.TagKey == "" {
		return fmt.Errorf("Config value 'TagKey' is unset")
	}

	return nil
}
