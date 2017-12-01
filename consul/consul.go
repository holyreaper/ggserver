package consul

import (
	consulapi "github.com/hashicorp/consul/api"
)

//RegisterService ...
func RegisterService() bool {
	client, err := consulapi.NewClient(&consulapi.Config{})
	if err != nil {

	}
	client.Agent()
	return true
}
