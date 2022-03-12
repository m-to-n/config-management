package main

import (
	"fmt"
	"github.com/m-to-n/config-management/tenant"
)

func main() {
	/* if err := tenant.CreateDummyTenant(); err != nil {
		panic(err)
	} */

	tcp, err := tenant.GetTenantConfig("tenant-123")
	if err != nil {
		panic(err)
	}

	tenantConfig := *tcp

	fmt.Println("I am done here %s", tenantConfig)

	fmt.Println("I am done here!!!")
}
