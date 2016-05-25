package main

type CoreConf struct {
	DbPath     string `json:"dbpath"`
	Debug      bool   `json:"debug"`
	TLS        bool   `json:"tls_enabled"`
	TLSKey     string `json:"tls_private_key"`
	TLSCert    string `json:"tls_certificate"`
	FQDN       string `json:fqdn"`
	ListenIP   string `json:"listen_ip"`
	ListenIPv6 string `json:"listen_ipv6"`
	ListenPort int64  `json:"listen_port"`
	BTCAddr    string `json:"btc_full_node_address"`
	BTCUser    string `json:"btc_rpc_username"`
	BTCPass    string `json:"btc_rpc_password"`
}

type TemplateConf struct {
	FQDN       string
	ListenPort int64
}

type APIResponse struct {
	StatusCode int    `json:"status_code"`
	Success    bool   `json:"success"`
	Message    string `json:"status_message"`
	Version    int64  `json:"version"`

	/*BTCAddr string `json:"btc_addr"`
	UserID  uint64 `json:"user_id"`
	CryptID string `json:"crypt_id"`*/
}

func NewAPIResponse() *APIResponse {
	return &APIResponse{StatusCode: 200, Version: 1, Success: true}
}

type APIRegisterResponse struct {
	StatusCode int    `json:"status_code"`
	Success    bool   `json:"success"`
	Message    string `json:"status_message"`
	Version    int64  `json:"version"`

	BTCAddr string `json:"btc_addr"`
	UserID  uint64 `json:"user_id"`
}
type APIChallengeResponse struct {
	Challenge   string `json:"challenge"`
	ChallengeID uint64 `json:"challenge_id"`
	UserID      uint64 `json:"user_id"`
	StatusCode  int    `json:"status_code"`
	Success     bool   `json:"success"`
	Message     string `json:"status_message"`
	Version     int64  `json:"version"`
}

type APICryptResponse struct {
	StatusCode int    `json:"status_code"`
	Success    bool   `json:"success"`
	Message    string `json:"status_message"`
	Version    int64  `json:"version"`

	CryptPayload Crypt `json:"crypt"`
}

type Crypt struct {
	UserID          uint64 `json:"-"`
	CryptID         string `json:"crypt_id"`
	CipherText      string `json:"ciphertext"`
	CreateTimeStamp int64  `json:"crypt_timestamp"`
	Description     string `json:"crypt_description"`
	IsDestroyed     bool   `json:"is_crypt_destroyed"`
	LastCheckIn     int64  `json:"last_checkin"`
	CheckInDuration int64  `json:"check_in_duration"`
	MissCount       int64  `json:"miss_count"`
}

type ClientRequest struct {
	UserID          uint64 `json:"user_id"`
	CryptContent    string `json:"crypt_content"`
	Challenge       string `json:"challenge"`
	ChallengeID     uint64 `json:"challenge_id"`
	CryptID         string `json:"crypt_id"`
	PublicKey       string `json:"public_key"`
	Fingerprint     string `json:"fingerprint"`
	Description     string `json:"description"`
	CheckInDuration int64  `json:"checkin_duration"`
	MissCount       int64  `json:"miss_count"`
}

type Account struct {
	UserID               uint64 `json:"user_id"`
	PublicKey            string `json:"public_key"`
	PublicKeyFingerprint string `json:"fingerprint"`
	BTCAddr              string `json:"btc_addr"`
	ByteSize             int64  `json:byte_size"`
	ByteBudget           int64  `json:byte_budget"`
}

type Challenge struct {
	ChallengeID uint64 `json:"challenge_id"`
	Challenge   string `json:"challenge"`
	UserID      uint64 `json:"user_id"`
}
