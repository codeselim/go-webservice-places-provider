package providers

import (
	"github.com/codeselim/go-webservice-places-provider/config"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUnitgetHttpClientFromConfig(t *testing.T) {
	//empty config
	config1 := ProviderConfig{}
	httpClient1 := getHttpClientFromConfig(&config1)
	//with value
	config2 := ProviderConfig{Timeout: time.Second * 10}
	httpClient2 := getHttpClientFromConfig(&config2)

	assert.Equal(t, config.DefaultProviderTimeout, httpClient1.Timeout)
	assert.Equal(t, time.Second*10, httpClient2.Timeout)
}

func TestUnitgetSearchRadiusFromConfig(t *testing.T) {
	//empty config
	config1 := ProviderConfig{}
	radius := getSearchRadiusFromConfig(&config1)
	//with value
	config2 := ProviderConfig{SearchRadius: 10}
	radius2 := getSearchRadiusFromConfig(&config2)
	//with value
	config3 := ProviderConfig{SearchRadius: config.MaxAllowedSearchRadius + 1}
	radius3 := getSearchRadiusFromConfig(&config3)

	assert.Equal(t, config.DefaultSearchRadius, radius)
	assert.Equal(t, 10, radius2)
	assert.Equal(t, config.DefaultSearchRadius, radius3)
}
