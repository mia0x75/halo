module github.com/mia0x75/halo

require (
	cloud.google.com/go v0.34.0 // indirect
	github.com/99designs/gqlgen v0.8.3
	github.com/akhenakh/statgo v0.0.0-20171021021904-3ae2cda264c5
	github.com/cznic/mathutil v0.0.0-20160613104831-78ad7f262603 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fatih/structs v1.1.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/builder v0.3.2
	github.com/go-xorm/core v0.6.2
	github.com/go-xorm/xorm v0.7.1
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.6.1
	github.com/gorilla/websocket v1.4.0
	github.com/hashicorp/golang-lru v0.5.1 // indirect
	github.com/machinebox/graphql v0.2.2
	github.com/matryer/is v1.2.0 // indirect
	github.com/mia0x75/antlr v0.0.0-20190323140341-bf6915c3dd7b // indirect
	github.com/mia0x75/parser v0.0.0-20190503015245-26d8edb95c53
	github.com/mia0x75/yql v0.0.0-20190325023231-8a8982f46522
	github.com/pingcap/check v0.0.0-20171206051426-1c287c953996 // indirect
	github.com/pingcap/errors v0.11.0 // indirect
	github.com/sirupsen/logrus v1.3.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/toolkits/file v0.0.0-20160325033739-a5b3c5147e07
	github.com/vektah/gqlparser v1.1.2
	golang.org/x/crypto v0.0.0-20190130090550-b01c7a725664
	golang.org/x/sys v0.0.0-20190222072716-a9d3bda3a223 // indirect
)

replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.36.0
	github.com/go-xorm/core => github.com/mia0x75/xorm-core v0.6.3

	golang.org/x/build => github.com/golang/build v0.0.0-20190228010158-44b79b8774a7
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190227175134-215aa809caaf
	golang.org/x/exp => github.com/golang/exp v0.0.0-20190221220918-438050ddec5e
	golang.org/x/lint => github.com/golang/lint v0.0.0-20190227174305-5b3e6a55c961
	golang.org/x/net => github.com/golang/net v0.0.0-20190227160552-c95aed5357e7
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20190226205417-e64efc72b421
	golang.org/x/perf => github.com/golang/perf v0.0.0-20190124201629-844a5f5b46f4
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190227155943-e225da77a7e6
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190226215855-775f8194d0f9
	golang.org/x/text => github.com/golang/text v0.3.0
	golang.org/x/time => github.com/golang/time v0.0.0-20181108054448-85acf8d2951c
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190227232517-f0a709d59f0f
	google.golang.org/api => github.com/googleapis/google-api-go-client v0.1.0
	google.golang.org/appengine => github.com/golang/appengine v1.4.0
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20190227213309-4f5b463f9597
	google.golang.org/grpc => github.com/grpc/grpc-go v1.19.0
)
