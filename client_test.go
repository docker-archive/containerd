package containerd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"testing"
)

func init() {
	flag.StringVar(&address, "address", "/run/containerd/containerd.sock", "The address to the containerd socket for use in the tests")
	flag.Parse()
}

const defaultRoot = "/var/lib/containerd-test"

func TestMain(m *testing.M) {
	// setup a new containerd daemon
	cmd := exec.Command("containerd", "--root", defaultRoot)
	// TODO: what todo with IO?
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	status := m.Run()

	// tear down the daemon and resources created
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		fmt.Println(err)
	}
	if _, err := cmd.Process.Wait(); err != nil {
		fmt.Println(err)
	}
	if err := os.RemoveAll(defaultRoot); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(status)
}

var address string

func TestNewClient(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}
	client, err := New(address)
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("New() returned nil client")
	}
	if err := client.Close(); err != nil {
		t.Errorf("client closed returned errror %v", err)
	}
}

func TestImagePull(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}
	client, err := New(address)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	const ref = "docker.io/library/alpine:latest"
	_, err = client.Pull(context.Background(), ref)
	if err != nil {
		t.Error(err)
		return
	}
}
