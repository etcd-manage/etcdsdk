package model

const (
	ETCD_VERSION_V2 = "v2"
	ETCD_VERSION_V3 = "v3"
)

// Config etcd 连接配置
type Config struct {
	Version   string   `json:"version,omitempty"`
	Address   []string `json:"address,omitempty"`
	TlsEnable bool     `json:"tls_enable,omitempty"`
	CertFile  string   `json:"cert_file,omitempty"`
	KeyFile   string   `json:"key_file,omitempty"`
	CaFile    string   `json:"ca_file,omitempty"`
	Username  string   `json:"username,omitempty"`
	Password  string   `json:"password,omitempty"`
}
