package models

// UpdateTemplateInput GraphQL API交互所需要的结构体
type UpdateTemplateInput struct {
	TemplateUUID string `valid:"required,length(36|36)"   gqlgen:"TemplateUUID"` //
	Subject      string `valid:"required,length(1|100)"   gqlgen:"Subject"`      //
	Body         string `valid:"required,length(1|65535)" gqlgen:"Body"`         //
	Description  string `valid:"required,length(1|50)" gqlgen:"Description"`     //
}

// CreateTicketInput GraphQL API交互所需要的结构体
type CreateTicketInput struct {
	ClusterUUID  string `valid:"required,length(36|36)"   gqlgen:"ClusterUUID"`  //
	Database     string `valid:"required,length(1|50)"    gqlgen:"Database"`     //
	Subject      string `valid:"required,length(1|75)"    gqlgen:"Subject"`      //
	Content      string `valid:"required,length(1|65535)" gqlgen:"Content"`      //
	ReviewerUUID string `valid:"required,length(36|36)"   gqlgen:"ReviewerUUID"` //
}

// UpdateTicketInput GraphQL API交互所需要的结构体
type UpdateTicketInput struct {
	TicketUUID   string `valid:"required,length(36|36)"   gqlgen:"TicketUUID"`   //
	ClusterUUID  string `valid:"required,length(36|36)"   gqlgen:"ClusterUUID"`  //
	Database     string `valid:"required,length(1|50)"    gqlgen:"Database"`     //
	Subject      string `valid:"required,length(1|75)"    gqlgen:"Subject"`      //
	Content      string `valid:"required,length(1|65535)" gqlgen:"Content"`      //
	ReviewerUUID string `valid:"required,length(36|36)"   gqlgen:"ReviewerUUID"` //
}

// PatchTicketStatusInput GraphQL API交互所需要的结构体
type PatchTicketStatusInput struct {
	TicketUUID string `valid:"required,length(36|36)" gqlgen:"TicketUUID"` //
	Status     string `valid:"required"               gqlgen:"Status"`     //
}

// ScheduleTicketInput GraphQL API交互所需要的结构体
type ScheduleTicketInput struct {
	TicketUUID string `valid:"required,length(36|36)" gqlgen:"TicketUUID"` //
	Schedule   string `valid:"required,int"           gqlgen:"Schedule"`   //
}

// CreateCommentInput GraphQL API交互所需要的结构体
type CreateCommentInput struct {
	TicketUUID string `valid:"required,length(36|36)"   gqlgen:"TicketUUID"` //
	Content    string `valid:"required,length(1|65535)" gqlgen:"Content"`    //
}

// CreateUserInput GraphQL API交互所需要的结构体
type CreateUserInput struct {
	Email         string   `valid:"required,length(3|75)"           gqlgen:"Email"`         //
	Password      string   `valid:"required,length(6|25)"           gqlgen:"Password"`      //
	Name          string   `valid:"required,length(1|15)"           gqlgen:"Name"`          //
	Phone         uint64   `valid:"-"                               gqlgen:"Phone"`         //
	RoleUUIDs     []string `valid:"required"                        gqlgen:"RoleUUIDs"`     //
	ClusterUUIDs  []string `valid:"optional"                        gqlgen:"ClusterUUIDs"`  //
	ReviewerUUIDs []string `valid:"optional"                        gqlgen:"ReviewerUUIDs"` //
	AvatarUUID    string   `valid:"required,length(36|36)"          gqlgen:"AvatarUUID"`    //
	Status        uint8    `valid:"required,int,matches(^(1|2|3)$)" gqlgen:"Status"`        //
}

// UpdateUserInput GraphQL API交互所需要的结构体
type UpdateUserInput struct {
	UserUUID   string `valid:"required,length(36|36)"          gqlgen:"UserUUID"`   //
	Email      string `valid:"required,length(3|75)"           gqlgen:"Email"`      //
	Password   string `valid:"required,length(3|75)"           gqlgen:"Password"`   //
	Status     uint8  `valid:"required,int,matches(^(1|2|3)$)" gqlgen:"Status"`     //
	Name       string `valid:"required,length(1|25)"           gqlgen:"Name"`       //
	Phone      uint64 `valid:"-"                               gqlgen:"Phone"`      //
	AvatarUUID string `valid:"required,length(36|36)"          gqlgen:"AvatarUUID"` //
}

