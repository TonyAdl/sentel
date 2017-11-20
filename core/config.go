//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use this file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package core

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/Unknwon/goconfig"
	"github.com/golang/glog"
)

// Config interface
type Config interface {
	Bool(section string, key string) (bool, error)
	Int(section string, key string) (int, error)
	String(section string, key string) (string, error)
	MustBool(section string, key string) bool
	MustInt(section string, key string) int
	MustString(section string, key string) string
	SetValue(section string, key string, val string)
}

type configSection struct {
	items map[string]string
}

type globalConfig struct{}

var _allConfigSections map[string]*configSection = make(map[string]*configSection)

var (
	ErrorInvalidConfiguration = errors.New("Invalid configuration")
)

// globalConfig implementations

// Bool return bool value for key
func (c *globalConfig) Bool(section string, key string) (bool, error) {
	if _allConfigSections[section] == nil {
		return false, ErrorInvalidConfiguration
	}
	val := _allConfigSections[section].items[key]
	switch val {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}
	return false, fmt.Errorf("Invalid configuration item %s:%s", key, val)
}

// Int return int value for key
func (c *globalConfig) Int(section string, key string) (int, error) {
	if _allConfigSections[section] == nil {
		return -1, ErrorInvalidConfiguration
	}
	val := _allConfigSections[section].items[key]
	return strconv.Atoi(val)
}

// String return string valu for key
func (c *globalConfig) String(section string, key string) (string, error) {
	if _allConfigSections[section] == nil {
		return "", ErrorInvalidConfiguration
	}
	return _allConfigSections[section].items[key], nil
}

func (c *globalConfig) MustBool(section string, key string) bool {
	if _allConfigSections[section] == nil {
		glog.Fatal("Invalid configuration item:%s:%s", section, key)
		os.Exit(0)
	}
	val := _allConfigSections[section].items[key]
	switch val {
	case "true":
		return true
	case "false":
		return false
	}
	os.Exit(0)
	return false
}
func (c *globalConfig) MustInt(section string, key string) int {
	if _, ok := _allConfigSections[section]; !ok {
		glog.Fatal("Invalid configuration item:%s:%s", section, key)
		os.Exit(0)
	}
	val := _allConfigSections[section].items[key]
	n, err := strconv.Atoi(val)
	if err != nil {
		glog.Fatalf("Invalid configuration item:%s:%s", section, key)
		os.Exit(0)
	}
	return n
}

func (c *globalConfig) MustString(section string, key string) string {
	if _, ok := _allConfigSections[section]; !ok {
		glog.Fatalf("Invalid configuration: %s not found", section)
		os.Exit(0)
	}
	return _allConfigSections[section].items[key]
}

func (c *globalConfig) SetValue(section string, key string, valu string) {
}

// NewWithConfigFile load configurations from files
func NewConfigWithFile(fileName string, moreFiles ...string) (Config, error) {
	// load all config sections in _allConfigSections, get section and item to overide
	cfg, err := goconfig.LoadConfigFile(fileName, moreFiles...)
	if err == nil {
		sections := cfg.GetSectionList()
		for _, name := range sections {
			// create section if it doesn't exist
			if _, ok := _allConfigSections[name]; !ok {
				_allConfigSections[name] = &configSection{items: make(map[string]string)}
			}
			items, err := cfg.GetSection(name)
			if err == nil {
				for key, val := range items {
					_allConfigSections[name].items[key] = val
				}
			}
		}
	}
	return &globalConfig{}, nil
}

// Config global functions
func RegisterConfig(sectionName string, items map[string]string) {
	if _allConfigSections[sectionName] != nil { // section already exist
		section := _allConfigSections[sectionName]
		for key, val := range items {
			if section.items[key] != "" {
				glog.Infof("Config item(%s) will overide existed item:%s", key, section.items[key])
			}
			section.items[key] = val
		}
	} else {
		section := new(configSection)
		section.items = make(map[string]string)
		for key, val := range items {
			section.items[key] = val
		}
		_allConfigSections[sectionName] = section
	}
}

// RegisterConfigGropu reigster all group's subconfigurations
func RegisterConfigGroup(configs map[string]map[string]string) {
	for group, values := range configs {
		RegisterConfig(group, values)
	}
}
