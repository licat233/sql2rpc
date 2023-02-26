/*
 * @Author: licat
 * @Date: 2023-01-15 14:07:25
 * @LastEditors: licat
 * @LastEditTime: 2023-02-16 16:04:18
 * @Description: licat233@gmail.com
 */

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/licat233/sql2rpc/cmd"
	"github.com/licat233/sql2rpc/db"
	"github.com/licat233/sql2rpc/update"

	"github.com/licat233/sql2rpc/cmd/api"
	"github.com/licat233/sql2rpc/cmd/model"
	"github.com/licat233/sql2rpc/cmd/pb"
	"github.com/licat233/sql2rpc/config"
)

var _startTime = time.Now()

func main() {
	Initialize()
	if err := db.InitConn(); err != nil {
		log.Fatal(err)
	}
	defer func(conn *sql.DB) {
		if conn == nil {
			return
		}
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db.Conn)

	if err := cmd.New().Register(pb.New(), api.New(), model.New()).Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done.")
	elapsed := time.Since(_startTime)
	fmt.Println("耗时：", elapsed)
}

func Initialize() {
	arges := os.Args
	apiArgStatus := false
	pbArgStatus := false
	modelArgStatus := false
	if len(arges) > 1 {
		argV := arges[1]
		if argV[0] != '-' {
			if argV == "upgrade" || argV == "up" {
				update.Update()
				os.Exit(0)
			} else if argV == config.ApiCoreName {
				apiArgStatus = true
			} else if argV == config.PbCoreName {
				pbArgStatus = true
			} else if argV == config.ModelCoreName {
				modelArgStatus = true
			} else {
				fmt.Println("Input command error, please enter -help or -h for help")
				os.Exit(1)
			}
		}
	}

	defaultConfig := config.NewDefaultConfig()
	//base flag
	initConfigFile := flag.Bool("init", false, "Create default config file，priority is given to the data in this file")
	version1 := flag.Bool("v", false, "Current version")
	version2 := flag.Bool("version", false, "Current version")
	upgrade1 := flag.Bool("up", false, "Upgrade sql2rpc to latest version")
	upgrade2 := flag.Bool("upgrade", false, "Upgrade sql2rpc to latest version")

	//common flag
	serviceName := flag.String(defaultConfig.ServiceName.FlagString())
	fileName := flag.String(defaultConfig.Filename.FlagString())
	//database flag
	dbType := flag.String(defaultConfig.DBType.FlagString())
	dbHost := flag.String(defaultConfig.DBHost.FlagString())
	dbPort := flag.Int(defaultConfig.DBPort.FlagInt())
	dbUser := flag.String(defaultConfig.DBUser.FlagString())
	dbPassword := flag.String(defaultConfig.DBPassword.FlagString())
	dbSchema := flag.String(defaultConfig.DBSchema.FlagString())
	dbTable := flag.String(defaultConfig.DBTable.FlagString())
	//ignore flag
	ignoreTableStr := flag.String(defaultConfig.IgnoreTableStr.FlagString())
	ignoreColumnStr := flag.String(defaultConfig.IgnoreColumnStr.FlagString())
	//pbStatus flag
	pbStatus := flag.Bool(defaultConfig.Pb.FlagBool())
	pbPackageName := flag.String(defaultConfig.PbPackage.FlagString())
	pbGoPackageName := flag.String(defaultConfig.PbGoPackage.FlagString())
	pbMultiple := flag.Bool(defaultConfig.PbMultiple.FlagBool())
	// apiStatus flag
	apiStatus := flag.Bool(defaultConfig.Api.FlagBool())
	apiStyle := flag.String(defaultConfig.ApiStyle.FlagString())
	apiJwt := flag.String(defaultConfig.ApiJwt.FlagString())
	apiMiddleware := flag.String(defaultConfig.ApiMiddleware.FlagString())
	apiPrefix := flag.String(defaultConfig.ApiPrefix.FlagString())
	apiMultiple := flag.Bool(defaultConfig.ApiMultiple.FlagBool())
	// modelStatus flag
	modelStatus := flag.Bool(defaultConfig.Model.FlagBool())
	flag.Parse()

	*pbStatus = *pbStatus || pbArgStatus
	*apiStatus = *apiStatus || apiArgStatus
	*modelStatus = *modelStatus || modelArgStatus

	if *version1 || *version2 {
		fmt.Println("Current version:", config.CurrentVersion)
		os.Exit(0)
	}

	if *upgrade1 || *upgrade2 {
		update.Update()
		os.Exit(0)
	}

	if *initConfigFile {
		if err := config.CreateConfigFile(config.DefaultFileName); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done.")
		os.Exit(0)
	}

	cmdConfig := &config.Config{
		DBType:          defaultConfig.DBType.Set(*dbType),
		DBHost:          defaultConfig.DBHost.Set(*dbHost),
		DBPort:          defaultConfig.DBPort.Set(*dbPort),
		DBUser:          defaultConfig.DBUser.Set(*dbUser),
		DBPassword:      defaultConfig.DBPassword.Set(*dbPassword),
		DBSchema:        defaultConfig.DBSchema.Set(*dbSchema),
		DBTable:         defaultConfig.DBTable.Set(*dbTable),
		IgnoreTableStr:  defaultConfig.IgnoreTableStr.Set(*ignoreTableStr),
		IgnoreColumnStr: defaultConfig.IgnoreColumnStr.Set(*ignoreColumnStr),
		ServiceName:     defaultConfig.ServiceName.Set(*serviceName),
		Filename:        defaultConfig.Filename.Set(*fileName),
		Pb:              defaultConfig.Pb.Set(*pbStatus),
		PbPackage:       defaultConfig.PbPackage.Set(*pbPackageName),
		PbGoPackage:     defaultConfig.PbGoPackage.Set(*pbGoPackageName),
		PbMultiple:      defaultConfig.PbMultiple.Set(*pbMultiple),
		Api:             defaultConfig.Api.Set(*apiStatus),
		ApiStyle:        defaultConfig.ApiStyle.Set(*apiStyle),
		ApiJwt:          defaultConfig.ApiJwt.Set(*apiJwt),
		ApiMiddleware:   defaultConfig.ApiMiddleware.Set(*apiMiddleware),
		ApiPrefix:       defaultConfig.ApiPrefix.Set(*apiPrefix),
		ApiMultiple:     defaultConfig.ApiMultiple.Set(*apiMultiple),
		Model:           defaultConfig.Model.Set(*modelStatus),
	}

	config.C = cmdConfig.Assignment(cmdConfig)
	if err := config.C.Validate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	config.C.Initialize()
	if config.C == nil {
		log.Fatal("- config.C is nil")
	}
}
