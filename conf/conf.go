package conf

import (
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf"
	kyaml "github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	kenv "github.com/knadh/koanf/providers/env"
	kfile "github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/sisukasco/commons/utils"
	"github.com/spf13/pflag"
)

type Confx struct {
	Konf  *koanf.Koanf
	Flags *pflag.FlagSet
}

func LoadConf(args []string, envPrefix string, addnl func(flags *pflag.FlagSet)) (*Confx, error) {
	confx := &Confx{}
	confx.Konf = koanf.New(".")

	confx.Flags = pflag.NewFlagSet("config", pflag.ContinueOnError)

	confx.Flags.String("conf", "conf.yaml", "path .yaml config file")
	addnl(confx.Flags)
	confx.Flags.Usage = func() {
		fmt.Println(confx.Flags.FlagUsages())
		os.Exit(0)
	}

	confx.Flags.Parse(args)

	cf, _ := confx.Flags.GetString("conf")

	if len(cf) > 0 && utils.FileExists(cf) {
		err := confx.Konf.Load(kfile.Provider(cf), kyaml.Parser())
		if err != nil {
			return nil, err
		}
	}

	err := confx.Konf.Load(kenv.Provider(envPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envPrefix)), "_", ".", -1)
	}), nil)
	if err != nil {
		return nil, err
	}

	err = confx.Konf.Load(posflag.Provider(confx.Flags, ".", confx.Konf), nil)
	if err != nil {
		return nil, err
	}

	if confx.Konf.Exists("db") {
		confx.initPostgresURL()
	}

	if confx.Konf.Exists("redis") {
		confx.initRedisURL()
	}
	return confx, nil
}

func (conf *Confx) initPostgresURL() {

	if len(conf.Konf.String("db.url")) > 0 {
		return
	}
	url := "postgres://" + conf.Konf.String("db.user") + ":" +
		conf.Konf.String("db.pass") + "@" +
		conf.Konf.String("db.server") + ":5432/" +
		conf.Konf.String("db.db") +
		"?sslmode=" + conf.Konf.String("db.sslmode")

	conf.Konf.Load(confmap.Provider(map[string]interface{}{
		"db.url": url,
	}, "."), nil)
}

func (conf *Confx) initRedisURL() {
	if len(conf.Konf.String("redis.url")) <= 0 {
		url := "redis://"
		pass := conf.Konf.String("redis.pass")
		if len(pass) > 0 {
			url += pass + "@"
		}
		url += conf.Konf.String("redis.host") +
			":6379/"
		conf.Konf.Load(confmap.Provider(map[string]interface{}{
			"redis.url":         url + conf.Konf.String("redis.db.common"),
			"machine.redis.url": url + conf.Konf.String("redis.db.machine"),
		}, "."), nil)
	}
}
