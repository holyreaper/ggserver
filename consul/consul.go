package consul

import (
	"errors"
	"fmt"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/holyreaper/ggserver/def"
	. "github.com/holyreaper/ggserver/glog"
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
		LogFatal("consul GetServices get NewClient fail  error %s ", err)
		return
	}
	agent := client.Agent()
	ret, err = agent.Services()
	if err != nil {
		LogFatal("consul GetServices get AllServices fail  error %s ", err)
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

//GetSingleServerInfo .
func GetSingleServerInfo(sid def.SID) (info *consulapi.AgentService, err error) {
	srv, err := GetServices()
	if err != nil {
		LogFatal("consul GetSingleServerInfo get %d fail  error %s ", sid, err)
		return nil, err
	}
	if v, ok := srv[strconv.Itoa(int(sid))]; ok {
		LogInfo("consul GetSingleServerInfo get %d succ ", sid)
		return v, nil
	}
	LogFatal("consul GetSingleServerInfo get %d fail  error %s ", sid, err)
	return nil, errors.New("consul have no id ")
}
