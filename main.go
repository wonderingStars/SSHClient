package main

import (
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	key           = "key.pem"
	publicIPNPort = "192.168.0.1:22"
)

func SendCommand(in io.WriteCloser, cmd string) error {
	if _, err := in.Write([]byte(cmd + "\n")); err != nil {
		return err
	}

	return nil
}
func main() {

	pemBytes, err := ioutil.ReadFile(key)
	if err != nil {
		log.Fatal(err)
	}
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		log.Fatalf("parse key failed:%v", err)
	}
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 5,
		User:            "ubuntu",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the client
	client, err := ssh.Dial("tcp", publicIPNPort, config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Create a session
	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Setup StdinPipe to send commands
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stdin.Close()

	// Route session Stdout/Stderr to system Stdout/Stderr
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// Start a shell
	if err := session.Shell(); err != nil {
		log.Fatal(err)
	}

	// Run configuration commands
	SendCommand(stdin, "docker run -p 9987:9987/udp -p 10011:10011 -p 30033:30033 -e TS3SERVER_LICENSE=accept teamspeak ")
	//SendCommand(stdin, "ls")
	//SendCommand(stdin, "ls")

	session.Wait()

}
