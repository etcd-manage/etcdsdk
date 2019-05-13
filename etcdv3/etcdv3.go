package etcdv3

import (
	"context"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/etcd-manage/etcdsdk/model"
)

var (
	// DefaultTimeout 默认查询超时
	DefaultTimeout = 5 * time.Second
)

// EtcdV3Sdk etcd v3版
type EtcdV3Sdk struct {
	cli *clientv3.Client
}

// NewClient 创建etcd连接
func NewClient(cfg *model.Config) (client model.EtcdSdk, err error) {
	if cfg == nil {
		err = model.ERR_CONFIG_ISNIL
		return
	}
	if cfg.TlsEnable == true && (cfg.CertFile == "" || cfg.KeyFile == "" || cfg.CaFile == "") {
		err = model.ERR_TLS_CONFIG_ISNIL
		return
	}
	if len(cfg.Address) == 0 {
		err = model.ERR_ETCD_ADDRESS_EMPTY
		return
	}

	var cli *clientv3.Client

	if cfg.TlsEnable == true {
		// tls 配置
		tlsInfo := transport.TLSInfo{
			CertFile:      cfg.CertFile,
			KeyFile:       cfg.KeyFile,
			TrustedCAFile: cfg.CaFile,
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			return nil, err
		}

		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   cfg.Address,
			DialTimeout: 10 * time.Second,
			TLS:         tlsConfig,
			Username:    cfg.Username,
			Password:    cfg.Password,
		})
	} else {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   cfg.Address,
			DialTimeout: 10 * time.Second,
			Username:    cfg.Username,
			Password:    cfg.Password,
		})
	}

	if err != nil {
		return
	}
	// 可操作etcd客户端对象
	client = &EtcdV3Sdk{
		cli: cli,
	}
	return
}

// List 显示当前path下所有key
func (sdk *EtcdV3Sdk) List(path string) (list []*model.Node, err error) {
	// 9 秒超时
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	// 获取指定前缀key列表
	resp, err := sdk.cli.Get(ctx, path,
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		return
	}
	/* 处理出当前目录层的key */
	if resp.Count == 0 {
		return
	}
	list, err = sdk.ConvertToPath(path, resp.Kvs)

	return
}

// Val 获取path的值
func (sdk *EtcdV3Sdk) Val(path string) (data []byte, err error) {
	return
}

// Add 添加key
func (sdk *EtcdV3Sdk) Add(path string, data []byte) (err error) {
	return
}

// Put 修改key
func (sdk *EtcdV3Sdk) Put(path string, data []byte) (err error) {
	return
}

// Del 删除key
func (sdk *EtcdV3Sdk) Del(path string) (err error) {
	return
}

// Members 获取节点列表
func (sdk *EtcdV3Sdk) Members() (members []*model.Member, err error) {
	return
}
