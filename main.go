package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {

	runProtoc()

}

func runProtoc() {
	os.Chdir("test")
	execCMD("mkdir", "-p", "todo")
	execCMD("protoc", "-I../", "-I.", "--go_out=plugins=grpc+qr:todo", "todo.proto")
}

func execCMD(name string, args ...string) {

	cmd := exec.Command(name, args...)
	// Combine stdout and stderr
	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printError(err)
	printOutput(output)

}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}
