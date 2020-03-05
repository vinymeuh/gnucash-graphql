// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the format of the application's configuration file
type Config struct {
	Listen struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"listen"`
	Gnucash struct {
		File string `yaml:"file"`
	} `yaml:"gnucash"`
}

// NewConfig creates a new Config with default values
func NewConfig() Config {
	var c Config
	c.Listen.Host = "localhost"
	return c
}

// LoadFromFile fills the Config from a YAML file
func (c *Config) LoadFromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(c)
	return err
}
