module github.com/mia0x75/halo

go 1.14

require (
	github.com/99designs/gqlgen v0.11.3
	github.com/akhenakh/statgo v0.0.0-20171021021904-3ae2cda264c5
	github.com/cznic/mathutil v0.0.0-20181122101859-297441e03548 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fatih/structs v1.1.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-xorm/builder v0.3.4
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.4
	github.com/gorilla/websocket v1.4.2
	github.com/mia0x75/antlr v0.0.0-20190323140341-bf6915c3dd7b // indirect
	github.com/mia0x75/parser v0.0.0-20190531113551-6fbc203ea218
	github.com/mia0x75/yql v0.0.0-20190325023231-8a8982f46522
	github.com/pingcap/errors v0.11.4 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 // indirect
	github.com/sirupsen/logrus v1.5.0
	github.com/spf13/cobra v1.0.0
	github.com/toolkits/file v0.0.0-20160325033739-a5b3c5147e07
	github.com/vektah/gqlparser v1.3.1
	github.com/vektah/gqlparser/v2 v2.0.1
	golang.org/x/crypto v0.0.0-20200427165652-729f1e841bcc
	xorm.io/core v0.7.3
	xorm.io/xorm v1.0.1
)

replace xorm.io/core => github.com/mia0x75/xorm-core v0.7.3
