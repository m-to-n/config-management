package main

import (
	"fmt"
	"github.com/m-to-n/config-management/tenant"
)

func main() {
	if err := tenant.CreateDummyTenant(); err != nil {
		panic(err)
	}

	fmt.Println("state stored!!!")
}
