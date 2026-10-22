// Package bootstrap performs the first-run setup behind the /welcome wizard:
// it connects to the metadata database, migrates it, loads the base data,
// writes the config file and creates the initial administrator.
package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"math"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/internal/app"
	"github.com/mia0x75/halo/internal/cache"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/platform/config"
	"github.com/mia0x75/halo/internal/platform/db"
	"github.com/mia0x75/halo/internal/platform/migrate"
	"github.com/mia0x75/halo/internal/platform/seed"
	"github.com/mia0x75/halo/migrations"
)

// Minimum lengths of the secrets the wizard collects.
const (
	MinJWTSecretLength = 32
	MinCryptoLength    = 16
	MinPasswordLength  = 8
)

// Request carries the values collected by the setup wizard.
type Request struct {
	Database DatabaseRequest `json:"database"`
	Secret   SecretRequest   `json:"secret"`
	Admin    AdminRequest    `json:"admin"`
}

// DatabaseRequest describes the metadata database connection.
type DatabaseRequest struct {
	Driver         string `json:"driver"`
	Addr           string `json:"addr"`
	MaxIdle        int    `json:"max_idle"`
	MaxConnections int    `json:"max_connections"`
	WaitTimeout    int    `json:"wait_timeout"`
}

// SecretRequest describes the keys the server signs and encrypts with.
type SecretRequest struct {
	TokenName string `json:"token_name"`
	Issuer    string `json:"issuer"`
	SecretKey string `json:"secret_key"`
	Crypto    string `json:"crypto"`
}

// AdminRequest describes the initial administrator account.
type AdminRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Validate checks the wizard input before anything is written.
func (r Request) Validate() error {
	if strings.TrimSpace(r.Database.Addr) == "" {
		return errors.New("请填写元数据库连接串")
	}
	if r.Database.Driver != "" && r.Database.Driver != "mysql" && r.Database.Driver != "postgres" {
		return fmt.Errorf("未知的元数据库类型: %s", r.Database.Driver)
	}
	if len(r.Secret.SecretKey) < MinJWTSecretLength {
		return fmt.Errorf("JWT 密钥至少需要 %d 个字符", MinJWTSecretLength)
	}
	if len(r.Secret.Crypto) < MinCryptoLength {
		return fmt.Errorf("加密密钥至少需要 %d 个字符", MinCryptoLength)
	}
	if !strings.Contains(r.Admin.Email, "@") {
		return errors.New("请填写管理员邮箱")
	}
	if len(r.Admin.Password) < MinPasswordLength {
		return fmt.Errorf("管理员密码至少需要 %d 个字符", MinPasswordLength)
	}
	return nil
}

// normalize fills in the defaults the wizard leaves empty.
func (r Request) normalize() Request {
	if r.Database.MaxIdle <= 0 {
		r.Database.MaxIdle = 10
	}
	if r.Database.MaxConnections <= 0 {
		r.Database.MaxConnections = 100
	}
	if r.Database.WaitTimeout <= 0 {
		r.Database.WaitTimeout = 300
	}
	if strings.TrimSpace(r.Secret.TokenName) == "" {
		r.Secret.TokenName = "Authentication"
	}
	if strings.TrimSpace(r.Secret.Issuer) == "" {
		r.Secret.Issuer = "halo"
	}
	if strings.TrimSpace(r.Admin.Name) == "" {
		r.Admin.Name = "管理员"
	}
	return r
}

// Run performs the setup and returns the wired application.
func Run(ctx context.Context, request Request) (*app.App, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	request = request.normalize()

	dbCfg := db.Config{
		Addr:           request.Database.Addr,
		MaxIdle:        request.Database.MaxIdle,
		MaxConnections: request.Database.MaxConnections,
		WaitTimeout:    request.Database.WaitTimeout,
		Driver:         request.Database.Driver,
	}

	log.Infof("[I] setup: 连接元数据库")
	pool, err := db.New(ctx, dbCfg)
	if err != nil {
		return nil, fmt.Errorf("连接元数据库失败: %w", err)
	}
	defer func() { _ = pool.Close() }()

	log.Infof("[I] setup: 执行数据库迁移")
	if err := migrateSchema(ctx, pool); err != nil {
		return nil, err
	}
	log.Infof("[I] setup: 写入基础数据")
	if err := loadBaseData(ctx, pool); err != nil {
		return nil, err
	}

	// The config file is written before the application starts so a restart
	// comes up with the same settings.
	if err := config.Save(func(c *config.GlobalConfig) {
		c.Database = &config.DatabaseConfig{
			Driver:         pool.Driver,
			Addr:           dbCfg.Addr,
			MaxIdle:        dbCfg.MaxIdle,
			MaxConnections: dbCfg.MaxConnections,
			WaitTimeout:    dbCfg.WaitTimeout,
		}
		c.Secret = &config.SecretConfig{
			Jwt: &config.JwtConfig{
				TokenName: request.Secret.TokenName,
				Issuer:    request.Secret.Issuer,
				SecretKey: request.Secret.SecretKey,
			},
			Crypto: request.Secret.Crypto,
		}
	}); err != nil {
		return nil, fmt.Errorf("写入配置文件失败: %w", err)
	}

	log.Infof("[I] setup: 启动应用")
	application, err := app.New(ctx, dbCfg)
	if err != nil {
		return nil, fmt.Errorf("初始化应用失败: %w", err)
	}
	log.Infof("[I] setup: 写入初始管理员")
	if err := ensureAdmin(ctx, application, request.Admin); err != nil {
		_ = application.Close()
		return nil, err
	}
	return application, nil
}

