package consul

import (
	"fmt"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
)

//Consul Consul manager
type Consul struct {
	Address    string
	Datacenter string
}

func init() {

}

//NewConsul ...
func NewConsul() *Consul {
	consul := &Consul{
		Address:    "192.168.1.177:8500",
		Datacenter: "ggserver",
	}
	return consul
}

//GetServices .
func GetServices() (ret map[string]*consulapi.AgentService, err error) {
	cl := NewConsul()
	return cl.GetServices()
}

//GetServices get all service
func (c *Consul) GetServices() (ret map[string]*consulapi.AgentService, err error) {
	conf := consulapi.DefaultConfig()
	conf.Address = c.Address
	conf.Datacenter = c.Datacenter
	client, err := consulapi.NewClient(conf)
	if err != nil {
		return
	}
	agent := client.Agent()
	ret, err = agent.Services()
	if err != nil {
		return
	}
	return
}

//AddKeyValue ...
func (c *Consul) AddKeyValue(key string, value string) (ret bool, err error) {
	conf := consulapi.DefaultConfig()
	conf.Address = c.Address
	conf.Datacenter = c.Datacenter
	client, err := consulapi.NewClient(conf)
	if err != nil {
		return
	}
	kv := client.KV()

	pair := consulapi.KVPair{Key: key}
	pair.Value = strconv.AppendQuote(pair.Value, value)
	wt, err := kv.Put(&pair, nil)
	if err != nil {
		return
	}
	fmt.Println("write key value ", wt.RequestTime)
	return
}

//GetKeyValue ..
func (c *Consul) GetKeyValue(key string) (value string) {
	conf := consulapi.DefaultConfig()
	conf.Address = c.Address
	conf.Datacenter = c.Datacenter
	client, err := consulapi.NewClient(conf)
	if err != nil {
		return
	}
	kv := client.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return
	}
	value = string(pair.Value)
	fmt.Printf("GetKeyValue key value %s ", pair.Value)
	return
}

//DelKey ...
func (c *Consul) DelKey(key string) bool {
	conf := consulapi.DefaultConfig()
	conf.Address = c.Address
	conf.Datacenter = c.Datacenter
	client, err := consulapi.NewClient(conf)
	if err != nil {
		fmt.Println("DelKey newclient fail ", err)
		return false
	}
	kv := client.KV()
	_, err = kv.Delete(key, nil)
	if err != nil {
		fmt.Println("DelKey fail ", err)
		return false
	}
	return true
}
