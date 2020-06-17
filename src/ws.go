package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/kardianos/service"
)

var (
	cmd *exec.Cmd
)

type program struct {
	Name      string   //应用程序名
	Args      []string //应用程序参数
	WorkSpace string
	Env       []string
	logPath   string
}

func (p *program) Start(s service.Service) error {
	log.Println("Start ...")
	go p.run()
	return nil
}

func (p *program) run() {
	cmd = exec.Command(p.Name, p.Args[1:]...)
	cmd.Dir = p.WorkSpace
	cmd.Env = append(os.Environ(), p.Env...)
	stdout, err := os.OpenFile(p.WorkSpace+"/"+p.logPath+"/out."+time.Now().Format("20200202")+".log", os.O_CREATE|os.O_WRONLY, 0600)

	if err != nil {
		log.Fatalf("create log file:%s error:%s", p.WorkSpace+"/"+p.logPath+"/out."+time.Now().Format("20200202")+".log", err.Error())
	}
	defer stdout.Close()
	// cmd.Stdout = stdout
	cmd.Stdout = io.MultiWriter(stdout, os.Stdout)
	log.Println("cd", p.WorkSpace)
	log.Println(p.Name, strings.Join(p.Args, " "))
	err = cmd.Run()
	if err != nil {
		log.Fatalln("service application has error:", err)
	}
}
func (p *program) Stop(s service.Service) error {
	log.Println("Stop ...", p.Name)
	cmd.Process.Kill()
	return nil
}
func usage() {
	s := fmt.Sprintf("Usage:\n\nws <command> [service name] options.. args ... \n")
	s += fmt.Sprintf("\nThe command should be:\n")
	s += fmt.Sprintf("  %s\t%s\n", "install", "install service")
	s += fmt.Sprintf("  %s\t%s\n", "uninstall", "remove service")
	s += fmt.Sprintf("  %s\t\t%s\n", "start", "start service")
	s += fmt.Sprintf("  %s\t\t%s\n", "stop", "stop service")
	s += fmt.Sprintf("  %s\t%s\n", "restart", "restart service")
	s += fmt.Sprintf("  %s\t%s\n", "status", "remove service")
	s += fmt.Sprintf("\nThe options are:\n")
	fmt.Print(s)
	cmdLine.PrintDefaults()
	fmt.Print("\nUsing args to specify args for application\n\n")
}

var cmdLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

func main() {
	dir := cmdLine.String("w", "", "specify service absolute `workspace` path")
	env := cmdLine.String("e", "", "specify serivce application `env`, eg. 'key1=value1,key2=value2'")
	path := cmdLine.String("p", "", "specify service absolute application `path`")
	logPath := cmdLine.String("l", "./", "specify service log `relative path`")

	if len(os.Args) < 3 {
		usage()
		return
	}
	err := cmdLine.Parse(os.Args[3:])
	if err != nil {
		usage()
		return
	}

	workspace := *dir
	if workspace == "" {
		workspace, _ = os.Getwd()
	}

	file, err := os.OpenFile(workspace+"/"+*logPath+"/ws."+time.Now().Format("20200202")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open info log file:", err.Error())
		return
	}
	log.SetOutput(io.MultiWriter(file, os.Stdout))

	cmd := os.Args[1]
	serviceName := os.Args[2]

	args := []string{
		"run",
		serviceName,
		"-p", *path,
		"-w", workspace,
		"-l", *logPath,
		"-e", *env,
	}
	args = append(args, cmdLine.Args()...)

	svcConfig := &service.Config{
		Name:             serviceName,                        //服务名称
		DisplayName:      serviceName,                        //显示名称
		Description:      "This service runs " + serviceName, //服务描述
		Arguments:        args,
		WorkingDirectory: workspace,
	}
	prg := &program{
		Name:      *path,
		Args:      cmdLine.Args(),
		WorkSpace: workspace,
		logPath:   *logPath,
		Env:       strings.Split(*env, ","),
	}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		usage()
		log.Fatalln("create new service error: ", err.Error())
	}

	switch cmd {
	case "install":
		err := s.Install()
		if err != nil {
			usage()
			log.Fatalf("Install service %s error: %s\n", serviceName, err.Error())
		}
		log.Printf("The service %s is installed", serviceName)
		return
	case "uninstall":
		err := s.Uninstall()
		if err != nil {
			usage()
			log.Fatalf("Uninstall service %s error: %s\n", serviceName, err.Error())
		}
		log.Printf("The service %s is removed", serviceName)
		return
	case "start":
		err := s.Start()
		if err != nil {
			usage()
			log.Fatalf("Start service %s error: %s\n", serviceName, err.Error())
		}
		return
	case "stop":
		err := s.Stop()
		if err != nil {
			usage()
			log.Fatalf("Stop service %s error: %s\n", serviceName, err.Error())
		}
		return
	case "restart":
		err := s.Restart()
		if err != nil {
			usage()
			log.Fatalf("Restart service %s error: %s\n", serviceName, err.Error())
		}
		return
	case "status":
		status, err := s.Status()
		if err != nil {
			usage()
			log.Fatalf("Status service %s error:%s\n", serviceName, err.Error())
			return
		}
		if status == 1 {
			log.Printf("Status:%d\nThe service %s is running", status, serviceName)
		} else {
			log.Printf("The service %s is not running", serviceName)
		}
		return
	case "run":
		err := s.Run()
		if err != nil {
			usage()
			log.Fatalf("Run service %s error:%s\n", serviceName, err.Error())
		}
		return
	}
	usage()
}
