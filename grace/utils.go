package grace

import (
	"os"
	"os/exec"
	"syscall"
)

// Fork starts a new process executing the same command
// as current process with same arguments,
// stdout, stderr and extraFiles are inherited by the new process.
func Fork(extraFiles []*os.File) error {
	path := os.Args[0]
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if len(extraFiles) > 0 {
		cmd.ExtraFiles = extraFiles
	}

	return cmd.Start()
}

// KillParent sends SIGTERM to parent process.
func KillParent() error {
	parent := syscall.Getppid()
	return syscall.Kill(parent, syscall.SIGTERM)
}
