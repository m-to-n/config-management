package tenant

import (
	"encoding/json"
	"fmt"
	"github.com/m-to-n/common/channels"
	"github.com/m-to-n/common/tenants"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestDaprStateSerialization(t *testing.T) {
	t.Parallel()

	tenantConfig := tenants.TenantConfig{
		TenantId: "tenant-123",
		Name:     "dummyTenant",
		Desc:     "dummy tenant for development & testing",
		Channels: make([]tenants.TenantChannelConfig, 0),
	}

	tenantChannelConfig := tenants.TenantChannelConfig{
		Channel: channels.CHANNELS_WHATSAPP,
	}

	tenantConfig.Channels = append(tenantConfig.Channels, tenantChannelConfig)

	daprState := DaprStateWrapper{
		Key:   tenantConfig.TenantId,
		Value: tenantConfig,
	}

	var daprStates []DaprStateWrapper
	daprStates = append(daprStates, daprState)

	data, err := json.Marshal(daprStates)

	if err != nil {
		log.Panic(fmt.Sprintf("error when marshaling daprState: %s", err.Error()))
	}

	datastr := string(data)
	expected := `[{"key":"tenant-123","value":{"tenantId":"tenant-123","name":"dummyTenant","desc":"dummy tenant for development \u0026 testing","channels":[{"channel":"whatsapp"}]}}]`
	assert.Equal(t, datastr, expected)

}
