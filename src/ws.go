package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/kardianos/service"
)

var cmd *exec.Cmd

type program struct {
	Name      string   //应用程序名
	Args      []string //应用程序参数
	WorkSpace string
	Env       []string
}

func (p *program) Start(s service.Service) error {
	Info.Println("Start ...")
	go p.run()
	return nil
}

func (p *program) run() {
	cmd = exec.Command(p.Name, p.Args[1:]...)
	cmd.Dir = p.WorkSpace
	cmd.Env = append(os.Environ(), p.Env...)
	stdout, err := os.OpenFile(p.WorkSpace+"/out.log", os.O_CREATE|os.O_WRONLY, 0600)
	// Info.Fatal(cmd.Env)
	if err != nil {
		log.Fatalln(err)
	}
	defer stdout.Close()
	cmd.Stdout = stdout
	Info.Println("cd", p.WorkSpace)
	Info.Println(p.Name, strings.Join(p.Args, " "))
	err = cmd.Run()
	if err != nil {
		Error.Fatal("service application has error:", err)
	}
}
func (p *program) Stop(s service.Service) error {
	Info.Println("Stop ...", p.Name)
	cmd.Process.Kill()
	return nil
}
func usage() {
	fmt.Printf("Usage:\n\nws <command> [service name] options.. args ... \n")
	fmt.Println("\nThe command should be:")
	fmt.Printf("  %s\t  %s\n", "install", "install service")
	fmt.Printf("  %s\t%s\n", "uninstall", "remove service")
	fmt.Printf("  %s\t    %s\n", "start", "start service")
	fmt.Printf("  %s\t     %s\n", "stop", "stop service")
	fmt.Printf("  %s\t  %s\n", "restart", "restart service")
	fmt.Printf("  %s\t   %s\n", "status", "remove service")

	fmt.Println("\nThe options are:")
	cmdLine.PrintDefaults()
	fmt.Println("\nUsing args to specify args for application")
}

var cmdLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

func main() {

	env := cmdLine.String("e", "", "specify serivce application `env`, eg. 'key1=value1,key2=value2'")
	dir := cmdLine.String("w", "", "specify service application  `workspace`")
	path := cmdLine.String("p", "", "specify service application `path`")

	if len(os.Args) < 3 {
		usage()
		return
	}
	err := cmdLine.Parse(os.Args[3:])
	if err != nil {
		usage()
		Error.Fatal("Parse args error: ", err.Error())
		return
	}
	cmd := os.Args[1]
	serviceName := os.Args[2]

	workspace := *dir
	if workspace == "" {
		workspace, _ = os.Getwd()
	}

	args := []string{
		"run",
		serviceName,
		"-p", *path,
		"-w", workspace,
		"-e", *env,
	}
	args = append(args, cmdLine.Args()...)

	svcConfig := &service.Config{
		Name:             serviceName, //服务名称
		DisplayName:      serviceName, //显示名称
		Description:      serviceName, //服务描述
		Arguments:        args,
		WorkingDirectory: workspace,
	}
	prg := &program{
		Name:      *path,
		Args:      cmdLine.Args(),
		WorkSpace: workspace,
		Env:       strings.Split(*env, ","),
	}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		usage()
		log.Fatal("create new service error: ", err.Error())
	}

	switch cmd {
	case "install":
		err := s.Install()
		if err != nil {
			usage()
			log.Fatal("install new service error: ", err.Error())
		}
		log.Println(serviceName, " is installed")
		return
	case "uninstall":
		err := s.Uninstall()
		if err != nil {
			usage()
			log.Fatal("remove service ", serviceName, " error: ", err.Error())
		}
		log.Println(serviceName + " is removed")
		return
	case "start":
		err := s.Start()
		if err != nil {
			usage()
			log.Fatal("run service error: ", err.Error())
		}
		return
	case "stop":
		err := s.Stop()
		if err != nil {
			usage()
			log.Fatal("run service error: ", err.Error())
		}
		return
	case "restart":
		err := s.Restart()
		if err != nil {
			usage()
			log.Fatal("run service error: ", err.Error())
		}
		return
	case "status":
		status, err := s.Status()
		if err != nil {
			usage()
			log.Fatal("run service error: ", err.Error())
			return
		}
		fmt.Printf("the status of service %s :%d", serviceName, status)
		return
	case "run":
		err := s.Run()
		if err != nil {
			usage()
			log.Fatal("run service error: ", err.Error())
		}
		return
	}
	usage()
}
