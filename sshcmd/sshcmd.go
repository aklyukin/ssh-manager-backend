package sshcmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"bytes"
	"log"
)

type SSHCommand struct {
	Path   string
	Env    []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type SSHClient struct {
	Config *ssh.ClientConfig
	Host   string
	Port   int
}

type sshcmd struct {
	client SSHClient
	cmd SSHCommand
}


func (client *SSHClient) RunCommand(cmd *SSHCommand) error {
	var (
		session *ssh.Session
		err     error
	)

	if session, err = client.newSession(); err != nil {
		return err
	}
	defer session.Close()

	if err = client.prepareCommand(session, cmd); err != nil {
		return err
	}

	err = session.Run(cmd.Path)
	return err
}

func (client *SSHClient) prepareCommand(session *ssh.Session, cmd *SSHCommand) error {
	for _, env := range cmd.Env {
		variable := strings.Split(env, "=")
		if len(variable) != 2 {
			continue
		}

		if err := session.Setenv(variable[0], variable[1]); err != nil {
			return err
		}
	}

	if cmd.Stdin != nil {
		stdin, err := session.StdinPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdin for session: %v", err)
		}
		go io.Copy(stdin, cmd.Stdin)
	}

	if cmd.Stdout != nil {
		stdout, err := session.StdoutPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdout for session: %v", err)
		}
		go io.Copy(cmd.Stdout, stdout)
	}

	if cmd.Stderr != nil {
		stderr, err := session.StderrPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stderr for session: %v", err)
		}
		go io.Copy(cmd.Stderr, stderr)
	}

	return nil
}

func (client *SSHClient) newSession() (*ssh.Session, error) {
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", client.Host, client.Port), client.Config)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial: %s", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %s", err)
	}

	modes := ssh.TerminalModes{
		// ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		return nil, fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	return session, nil
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func RunCmd(hostname string){

	sshConfig := &ssh.ClientConfig{
		User: "root",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			PublicKeyFile("/Users/andrey/.ssh/fbs_new_id_rsa"),
		},
	}

	client := &SSHClient{
		Config: sshConfig,
		Host:   hostname,
		Port:   22,
	}

	var stdoutBuf bytes.Buffer

	cmd := &SSHCommand{
		Path:   "ls -l /",
		Env:    []string{""},
		Stdin:  os.Stdin,
		Stdout: &stdoutBuf,
		Stderr: os.Stderr,
	}

	fmt.Printf("Running command: %s\n", cmd.Path)

	if err := client.RunCommand(cmd); err != nil {
		fmt.Fprintf(os.Stderr, "command run error: %s\n", err)
		os.Exit(1)
	}
	log.Printf(stdoutBuf.String())
}

func getCmd(hostname string) sshcmd {

	sshConfig := &ssh.ClientConfig{
		User:            "root",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			PublicKeyFile("/Users/andrey/.ssh/fbs_new_id_rsa"),
		},
	}

	sshcmd := &sshcmd{
		client: SSHClient{
			Config: sshConfig,
			Host:   hostname,
			Port:   22,
		},
		cmd: SSHCommand{
			Path:   "ls -l /",
			Env:    []string{""},
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		},
	}
	return *sshcmd
}

func GetUsers(hostname string) []string {
	cmd := getCmd(hostname)
	var stdoutBuf bytes.Buffer
    cmd.cmd.Path = "grep -h -v -e ':/sbin/nologin$' -e ':/bin/sync$' -e ':/bin/false$' -e ':/sbin/shutdown$' -e ':/sbin/halt$' /etc/passwd | awk -F: '$3 >= 500 {print $1}'"
	cmd.cmd.Stdout = &stdoutBuf
	log.Printf("Running command: %s\n", cmd.cmd.Path)

	if err := cmd.client.RunCommand(&cmd.cmd); err != nil {
		fmt.Fprintf(os.Stderr, "command run error: %s\n", err)
		os.Exit(1)
	}

	var userList []string
	for _, str := range strings.Split(stdoutBuf.String(), "\r\n") {
		if str != "" {
			userList = append(userList, str)
		}
	}
	return userList
}

//func GetKeysForUser(hostName string,userName string) []string {
//	cmd := getCmd(hostName)
//	var stdoutBuf bytes.Buffer
//	cmd.cmd.Path = "grep -h -v -e ':/sbin/nologin$' -e ':/bin/sync$' -e ':/bin/false$' -e ':/sbin/shutdown$' -e ':/sbin/halt$' /etc/passwd | awk -F: '$3 >= 500 {print $1}'"
//	cmd.cmd.Stdout = &stdoutBuf
//	fmt.Printf("Running command: %s\n", cmd.cmd.Path)
//
//	if err := cmd.client.RunCommand(&cmd.cmd); err != nil {
//		fmt.Fprintf(os.Stderr, "command run error: %s\n", err)
//		os.Exit(1)
//	}
//	log.Printf(stdoutBuf.String())
//	var userList []string
//	for _, str := range strings.Split(stdoutBuf.String(), "\n") {
//		if str != "" {
//			userList = append(userList, str)
//		}
//	}
//	return userList
//}