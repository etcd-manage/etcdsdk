package model

import "github.com/coreos/etcd/etcdserver/etcdserverpb"

// 这里保存sdk使用到的所有结构图

// Node 一个key 目录或文件
type Node struct {
	IsDir   bool   `json:"is_dir,omitempty"`
	Path    string `json:"path,omitempty"`
	Name    string `json:"name,omitempty"`
	Value   string `json:"value,string,omitempty"`
	Version int64  `json:"version,omitempty"`
}

const (
	ROLE_LEADER   = "leader"
	ROLE_FOLLOWER = "follower"

	STATUS_HEALTHY   = "healthy"
	STATUS_UNHEALTHY = "unhealthy"
)

// Member 节点信息
type Member struct {
	*etcdserverpb.Member
	Role   string `json:"role"`
	Status string `json:"status"`
	DbSize int64  `json:"db_size"`
}
