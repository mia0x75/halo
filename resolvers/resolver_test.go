package resolvers

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"testing"

	"github.com/machinebox/graphql"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/models"
)

var LoginData struct {
	Login *gqlapi.LoginPayload
}

func undo() {

	command := "mysql -utickets -ptickets -h 10.100.11.81 -P3309 tickets <../undo.sql"
	cmd := exec.Command("/bin/sh", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("undo out : %v", string(out))
}

var cfg = flag.String("c", "../cfg.json", "configuration file")

func init() {
	undo()
	flag.Parse()
	g.ParseConfig(*cfg)

	if err := g.InitDB(); err != nil {
		os.Exit(0)
	}
	caches.Init()

	Login()
}

func Do(req *graphql.Request, resp interface{}) (err error) {
	option := graphql.WithHTTPClient(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	// create a client (safe to share across requests)
	client := graphql.NewClient("https://127.0.0.1:3800/api/query", option)
	// define a Context for the request
	ctx := context.Background()

	err = client.Run(ctx, req, resp)

	return
}

func TestTicketResolverTicket(t *testing.T) {
	req := graphql.NewRequest(`
    query ($ticketUUID: ID!) {
		ticket (
			id: $ticketUUID
		) {
			Database
			Subject
		}
	}`)

	req.Var("ticketUUID", "0aa345ab-3235-43f0-a4c3-4d7f10b62e2f")

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authentication", LoginData.Login.Token)

	var respData struct {
		Ticket *models.Ticket
	}
	if err := Do(req, &respData); err != nil {

	}
	require.Equal(t, "sbtest", respData.Ticket.Database)
}

func TestTicketResolverTickets(t *testing.T) {
	req := graphql.NewRequest(`
    query ($first: Int!) {
		tickets (
			first: $first
		) {
			edges {
				node {
					Subject
				}
			}
		}
	}`)

	req.Var("first", 10)

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authentication", LoginData.Login.Token)

	var respData struct {
		Tickets *gqlapi.TicketConnection
	}
	if err := Do(req, &respData); err != nil {

	}
}

func TestCreateTicket(t *testing.T) {
	req := graphql.NewRequest(`
    mutation ($clusterid: String!,$subject:String!,$db:String!,$content:String!,$reviewerid:String!){
		createTicket(
			input: {
				Subject: $subject,
				Database:$db,
				Content: $content,
				ClusterUUID:$clusterid,
				ReviewerUUID:$reviewerid,
			}
		){
			Database
			Subject
			TicketUUID
		}
	}`)

	type Input struct {
		Subject      string
		Database     string
		Content      string
		ClusterUUID  string
		ReviewerUUID string
		Error        string
	}
	//
	cases := [7]Input{}

	//群集不存在
	//错误代码: %s, 错误信息: 群集(uuid=%s)不存在。
	cases[0] = Input{
		Subject:      "graphAPI 单元测试 update",
		Database:     "sbtest",
		Content:      "create table t_create(id2 int primary key)",
		ClusterUUID:  "zzzz",
		ReviewerUUID: "e70e78bb-9d08-405d-a0ed-266ec703de19",
	}
	cases[0].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 群集(uuid=%s)不存在。", cases[0].ClusterUUID)

	//群集不可用
	cases[1] = Input{
		Subject:      "graphAPI 单元测试 update",
		Database:     "sbtest",
		Content:      "create table t_create(id2 int primary key)",
		ClusterUUID:  "cc3d1347-aa8a-4a41-bbdb-8aa5ba0c78a1", //当前cluster 的status 为6
		ReviewerUUID: "e70e78bb-9d08-405d-a0ed-266ec703de19",
	}
	cases[1].Error = fmt.Sprintf("graphql: 错误代码: 2001, 错误信息: 群集(uuid=%s)不可用。", cases[1].ClusterUUID)

	//群集和用户没有关联
	cases[2] = Input{
		Subject:      "graphAPI 单元测试 update",
		Database:     "sbtest",
		Content:      "create table t_create(id2 int primary key)",
		ClusterUUID:  "a42d9b05-eab7-4e95-8d9c-c4e37e12383c", //当前cluster 的status 为1
		ReviewerUUID: "e70e78bb-9d08-405d-a0ed-266ec703de19",
	}
	cases[2].Error = fmt.Sprintf("graphql: 错误代码: 1403, 错误信息: 用户(uuid=e70e78bb-9d08-405d-a0ed-266ec703de19)没有关联群集(uuid=%s)。", cases[2].ClusterUUID)

	// 审核人不存在
	cases[3] = Input{
		Subject:      "graphAPI 单元测试 update",
		Database:     "sbtest",
		Content:      "create table t_create(id2 int primary key)",
		ClusterUUID:  "dbf8ead1-c633-4ac4-9f30-dbf889a197aa", //当前cluster 的status 为1
		ReviewerUUID: "zzzzz",
	}
	cases[3].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 审核用户(uuid=%s)不存在。", cases[3].ReviewerUUID)

	//审核人状态异常
	cases[4] = Input{
		Subject:      "graphAPI 单元测试 update",
		Database:     "sbtest",
		Content:      "create table t_create(id2 int primary key)",
		ClusterUUID:  "dbf8ead1-c633-4ac4-9f30-dbf889a197aa", //当前cluster 的status 为1
		ReviewerUUID: "fcd3ece1-ddc9-423d-809b-0f00d9508bae",
	}
	cases[4].Error = fmt.Sprintf("graphql: 错误代码: 2026, 错误信息: 审核用户(uuid=%s)状态异常。", cases[4].ReviewerUUID)

	// 用户没有和审核人关联

	cases[5] = Input{
		Subject:      "graphAPI 单元测试 update",
		Database:     "sbtest",
		Content:      "create table t_create(id2 int primary key)",
		ClusterUUID:  "dbf8ead1-c633-4ac4-9f30-dbf889a197aa", //当前cluster 的status 为1
		ReviewerUUID: "f0ebc1b4-c776-4613-8786-a40c1b02d3cb",
	}
	cases[5].Error = fmt.Sprintf("graphql: 错误代码: 2026, 错误信息: 用户(uuid=e70e78bb-9d08-405d-a0ed-266ec703de19)没有关联审核用户(uuid=%s)。", cases[5].ReviewerUUID)

	//正常创建工单
	cases[6] = Input{
		Subject:      "graphAPI 单元测试23",
		Database:     "sbtest",
		Content:      "create table t_create(id int primary key)",
		ClusterUUID:  "a42d9b05-eab7-4e95-8d9c-c4e37e12383c",
		ReviewerUUID: "e70e78bb-9d08-405d-a0ed-266ec703de19",
	}
	cases[6].Error = ""

	var respData struct {
		CreateTicket struct {
			Database   string
			Subject    string
			TicketUUID string
		}
	}
	if err := Do(req, &respData); err != nil {
		log.Printf("err: %v", err)
	}
	for _, k := range cases {
		// mutation ($clusterid: String!,$subject:String!,$db:String!,$content:String!,%reviewerid:String!){
		req.Var("clusterid", k.ClusterUUID)
		req.Var("subject", k.Subject)
		req.Var("content", k.Content)
		req.Var("reviewerid", k.ReviewerUUID)
		req.Var("db", k.Database)
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authentication", LoginData.Login.Token)
		err := Do(req, &respData)

		if err != nil {
			require.Equal(t, k.Error, err.Error())
		} else {
			// fmt.Printf("respData.UpdateTicket :%v", respData.CreateTicket)
			require.Equal(t, k.Subject, respData.CreateTicket.Subject)
		}

	}
}

func TestUpdateTicket(t *testing.T) {
	req := graphql.NewRequest(`
    mutation ($id:ID!,$clusterid: String!,$subject:String!,$db:String!,$content:String!,$reviewerid:String!){
		updateTicket(
			input: {
				TicketUUID:$id,
				Subject: $subject,
				Database: $db,
				Content:$content,
				ClusterUUID:$clusterid,
				ReviewerUUID:$reviewerid 
			}
		){
			Database
			Subject
			TicketUUID
		}
	}`)

	type Input struct {
		TicketUUID   string
		Subject      string
		Database     string
		Content      string
		ClusterUUID  string
		ReviewerUUID string
		Error        string
	}

	cases := [9]Input{}
	//群集不存在
	//错误代码: %s, 错误信息: 群集(uuid=%s)不存在。
	cases[0] = Input{
		TicketUUID:   "e70e78bb-9d08-405d-a0ed-266ec703de19",
		Subject:      "graphAPI 单元测试 update",
		Database:     "sbtest",
		Content:      "create table t_create(id2 int primary key)",
		ClusterUUID:  "zzzz",
		ReviewerUUID: "e70e78bb-9d08-405d-a0ed-266ec703de19",
	}
	cases[0].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 群集(uuid=%s)不存在。", cases[0].ClusterUUID)

	//群集不可用
	cases[1] = Input{
		TicketUUID:   "e70e78bb-9d08-405d-a0ed-266ec703de19",
		Subject:      "graphAPI 单元测试 update",
		Database:     "sbtest",
		Content:      "create table t_create(id2 int primary key)",
		ClusterUUID:  "cc3d1347-aa8a-4a41-bbdb-8aa5ba0c78a1", //当前cluster 的status 为6
		ReviewerUUID: "e70e78bb-9d08-405d-a0ed-266ec703de19",
	}
	cases[1].Error = fmt.Sprintf("graphql: 错误代码: 2001, 错误信息: 群集(uuid=%s)不可用。", cases[1].ClusterUUID)

	//群集和用户没有关联
	cases[2] = Input{
		TicketUUID:   "12d23647-352e-4ba6-8d7a-ed9364c89050",
		Subject:      "graphAPI 单元测试 update",
		Database:     "sbtest",
		Content:      "create table t_create(id2 int primary key)",
		ClusterUUID:  "a42d9b05-eab7-4e95-8d9c-c4e37e12383c", //当前cluster 的status 为1
		ReviewerUUID: "e70e78bb-9d08-405d-a0ed-266ec703de19",
	}
	cases[2].Error = fmt.Sprintf("graphql: 错误代码: 1403, 错误信息: 用户(uuid=e70e78bb-9d08-405d-a0ed-266ec703de19)没有关联群集(uuid=%s)。", cases[2].ClusterUUID)

	// 审核人不存在
	cases[3] = Input{
		TicketUUID:   "e70e78bb-9d08-405d-a0ed-266ec703de19",
		Subject:      "graphAPI 单元测试 update",
		Database:     "sbtest",
		Content:      "create table t_create(id2 int primary key)",
		ClusterUUID:  "dbf8ead1-c633-4ac4-9f30-dbf889a197aa", //当前cluster 的status 为1
		ReviewerUUID: "zzzzz",
	}
	cases[3].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 审核用户(uuid=%s)不存在。", cases[3].ReviewerUUID)

	//审核人状态异常
	cases[4] = Input{
		TicketUUID:   "e70e78bb-9d08-405d-a0ed-266ec703de19",
		Subject:      "graphAPI 单元测试 update",
		Database:     "sbtest",
		Content:      "create table t_create(id2 int primary key)",
		ClusterUUID:  "dbf8ead1-c633-4ac4-9f30-dbf889a197aa", //当前cluster 的status 为1
		ReviewerUUID: "fcd3ece1-ddc9-423d-809b-0f00d9508bae",
	}
	cases[4].Error = fmt.Sprintf("graphql: 错误代码: 2026, 错误信息: 审核用户(uuid=%s)的当前状态异常。", cases[4].ReviewerUUID)

	// 用户没有和审核人关联

	cases[5] = Input{
		TicketUUID:   "e70e78bb-9d08-405d-a0ed-266ec703de19",
		Subject:      "graphAPI 单元测试 update",
		Database:     "sbtest",
		Content:      "create table t_create(id2 int primary key)",
		ClusterUUID:  "dbf8ead1-c633-4ac4-9f30-dbf889a197aa", //当前cluster 的status 为1
		ReviewerUUID: "f0ebc1b4-c776-4613-8786-a40c1b02d3cb",
	}
	cases[5].Error = fmt.Sprintf("graphql: 错误代码: 2026, 错误信息: 用户(uuid=e70e78bb-9d08-405d-a0ed-266ec703de19)没有关联审核用户(uuid=%s)。", cases[5].ReviewerUUID)

	//正常update

	cases[6] = Input{
		TicketUUID:   "82cd4f3d-2995-436d-b797-d68d622ce617",
		Subject:      "graphAPI 单元测试12267",
		Database:     "sbtest",
		Content:      "create table t_create(id int primary key)",
		ClusterUUID:  "dbf8ead1-c633-4ac4-9f30-dbf889a197aa",
		ReviewerUUID: "e70e78bb-9d08-405d-a0ed-266ec703de19",
	}
	cases[6].Error = ""

	//工单不存在
	cases[7] = Input{
		TicketUUID:   "82cd4f3d-2995-436d-b797-d68d622ce617-----",
		Subject:      "graphAPI 单元测试122",
		Database:     "sbtest",
		Content:      "create table t_create(id int primary key)",
		ClusterUUID:  "dbf8ead1-c633-4ac4-9f30-dbf889a197aa",
		ReviewerUUID: "e70e78bb-9d08-405d-a0ed-266ec703de19",
	}
	cases[7].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 工单(uuid=%s)不存在。", cases[7].TicketUUID)

	//错误代码: %s, 错误信息: 已执行或执行失败的工单不可编辑。
	//626d8fb7-542b-4e20-ae1e-2676014e3c5c
	cases[8] = Input{
		TicketUUID:   "d9f81cc2-1794-4d9d-a5fb-f77e3c7f5733",
		Subject:      "graphAPI 单元测试123332",
		Database:     "sbtest",
		Content:      "create table t_create(id int primary key)",
		ClusterUUID:  "dbf8ead1-c633-4ac4-9f30-dbf889a197aa",
		ReviewerUUID: "e70e78bb-9d08-405d-a0ed-266ec703de19",
	}
	cases[8].Error = fmt.Sprintf("graphql: 错误代码: 1200, 错误信息: 已执行或执行失败的工单不可编辑。")

	// cases[9] = Input{
	// 	TicketUUID:   "626d8fb7-542b-4e20-ae1e-2676014e3c5c",
	// 	Subject:      "graphAPI 单元测试123332",
	// 	Database:     "sbtest",
	// 	Content:      "create table t_create(id int primary key)",
	// 	ClusterUUID: "dbf8ead1-c633-4ac4-9f30-dbf889a197aa",
	// 	ReviewerUUID: "e70e78bb-9d08-405d-a0ed-266ec703de19",
	// }
	// cases[9].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 已执行或执行失败的工单不可编辑。")

	var respData struct {
		UpdateTicket struct {
			Database string
			Subject  string
		}
	}
	for _, k := range cases {
		//($id:ID!,$clusterid: String!,$subject:String!,$db:String!,$content:String!,$reviewerid:String)
		req.Var("clusterid", k.ClusterUUID)
		req.Var("subject", k.Subject)
		req.Var("content", k.Content)
		req.Var("id", k.TicketUUID)
		req.Var("reviewerid", k.ReviewerUUID)
		req.Var("db", k.Database)
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authentication", LoginData.Login.Token)
		err := Do(req, &respData)

		if err != nil {
			require.Equal(t, k.Error, err.Error())
		} else {
			fmt.Printf("respData.UpdateTicket :%v", respData.UpdateTicket)
			require.Equal(t, k.Subject, respData.UpdateTicket.Subject)
		}

	}
}

func TestPatchTicket(t *testing.T) {
	req := graphql.NewRequest(`
    mutation ($id: ID!,$status: Int!) {
		patchTicketStatus(
			input :{
				TicketUUID: $id,
				Status:     $status
			}
		) 
	}`)

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authentication", LoginData.Login.Token)
	var respData struct {
		ok  bool
		err error
	}

	/* 1、别人的工单
	   2、我的工单，但是status 初始状态 1
	   3、我的工单，status初始状态为2，目标状态1
	   4、合法值
	*/
	cases := make(map[string]int)
	cases["626d8fb7-542b-4e20-ae1e-2676014e3c5c"] = 6 //合法值
	cases["0bc1f001-e9ff-4e1a-9635-b840ad069154"] = 3 //status初始状态1
	cases["57ac663b-d02e-4cf1-81a8-d65f1b4308f6"] = 5 //别人的工单
	cases["6c93ab67-3c63-4909-a8d4-f2d9b5ae2753"] = 1 //目标status 1

	var err error
	for k, v := range cases {
		req.Var("id", k)
		req.Var("status", v)
		if err = Do(req, &respData); err != nil {

		}
		if k != "626d8fb7-542b-4e20-ae1e-2676014e3c5c" {
			require.NotEqual(t, nil, err)
		} else {
			log.Println(respData)
			require.Equal(t, false, respData.ok)
		}

	}
}

func TestRemoveTickets(t *testing.T) {
	req := graphql.NewRequest(`
    mutation ($id: ID!) {
		removeTicket (
			id: $id
		)
	}`)
	//case : 错误信息
	cases := []map[string]string{
		{"57ac663b-d02e-4cf1-81a8-d65f1b436": "graphql: 错误代码: 1404, 错误信息: 工单(uuid=57ac663b-d02e-4cf1-81a8-d65f1b436)不存在。"},
		{"6c93ab67-3c63-4909-a8d4-f2d9b5ae2753": "graphql: 错误代码: 1403, 错误信息: 只有工单(uuid=6c93ab67-3c63-4909-a8d4-f2d9b5ae2753)的发起人可以删除工单。"},
		{"626d8fb7-542b-4e20-ae1e-2676014e3c5c": "graphql: 错误代码: 1200, 错误信息: 已执行或执行失败的工单不可删除。"},
		{"da5be394-cf8b-40d5-be81-89ac5a8ce4ae": ""},
	}

	var respData struct {
		RemoveTicket bool
	}
	for _, m := range cases {
		for k, v := range m {
			req.Var("id", k)
			req.Header.Set("Cache-Control", "no-cache")
			req.Header.Set("Authentication", LoginData.Login.Token)
			err := Do(req, &respData)
			if err != nil {
				require.Equal(t, v, err.Error())
			} else {
				require.Equal(t, true, respData.RemoveTicket)
			}
		}
	}

}

//user
func TestRegister(t *testing.T) {
	req := graphql.NewRequest(`
    mutation ($email: String!,$password: String!) {
		register (
			input :{
				Email: $email,
				Password: $password
			}
		) {
  			Email
 		}
	}`)

	cases := []map[string]string{
		{"tuandaidba178@tuandai.com": "graphql: 错误代码: 2005, 错误信息: 账号(email=tuandaidba178@tuandai.com)当前状态是等待验证。"}, //
		{"tzzz232322@qq.com": "tzzz232322@qq.com"}, // 正常注册
		{"root@dba.com": "graphql: 错误代码: 2009, 错误信息: 账号(email=root@dba.com)已经被注册。"},                             // 已经被注册
		{"locked@localhost.localhost": "graphql: 错误代码: 2016, 错误信息: 账号(email=locked@localhost.localhost)已经被禁用。"}, // 已经被注册
	}
	var respData struct {
		Register struct {
			Email string
		}
	}
	for _, m := range cases {
		for k, v := range m {
			req.Var("email", k)
			req.Var("password", k)
			req.Header.Set("Cache-Control", "no-cache")
			req.Header.Set("Authentication", LoginData.Login.Token)
			err := Do(req, &respData)
			if err != nil {
				require.Equal(t, v, err.Error())
			} else {
				require.Equal(t, k, respData.Register.Email)
			}
		}
	}
}

func TestCreateUser(t *testing.T) {
	req := graphql.NewRequest(`
    mutation createUser($email:String!,$password:String!,$name:String!,$roleid:[String!]!,$clusterid:[String!]!,
		$reviewerid:[String!]!,$avatarid:String!,$status:Int!){
		createUser(
		  input:{
			Email:$email,
			Password:$password,
			Name:$name,
			RoleUUIDs:$roleid,
			ClusterUUIDs:$clusterid,
			ReviewerUUIDs:$reviewerid,
			AvatarUUID:$avatarid,
			Status:$status
			}
		){
		  Name
		  Email
		}
	  }`)
	type Input struct {
		Email         string
		Password      string
		Name          string
		RoleUUIDs     []string
		ClusterUUIDs  []string
		ReviewerUUIDs []string
		AvatarUUID    string
		Status        int
		Error         string
	}

	cases := [6]Input{}

	//邮箱已经被注册
	cases[0] = Input{
		Email:         "root@root.com",
		Password:      "create",
		Name:          "create",
		RoleUUIDs:     []string{"286c730f-5c76-4280-91c5-2a172f782b84"},
		ClusterUUIDs:  []string{"335ddb23-1a43-4c73-9c14-6f2900ed627d"},
		ReviewerUUIDs: []string{"e70e78bb-9d08-405d-a0ed-266ec703de19"},
		AvatarUUID:    "8e79d8ef-1fbf-496b-95af-b1d790ee03d3",
		Status:        1,
		Error:         "graphql: 错误代码: 2009, 错误信息: 账号(email=root@root.com)已经注册。",
	}
	//roleid 不存在
	cases[1] = Input{
		Email:         "create3@dba.com",
		Password:      "create",
		Name:          "create",
		RoleUUIDs:     []string{"286c730f-5c76-4280-91c5-2a172f782b84zzz"},
		ClusterUUIDs:  []string{"335ddb23-1a43-4c73-9c14-6f2900ed627d"},
		ReviewerUUIDs: []string{"e70e78bb-9d08-405d-a0ed-266ec703de19"},
		AvatarUUID:    "8e79d8ef-1fbf-496b-95af-b1d790ee03d3",
		Status:        1,
		Error:         "graphql: 错误代码: 1404, 错误信息: 角色(uuid=286c730f-5c76-4280-91c5-2a172f782b84zzz)不存在。",
	}
	//clusterid 不存在
	cases[2] = Input{
		Email:         "create3@dba.com",
		Password:      "create",
		Name:          "create",
		RoleUUIDs:     []string{"286c730f-5c76-4280-91c5-2a172f782b84"},
		ClusterUUIDs:  []string{"335ddb23-1a43-4c73-9c14-6f2900ed627dzzzzzzzzzzzzzzzz"},
		ReviewerUUIDs: []string{"e70e78bb-9d08-405d-a0ed-266ec703de19"},
		AvatarUUID:    "8e79d8ef-1fbf-496b-95af-b1d790ee03d3",
		Status:        1,
		Error:         "graphql: 错误代码: 1404, 错误信息: 群集(uuid=335ddb23-1a43-4c73-9c14-6f2900ed627dzzzzzzzzzzzzzzzz)不存在。",
	}
	//reviewerid 不存在
	cases[3] = Input{
		Email:         "create3@dba.com",
		Password:      "create",
		Name:          "create",
		RoleUUIDs:     []string{"286c730f-5c76-4280-91c5-2a172f782b84"},
		ClusterUUIDs:  []string{"335ddb23-1a43-4c73-9c14-6f2900ed627d"},
		ReviewerUUIDs: []string{"e70e78bb-9d08-405d-a0ed-266ec703de19222"},
		AvatarUUID:    "8e79d8ef-1fbf-496b-95af-b1d790ee03d3",
		Status:        1,
		Error:         "graphql: 错误代码: 1404, 错误信息: 账号(uuid=e70e78bb-9d08-405d-a0ed-266ec703de19222)不存在。",
	}

	//reviewerid 状态异常
	cases[4] = Input{
		Email:         "create4@dba.com",
		Password:      "create",
		Name:          "create",
		RoleUUIDs:     []string{"286c730f-5c76-4280-91c5-2a172f782b84"},
		ClusterUUIDs:  []string{"335ddb23-1a43-4c73-9c14-6f2900ed627d"},
		ReviewerUUIDs: []string{"e70e78bb-9d08-405d-a0ed-266ec703de19"},
		AvatarUUID:    "8e79d8ef-1fbf-496b-95af-b1d790ee03d3",
		Status:        1,
		Error:         "graphql: 错误代码: 2026, 错误信息: 审核用户(uuid=e70e78bb-9d08-405d-a0ed-266ec703de19)的当前状态异常。",
	}

	//正常创建用户
	cases[5] = Input{
		Email:         "create3@dba.com",
		Password:      "create",
		Name:          "create",
		RoleUUIDs:     []string{"286c730f-5c76-4280-91c5-2a172f782b84"},
		ClusterUUIDs:  []string{"335ddb23-1a43-4c73-9c14-6f2900ed627d"},
		ReviewerUUIDs: []string{"e70e78bb-9d08-405d-a0ed-266ec703de19"},
		AvatarUUID:    "8e79d8ef-1fbf-496b-95af-b1d790ee03d3",
		Status:        1,
		Error:         "",
	}
	var respData struct {
		CreateUser struct {
			Email string
			Name  string
		}
	}
	for _, m := range cases {
		req.Var("email", m.Email)
		req.Var("password", m.Password)
		req.Var("name", m.Name)
		req.Var("roleid", m.RoleUUIDs)
		req.Var("clusterid", m.ClusterUUIDs)
		req.Var("reviewerid", m.ReviewerUUIDs)
		req.Var("avatarid", m.AvatarUUID)
		req.Var("status", m.Status)
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authentication", LoginData.Login.Token)
		err := Do(req, &respData)
		if err != nil {
			require.Equal(t, m.Error, err.Error())
		} else {
			require.Equal(t, m.Email, respData.CreateUser.Email)
		}
		// fmt.Println(m)
	}
}

func TestUpdateUser(t *testing.T) {
	req := graphql.NewRequest(`
	mutation updateUser($id:ID!,$email:String!,$password:String!,$name:String!,$status:Int!,$avatarid:String!){
		updateUser(
		input:{
		  UserUUID:$id,
		  Email:$email,
		  Password:$password,
		  Name:$name,
		  Status:$status,
		  AvatarUUID:$avatarid
		}
	  ){
		Email
		Name
		Status
	  }
	}`)
	type Input struct {
		UUID       string
		Email      string
		Password   string
		Name       string
		AvatarUUID string
		Status     int
		Error      string
	}

	cases := [3]Input{}
	//正常更新用户
	cases[0] = Input{
		Email:      "update@dba.com",
		Password:   "update",
		Name:       "update",
		UUID:       "09c5e0f9-d449-4e37-8008-8d97a8b8609a",
		AvatarUUID: "8e79d8ef-1fbf-496b-95af-b1d790ee03d3",
		Status:     1,
		Error:      "",
	}
	//邮箱已经被注册
	cases[1] = Input{
		Email:      "root@root.com",
		Password:   "udpate",
		Name:       "update",
		UUID:       "09c5e0f9-d449-4e37-8008-8d97a8b8609a",
		AvatarUUID: "8e79d8ef-1fbf-496b-95af-b1d790ee03d3",
		Status:     1,
		Error:      "graphql: 错误代码: 2009, 错误信息: 账号(email=root@root.com)已经注册。",
	}
	//用户不存在
	cases[2] = Input{
		Email:      "create3@dba.com",
		Password:   "create",
		Name:       "create",
		UUID:       "09c5e0f9-d449-4e37-8008-8d97a8b8609az",
		AvatarUUID: "8e79d8ef-1fbf-496b-95af-b1d790ee03d3",
		Status:     1,
		Error:      "graphql: 错误代码: 1404, 错误信息: 账号(uuid=09c5e0f9-d449-4e37-8008-8d97a8b8609az)不存在。",
	}
	var respData struct {
		UpdateUser struct {
			Email string
			Name  string
		}
	}
	for _, m := range cases {

		req.Var("email", m.Email)
		req.Var("password", m.Password)
		req.Var("name", m.Name)
		req.Var("id", m.UUID)
		req.Var("avatarid", m.AvatarUUID)
		req.Var("status", m.Status)
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authentication", LoginData.Login.Token)
		err := Do(req, &respData)
		if err != nil {
			require.Equal(t, m.Error, err.Error())
		} else {
			require.Equal(t, m.Email, respData.UpdateUser.Email)
		}
		// fmt.Println(m)
	}
}

func TestPatchUserStatus(t *testing.T) {
	req := graphql.NewRequest(`
	mutation patchUserStatus($id:ID!,$status:Int!){
		patchUserStatus(
			input:{
			UserUUID:$id,
			Status:$status
		  } 
	  )
  }
  `)
	type Input struct {
		UUID   string
		Status int
		Error  string
	}

	cases := [3]Input{}
	//参数无效
	cases[0] = Input{
		UUID:   "09c5e0f9-d449-4e37-8008-8d97a8b8609a",
		Status: -1,
		Error:  fmt.Sprintf("graphql: 错误代码: 1400, 错误信息: 参数(status=%d)无效。", -1),
	}
	//无效uuid
	cases[1] = Input{
		UUID:   "09c5e0f9-d449-4e37-8008-8d97a8b8609azzzzzz",
		Status: 1,
		Error:  fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 账号(uuid=%s)不存在。", "09c5e0f9-d449-4e37-8008-8d97a8b8609azzzzzz"),
	}
	//用户状态无需更新
	cases[2] = Input{
		UUID:   "09c5e0f9-d449-4e37-8008-8d97a8b8609a",
		Status: 1,
		Error:  "graphql: 错误代码: 1200, 错误信息: 账号(uuid=09c5e0f9-d449-4e37-8008-8d97a8b8609a)状态无需更新。",
	}
	var respData struct {
		PatchUserStatus bool
	}
	for _, m := range cases {
		req.Var("id", m.UUID)
		req.Var("status", m.Status)
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authentication", LoginData.Login.Token)
		err := Do(req, &respData)
		if err != nil {
			require.Equal(t, m.Error, err.Error())
		} else {
			require.Equal(t, true, respData.PatchUserStatus)
		}
		// fmt.Println(m)
	}
}

//UpdateProfile
func TestUpdateProfile(t *testing.T) {
	req := graphql.NewRequest(`
    mutation ($id: String!,$name: String!) {
		updateProfile(
			input:{
			  AvatarUUID: $id
			  Name: $name
			}){
				Email
				Name
		  }
	}`)

	var respData struct {
		UpdateProfile *models.User
	}

	/*
	   1、非法avatarUUID,也会更新
	   2、合法值
	   3、空name
	*/
	cases := make(map[string]string)
	cases["zzz"] = "avatar测试2"
	cases["57ac663b-d02e-4cf1-81a8-d65f1b4308f6"] = "测试用户avatar5"
	cases["ef5078a7-b664-4a4a-935d-83ecc763e611"] = "6" //找不到

	for k, v := range cases {
		req.Var("id", k)
		req.Var("name", v)
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authentication", LoginData.Login.Token)

		if err := Do(req, &respData); err != nil {
			log.Println(err)
		}

		require.Equal(t, v, respData.UpdateProfile.Name)
	}
}

func TestUpdateEmail(t *testing.T) {
	req := graphql.NewRequest(`
    mutation ($email:String!) {
		updateEmail(
			input:{
			  NewEmail:$email
			})
	}`)

	cases := []map[string]string{
		{"root@root.com": "graphql: 错误代码: 1200, 错误信息: 账号(email=root@root.com)无需修改。"},
		{"111tzzz@qq.com": "graphql: 错误代码: 2009, 错误信息: 账号(email=111tzzz@qq.com)已经存在。"}, //账号已经存在
		{"root@root2.com": ""},
	}
	var respData struct {
		UpdateEmail bool
	}
	for _, m := range cases {
		for k, v := range m {
			req.Var("email", k)
			req.Header.Set("Cache-Control", "no-cache")
			req.Header.Set("Authentication", LoginData.Login.Token)
			err := Do(req, &respData)
			if err != nil {
				require.Equal(t, v, err.Error())
			} else {
				require.Equal(t, true, respData.UpdateEmail)
			}
		}
	}
}

//GrantReviewers

func TestGrantReviewers(t *testing.T) {
	req := graphql.NewRequest(`
    mutation ($id:ID!,$reviewerid: [String!]!) {
		grantReviewers(
			input:{
			  UserUUID:$id,
			  ReviewerUUIDs: $reviewerid
			})
	}`)

	type Input struct {
		UserUUID      string
		ReviewerUUIDS []string
		Error         string
	}

	cases := [4]Input{}
	// 用户uuid非法
	// cases["ef5078a7-b664-4a4a-935d-83ecc763e61222"] = []string{}
	cases[0] = Input{
		UserUUID:      "ef5078a7-b664-4a4a-935d-83ecc763e61222",
		ReviewerUUIDS: []string{"e70e78bb-9d08-405d-a0ed-266ec703de19"},
	}
	cases[0].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 账号(uuid=%s)不存在。", cases[0].UserUUID)

	//审核人uuid 非法、
	cases[1] = Input{
		UserUUID:      "c655a587-33b2-4efc-b39d-a258eff6308c",
		ReviewerUUIDS: []string{"6bae321e-74dd-46df-a86c-3fe1867444bazzzz"},
	}
	cases[1].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 账号(uuid=%s)不存在。", cases[1].ReviewerUUIDS[0])

	//状态异常
	cases[2] = Input{
		UserUUID:      "f0ebc1b4-c776-4613-8786-a40c1b02d3cb",
		ReviewerUUIDS: []string{"5f53bbda-0084-494e-b6b0-dfe60f4615a0"},
	}
	cases[2].Error = fmt.Sprintf("graphql: 错误代码: 2026, 错误信息: 审核用户(uuid=%s)的当前状态异常。", cases[2].ReviewerUUIDS[0])

	//状态正常
	cases[3] = Input{
		UserUUID:      "f0ebc1b4-c776-4613-8786-a40c1b02d3cb",
		ReviewerUUIDS: []string{"672d8d25-9448-4388-9aa1-d32e674063de"},
	}
	cases[3].Error = ""

	var respData struct {
		GrantReviewers bool
	}
	for _, k := range cases {
		req.Var("id", k.UserUUID)
		req.Var("reviewerid", k.ReviewerUUIDS)

		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authentication", LoginData.Login.Token)

		err := Do(req, &respData)

		if err != nil {
			require.Equal(t, k.Error, err.Error())
		} else {
			require.Equal(t, true, respData.GrantReviewers)
		}

	}

}

func TestRevokeReviewers(t *testing.T) {
	req := graphql.NewRequest(`
    mutation ($id:ID!,$reviewerid: [String!]!) {
		revokeReviewers(
			input:{
			  UserUUID:$id,
			  ReviewerUUIDs: $reviewerid
			})
	}`)
	type Input struct {
		UserUUID      string
		ReviewerUUIDS []string
		Error         string
	}
	cases := [4]Input{}
	// 用户uuid非法
	// cases["ef5078a7-b664-4a4a-935d-83ecc763e61222"] = []string{}
	cases[0] = Input{
		UserUUID:      "ef5078a7-b664-4a4a-935d-83ecc763e61222",
		ReviewerUUIDS: []string{"e70e78bb-9d08-405d-a0ed-266ec703de19"},
	}
	cases[0].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 账号(uuid=%s)不存在。", cases[0].UserUUID)

	//审核人uuid 非法、
	cases[1] = Input{
		UserUUID:      "c655a587-33b2-4efc-b39d-a258eff6308c",
		ReviewerUUIDS: []string{"6bae321e-74dd-46df-a86c-3fe1867444bazzzz"},
	}
	cases[1].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 账号(uuid=%s)不存在。", cases[1].ReviewerUUIDS[0])

	//状态异常
	cases[2] = Input{
		UserUUID:      "f0ebc1b4-c776-4613-8786-a40c1b02d3cb",
		ReviewerUUIDS: []string{"5f53bbda-0084-494e-b6b0-dfe60f4615a0"},
	}
	cases[2].Error = fmt.Sprintf("graphql: 错误代码: 2026, 错误信息: 审核用户(uuid=%s)的当前状态异常。", cases[2].ReviewerUUIDS[0])

	//状态正常
	cases[3] = Input{
		UserUUID:      "f0ebc1b4-c776-4613-8786-a40c1b02d3cb",
		ReviewerUUIDS: []string{"672d8d25-9448-4388-9aa1-d32e674063de"},
	}
	cases[3].Error = ""

	var respData struct {
		RevokeReviewers bool
	}

	for _, k := range cases {
		req.Var("id", k.UserUUID)
		req.Var("reviewerid", k.ReviewerUUIDS)
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authentication", LoginData.Login.Token)

		err := Do(req, &respData)

		if err != nil {
			require.Equal(t, k.Error, err.Error())
		} else {
			require.Equal(t, true, respData.RevokeReviewers)
		}
	}

}

func TestGrantClusters(t *testing.T) {
	req := graphql.NewRequest(`
    mutation ($id:ID!,$clusterid: [String!]!) {
		grantClusters(
			input:{
			  UserUUID:$id,
			  ClusterUUIDs: $clusterid
			})
	}`)

	type Input struct {
		UserUUID     string
		ClusterUUIDs []string
		Error        string
	}
	cases := [4]Input{}
	// 用户uuid非法
	cases[0] = Input{
		UserUUID:     "ef5078a7-b664-4a4a-935d-83ecc763e61222",
		ClusterUUIDs: []string{"e70e78bb-9d08-405d-a0ed-266ec703de19"},
	}
	cases[0].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 账号(uuid=%s)不存在。", cases[0].UserUUID)

	//群集uuid
	cases[1] = Input{
		UserUUID:     "c655a587-33b2-4efc-b39d-a258eff6308c",
		ClusterUUIDs: []string{"6bae321e-74dd-46df-a86c-3fe1867444bazzzz"},
	}
	cases[1].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 群集(uuid=%s)不存在。", cases[1].ClusterUUIDs[0])

	//正常状态的群集
	cases[2] = Input{
		UserUUID:     "f0ebc1b4-c776-4613-8786-a40c1b02d3cb",
		ClusterUUIDs: []string{"f4dc547b-28e6-499a-b375-9dd97cdb35b2"},
	}
	cases[2].Error = ""
	//错误代码: %s, 错误信息: 群集(uuid=%s)不可用。

	cases[3] = Input{
		UserUUID:     "f0ebc1b4-c776-4613-8786-a40c1b02d3cb",
		ClusterUUIDs: []string{"cc3d1347-aa8a-4a41-bbdb-8aa5ba0c78a1"},
	}
	cases[3].Error = fmt.Sprintf("graphql: 错误代码: 2001, 错误信息: 群集(uuid=%s)不可用。", cases[3].ClusterUUIDs[0])

	var respData struct {
		GrantClusters bool
	}
	for _, k := range cases {
		req.Var("id", k.UserUUID)
		req.Var("clusterid", k.ClusterUUIDs)

		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authentication", LoginData.Login.Token)

		err := Do(req, &respData)

		if err != nil {
			require.Equal(t, k.Error, err.Error())
		} else {
			require.Equal(t, true, respData.GrantClusters)
		}
	}

}

func TestRevokeClusters(t *testing.T) {
	req := graphql.NewRequest(`
    mutation ($id:ID!,$clusterid: [String!]!) {
		revokeClusters(
			input:{
			  UserUUID:$id,
			  ClusterUUIDs: $clusterid
			})
	}`)

	type Input struct {
		UserUUID     string
		ClusterUUIDs []string
		Error        string
	}
	cases := [4]Input{}
	// 用户uuid非法
	cases[0] = Input{
		UserUUID:     "ef5078a7-b664-4a4a-935d-83ecc763e61222",
		ClusterUUIDs: []string{"e70e78bb-9d08-405d-a0ed-266ec703de19"},
	}
	cases[0].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 账号(uuid=%s)不存在。", cases[0].UserUUID)

	//群集uuid
	cases[1] = Input{
		UserUUID:     "c655a587-33b2-4efc-b39d-a258eff6308c",
		ClusterUUIDs: []string{"6bae321e-74dd-46df-a86c-3fe1867444bazzzz"},
	}
	cases[1].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 群集(uuid=%s)不存在。", cases[1].ClusterUUIDs[0])

	//错误代码: %s, 错误信息: 群集(uuid=%s)不可用。
	cases[2] = Input{
		UserUUID:     "f0ebc1b4-c776-4613-8786-a40c1b02d3cb",
		ClusterUUIDs: []string{"cc3d1347-aa8a-4a41-bbdb-8aa5ba0c78a1"},
	}
	cases[2].Error = fmt.Sprintf("graphql: 错误代码: 2001, 错误信息: 群集(uuid=%s)不可用。", cases[2].ClusterUUIDs[0])

	//正常状态的群集
	cases[3] = Input{
		UserUUID:     "f0ebc1b4-c776-4613-8786-a40c1b02d3cb",
		ClusterUUIDs: []string{"f4dc547b-28e6-499a-b375-9dd97cdb35b2"},
	}
	cases[3].Error = ""

	var respData struct {
		RevokeClusters bool
	}
	for _, k := range cases {
		req.Var("id", k.UserUUID)
		req.Var("clusterid", k.ClusterUUIDs)

		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authentication", LoginData.Login.Token)

		err := Do(req, &respData)

		if err != nil {
			require.Equal(t, k.Error, err.Error())
		} else {
			require.Equal(t, true, respData.RevokeClusters)
		}
	}

}

func TestGrantRoles(t *testing.T) {
	req := graphql.NewRequest(`
    mutation ($id:ID!,$roleids: [String!]!) {
		grantRoles(
			input:{
			  UserUUID:$id,
			  RoleUUIDs: $roleids
			})
	}`)
	type Input struct {
		UserUUID  string
		RoleUUIDs []string
		Error     string
	}
	cases := [3]Input{}
	// 用户uuid非法
	cases[0] = Input{
		UserUUID:  "ef5078a7-b664-4a4a-935d-83ecc763e61222",
		RoleUUIDs: []string{"e70e78bb-9d08-405d-a0ed-266ec703de19"},
	}
	cases[0].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 账号(uuid=%s)不存在。", cases[0].UserUUID)

	//角色uuid非法
	cases[1] = Input{
		UserUUID:  "c655a587-33b2-4efc-b39d-a258eff6308c",
		RoleUUIDs: []string{"6bae321e-74dd-46df-a86c-3fe1867444bazzzz"},
	}
	cases[1].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 角色(uuid=%s)不存在。", cases[1].RoleUUIDs[0])

	//正常授权
	cases[2] = Input{
		UserUUID:  "f0ebc1b4-c776-4613-8786-a40c1b02d3cb",
		RoleUUIDs: []string{"dba2929d-19d5-44ad-9d6a-71c32139f36c"},
	}
	cases[2].Error = ""

	var respData struct {
		GrantRoles bool
	}
	for _, k := range cases {
		req.Var("id", k.UserUUID)
		req.Var("roleids", k.RoleUUIDs)

		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authentication", LoginData.Login.Token)

		err := Do(req, &respData)

		if err != nil {
			require.Equal(t, k.Error, err.Error())
		} else {
			require.Equal(t, true, respData.GrantRoles)
		}
	}

}

func TestRevokeRoles(t *testing.T) {
	req := graphql.NewRequest(`
    mutation ($id:ID!,$roleid:[String!]!) {
		revokeRoles(
			input:{
			  UserUUID:$id,
			  RoleUUIDs:  $roleid
			}) 
	}`)

	var respData struct {
		RevokeRoles bool
	}

	type Input struct {
		UserUUID  string
		RoleUUIDs []string
		Error     string
	}
	cases := [3]Input{}
	// 用户uuid非法
	cases[0] = Input{
		UserUUID:  "ef5078a7-b664-4a4a-935d-83ecc763e61222",
		RoleUUIDs: []string{"e70e78bb-9d08-405d-a0ed-266ec703de19"},
	}
	cases[0].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 账号(uuid=%s)不存在。", cases[0].UserUUID)
	//角色uuid非法
	cases[1] = Input{
		UserUUID:  "f0ebc1b4-c776-4613-8786-a40c1b02d3cb",
		RoleUUIDs: []string{"6bae321e-74dd-46df-a86c-3fe1867444bazzzz"},
	}
	cases[1].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 角色(uuid=%s)不存在。", cases[1].RoleUUIDs[0])

	//正常授权
	cases[2] = Input{
		UserUUID:  "f0ebc1b4-c776-4613-8786-a40c1b02d3cb",
		RoleUUIDs: []string{"dba2929d-19d5-44ad-9d6a-71c32139f36c"},
	}
	cases[2].Error = ""

	for _, k := range cases {
		req.Var("id", k.UserUUID)
		req.Var("roleid", k.RoleUUIDs)
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authentication", LoginData.Login.Token)

		err := Do(req, &respData)

		if err != nil {
			require.Equal(t, k.Error, err.Error())
		} else {
			require.Equal(t, true, respData.RevokeRoles)
		}
	}

}

func Login() {
	// make a request
	req := graphql.NewRequest(`
    mutation ($email: String! $password: String!) {
		login (
			input: {
				Email: $email
				Password: $password
			}
		) {
			Me {
				UserUUID
			}
			Token
		}
	}`)

	// set any variables
	req.Var("email", "root@root.com")
	req.Var("password", "create")

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")

	if err := Do(req, &LoginData); err != nil {
		log.Println(err)
	}
}

//clusters

func TestCreateCluster(t *testing.T) {
	req := graphql.NewRequest(`
    mutation createCluster($host:String!,$alias:String!,$ip:String!,$port:Int!,$user:String!,$password:String!,$status:Int!){
		createCluster(
		  input:{
			Host:$host,
			Alias:$alias,
			IP:$ip,
			Port:$port,
			User:$user,
			Password:$password,
			Status:$status
		  }
		){
		  Host
		  IP
		  Port
		  User
		  Alias
		}
	  }
	  `)

	type Input struct {
		Host     string
		Password string
		User     string
		Alias    string
		IP       string
		Port     int
		Status   int
		Error    string
	}

	cases := [2]Input{}
	//群集别名已经存在
	cases[0] = Input{
		Host:     "10.100.11.8734",
		Password: "create",
		Status:   1,
		Alias:    "11873453",
		IP:       "10.100.11.27",
		User:     "create",
		Port:     3309,
	}
	cases[0].Error = fmt.Sprintf("graphql: 错误代码: 2003, 错误信息: 群集(alias=%s)已经存在。", cases[0].Alias)
	//群集ip、port 重复
	cases[1] = Input{
		Host:     "10.100.11.28",
		Password: "create",
		Alias:    "11828",
		Status:   1,
		IP:       "10.100.11.28",
		User:     "create",
		Port:     3309,
	}
	cases[1].Error = fmt.Sprintf("graphql: 错误代码: 2006, 错误信息: 群集(cluster=%s:%d 或 cluster=%s:%d)已经存在", cases[1].Host, cases[1].Port, cases[1].IP, cases[1].Port)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authentication", LoginData.Login.Token)

	var respData struct {
		CreateCluster struct {
			Host  string
			IP    string
			Port  int
			User  string
			Alias string
		}
	}

	for _, k := range cases {
		req.Var("host", k.Host)
		req.Var("alias", k.Alias)
		req.Var("status", k.Status)
		req.Var("user", k.User)
		req.Var("port", k.Port)
		req.Var("password", k.Password)
		req.Var("ip", k.IP)
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authentication", LoginData.Login.Token)

		err := Do(req, &respData)

		if err != nil {
			require.Equal(t, k.Error, err.Error())
		} else {
			require.Equal(t, k.Alias, respData.CreateCluster.Alias)
		}
	}
}

func TestUpdateCluster(t *testing.T) {
	req := graphql.NewRequest(`
    mutation updateCluster($id:ID!,$host:String!,$alias:String!,$ip:String!,$port:Int!,$user:String!,$password:String!,$status:Int!){
		updateCluster(
		  input:{
			ClusterUUID:$id,
			Host:$host,
			Alias:$alias,
			IP:$ip,
			Port:$port,
			User:$user,
			Password:$password,
			Status:$status
		  }
		){
		  Host
		  IP
		  Port
		  User
		  Alias
		}
	  }
	  `)

	type Input struct {
		ID       string
		Host     string
		Password string
		User     string
		Alias    string
		IP       string
		Port     int
		Status   int
		Error    string
	}

	cases := [2]Input{}
	//群集uuid 不存在
	cases[0] = Input{
		ID:       "zzzzzzzzzzzzzzzzzzzzzzzzz",
		Host:     "10.100.11.87",
		Password: "create",
		Status:   1,
		Alias:    "1187",
		IP:       "10.100.11.87",
		User:     "create",
		Port:     3309,
	}
	cases[0].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 群集(uuid=%s)不存在。", cases[0].ID)

	//群集ip、port 重复
	cases[1] = Input{
		ID:       "77874b86-3f15-4541-86f6-d84855d3f59a",
		Host:     "10.100.11.79",
		Password: "create",
		Alias:    "118823242",
		Status:   1,
		IP:       "10.100.11.79",
		User:     "create",
		Port:     3309,
	}
	cases[1].Error = fmt.Sprintf("graphql: 错误代码: 2006, 错误信息: 群集(cluster=%s:%d 或 cluster=%s:%d)已经存在。", cases[1].Host, cases[1].Port, cases[1].IP, cases[1].Port)

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authentication", LoginData.Login.Token)

	var respData struct {
		UpdateCluster struct {
			Host  string
			IP    string
			Port  int
			User  string
			Alias string
		}
	}

	for _, k := range cases {
		req.Var("host", k.Host)
		req.Var("alias", k.Alias)
		req.Var("status", k.Status)
		req.Var("user", k.User)
		req.Var("port", k.Port)
		req.Var("password", k.Password)
		req.Var("ip", k.IP)
		req.Var("id", k.ID)
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authentication", LoginData.Login.Token)

		err := Do(req, &respData)

		if err != nil {
			require.Equal(t, k.Error, err.Error())
		} else {
			require.Equal(t, k.Alias, respData.UpdateCluster.Alias)
		}
	}
}

func TestPatchClusterStatus(t *testing.T) {
	req := graphql.NewRequest(`
    mutation patchClusterStatus($id:ID!,$status:Int!){
		patchClusterStatus(
		  input:{
			ClusterUUID:$id,
			Status:$status
		  }
		)}`)

	type Input struct {
		ID     string
		Status int
		Error  string
	}

	cases := [2]Input{}
	//无效的status
	cases[0] = Input{
		ID:     "686e53f3-47f0-4b2b-a7af-82386af1dfe1",
		Status: -1,
	}
	cases[0].Error = fmt.Sprintf("graphql: 错误代码: 1400, 错误信息: 参数(status=%d)无效。", cases[0].Status)

	//群集不存在
	cases[1] = Input{
		ID:     "77874b86-3f15-4541-86f6-d84855d3f59a",
		Status: 1,
	}
	cases[1].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 群集(uuid=%s)不存在。", cases[1].ID)

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authentication", LoginData.Login.Token)

	var respData struct {
		PatchClusterStatus bool
	}

	for _, k := range cases {

		req.Var("status", k.Status)
		req.Var("id", k.ID)
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authentication", LoginData.Login.Token)

		err := Do(req, &respData)

		if err != nil {
			require.Equal(t, k.Error, err.Error())
		} else {
			require.Equal(t, true, respData.PatchClusterStatus)
		}
	}
}

//rule

func TestPatchRuleValues(t *testing.T) {
	req := graphql.NewRequest(`
    mutation patchRuleValues($id:ID!,$values:String!){
		patchRuleValues(
		  input:{
			RuleUUID:$id,
			Values:$values
		  }
		)}`)

	type Input struct {
		ID     string
		Values string
		Error  string
	}

	cases := [2]Input{}
	//无效的status
	cases[0] = Input{
		ID:     "686e53f3-47f0-4b2b-a7af-82386af1dfe1",
		Values: "1",
	}
	cases[0].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 规则(uuid=%s)不存在。", cases[0].ID)
	//规则不允许更新
	//ceb63263-761b-4904-ba6c-993281924daf
	// 禁止创建数据库规则
	cases[1] = Input{
		ID:     "ceb63263-761b-4904-ba6c-993281924daf",
		Values: "2",
	}
	cases[1].Error = fmt.Sprintf("graphql: 错误代码: 2001, 错误信息: 规则(uuid=%s)不允许更新。", cases[1].ID)

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authentication", LoginData.Login.Token)

	var respData struct {
		PatchClusterStatus bool
	}
	for _, k := range cases {
		req.Var("values", k.Values)
		req.Var("id", k.ID)
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authentication", LoginData.Login.Token)
		err := Do(req, &respData)

		if err != nil {
			require.Equal(t, k.Error, err.Error())
		} else {
			require.Equal(t, true, respData.PatchClusterStatus)
		}
	}
}

func TestPatchRuleBitwise(t *testing.T) {
	req := graphql.NewRequest(`
    mutation patchRuleBitwise($id:ID!,$enabled:String!){
		patchRuleBitwise(
		  input:{
			RuleUUID:$id,
			Enabled:$enabled
		  }
		)}`)

	type Input struct {
		ID      string
		Enabled string
		Error   string
	}

	cases := [2]Input{}
	//无效的status
	cases[0] = Input{
		ID:      "686e53f3-47f0-4b2b-a7af-82386af1dfe1",
		Enabled: "true",
	}
	cases[0].Error = fmt.Sprintf("graphql: 错误代码: 1404, 错误信息: 规则(uuid=%s)不存在。", cases[0].ID)
	//规则不允许更新
	//ceb63263-761b-4904-ba6c-993281924daf
	// 禁止创建数据库规则
	cases[1] = Input{
		ID:      "ceb63263-761b-4904-ba6c-993281924daf",
		Enabled: "false",
	}
	cases[1].Error = fmt.Sprintf("graphql: 错误代码: 2001, 错误信息: 规则(uuid=%s)不允许更新。", cases[1].ID)

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authentication", LoginData.Login.Token)

	var respData struct {
		PatchClusterStatus bool
	}
	for _, k := range cases {
		req.Var("enabled", k.Enabled)
		req.Var("id", k.ID)
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authentication", LoginData.Login.Token)
		err := Do(req, &respData)

		if err != nil {
			require.Equal(t, k.Error, err.Error())
		} else {
			require.Equal(t, true, respData.PatchClusterStatus)
		}
	}
}

// func TestLogout(t *testing.T) {
// 	req := graphql.NewRequest(`
//     mutation  {
// 		logout
// 	}`)

// 	req.Header.Set("Cache-Control", "no-cache")
// 	req.Header.Set("Authentication", LoginData.Login.Token)

// 	var respData struct {
// 		Logout bool
// 	}
// 	if err := Do(req, &respData); err != nil {

// 	}
// 	require.Equal(t, true, respData.Logout)
// }

func TestUpdatePassword(t *testing.T) {
	req := graphql.NewRequest(`
    mutation ($new:String!,$old:String!) {
		updatePassword(
			input:{
			  NewPassword: $new
			  OldPassword: $old
			})
	}`)

	cases := []map[string]string{
		{"root1": "root1"},  //新旧密码相同无法更新
		{"root": "dba001"},  //正常更新
		{"wrong": "wrong2"}, //错误的就密码
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authentication", LoginData.Login.Token)

	var respData struct {
		UpdatePassword *bool
	}
	for _, m := range cases {
		for k, v := range m {
			req.Var("new", k)
			req.Var("old", v)
			req.Header.Set("Cache-Control", "no-cache")
			req.Header.Set("Authentication", LoginData.Login.Token)

			err := Do(req, &respData)

			if k == "root1" {
				require.Equal(t, fmt.Sprintf("graphql: 错误代码: 1200, 错误信息: 新密码(password=%s)和旧密码(password=%s)相同。", k, v), err.Error())
			} else if k == "wrong" {
				require.Equal(t, fmt.Sprintf("graphql: 错误代码: 2010, 错误信息: 旧密码(password=%s)不正确。", v), err.Error())
			} else {
				require.Equal(t, true, *respData.UpdatePassword, k)
			}
		}
	}
}

func TestUnautherized(t *testing.T) {

}
