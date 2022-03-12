package tenant

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/m-to-n/common/channels"
	"github.com/m-to-n/common/tenants"
	"github.com/m-to-n/config-management/dapr"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

// https://www.isolineltd.com/blog/2020/implementing-custom-dapr-state.html
// https://golangtutorial.dev/tips/http-post-json-go/
const MONGODB_STATE_STORE_TENANTS = "statestore-mongodb"

var (
	once                            sync.Once
	MONGODB_STATE_STORE_TENANTS_URL string
)

// see sample state data here: https://docs.dapr.io/developing-applications/building-blocks/state-management/howto-state-query-api/
type DaprStateWrapper struct {
	Key   string               `json:"key"`
	Value tenants.TenantConfig `json:"value"`
}

func saveState(state []byte) error {
	client := &http.Client{}

	once.Do(func() {
		MONGODB_STATE_STORE_TENANTS_URL = fmt.Sprintf("http://localhost:%s/v1.0/state/%s", dapr.DAPR_HTTP_PORT, MONGODB_STATE_STORE_TENANTS)
	})

	req, err := http.NewRequest("POST", MONGODB_STATE_STORE_TENANTS_URL, bytes.NewBuffer(state))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	responseString := string(responseData)

	fmt.Printf("http response: %s Body: %s", resp.Status, responseString)

	return nil
}

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
		Channel: channels.CHANNELS_WHATSAPP,
	}

	tenantConfig.Channels = append(tenantConfig.Channels, tenantChannelConfig)

	if 1 == 11 {
		// saving state via SDK, not sure how to save struct as json (will be stringified internally)
		data, err := json.Marshal(tenantConfig)
		if err != nil {
			fmt.Printf("error when marshaling tenantConfig: %s", err.Error())
			return err
		}

		if err := client.SaveState(ctx, MONGODB_STATE_STORE_TENANTS, tenantConfig.TenantId, data); err != nil {
			fmt.Printf("error when saving to %s client %s", MONGODB_STATE_STORE_TENANTS, err.Error())
			return err
		}
	} else {
		// saving state "manually" via HTTP API
		daprState := DaprStateWrapper{
			Key:   tenantConfig.TenantId,
			Value: tenantConfig,
		}

		var daprStates []DaprStateWrapper
		daprStates = append(daprStates, daprState)

		data, err := json.Marshal(daprStates)
		if err != nil {
			fmt.Printf("error when marshaling daprState: %s", err.Error())
			return err
		}
		fmt.Printf("serialized daprStates: %s", string(data))

		if err := saveState(data); err != nil {
			fmt.Printf("error when saving to %s client %s", MONGODB_STATE_STORE_TENANTS, err.Error())
			return err
		}
	}

	return nil
}
