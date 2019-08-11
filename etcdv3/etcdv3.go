package etcdv3

import (
	"context"
	"log"
	"sort"
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
		clientv3.WithPrefix(), clientv3.WithKeysOnly()) // 此处只查询出key，不查询值
	if err != nil {
		return
	}
	/* 处理出当前目录层的key */
	if resp.Count == 0 {
		return
	}
	list, err = sdk.ConvertToPath(path, resp.Kvs)

	// etcd 排序无效，自己实现
	sort.Slice(list, func(i, j int) bool {
		return list[i].Path < list[j].Path
	})

	// 如果是值，则查询值内容
	for _, v := range list {
		rv, err := sdk.cli.Get(ctx, v.Path)
		if err != nil {
			log.Println("读取值错误")
			continue
		}
		if len(rv.Kvs) > 0 {
			v.Value = string(rv.Kvs[0].Value)
		}
	}

	return
}

// Val 获取path的值
func (sdk *EtcdV3Sdk) Val(path string) (data *model.Node, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	resp, err := sdk.cli.Get(ctx, path)
	if err != nil {
		return
	}
	if len(resp.Kvs) == 0 {
		err = model.ERR_KEY_NOT_FOUND
		return
	}
	// 返回一个node结构
	list, err := sdk.ConvertToPath(path, resp.Kvs)
	if err != nil {
		return
	}
	data = list[0]
	return
}

// Add 添加key
func (sdk *EtcdV3Sdk) Add(path string, data []byte) (err error) {
	// 使用事物，防止覆盖，添加就是添加，不可以覆盖
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	txn := sdk.cli.Txn(ctx)
	txn.If(
		clientv3.Compare(
			clientv3.Version(path),
			"=",
			0,
		),
	).Then(
		clientv3.OpPut(path, string(data)),
	)

	txnResp, err := txn.Commit()
	if err != nil {
		return err
	}

	if !txnResp.Succeeded {
		return model.ERR_ADD_KEY
	}
	return
}

// Put 修改key
func (sdk *EtcdV3Sdk) Put(path string, data []byte) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	_, err = sdk.cli.Put(ctx, path, string(data))
	if err != nil {
		return
	}
	return
}

// Del 删除key - 虚拟目录不允许删除
func (sdk *EtcdV3Sdk) Del(path string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	_, err = sdk.cli.Delete(ctx, path)
	if err != nil {
		return
	}
	return
}

// Members 获取节点列表
func (sdk *EtcdV3Sdk) Members() (members []*model.Member, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	resp, err := sdk.cli.MemberList(ctx)
	if err != nil {
		return nil, err
	}

	for _, member := range resp.Members {
		if len(member.ClientURLs) > 0 {
			m := &model.Member{Member: member, Role: model.ROLE_FOLLOWER, Status: model.STATUS_UNHEALTHY}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			resp, err := sdk.cli.Status(ctx, m.ClientURLs[0])
			if err == nil {
				m.Status = model.STATUS_HEALTHY
				m.DbSize = resp.DbSize
				if resp.Leader == resp.Header.MemberId {
					m.Role = model.ROLE_LEADER
				}
			}
			members = append(members, m)
		}
	}
	return
}