// TestConnection tries to reach the metadata database without writing
// anything, so the wizard can verify the connection before the setup runs.
func TestConnection(ctx context.Context, request DatabaseRequest) error {
	if strings.TrimSpace(request.Addr) == "" {
		return errors.New("请填写元数据库连接串")
	}
	pool, err := db.New(ctx, db.Config{
		Addr:           request.Addr,
		MaxIdle:        request.MaxIdle,
		MaxConnections: request.MaxConnections,
		WaitTimeout:    request.WaitTimeout,
		Driver:         request.Driver,
	})
	if err != nil {
		return fmt.Errorf("连接元数据库失败: %w", err)
	}
	return pool.Close()
}

// migrateSchema applies the pending schema migrations.
func migrateSchema(ctx context.Context, pool *db.DB) error {
	fsys, err := migrationFS(pool.Driver)
	if err != nil {
		return err
	}
	migrator, err := migrate.New(pool.DB, migrate.Engine(pool.Driver), fsys)
	if err != nil {
		return fmt.Errorf("准备迁移失败: %w", err)
	}
	if _, err := migrator.Up(ctx); err != nil {
		return fmt.Errorf("执行迁移失败: %w", err)
	}
	return nil
}

// migrationFS returns the migration set of the given engine.
func migrationFS(driver string) (fsys fs.FS, err error) {
	if driver == "postgres" {
		return migrations.Postgres()
	}
	return migrations.MySQL()
}

// loadBaseData loads the roles, rules and templates when the database has
// none, so a fresh install can log in and validate SQL.
func loadBaseData(ctx context.Context, pool *db.DB) error {
	var count int
	if err := pool.GetContext(ctx, &count, "SELECT COUNT(*) FROM mm_roles"); err != nil {
		return fmt.Errorf("读取角色表失败: %w", err)
	}
	if count > 0 {
		return nil
	}
	if err := seed.Apply(ctx, pool, pool.Driver); err != nil {
		return fmt.Errorf("写入基础数据失败: %w", err)
	}
	return nil
}

// ensureAdmin makes sure an administrator account exists with the password the
// wizard collected. The seeded account is reused so no default password is
// left behind.
func ensureAdmin(ctx context.Context, application *app.App, admin AdminRequest) error {
	// The role list is keyset paginated on role_id.
	roles, err := application.Services.Role.List(ctx, math.MaxUint32, 100)
	if err != nil {
		return fmt.Errorf("读取角色表失败: %w", err)
	}
	var role *domain.Role
	for _, item := range roles {
		if strings.EqualFold(item.Name, "Root") {
			role = item
			break
		}
	}
	if role == nil {
		return errors.New("基础数据里没有 Root 角色，无法创建管理员")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("加密管理员密码失败: %w", err)
	}

	user, err := findAdmin(ctx, application, role.RoleID)
	switch {
	case err != nil:
		return err
	case user == nil:
		created := &domain.User{
			UUID:     uuid.NewString(),
			Email:    admin.Email,
			Password: string(hashed),
			Status:   domain.UserStatusNormal,
			Name:     admin.Name,
			Version:  1,
		}
		if _, err := application.Services.User.Create(ctx, created, 0); err != nil {
			return fmt.Errorf("创建管理员失败: %w", err)
		}
	default:
		user.Email = admin.Email
		user.Password = string(hashed)
		user.Name = admin.Name
		user.Status = domain.UserStatusNormal
		user.Version++
		if err := application.Services.User.Update(ctx, user, 0); err != nil {
			return fmt.Errorf("更新管理员失败: %w", err)
		}
	}

	if err := cache.UsersInit(); err != nil {
		return fmt.Errorf("刷新用户缓存失败: %w", err)
	}
	if err := cache.EdgesInit(); err != nil {
		return fmt.Errorf("刷新授权缓存失败: %w", err)
	}
	return nil
}

// findAdmin returns the user that already holds the Root role.
func findAdmin(ctx context.Context, application *app.App, roleID uint) (*domain.User, error) {
	// The edge list is keyset paginated on edge_id, so the cursor has to sort
	// after every row.
	edges, err := application.Services.Edge.List(ctx, math.MaxUint32, 1000)
	if err != nil {
		return nil, fmt.Errorf("读取授权关系失败: %w", err)
	}
	for _, edge := range edges {
		if edge.Type != domain.EdgeTypeUserToRole || edge.DescendantID != roleID {
			continue
		}
		user, err := application.Services.User.Get(ctx, edge.AncestorID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				continue
			}
			return nil, fmt.Errorf("读取管理员失败: %w", err)
		}
		return user, nil
	}
	return nil, nil
}
