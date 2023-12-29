package main

import (
	"flag"
	"fmt"
	goreMysql "git.tenvine.cn/backend/gore/db/mysql"
	"git.tenvine.cn/backend/gore/gonfig"
	"git.tenvine.cn/backend/gore/log"
	"gorm.io/gen"
	"os"
	"strings"
	"time"
)

func errExit(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

var (
	flagEnv     string
	flagName    string
	flagConsul  string
	flagTables  string
	flagOutPath string
)

func main() {
	flag.StringVar(&flagOutPath, "o", "./internal/query", "output path")
	flag.StringVar(&flagEnv, "env", "sdev0", "Environment")
	flag.StringVar(&flagName, "name", "", "Application Name")
	flag.StringVar(&flagConsul, "consul", "i-consul-${profile}.xk5.com:8500", "consul host:port")
	flag.StringVar(&flagTables, "tables", "", "Table names. Multiple table names are separated by commas ','. e.g. tb_member,tb_character")

	flag.Parse()

	if len(flagName) == 0 {
		errExit("invalid name")
	}

	gonfig.Instance().Set("env", flagEnv)
	gonfig.Instance().Set("name", flagName)
	gonfig.Instance().Set("consul", flagConsul)

	if err := gonfig.Setup(); err != nil {
		errExit("gonfig setup err: %s", err.Error())
	}
	if err := log.SetupDefault(); err != nil {
		errExit("log setup err: %s", err.Error())
	}

	if err := goreMysql.SetupDefault(); err != nil {
		panic(err)
	}
	if err := goreMysql.SetupGorm(); err != nil {
		panic(err)
	}

	g := gen.NewGenerator(gen.Config{OutPath: flagOutPath, Mode: gen.WithDefaultQuery})

	g.UseDB(goreMysql.GormDB())

	tableModels := make([]any, 0)
	if len(flagTables) == 0 {
		tableModels = g.GenerateAllTable()
	} else {
		for _, tableName := range strings.Split(flagTables, ",") {
			tableModels = append(tableModels, g.GenerateModel(tableName))
		}
	}
	g.ApplyBasic(tableModels...)
	g.ApplyInterface(func(Querier) {}, tableModels...)

	g.Execute()
}

type Querier interface {
	// FilterWithTime
	//	SELECT * FROM @@table
	// 		{{where}}
	//			{{if !begin.IsZero()}}
	//				created_time > @begin
	//			{{end}}
	//			{{if !end.IsZero()}}
	//				AND created_time < @end
	//			{{end}}
	//		{{end}}
	FilterWithTime(begin, end time.Time) ([]gen.T, error)
}
