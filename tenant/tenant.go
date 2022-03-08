package tenant

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/m-to-n/common/channels"
	"github.com/m-to-n/common/tenants"
	"github.com/m-to-n/config-management/dapr"
)

const MONGODB_STATE_STORE_TENANTS = "statestore-mongodb"

func CreateDummyTenant() error {
	client := dapr.DaprClient()
	if client == nil {
		return errors.New("unable to get DaprClient")
	}
	defer client.Close()

	ctx := context.Background()

	tenantConfig := tenants.TenantConfig{
		TenantId: "tenant-123",
		Name:     "dummyTenant",
		Desc:     "dummy tenant for development & testing",
		Channels: make([]tenants.TenantChannelConfig, 0),
	}

	tenantChannelConfig := tenants.TenantChannelConfig{
		Channel: channels.WhatsApp,
	}

	tenantConfig.Channels = append(tenantConfig.Channels, tenantChannelConfig)

	data, err := json.Marshal(tenantConfig)
	if err != nil {
		fmt.Printf("error when marshaling tenantConfig: %s", err.Error())
		return err
	}

	if err := client.SaveState(ctx, MONGODB_STATE_STORE_TENANTS, tenantConfig.TenantId, data); err != nil {
		fmt.Printf("error when saving to %s client %s", MONGODB_STATE_STORE_TENANTS, err.Error())
		return err
	}

	return nil
}
