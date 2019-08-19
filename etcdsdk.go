package etcdsdk

import (
	"log"
	"sync"

	"github.com/etcd-manage/etcdsdk/etcdv2"
	"github.com/etcd-manage/etcdsdk/etcdv3"
	"github.com/etcd-manage/etcdsdk/model"
)

/*
 golang 操作etcd sdk，可以兼容v2和v3版本etcd api
*/

// 保存所有etcd客户端连接
var (
	v2ClientMap = new(sync.Map)
	v3ClientMap = new(sync.Map)
)

// NewClientByConfig 创建一个etcd客户端
// 可重复调用，不会重复和etcd建立连接
func NewClientByConfig(cfgObj *model.Config) (client model.EtcdSdk, err error) {
	if cfgObj == nil {
		err = model.ERR_CONFIG_ISNIL
		return
	}
	cfg := cfgObj.String()
	if cfgObj.Version == model.ETCD_VERSION_V2 {
		if val, ok := v2ClientMap.Load(cfg); ok == true {
			client, ok = val.(model.EtcdSdk)
			if ok == false {
				client, err = etcdv2.NewClient(cfgObj)
			} else {
				return
			}
		} else {
			client, err = etcdv2.NewClient(cfgObj)
		}
		v2ClientMap.Store(cfg, client)
	} else if cfgObj.Version == model.ETCD_VERSION_V3 {
		if val, ok := v3ClientMap.Load(cfg); ok == true {
			client, ok = val.(model.EtcdSdk)
			if ok == false {
				log.Println("创建连接v3")
				client, err = etcdv3.NewClient(cfgObj)
			} else {
				return
			}
		} else {
			log.Println("创建连接v3")
			client, err = etcdv3.NewClient(cfgObj)
		}
		v3ClientMap.Store(cfg, client)
	} else {
		err = model.ERR_UNSUPPORTED_VERSION
	}
	return
}
