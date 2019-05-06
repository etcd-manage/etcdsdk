package etcdsdk

import (
	"github.com/etcd-manage/etcdsdk/etcdv2"
	"github.com/etcd-manage/etcdsdk/etcdv3"
	"github.com/etcd-manage/etcdsdk/model"
)

/*
 golang 操作etcd sdk，可以兼容v2和v3版本etcd api
*/

// NewClientByConfig 创建一个etcd客户端
func NewClientByConfig(cfg *model.Config) (client model.EtcdSdk, err error) {
	if cfg == nil {
		err = model.ERR_CONFIG_ISNIL
		return
	}
	if cfg.Version == model.ETCD_VERSION_V2 {
		client, err = etcdv2.NewClient(cfg)
	} else if cfg.Version == model.ETCD_VERSION_V3 {
		client, err = etcdv3.NewClient(cfg)
	} else {
		err = model.ERR_UNSUPPORTED_VERSION
	}
	return
}
