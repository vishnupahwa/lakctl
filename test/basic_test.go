package test

import (
	"bytes"
	. "github.com/vishnupahwa/lakctl/test/short"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"
)

const port = ":8008"
const commandPort = ":9999"

var lakctl *exec.Cmd

func init() {
	_ = os.Chdir("..")
}
func setup(t *testing.T) {
	lakctl = exec.Command("./lakctl", "start", "-r", "go run ./testdata/testapp/main.go World")
	lakctl.Stdout = os.Stdout
	lakctl.Stderr = os.Stderr
	lakctl.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	Must(t, lakctl.Start())
	time.Sleep(1 * time.Second)
}

func cleanup() {
	pid := lakctl.Process.Pid
	_ = syscall.Kill(-pid, syscall.SIGTERM)
	_ = lakctl.Process.Signal(syscall.SIGTERM)
	_ = lakctl.Wait()
}

func TestStoppingCommand(t *testing.T) {
	setup(t)
	t.Cleanup(cleanup)
	_, err := http.Get("http://localhost" + port + "/stop")
	Ok(t, err)
	time.Sleep(1 * time.Second)

	_, err = http.Get("http://localhost" + commandPort)
	Assert(t, err != nil, "expected error from get request")
}

func TestStoppingThenStartingCommand(t *testing.T) {
	setup(t)
	t.Cleanup(cleanup)
	_, err := http.Get("http://localhost" + port + "/stop")
	Ok(t, err)
	_, err = http.Get("http://localhost" + port + "/start")
	Ok(t, err)
	time.Sleep(1 * time.Second)

	response, _ := http.Get("http://localhost" + commandPort)
	body := toString(t, response)
	Equals(t, "Hello World", body)
}

func TestStoppingThenStartingThenStoppingCommand(t *testing.T) {
	setup(t)
	t.Cleanup(cleanup)
	_, err := http.Get("http://localhost" + port + "/stop")
	Ok(t, err)
	_, err = http.Get("http://localhost" + port + "/start")
	Ok(t, err)
	_, err = http.Get("http://localhost" + port + "/stop")
	Ok(t, err)
	time.Sleep(1 * time.Second)

	_, err = http.Get("http://localhost" + commandPort)
	Assert(t, err != nil, "expected error from get request")
}

func TestRestart(t *testing.T) {
	setup(t)
	t.Cleanup(cleanup)
	replaceJson := strings.NewReader(`{
		"old": "World",
		"new": "世界"
	}`)
	_, err := http.Post("http://localhost"+port+"/restart/", "application/json", replaceJson)
	Ok(t, err)
	time.Sleep(1 * time.Second)

	response, _ := http.Get("http://localhost" + commandPort)
	body := toString(t, response)
	Equals(t, "Hello 世界", body)
}

func TestRestartWithNoSubstitution(t *testing.T) {
	setup(t)
	t.Cleanup(cleanup)
	replaceJson := strings.NewReader(`{}`)
	_, err := http.Post("http://localhost"+port+"/restart/", "application/json", replaceJson)
	Ok(t, err)
	time.Sleep(1 * time.Second)

	response, _ := http.Get("http://localhost" + commandPort)
	body := toString(t, response)
	Equals(t, "Hello World", body)
}

func toString(t *testing.T, response *http.Response) string {
	if response == nil {
		t.Fatal("No response was returned")
	}
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(response.Body)
	r := buf.String()
	return r
}
