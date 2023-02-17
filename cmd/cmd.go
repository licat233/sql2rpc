/*
 * @Author: licat
 * @Date: 2023-02-07 09:46:20
 * @LastEditors: licat
 * @LastEditTime: 2023-02-15 16:49:07
 * @Description: licat233@gmail.com
 */

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/licat233/sql2rpc/config"
)

type Core interface {
	Allow() bool  //是否可以运行
	Run() error   //运行
	Name() string //名称
}

type Server struct {
	Cores []Core
}

func New() *Server {
	return &Server{
		Cores: []Core{},
	}
}

func (cs *Server) Register(cores ...Core) *Server {
	cs.Cores = append(cs.Cores, cores...)
	return cs
}

func (cs *Server) Run() error {
	has := false
	var listName []string
	for _, c := range cs.Cores {
		listName = append(listName, c.Name())
		if c.Allow() {
			has = true
			if err := c.Run(); err != nil {
				return err
			}
		}
	}
	if !has {
		fmt.Printf("\n- please choose %s\n- Run the \"%s -h\" command for help", strings.Join(listName, "or"), config.ProjectName)
		os.Exit(1)
	}
	return nil
}
