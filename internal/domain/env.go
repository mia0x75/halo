package domain

// CPUStats is the host CPU sample produced by the collector layer.
type CPUStats struct {
	User      float64
	Kernel    float64
	Idle      float64
	IOWait    float64
	Swap      float64
	Nice      float64
	LoadMin1  float64
	LoadMin5  float64
	LoadMin15 float64
}

// MemStats is the memory and swap sample, values in bytes.
type MemStats struct {
	Total     int
	Free      int
	Used      int
	Cache     int
	SwapTotal int
	SwapUsed  int
	SwapFree  int
}

// HostInfos describes the runtime host environment.
type HostInfos struct {
	NCPUs     int
	MaxCPUs   int
	BitWidth  int
	OSName    string
	OSRelease string
	OSVersion string
	Platform  string
	HostName  string
}

// ProcessStats counts processes grouped by OS state.
type ProcessStats struct {
	Total    int
	Running  int
	Sleeping int
	Stopped  int
	Zombie   int
}
