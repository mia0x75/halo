package tools

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

func init() {
	cfg := flag.String("c", "../cfg.json", "configuration file")
	flag.Parse()

	g.ParseConfig(*cfg)

	if err := g.InitDB(); err != nil {
		os.Exit(0)
	}
	caches.Init()
}

func TestGenerator_ParseToken(t *testing.T) {
	//token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTExNDYwOTcsImFkbWluIjp0cnVlLCJ1c2VySWQiOjF9.jKF2b0B-stLwSwGuhq4Yf7mqc19QfHesdojbEQc_N0k"
	tk := New()
	tokenStr, err := tk.CreateToken(&models.User{})
	fmt.Printf("token:%s\n", tokenStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	claims, err := tk.ParseToken(tokenStr)
	if err == nil {
		// 验证通过
		uuid := claims.UserUUID
		fmt.Printf("claims :%v\n", claims)
		fmt.Printf("uuid :%v\n", uuid)
	} else {
		fmt.Printf("err :%v\n", err.Error())
	}

}

func Test_role(t *testing.T) {
	user := caches.UsersMap.GetById(1)
	fmt.Print(user)
	roleIDs := caches.EdgesMap.GetDescendants(g.EdgeTypeUserToRole, user.UserID)
	fmt.Println(roleIDs)
	roles := caches.RolesMap.GetByIds(roleIDs)
	fmt.Print(roles)
	for _, role := range roles {
		fmt.Printf("role:%v\n", role)
	}
}
