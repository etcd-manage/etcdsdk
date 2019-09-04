package etcdv2

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/etcd-manage/etcdsdk/model"
)

// 存储证书文件
func writeCa(cfg *model.Config, etcdId int32) (certFilePath, keyFilePath, caFilePath string, err error) {
	certFilePath = fmt.Sprintf("./%d_cert.pem", etcdId)
	certBody, err := base64.StdEncoding.DecodeString(cfg.CertFile)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(certFilePath, certBody, 0755)
	if err != nil {
		return
	}
	keyFilePath = fmt.Sprintf("./%d_key.pem", etcdId)
	keyBody, err := base64.StdEncoding.DecodeString(cfg.KeyFile)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(keyFilePath, keyBody, 0755)
	if err != nil {
		return
	}
	caFilePath = fmt.Sprintf("./%d_ca.pem", etcdId)
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
