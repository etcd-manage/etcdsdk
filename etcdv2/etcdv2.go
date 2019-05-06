package etcdv2

import (
	"github.com/etcd-manage/etcdsdk/model"
)

// EtcdV2Sdk etcd v2版
type EtcdV2Sdk struct {
}

// NewClient 创建etcd连接
func NewClient(cfg *model.Config) (client model.EtcdSdk, err error) {

	return
}

// List 显示当前path下所有key
func (sdk *EtcdV2Sdk) List(path string) (list []*model.Node, err error) {

	return
}

// Val 获取path的值
func (sdk *EtcdV2Sdk) Val(path string) (data []byte, err error) {
	return
}

// Add 添加key
func (sdk *EtcdV2Sdk) Add(path string, data []byte) (err error) {
	return
}

// Put 修改key
func (sdk *EtcdV2Sdk) Put(path string, data []byte) (err error) {
	return
}

// Del 删除key
func (sdk *EtcdV2Sdk) Del(path string) (err error) {
	return
}

// Members 获取节点列表
func (sdk *EtcdV2Sdk) Members() (members []*model.Member, err error) {
	return
}
