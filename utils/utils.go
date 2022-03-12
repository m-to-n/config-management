package utils

import (
	"github.com/m-to-n/common/channels"
	"github.com/m-to-n/common/tenants"
)

func GetDummyTenant() tenants.TenantConfig {
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

	return tenantConfig
}