// UpdateProfileInput GraphQL API交互所需要的结构体
type UpdateProfileInput struct {
	AvatarUUID string `valid:"required,length(36|36)" gqlgen:"AvatarUUID"` //
	Name       string `valid:"required,length(1|25)"  gqlgen:"Name"`       //
	Phone      uint64 `valid:"-"                      gqlgen:"Phone"`      //
}

// PatchPasswordInput GraphQL API交互所需要的结构体
type PatchPasswordInput struct {
	OldPassword string `valid:"required，length(6|25)" gqlgen:"OldPassword"` //
	NewPassword string `valid:"required，length(6|25)" gqlgen:"NewPassword"` //
}

// PatchEmailInput GraphQL API交互所需要的结构体
type PatchEmailInput struct {
	NewEmail string `valid:"required,email" gqlgen:"NewEmail"` //
}

// GrantRolesInput GraphQL API交互所需要的结构体
type GrantRolesInput struct {
	UserUUID  string   `valid:"required,length(36|36)" gqlgen:"UserUUID"`  //
	RoleUUIDs []string `valid:"required"               gqlgen:"RoleUUIDs"` //
}

// GrantReviewersInput GraphQL API交互所需要的结构体
type GrantReviewersInput struct {
	UserUUID      string   `valid:"required,length(36|36)" gqlgen:"UserUUID"`      //
	ReviewerUUIDs []string `valid:"required"               gqlgen:"ReviewerUUIDs"` //
}

// RevokeReviewersInput GraphQL API交互所需要的结构体
type RevokeReviewersInput struct {
	UserUUID      string   `valid:"required,length(36|36)" gqlgen:"UserUUID"`      //
	ReviewerUUIDs []string `valid:"required"               gqlgen:"ReviewerUUIDs"` //
}

// RevokeRolesInput GraphQL API交互所需要的结构体
type RevokeRolesInput struct {
	UserUUID  string   `valid:"required,length(36|36)" gqlgen:"UserUUID"`  //
	RoleUUIDs []string `valid:"required"               gqlgen:"RoleUUIDs"` //
}

// UserRegisterInput GraphQL API交互所需要的结构体
type UserRegisterInput struct {
	Email    string `valid:"required,email"        gqlgen:"Email"`    //
	Password string `valid:"required,length(6|25)" gqlgen:"Password"` //
}

// UserLoginInput GraphQL API交互所需要的结构体
type UserLoginInput struct {
	Email    string `valid:"required,email"        gqlgen:"Email"`    //
	Password string `valid:"required,length(6|25)" gqlgen:"Password"` //
}

// GrantClustersInput GraphQL API交互所需要的结构体
type GrantClustersInput struct {
	UserUUID     string   `valid:"required,length(36|36)" gqlgen:"UserUUID"`     //
	ClusterUUIDs []string `valid:"required"               gqlgen:"ClusterUUIDs"` //
}

// RevokeClustersInput GraphQL API交互所需要的结构体
type RevokeClustersInput struct {
	UserUUID     string   `valid:"required,length(36|36)" gqlgen:"UserUUID"`     //
	ClusterUUIDs []string `valid:"required"               gqlgen:"ClusterUUIDs"` //
}

// PatchUserStatusInput GraphQL API交互所需要的结构体
type PatchUserStatusInput struct {
	UserUUID string `valid:"required,length(36|36)"      gqlgen:"UserUUID"` //
	Status   uint8  `valid:"required,matches(^(0|1|2)$)" gqlgen:"Status"`   //
}

// PatchRuleValuesInput GraphQL API交互所需要的结构体
type PatchRuleValuesInput struct {
	RuleUUID string `valid:"required,length(36|36)" gqlgen:"RuleUUID"` //
	Values   string `valid:"required,length(1|150)" gqlgen:"Values"`   //
}

// PatchRuleBitwiseInput GraphQL API交互所需要的结构体
type PatchRuleBitwiseInput struct {
	RuleUUID string `valid:"required,length(36|36)"           gqlgen:"RuleUUID"` //
	Enabled  string `valid:"required,matches(^(true|false)$)" gqlgen:"Enabled"`  //
}

