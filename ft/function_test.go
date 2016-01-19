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

	"github.com/henryscala/sipq/util"
)

const (
	clientPortStr = "50700"
	serverPortStr = "50600"
	clientIP      = "127.0.0.1"
	serverIP      = clientIP
)

var (
	workDir     string
	exeName     string
	exeFullName string
)

func buildExe() {
	buildCommand := exec.Command("go", "build")

	os.Chdir(workDir)
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
	serverScenarioFile := filepath.Join(workDir, "scenario", "example", "call_establish_server.js")
	clientScenarioFile := filepath.Join(workDir, "scenario", "example", "call_establish_client.js")
	serverCommand := exec.Command(exeFullName, "-local-port", serverPortStr,
		"-remote-port", clientPortStr,
		"-scenario-file", serverScenarioFile,
	)

	clientCommand := exec.Command(exeFullName, "-local-port", clientPortStr,
		"-remote-port", serverPortStr,
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
