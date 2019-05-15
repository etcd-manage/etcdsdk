package etcdsdk

import (
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
func NewClientByConfig(cfg *model.Config) (client model.EtcdSdk, err error) {
	if cfg == nil {
		err = model.ERR_CONFIG_ISNIL
		return
	}
	if cfg.Version == model.ETCD_VERSION_V2 {
		if val, ok := v2ClientMap.Load(cfg); ok == true {
			client, ok = val.(model.EtcdSdk)
			if ok == false {
				client, err = etcdv2.NewClient(cfg)
			} else {
				return
			}
		} else {
			client, err = etcdv2.NewClient(cfg)
		}
		v2ClientMap.Store(cfg, client)
	} else if cfg.Version == model.ETCD_VERSION_V3 {
		if val, ok := v3ClientMap.Load(cfg); ok == true {
			client, ok = val.(model.EtcdSdk)
			if ok == false {
				client, err = etcdv3.NewClient(cfg)
			} else {
				return
			}
		} else {
			client, err = etcdv3.NewClient(cfg)
		}
		v3ClientMap.Store(cfg, client)
	} else {
		err = model.ERR_UNSUPPORTED_VERSION
	}
	return
}