// CreateClusterInput GraphQL API交互所需要的结构体
type CreateClusterInput struct {
	Host     string `valid:"required,length(1|75)"         gqlgen:"Host"`     //
	IP       string `valid:"required,ipv4"                 gqlgen:"IP"`       //
	Port     uint16 `valid:"required,port"                 gqlgen:"Port"`     //
	Alias    string `valid:"required,length(4|75)"         gqlgen:"Alias"`    //
	User     string `valid:"required,length(4|40)"         gqlgen:"User"`     //
	Password string `valid:"required,length(4|40)"         gqlgen:"Password"` // 注意：密码最大长度40位，超过40位会导致数据库截断
	Status   uint8  `valid:"required,int,matches(^(1|2)$)" gqlgen:"Status"`   //
}

// UpdateClusterInput GraphQL API交互所需要的结构体
type UpdateClusterInput struct {
	ClusterUUID string `valid:"required,length(36|36)"        gqlgen:"ClusterUUID"` //
	Host        string `valid:"required,length(1|75)"         gqlgen:"Host"`        //
	IP          string `valid:"required,ipv4"                 gqlgen:"IP"`          //
	Port        uint16 `valid:"required,port"                 gqlgen:"Port"`        //
	Alias       string `valid:"required,length(4|75)"         gqlgen:"Alias"`       //
	User        string `valid:"required,length(4|40)"         gqlgen:"User"`        //
	Status      uint8  `valid:"required,int,matches(^(1|2)$)" gqlgen:"Status"`      //
	Password    string `valid:"required,length(4|40)"         gqlgen:"Password"`    //
}

// PatchClusterStatusInput GraphQL API交互所需要的结构体
type PatchClusterStatusInput struct {
	ClusterUUID string `valid:"required,length(36|36)"        gqlgen:"ClusterUUID"` //
	Status      uint8  `valid:"required,int,matches(^(1|2)$)" gqlgen:"Status"`      //
}

// ValidateConnectionInput GraphQL API交互所需要的结构体
type ValidateConnectionInput struct {
	IP       string `valid:"required,ipv4"         gqlgen:"IP"`       //
	Port     uint16 `valid:"required,port"         gqlgen:"Port"`     //
	User     string `valid:"required,length(4|40)" gqlgen:"User"`     //
	Password string `valid:"required,length(4|40)" gqlgen:"Password"` // 注意：密码最大长度40位，超过40位会导致数据库截断
}

// ValidatePatternInput GraphQL API交互所需要的结构体
type ValidatePatternInput struct {
	Pattern string `valid:"required,length(1|255)" gqlgen:"Pattern"` //
}

// PatchOptionValueInput GraphQL API交互所需要的结构体
type PatchOptionValueInput struct {
	OptionUUID string `valid:"required,length(36|36)" gqlgen:"OptionUUID"` //
	Value      string `valid:"required,length(1|40)"  gqlgen:"Value"`      //
}

// CreateQueryInput GraphQL API交互所需要的结构体
type CreateQueryInput struct {
	ClusterUUID string `valid:"required,length(36|36)"   gqlgen:"ClusterUUID"` //
	Database    string `valid:"required,length(1|50)"    gqlgen:"Database"`    //
	Content     string `valid:"required,length(1|65535)" gqlgen:"Content"`     //
}

// SoarQueryInput GraphQL API交互所需要的结构体
type SoarQueryInput struct {
	ClusterUUID string `valid:"required,length(36|36)"   gqlgen:"ClusterUUID"` //
	Database    string `valid:"required,length(1|50)"    gqlgen:"Database"`    //
	Content     string `valid:"required,length(1|65535)" gqlgen:"Content"`     //
}

// ActivateInput GraphQL API交互所需要的结构体
type ActivateInput struct {
	Code string `valid:"required" gqlgen:"Code" json:"code"`
}

// LostPasswdInput GraphQL API交互所需要的结构体
type LostPasswdInput struct {
	Email string `valid:"required,length(1|75)" gqlgen:"Email" json:"email"`
}

// ResetPasswdInput GraphQL API交互所需要的结构体
type ResetPasswdInput struct {
	Code     string `valid:"required" gqlgen:"Code" json:"code"`
	Password string `valid:"required,length(6|20)" gqlgen:"Password" json:"password"`
}

// ResendActivationMailInput GraphQL API交互所需要的结构体
type ResendActivationMailInput struct {
	Email string `valid:"required,length(1|75)" gqlgen:"Email" json:"email"`
}
