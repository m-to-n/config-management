package tenant

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	common_dapr "github.com/m-to-n/common/dapr"
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
	once                                  sync.Once
	MONGODB_STATE_STORE_TENANTS_URL       string
	MONGODB_STATE_STORE_TENANTS_URL_ALPHA string
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

func GetTenantByAccIdAndPhoneNum(ctx context.Context, tenantReq *tenants.TenantConfigByTwilioAccIDAndReceiverNumReq) (*tenants.TenantConfig, error) {
	client := &http.Client{}

	once.Do(func() {
		MONGODB_STATE_STORE_TENANTS_URL_ALPHA = fmt.Sprintf("http://localhost:%s/v1.0-alpha1/state/%s/query", dapr.DAPR_HTTP_PORT, MONGODB_STATE_STORE_TENANTS)
	})

	daprQuery := fmt.Sprintf(`
		{
			"filter": {
				"AND": [
					{
						"EQ": { "value.channels.data.whatsapp.accountSid": "%s" }
					},
					{
						"EQ": { "value.channels.data.whatsapp.numbers.phoneNumber": "%s" }
					}
				]
			},
			"page": {
				"limit": 1
			}
		}
	`, tenantReq.AccountSid, tenantReq.ReceiverPhoneNumber)

	req, err := http.NewRequest("POST", MONGODB_STATE_STORE_TENANTS_URL_ALPHA, bytes.NewBuffer([]byte(daprQuery)))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("dapr-app-id", "config-management") // is this needed?

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseString := string(responseData)

	/**
	will return 204 if no data found, 200 otherwise with payload like this:

		{
			"results": [
				{
					"key": "tenant-456",
					"data": {
						"tenantId": "tenant-456",
						"name": "dummyTenant456",
						"desc": "dummy tenant for dev",
						"channels": [
							{
								"data": {
									"whatsapp": {
										"accountSid": "AC4...",
										"authToken": "b34...",
										"numbers": [
											{
												"phoneNumber": "+420123456789",
												"language": "en"
											}
										]
									}
								},
								"channel": "whatsapp"
							}
						]
					},
					"etag": "66d260bf-bc1f-4f3e-bfb0-48f2dd3e14d7"
				}
			],
			"token": "1"
		}
	*/

	fmt.Printf("http response: %s Body: %s", resp.Status, responseString)

	if resp.StatusCode == 204 {
		log.Printf("No tenant config found: %s", resp.StatusCode)
		return nil, nil
	}

	if resp.StatusCode != 200 {
		errMsg := fmt.Sprintf("Error http status code received: %s", resp.StatusCode)
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	var tenantConfigResult tenants.TenantConfigDaprQueryResult
	if err := json.Unmarshal(responseData, &tenantConfigResult); err != nil {
		return nil, err
	}

	return &tenantConfigResult.Results[0].Data, nil

}

func SaveTenantConfig(ctx context.Context, tenantConfig tenants.TenantConfig) error {
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

func GetTenantConfig(ctx context.Context, tenantId string) (*tenants.TenantConfig, error) {
	client := common_dapr.DaprClient(dapr.DAPR_GRPC_PORT)

	fmt.Printf("getting state store: %s  state key: %s", MONGODB_STATE_STORE_TENANTS, tenantId)
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
	client := common_dapr.DaprClient(dapr.DAPR_GRPC_PORT)
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
		if err := SaveTenantConfig(ctx, tenantConfig); err != nil {
			fmt.Printf("error when saving to %s client %s", MONGODB_STATE_STORE_TENANTS, err.Error())
			return err
		}
	}

	return nil
}
