package main

import (
	"fmt"
	"order-service-wb/pkg/config"
)

func main() {
	conf := config.NewConfig()
	fmt.Printf("%+v\n", conf)
}
