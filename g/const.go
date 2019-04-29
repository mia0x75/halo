package g

import (
	"github.com/akhenakh/statgo"
)

const CREDENTIAL_KEY = "Credential"

var (
	// GlobalStat 获取服务器当前系统环境的全局对象
	GlobalStat = statgo.NewStat()
)

// Template 邮件模板
type Template string

// 邮件模板
const (
	TplTicketCreated    Template = "b5c5ac9c-2071-4dd3-af5c-69ec149ee682"
	TplTicketUpdated    Template = "5043d567-02ef-4f5e-be9a-13df9f5fde11"
	TplTicketRemoved    Template = "6e861f70-5d8c-4042-879c-9ca932fb792b"
	TplTicketExecuted   Template = "03714a3f-eafe-4836-8e85-d360ee29a70f"
	TplTicketFailed     Template = "0a55142a-e336-4a97-b655-94ecac454da2"
	TplTicketScheduled  Template = "0c7bf7ab-8e39-464e-b0b0-6a209842058a"
	TplTicketClosed     Template = "5a36648d-0c97-4aa5-b753-2872ea2e0ac6"
	TplTicketMrvFailure Template = "9676f8e5-988c-4d5f-802b-a92f619a7ef0"
	TplTicketLgtm       Template = "33a3d82e-bb2b-4428-8c45-6a8e50c0ed0c"
	TplUserRegistered   Template = "30f37d4f-2cfa-40f4-8b44-4b660f9c613d"
	TplPasswordUpdated  Template = "69d05ebf-7626-433f-b906-cb69a596f78e"
	TplEmailUpdated     Template = "aa5404c5-ce37-4e01-a41c-75833028e122"
	TplProfileUpdated   Template = "64c110cf-18d5-494c-917e-fc61322c98e0"
	TplUserCreated      Template = "7de2bf1a-c03a-49d0-822a-a1dd1c98bdc1"
	TplCommentCreated   Template = "ac156eb3-9948-4e2f-997f-77fdeceb12ca"
	TplCronCancelled    Template = "ff0a4c66-9356-498a-afff-40a4407d9d8a"
)
