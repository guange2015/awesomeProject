package main

import (
	"awesomeProject/sync"
	"flag"
	"fmt"
	"github.com/howeyc/gopass"
	"log"
	"os"
	"io/ioutil"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, `gsync version: gsync/0.0.1
Usage: gsync [-hvVtTq] [-u username] [-p password] [-s sshaddr] localdir remotedir
 
Options:
`)

	flag.PrintDefaults()
}

func main() {
	username := flag.String("u", "", "ssh用户名")
	password := flag.String("p", "", "ssh密码")
	sshAddr := flag.String("s", "", "ssh地址: 127.0.0.1:2222")
	ignoreFile := flag.String("i", "", "忽略文件列表: ignore.txt，支持绝对匹配和正则匹配")

	flag.Parse()

	if len(*sshAddr) <= 0 {
		usage()
		return
	}
	if len(os.Args) < 3 {
		usage()
		return
	}

	if len(*username) <= 0 {
		usage()
		return
	}

	if len(*password) <= 0 {
		fmt.Printf("Password: ")
		_pwd, err := gopass.GetPasswd()
		if err != nil {
			log.Fatal("输入密码错误", err)
		}
		*password = string(_pwd)
	}

	var ignoreList []string
	if len(*ignoreFile)>0 {
		bytes, err := ioutil.ReadFile(*ignoreFile)
		if err != nil {
			log.Fatal("读取文件错误", err)
		}
		ignoreList = strings.Split(string(bytes), "\n")
	}

	g := sync.NewGsync(*sshAddr, *username, *password, ignoreList)
	g.SyncDir(flag.Arg(0), flag.Arg(1))
}
