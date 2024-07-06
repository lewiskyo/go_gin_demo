package bo

type Msg struct {
	Timestamp int64  `json:"timestamp"`
	Uid       int64  `json:"uid"`
	Type      uint8  `json:"type"`
	Region    string `json:"region"`
	Device    string `json:"device"`
	Ip        string `json:"ip"`
	Network   string `json:"network"`
	Version   uint64 `json:"version"`
}
