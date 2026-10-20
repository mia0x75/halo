package cli

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/mia0x75/halo/internal/platform/config"
	"github.com/mia0x75/halo/internal/platform/db"
	"github.com/mia0x75/halo/internal/platform/migrate"
	"github.com/mia0x75/halo/migrations"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "管理元数据库 schema 迁移",
	Long:  "对元数据库执行版本化 schema 迁移，并维护 mm_migrations 台账",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "执行尚未应用的迁移",
	Run:   migrateUp,
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看迁移与台账状态",
	Run:   migrateStatus,
}

func init() {
	RootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
}

// migrationFS returns the embedded migration set of the given engine.
func migrationFS(driver string) (fs.FS, error) {
	if driver == "postgres" {
		return migrations.Postgres()
	}
	return migrations.MySQL()
}

// openMigrator opens the metadata database and builds a migrator for it.
func openMigrator() (*migrate.Migrator, *db.DB, error) {
	if cfg, ok := os.LookupEnv("HALO_CFG"); ok {
		config.ParseConfig(cfg)
	} else {
		config.ParseConfig("")
	}

	dbCfg := config.Config().Database
	if dbCfg == nil || dbCfg.Addr == "" {
		return nil, nil, fmt.Errorf("缺少元数据库配置（database.addr）")
	}

	pool, err := db.New(context.Background(), db.Config{
		Addr:           dbCfg.Addr,
		MaxIdle:        dbCfg.MaxIdle,
		MaxConnections: dbCfg.MaxConnections,
		WaitTimeout:    dbCfg.WaitTimeout,
		Driver:         dbCfg.Driver,
	})
	if err != nil {
		return nil, nil, err
	}

	fsys, err := migrationFS(pool.Driver)
	if err != nil {
		_ = pool.Close()
		return nil, nil, err
	}

	migrator, err := migrate.New(pool.DB, migrate.Engine(pool.Driver), fsys)
	if err != nil {
		_ = pool.Close()
		return nil, nil, err
	}

	return migrator, pool, nil
}

func migrateUp(cmd *cobra.Command, args []string) {
	migrator, pool, err := openMigrator()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() { _ = pool.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	applied, err := migrator.Up(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(applied) == 0 {
		fmt.Println("没有需要执行的迁移")
		return
	}
	for _, mig := range applied {
		fmt.Printf("已应用迁移 %s\n", mig.Filename)
	}
}

func migrateStatus(cmd *cobra.Command, args []string) {
	migrator, pool, err := openMigrator()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() { _ = pool.Close() }()

	status, err := migrator.Status(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, item := range status {
		switch item.State {
		case migrate.StateApplied:
			fmt.Printf("已应用  %s\n", item.Migration.Filename)
		case migrate.StateFailed:
			fmt.Printf("失败    %s\n", item.Migration.Filename)
		default:
			fmt.Printf("待执行  %s\n", item.Migration.Filename)
		}
	}
}
