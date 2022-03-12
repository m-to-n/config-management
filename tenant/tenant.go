package tenant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/m-to-n/common/tenants"
	"github.com/m-to-n/config-management/dapr"
	"github.com/m-to-n/config-management/utils"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

/**
some useful links:
	https://www.isolineltd.com/blog/2020/implementing-custom-dapr-state.html
	https://golangtutorial.dev/tips/http-post-json-go/
	https://docs.dapr.io/developing-applications/building-blocks/state-management/howto-state-query-api/
*/

const MONGODB_STATE_STORE_TENANTS = "statestore-mongodb"

var (
	once                            sync.Once
	MONGODB_STATE_STORE_TENANTS_URL string
)

type DaprStateWrapper struct {
	Key   string               `json:"key"`
	Value tenants.TenantConfig `json:"value"`
}

func saveStateHttpApi(state []byte) error {
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

func SaveTenantConfig(tenantConfig tenants.TenantConfig) error {
	daprState := DaprStateWrapper{
		Key:   tenantConfig.TenantId,
		Value: tenantConfig,
	}

	var daprStates []DaprStateWrapper
	daprStates = append(daprStates, daprState)

	data, err := json.Marshal(daprStates)
	if err != nil {
		return err
	}

	if err := saveStateHttpApi(data); err != nil {
		fmt.Printf("error when saving to %s client %s", MONGODB_STATE_STORE_TENANTS, err.Error())
		return err
	}

	return nil
}

func GetTenantConfig(tenantId string) (*tenants.TenantConfig, error) {
	client := dapr.DaprClient()
	ctx := context.Background()

	stateItem, err := client.GetState(ctx, MONGODB_STATE_STORE_TENANTS, tenantId)
	if err != nil {
		return nil, err
	}

	var tenantConfig tenants.TenantConfig
	if err := json.Unmarshal(stateItem.Value, &tenantConfig); err != nil {
		return nil, err
	}

	return &tenantConfig, nil
}

// this functions just serves for quick testing of saving state
// via DAPR SDK API (which always serializes struct into json string in mongodb)
// and HTTP API (which works better and can serialize state value s nested json object in mongodb)
func CreateDummyTenant() error {
	client := dapr.DaprClient()
	ctx := context.Background()

	tenantConfig := utils.GetDummyTenant()

	if 1 == 2 /* never go this way :) */ {
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
		if err := SaveTenantConfig(tenantConfig); err != nil {
			fmt.Printf("error when saving to %s client %s", MONGODB_STATE_STORE_TENANTS, err.Error())
			return err
		}
	}

	return nil
}
