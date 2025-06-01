// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kr/pretty"
	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v2"
)

var (
	conf *Config
	once sync.Once
)

type ProductService struct {
	Address string `yaml:"address"`
}

type Config struct {
	Env          string
	Kitex        Kitex        `yaml:"kitex"`
	MySQL        MySQL        `yaml:"mysql"`
	Redis        Redis        `yaml:"redis"`
	RateLimiter  RateLimiter  `yaml:"rate_limiter"`
	Registry     Registry     `yaml:"registry"`
	ProductService ProductService `yaml:"product_service"`
}

type Registry struct {
	ETCD ETCD `yaml:"etcd"`
}

type ETCD struct {
	Address  []string `yaml:"address"`
	TTL      int      `yaml:"ttl"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}

type RateLimiter struct {
	Enabled    bool `yaml:"enabled"`
	Rate       int  `yaml:"rate"`
	BucketSize int  `yaml:"bucket_size"`
}

type MySQL struct {
	DSN      string `yaml:"dsn"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Kitex struct {
	Service         string `yaml:"service"`
	Address         string `yaml:"address"`
	MetricsPort     string `yaml:"metrics_port"`
	EnablePprof     bool   `yaml:"enable_pprof"`
	EnableGzip      bool   `yaml:"enable_gzip"`
	EnableAccessLog bool   `yaml:"enable_access_log"`
	LogLevel        string `yaml:"log_level"`
	LogFileName     string `yaml:"log_file_name"`
	LogMaxSize      int    `yaml:"log_max_size"`
	LogMaxBackups   int    `yaml:"log_max_backups"`
	LogMaxAge       int    `yaml:"log_max_age"`
	UseK8sDiscovery bool   `yaml:"use_k8s_discovery"`
	K8sNamespace    string `yaml:"k8s_namespace"`
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	// Try multiple possible config paths
	possiblePaths := []string{
		filepath.Join("conf", GetEnv(), "conf.yaml"), // Relative to working dir
		filepath.Join("app", "cart", "conf", GetEnv(), "conf.yaml"), // Relative to project root
	}

	var content []byte
	var err error
	for _, path := range possiblePaths {
		content, err = os.ReadFile(path)
		if err == nil {
			break
		}
		klog.Warnf("failed to read config from %s: %v", path, err)
	}
	
	if err != nil {
		panic(fmt.Sprintf("failed to find config file in any of: %v, last error: %v", possiblePaths, err))
	}
	conf = new(Config)
	err = yaml.Unmarshal(content, conf)
	if err != nil {
		klog.Error("parse yaml error - %v", err)
		panic(err)
	}
	if err := validator.Validate(conf); err != nil {
		klog.Error("validate config error - %v", err)
		panic(err)
	}
	conf.Env = GetEnv()
	pretty.Printf("%+v\n", conf)
}

func GetEnv() string {
	e := os.Getenv("GO_ENV")
	if len(e) == 0 {
		return "dev"
	}
	return e
}

func LogLevel() klog.Level {
	level := GetConf().Kitex.LogLevel
	switch level {
	case "trace":
		return klog.LevelTrace
	case "debug":
		return klog.LevelDebug
	case "info":
		return klog.LevelInfo
	case "notice":
		return klog.LevelNotice
	case "warn":
		return klog.LevelWarn
	case "error":
		return klog.LevelError
	case "fatal":
		return klog.LevelFatal
	default:
		return klog.LevelInfo
	}
}
