package tools

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

var cfg = flag.String("c", "../cfg.json", "configuration file")

// configAvailable is set by TestMain when the config file exists; tests that
// need the JWT secret skip otherwise.
var configAvailable bool

func TestMain(m *testing.M) {
	flag.Parse()
	if _, err := os.Stat(*cfg); err == nil {
		g.ParseConfig(*cfg)
		configAvailable = true
	}
	os.Exit(m.Run())
}

func TestGenerator_ParseToken(t *testing.T) {
	if !configAvailable {
		t.Skip("no configuration file")
	}
	tk := New()
	tokenStr, err := tk.CreateToken(&models.User{})
	fmt.Printf("token:%s\n", tokenStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	claims, err := tk.ParseToken(tokenStr)
	if err != nil {
		fmt.Printf("err :%v\n", err.Error())
		return
	}
	// 验证通过
	var uuid string
	if claims.User != nil {
		uuid = claims.User.UUID
	}
	fmt.Printf("claims :%v\n", claims)
	fmt.Printf("uuid :%v\n", uuid)
}
