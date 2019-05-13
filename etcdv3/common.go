package etcdv3

import (
	"strings"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/etcd-manage/etcdsdk/model"
)

// ConvertToPath 处理etcd3 的key为目录形式
func (sdk *EtcdV3Sdk) ConvertToPath(path string, keys []*mvccpb.KeyValue) (list []*model.Node, err error) {
	path = strings.TrimRight(path, "/")
	keyMap := make(map[string]*model.Node, 0)
	for _, val := range keys {
		if ok := strings.HasPrefix(string(val.Key), path); ok == true {
			key := string(val.Key)[len(path):]
			// 查找下一个/位置
			i := strings.Index(key, "/")
			if i == -1 { // 未查询到则证明是key，而不是目录
				keyMap[key] = &model.Node{
					IsDir:   false,
					Path:    string(val.Key),
					Name:    key,
					Value:   string(val.Value),
					Version: val.Version,
				}
			} else {
				key = key[:i]
				// 等于当前path，不返回
				if key == "" {
					continue
				}
				keyMap[key] = &model.Node{
					IsDir:   true,
					Path:    path + key,
					Name:    key,
					Value:   "",
					Version: 0,
				}
			}
		}
	}
	for _, val := range keyMap {
		list = append(list, val)
	}
	return
}
