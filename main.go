package main

import (
	"context"
	"fmt"
	daprc "github.com/dapr/go-sdk/client"
)

func main() {
	client, err := daprc.NewClient()
	if err != nil {
		panic("no good!")
	}
	defer client.Close()

	ctx := context.Background()
	data := []byte("hello")
	store := "statestore-mongodb" // defined in the component YAML

	// save state with the key key1, default options: strong, last-write
	if err := client.SaveState(ctx, store, "key1", data); err != nil {
		panic(err)
	}

	fmt.Println("state stored!!!")
}
