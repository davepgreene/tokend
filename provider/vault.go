package provider

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

func vaultConfig() *api.Config {
	if viper.GetBool("vault.tls") == true {
		return api.DefaultConfig()
	}

	endpoint := fmt.Sprintf("http://%s:%d",
		viper.GetString("vault.host"),
		viper.GetInt("vault.port"))

	config := &api.Config{
		Address:    endpoint,
		HttpClient: cleanhttp.DefaultClient(),
	}
	config.HttpClient.Timeout = time.Second * 60

	return config
}
