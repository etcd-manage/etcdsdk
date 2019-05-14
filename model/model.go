package model

import "github.com/coreos/etcd/etcdserver/etcdserverpb"

// 这里保存sdk使用到的所有结构图

// Node 一个key 目录或文件
type Node struct {
	IsDir   bool   `json:"is_dir"`
	Path    string `json:"path"`
	Name    string `json:"name"`
	Value   string `json:"value,string"`
	Version int64  `json:"version"`
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
