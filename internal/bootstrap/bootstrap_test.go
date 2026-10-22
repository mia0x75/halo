package bootstrap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func validRequest() Request {
	return Request{
		Database: DatabaseRequest{Driver: "mysql", Addr: "root@tcp(127.0.0.1:3306)/halo"},
		Secret:   SecretRequest{SecretKey: "0123456789abcdef0123456789abcdef", Crypto: "0123456789abcdef"},
		Admin:    AdminRequest{Email: "admin@example.com", Password: "halo-admin"},
	}
}

func TestValidateAcceptsCompleteRequest(t *testing.T) {
	require.NoError(t, validRequest().Validate())
}

func TestValidateRejectsMissingAddress(t *testing.T) {
	request := validRequest()
	request.Database.Addr = "  "
	require.ErrorContains(t, request.Validate(), "元数据库连接串")
}

func TestValidateRejectsUnknownDriver(t *testing.T) {
	request := validRequest()
	request.Database.Driver = "sqlite"
	require.ErrorContains(t, request.Validate(), "元数据库类型")
}

func TestValidateRejectsShortSecrets(t *testing.T) {
	request := validRequest()
	request.Secret.SecretKey = "short"
	require.ErrorContains(t, request.Validate(), "JWT 密钥")

	request = validRequest()
	request.Secret.Crypto = "short"
	require.ErrorContains(t, request.Validate(), "加密密钥")
}

func TestValidateRejectsBadAdmin(t *testing.T) {
	request := validRequest()
	request.Admin.Email = "not-an-email"
	require.ErrorContains(t, request.Validate(), "管理员邮箱")

	request = validRequest()
	request.Admin.Password = "short"
	require.ErrorContains(t, request.Validate(), "管理员密码")
}

func TestNormalizeFillsDefaults(t *testing.T) {
	request := validRequest()
	request.Database.MaxIdle = 0
	request.Database.MaxConnections = 0
	request.Database.WaitTimeout = 0
	request.Secret.TokenName = ""
	request.Secret.Issuer = ""
	request.Admin.Name = ""

	normalized := request.normalize()
	require.Equal(t, 10, normalized.Database.MaxIdle)
	require.Equal(t, 100, normalized.Database.MaxConnections)
	require.Equal(t, 300, normalized.Database.WaitTimeout)
	require.Equal(t, "Authentication", normalized.Secret.TokenName)
	require.Equal(t, "halo", normalized.Secret.Issuer)
	require.Equal(t, "管理员", normalized.Admin.Name)
}
