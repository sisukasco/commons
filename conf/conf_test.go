package conf_test

import (
	"os"
	"github.com/sisukas/commons/conf"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestLoading(t *testing.T) {
	args := []string{}
	setting := faker.RandomString(12)
	os.Setenv("SIM_TEST_MY_SETTING_OPTION", setting)
	confx, err := conf.LoadConf(args, "SIM_TEST_", func(flags *pflag.FlagSet) {})
	assert.Nil(t, err)

	assert.Equal(t, confx.Konf.String("my.setting.option"), setting)

	assert.Equal(t, confx.Konf.String("my.setting.value"), "test value here")
}
