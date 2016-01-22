package ft

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"testing"

	"github.com/henryscala/sipq/trace"
	"github.com/henryscala/sipq/util"
)

const (
	clientIP = "127.0.0.1"
	serverIP = clientIP
)

var (
	workDir     string
	exeName     string
	exeFullName string
)

func buildExe() {
	trace.Trace("enter buildExe")
	defer trace.Trace("exit buildExe")
	os.Chdir(workDir)
	buildCommand := exec.Command("go", "build")
	trace.Debug("buildCommand", buildCommand)
	err := buildCommand.Run()
	util.ErrorPanic(err)
}

func init() {
	workDir, _ = filepath.Abs("..")
	trace.Debug("working dir", workDir)
	buildExe()
	trace.Debug("operating system", runtime.GOOS)
	if runtime.GOOS == "windows" {
		exeName = "sipq.exe"
	} else {
		exeName = "sipq"
	}
	exeFullName = filepath.Join(workDir, exeName)

}

func dumpReader(r io.Reader, prefix string, wg *sync.WaitGroup) {
	bs, _ := ioutil.ReadAll(r)
	fmt.Println("##########begin############")
	fmt.Println("[", prefix, "[")
	fmt.Println(string(bs))
	fmt.Println("]]")
	fmt.Println("##########end############")
	wg.Done()
}

func TestCallEstablishUDP(t *testing.T) {

	ports, err := util.FindFreePort("udp", clientIP, 2)
	if err != nil {
		t.Fatal("cannot find free ports", err)
	}
	serverPort := ports[0]
	clientPort := ports[1]
	t.Log("serverPort", serverPort)
	t.Log("clientPort", clientPort)

	serverScenarioFile := filepath.Join(workDir, "scenario", "example", "call_establish_server.js")
	clientScenarioFile := filepath.Join(workDir, "scenario", "example", "call_establish_client.js")
	serverCommand := exec.Command(exeFullName, "-local-port", fmt.Sprint(serverPort),
		"-remote-port", fmt.Sprint(clientPort),
		"-scenario-file", serverScenarioFile,
		"-transport-type", "udp",
	)

	clientCommand := exec.Command(exeFullName, "-local-port", fmt.Sprint(clientPort),
		"-remote-port", fmt.Sprint(serverPort),
		"-scenario-file", clientScenarioFile,
		"-transport-type", "udp",
	)
	t.Log("server command", serverCommand)
	t.Log("client command", clientCommand)

	serverStderrPipe, err := serverCommand.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}
	serverStdoutPipe, err := serverCommand.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	clientStderrPipe, err := clientCommand.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}
	clientStdoutPipe, err := clientCommand.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	wg.Add(4)
	go dumpReader(serverStderrPipe, "server-stderr", &wg)
	go dumpReader(serverStdoutPipe, "server-stdout", &wg)
	go dumpReader(clientStderrPipe, "client-stderr", &wg)
	go dumpReader(clientStdoutPipe, "client-stdout", &wg)

	err = serverCommand.Start()
	if err != nil {
		t.Fatal(err)
	}

	err = clientCommand.Start()
	if err != nil {
		t.Fatal(err)
	}

	wg.Wait()

	err = clientCommand.Wait()
	if err != nil {
		t.Fatal(err)
	}

	err = serverCommand.Wait()
	if err != nil {
		t.Fatal(err)
	}

}

func TestCallEstablishTCP(t *testing.T) {

	ports, err := util.FindFreePort("tcp", clientIP, 2)
	if err != nil {
		t.Fatal("cannot find free ports", err)
	}
	t.Log("ports", ports)
	serverPort := ports[0]
	clientPort := ports[1]
	t.Log("serverPort", serverPort)
	t.Log("clientPort", clientPort)

	serverScenarioFile := filepath.Join(workDir, "scenario", "example", "call_establish_server.js")
	clientScenarioFile := filepath.Join(workDir, "scenario", "example", "call_establish_client.js")
	serverCommand := exec.Command(exeFullName, "-local-port", fmt.Sprint(serverPort),
		"-remote-port", fmt.Sprint(clientPort),
		"-scenario-file", serverScenarioFile,
		"-transport-type", "tcp",
	)

	clientCommand := exec.Command(exeFullName, "-local-port", fmt.Sprint(clientPort),
		"-remote-port", fmt.Sprint(serverPort),
		"-scenario-file", clientScenarioFile,
		"-transport-type", "tcp",
	)
	t.Log("server command", serverCommand)
	t.Log("client command", clientCommand)

	serverStderrPipe, err := serverCommand.StderrPipe()

	if err != nil {
		t.Fatal(err)
	}
	serverStdoutPipe, err := serverCommand.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	clientStderrPipe, err := clientCommand.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}

	clientStdoutPipe, err := clientCommand.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(4)
	go dumpReader(serverStderrPipe, "server-stderr", &wg)
	go dumpReader(serverStdoutPipe, "server-stdout", &wg)
	go dumpReader(clientStderrPipe, "client-stderr", &wg)
	go dumpReader(clientStdoutPipe, "client-stdout", &wg)

	err = serverCommand.Start()
	if err != nil {
		t.Fatal(err)
	}

	err = clientCommand.Start()
	if err != nil {
		t.Fatal(err)
	}

	wg.Wait()

	err = clientCommand.Wait()
	if err != nil {
		t.Fatal(err)
	}

	err = serverCommand.Wait()
	if err != nil {
		t.Fatal(err)
	}

}
