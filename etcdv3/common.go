package etcdv3

import (
	"encoding/base64"
	"io/ioutil"
	"strings"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/etcd-manage/etcdsdk/model"
)

// ConvertToPath 处理etcd3 的key为目录形式 - path只能是/结尾或为空
func (sdk *EtcdV3Sdk) ConvertToPath(path string, keys []*mvccpb.KeyValue) (list []*model.Node, err error) {
	keyMapVal := make(map[string]*model.Node, 0)
	keyMapPath := make(map[string]*model.Node, 0)

	for _, val := range keys {
		if ok := strings.HasPrefix(string(val.Key), path); ok == true {
			key := string(val.Key)[len(path):]
			// 判断是否是//开头，如果是，则本级目录为//
			if strings.HasPrefix(key, "//") == true {
				fullKey := path + "/"
				keyMapPath["/"] = &model.Node{
					IsDir:   true,
					Path:    fullKey,
					Name:    "/",
					Value:   string(val.Value),
					Version: 0,
				}
				continue
			}
			// 处理path为/情况
			if path == "" && strings.HasPrefix(key, "/") {
				keyMapPath["/"] = &model.Node{
					IsDir:   true,
					Path:    "/",
					Name:    "/",
					Value:   string(val.Value),
					Version: 0,
				}
				continue
			}
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
				if path == "/" {
					fullKey = path + key
				} else if path == "" {
					fullKey = key
				}
				lastIndex := strings.LastIndex(fullKey, "/")
				key = fullKey[lastIndex+1:]
				// log.Println(path, " -- ", key)
				keyMapPath[key] = &model.Node{
					IsDir:   true,
					Path:    fullKey,
					Name:    key,
					Value:   string(val.Value),
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

// 存储证书文件
func writeCa(cfg *model.Config) (certFilePath, keyFilePath, caFilePath string, err error) {
	certFilePath = "./cert.pem"
	certBody, err := base64.StdEncoding.DecodeString(cfg.CertFile)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(certFilePath, certBody, 0755)
	if err != nil {
		return
	}
	keyFilePath = "./key.pem"
	keyBody, err := base64.StdEncoding.DecodeString(cfg.KeyFile)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(keyFilePath, keyBody, 0755)
	if err != nil {
		return
	}
	caFilePath = "./ca.pem"
	caBody, err := base64.StdEncoding.DecodeString(cfg.CaFile)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(caFilePath, caBody, 0755)
	if err != nil {
		return
	}
	return
}
