package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/cristalhq/aconfig"
	// "github.com/cristalhq/aconfig/aconfigyaml"
)

type Config struct {
	LogLevel    string `default:"info" usage:"Log level: panic, fatal, error, warn, info, debug, trace"`
	S           string `default:"" usage:"Source path"`
	D           string `default:"" usage:"Destination path"`
	Source      string `default:"" usage:"Source path"`
	Destination string `default:"" usage:"Destination path"`
	Layout      string `default:"2006-01-02" usage:"Datetime layout"`
}

func New() *Config {
	// fix for https://github.com/cristalhq/aconfig/issues/82
	args := []string{}
	for _, a := range os.Args {
		if !strings.HasPrefix(a, "-test.") {
			args = append(args, a)
		}
	}
	// fix for https://github.com/cristalhq/aconfig/issues/82

	var cfg Config
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		// feel free to skip some steps :)
		// SkipEnv:      true,
		// MergeFiles: true,
		SkipFiles: true,
		// SkipFiles:          false,
		AllowUnknownFlags:  true,
		AllowUnknownFields: true,
		SkipDefaults:       false,
		SkipFlags:          false,
		FailOnFileNotFound: false,
		EnvPrefix:          "MV2A",
		FlagPrefix:         "",
		// FileFlag:           "config",
		// Files: []string{
		// 	"./config.yaml",
		// 	"./config/config.yaml",
		// 	"/etc/tgbotfv/config.yaml",
		// 	"/etc/tgbotfv/config/config.yaml",
		// 	"/usr/local/tgbotfv/config.yaml",
		// 	"/usr/local/tgbotfv/config/config.yaml",
		// 	"/opt/tgbotfv/config.yaml",
		// 	"/opt/tgbotfv/config/config.yaml",
		// },
		// FileDecoders: map[string]aconfig.FileDecoder{
		// 	// from `aconfigyaml` submodule
		// 	// see submodules in repo for more formats
		// 	".yaml": aconfigyaml.New(),
		// },
		Args: args[1:], // [1:] важно, см. доку к FlagSet.Parse
	})
	if err := loader.Load(); err != nil {
		fmt.Println(err)
	}
	// if cfg.Path == "" {
	// 	cfg.Path = "."
	// }

	return &cfg
}
