package etcdv3

import (
	"strings"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/etcd-manage/etcdsdk/model"
)

// ConvertToPath 处理etcd3 的key为目录形式
func (sdk *EtcdV3Sdk) ConvertToPath(path string, keys []*mvccpb.KeyValue) (list []*model.Node, err error) {
	path = strings.TrimRight(path, "/")
	keyMapVal := make(map[string]*model.Node, 0)
	keyMapPath := make(map[string]*model.Node, 0)

	for _, val := range keys {
		if ok := strings.HasPrefix(string(val.Key), path); ok == true {
			key := string(val.Key)[len(path):]
			// 查找下一个/位置
			i := strings.Index(key, "/")
			// 截取后第一个字符是/的情况
			if i == 0 {
				key = key[1:]
				i = strings.Index(key, "/")
			}
			// 未查找到，则为key，而不是目录
			if i == -1 { // 未查询到则证明是key，而不是目录 则取完整路径最后一个/之后的部分作为name
				lastIndex := strings.LastIndex(string(val.Key), "/")
				key = string(val.Key)[lastIndex+1:]
				keyMapVal[key] = &model.Node{
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
				fullKey := path + "/" + key
				lastIndex := strings.LastIndex(fullKey, "/")
				key = fullKey[lastIndex+1:]
				keyMapPath[key] = &model.Node{
					IsDir:   true,
					Path:    fullKey,
					Name:    key,
					Value:   "",
					Version: 0,
				}
			}
		}
	}
	for _, val := range keyMapPath {
		list = append(list, val)
	}
	for _, val := range keyMapVal {
		list = append(list, val)
	}

	return
}
