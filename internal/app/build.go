package app

import "runtime"

const (
	Version = "2.0.1248"
	Git     = "2026-10-15 13:03:43 +0800 2c65dd1"
	Compile = "2026-10-15 13:04:44 +0800 by go version go1.26.8-X:nodwarf5 linux/amd64"
	Branch  = "master"
	Distro  = "Fedora release 43 (Forty Three)"
	Kernel  = "7.2.4-100.fc43.x86_64"
	Banner  = ` __  __     ______     __         ______    
/\ \_\ \   /\  __ \   /\ \       /\  __ \   
\ \  __ \  \ \  __ \  \ \ \____  \ \ \/\ \  
 \ \_\ \_\  \ \_\ \_\  \ \_____\  \ \_____\ 
  \/_/\/_/   \/_/\/_/   \/_____/   \/_____/ 
`
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
