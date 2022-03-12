package tenant

import (
	"encoding/json"
	"fmt"
	"github.com/m-to-n/config-management/utils"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestDaprStateSerialization(t *testing.T) {
	t.Parallel()

	tenantConfig := utils.GetDummyTenant()

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
