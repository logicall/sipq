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
	trace.Trace.Println("enter buildExe")
	defer trace.Trace.Println("exit buildExe")
	os.Chdir(workDir)
	buildCommand := exec.Command("go", "build")
	trace.Trace.Println("buildCommand", buildCommand)
	err := buildCommand.Run()
	util.ErrorPanic(err)
}

func init() {
	workDir, _ = filepath.Abs("..")
	fmt.Println("working dir", workDir)
	buildExe()
	fmt.Println("operating system", runtime.GOOS)
	if runtime.GOOS == "windows" {
		exeName = "sipq.exe"
	} else {
		exeName = "sipq"
	}
	exeFullName = filepath.Join(workDir, exeName)

}

func dumpReader(r io.Reader, wg *sync.WaitGroup) {

	bs, _ := ioutil.ReadAll(r)
	fmt.Println(string(bs))
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
	)

	clientCommand := exec.Command(exeFullName, "-local-port", fmt.Sprint(clientPort),
		"-remote-port", fmt.Sprint(serverPort),
		"-scenario-file", clientScenarioFile,
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
	var waitGroup sync.WaitGroup
	waitGroup.Add(4)
	go dumpReader(serverStderrPipe, &waitGroup)
	go dumpReader(serverStdoutPipe, &waitGroup)
	go dumpReader(clientStderrPipe, &waitGroup)
	go dumpReader(clientStdoutPipe, &waitGroup)

	err = serverCommand.Start()
	if err != nil {
		t.Fatal(err)
	}

	err = clientCommand.Start()
	if err != nil {
		t.Fatal(err)
	}

	waitGroup.Wait() //wait goroutine done

	err = clientCommand.Wait()
	if err != nil {
		t.Fatal(err)
	}

	err = serverCommand.Wait()
	if err != nil {
		t.Fatal(err)
	}

}
