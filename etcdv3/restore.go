package etcdv3

import (
	"context"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
)

const (
	DEFAULT_DIR_VALUE = "etcdv3_dir_$2H#%gRe3*t"
)

// Restore 修复v1版本和e3w标记目录用的特殊key
func (sdk *EtcdV3Sdk) Restore() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	// 查找出所有标记目录key
	resp, err := sdk.cli.Get(ctx, "", clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		return
	}
	for _, v := range resp.Kvs {
		key := string(v.Key)
		log.Println(key)
		txn := sdk.cli.Txn(ctx)
		txn.If(
			clientv3.Compare(
				clientv3.Value(key),
				"=",
				DEFAULT_DIR_VALUE,
			),
		).Then(
			clientv3.OpDelete(key),
		)
		rr, err := txn.Commit()
		if err != nil {
			return err
		}
		log.Println(rr.Succeeded)
	}

	return
}
