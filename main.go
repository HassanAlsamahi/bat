package main

import (
	"fmt"
	"os"
	"syscall"
	"os/exec"
	"io/ioutil"
)
// docker         run image <cmd> <params>
// go run main.go run       <cmd> <params>

func main() {
	var fs string
	switch os.Args[2] {
	case "ubuntu":
		fs = "ubuntu-xenial"

	case "centos":
		fs = "centos-fs"
	}

	switch os.Args[1] {
	case "run":
		run()

	case "child":
		child(fs)
	default:
		panic("bad command")
	}

}


func run() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}
	cmd.Run()
}


func child(fs string) {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	syscall.Sethostname([]byte("container"))

	rootFS(fs)
  syscall.Mount("proc", "/proc", "proc", 0, "")
	cmd := exec.Command("/bin/bash")
	//cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	syscall.Unmount("/proc", 0)
}

func rootFS(fs string) {
	rootfs := "/home/hassanalsamahi/Data/golang/src/goto2018/" + fs
	syscall.Chroot(rootfs)
	syscall.Chdir("/root")
}


func cg() {
	cgroups := "/sys/fs/cgroup"
	pids := filepath.Join(cgroups, "pids")
	err := os.Mkdir(filepath.Join(pids, "pups"), 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	must(ioutil.WriteFile(filepath.Join(pids, "pups/pids.max"), []byte("20"), 0700))
	// Removes the new cgroup in place after the container exits
	must(ioutil.WriteFile(filepath.Join(pids, "liz/notify_on_release"), []byte("1"), 0700))
	must(ioutil.WriteFile(filepath.Join(pids, "liz/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}
