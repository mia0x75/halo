package models

// CPUStats CPU 状态信息，与 GraphQL schema.graphql 中的 CPUStats 对应。
// 原先复用 github.com/akhenakh/statgo.CPUStats（CGO + libstatgrab），现改为纯 Go 实现，字段保持兼容。
type CPUStats struct {
	User      float64 `json:"User"`
	Kernel    float64 `json:"Kernel"`
	Idle      float64 `json:"Idle"`
	IOWait    float64 `json:"IOWait"`
	Swap      float64 `json:"Swap"`
	Nice      float64 `json:"Nice"`
	LoadMin1  float64 `json:"LoadMin1"`
	LoadMin5  float64 `json:"LoadMin5"`
	LoadMin15 float64 `json:"LoadMin15"`
}

// HostInfos 主机信息，与 GraphQL schema.graphql 中的 HostInfos 对应。
type HostInfos struct {
	OSName    string `json:"OSName"`
	OSRelease string `json:"OSRelease"`
	OSVersion string `json:"OSVersion"`
	Platform  string `json:"Platform"`
	HostName  string `json:"HostName"`
	NCPUs     int    `json:"NCPUs"`
	MaxCPUs   int    `json:"MaxCPUs"`
	BitWidth  int    `json:"BitWidth"`
}

// ProcessStats 进程统计信息，与 GraphQL schema.graphql 中的 ProcessStats 对应.
type ProcessStats struct {
	Total    int `json:"Total"`
	Running  int `json:"Running"`
	Sleeping int `json:"Sleeping"`
	Stopped  int `json:"Stopped"`
	Zombie   int `json:"Zombie"`
}

// MemStats 内存与交换分区信息（单位：字节），与 GraphQL schema.graphql 中的 MemStats 对应.
type MemStats struct {
	Total     int `json:"Total"`
	Free      int `json:"Free"`
	Used      int `json:"Used"`
	Cache     int `json:"Cache"`
	SwapTotal int `json:"SwapTotal"`
	SwapUsed  int `json:"SwapUsed"`
	SwapFree  int `json:"SwapFree"`
}
