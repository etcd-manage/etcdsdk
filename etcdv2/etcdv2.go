package etcdv2

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/etcd-manage/etcdsdk/model"
)

var (
	// DefaultTimeout 默认查询超时
	DefaultTimeout = 5 * time.Second
	sm             = new(sync.Mutex)
)

// EtcdV2Sdk etcd v2版
type EtcdV2Sdk struct {
	cli        client.Client
	keysAPI    client.KeysAPI
	membersAPI client.MembersAPI
}

// NewClient 创建etcd连接
func NewClient(cfg *model.Config) (clientv2 model.EtcdSdk, err error) {
	sm.Lock()
	defer func() {
		sm.Unlock()
	}()
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
	var cli client.Client

	if cfg.TlsEnable == true {
		// // 数据库配置存储为key文件内容，此处每次都将内容写入文件
		// certFilePath, keyFilePath, _, err := writeCa(cfg, cfg.EtcdId)
		// if err != nil {
		// 	return clientv2, err
		// }
		// // tls 配置
		// cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
		// if err != nil {
		// 	return clientv2, err
		// }

		transportTls := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				// Certificates: []tls.Certificate{cert},
			},
			TLSHandshakeTimeout: 2 * DefaultTimeout,
		}

		// 地址加前缀
		address := make([]string, 0)
		for _, v := range cfg.Address {
			address = append(address, "https://"+v)
		}

		cfg := client.Config{
			Endpoints:               address,
			Transport:               transportTls,
			HeaderTimeoutPerRequest: DefaultTimeout,
			Username:                cfg.Username,
			Password:                cfg.Password,
		}
		cli, err = client.New(cfg)
	} else {
		// 地址加前缀
		address := make([]string, 0)
		for _, v := range cfg.Address {
			address = append(address, "http://"+v)
		}
		cfg := client.Config{
			Endpoints:               address,
			Transport:               client.DefaultTransport,
			HeaderTimeoutPerRequest: DefaultTimeout,
			Username:                cfg.Username,
			Password:                cfg.Password,
		}
		cli, err = client.New(cfg)
	}
	if err != nil {
		return
	}
	// key操作api对象
	keysAPI := client.NewKeysAPI(cli)
	membersAPI := client.NewMembersAPI(cli)

	// 返回etcd v2客户端对象
	clientv2 = &EtcdV2Sdk{
		cli:        cli,
		keysAPI:    keysAPI,
		membersAPI: membersAPI,
	}
	return
}

// List 显示当前path下所有key
func (sdk *EtcdV2Sdk) List(path string) (list []*model.Node, err error) {
	ctx, cancel := sdk.newContext()
	defer cancel()
	rsp, err := sdk.keysAPI.Get(ctx, path, nil)
	if err != nil {
		return
	}
	if rsp.Node == nil || rsp.Node.Dir == false {
		err = model.ERR_KEY_NOT_DIR
		return
	}
	sort.Sort(rsp.Node.Nodes) // 排个序
	for _, v := range rsp.Node.Nodes {
		name := v.Key
		names := strings.Split(v.Key, "/")
		if len(names) > 0 {
			name = names[len(names)-1]
		}
		list = append(list, &model.Node{
			IsDir:   v.Dir,
			Path:    v.Key,
			Name:    name,
			Value:   v.Value,
			Version: 0,
		})
	}
	return
}

// Val 获取path的值
func (sdk *EtcdV2Sdk) Val(path string) (data *model.Node, err error) {
	ctx, cancel := sdk.newContext()
	defer cancel()
	rsp, err := sdk.keysAPI.Get(ctx, path, nil)
	if err != nil {
		return
	}
	if rsp.Node == nil || rsp.Node.Dir == true {
		err = model.ERR_KEY_NOT_DIR
		return
	}
	data = &model.Node{
		IsDir:   rsp.Node.Dir,
		Path:    rsp.Node.Key,
		Name:    rsp.Node.Key,
		Value:   rsp.Node.Value,
		Version: 0,
	}
	return
}

// Add 添加key
func (sdk *EtcdV2Sdk) Add(path string, data []byte) (err error) {
	ctx, cancel := sdk.newContext()
	defer cancel()
	_, err = sdk.keysAPI.Create(ctx, path, string(data))
	if err != nil {
		return
	}
	return
}

// Put 修改key
func (sdk *EtcdV2Sdk) Put(path string, data []byte) (err error) {
	ctx, cancel := sdk.newContext()
	defer cancel()
	_, err = sdk.keysAPI.Update(ctx, path, string(data))
	if err != nil {
		return
	}
	return
}

// Del 删除key
func (sdk *EtcdV2Sdk) Del(path string) (err error) {
	ctx, cancel := sdk.newContext()
	defer cancel()
	_, err = sdk.keysAPI.Delete(ctx, path, nil)
	if err != nil {
		return
	}
	return
}

// Members 获取节点列表
func (sdk *EtcdV2Sdk) Members() (members []*model.Member, err error) {
	ctx, cancel := sdk.newContext()
	defer cancel()
	list, err := sdk.membersAPI.List(ctx)
	if err != nil {
		return
	}
	for _, v := range list {
		members = append(members, &model.Member{
			ID:         v.ID,
			Name:       v.Name,
			PeerURLs:   v.PeerURLs,
			ClientURLs: v.ClientURLs,
			Status:     "",
		})
	}
	return
}

// Close 关闭连接
func (sdk *EtcdV2Sdk) Close() error {
	return nil
}

// 获取请求上下文，有默认超时
func (sdk *EtcdV2Sdk) newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}
