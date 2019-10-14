package pm

import (
	"bufio"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PmTestSuite struct {
	suite.Suite
	cmd *exec.Cmd
}

// Runs before each test in the PmTestSuite.
func (suite *PmTestSuite) SetupTest() {
	u, err := user.Current()
	assert.NoError(suite.T(), err)

	uid, err := strconv.Atoi(u.Uid)
	assert.NoError(suite.T(), err)

	gid, err := strconv.Atoi(u.Gid)
	assert.NoError(suite.T(), err)

	var cred = &syscall.Credential{
		Uid:         uint32(uid),
		Gid:         uint32(gid),
		NoSetGroups: true,
	}

	var attr = &syscall.SysProcAttr{
		Credential: cred,
		Noctty:     false,
	}

	dir, err := os.Getwd()
	suite.T().Logf("dir: %s\n", dir)
	assert.NoError(suite.T(), err)

	suite.cmd = &exec.Cmd{}
	suite.cmd.Dir = dir
	suite.cmd.Path = "./bin/test"
	suite.cmd.Env = os.Environ()
	suite.cmd.Stdin = os.Stdin
	suite.cmd.SysProcAttr = attr
}

func (suite *PmTestSuite) TestPm() {
	// Spawn a child process ruunning `./bin/test`, and attempt to read its memory.
	stdout, err := suite.cmd.StdoutPipe()
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.cmd.Start())

	pid := suite.cmd.Process.Pid
	suite.T().Logf("pid (parent from os): %d\n", os.Getpid())
	suite.T().Logf("pid (child from os): %d\n", pid)
	gid, _ := syscall.Getpgid(pid)
	suite.T().Logf("gid (parent from os): %d\n", os.Getegid())
	suite.T().Logf("gid (child from syscall): %d\n", gid)

	scanner := bufio.NewScanner(stdout)
	lines := make([]byte, 0)
	for scanner.Scan() {
		for _, b := range scanner.Bytes() {
			lines = append(lines, b)
		}
		lines = append(lines, ',', ' ')
	}
	suite.T().Logf("lines: %s\n", string(lines))
	address := uintptr(unsafe.Pointer(&lines[0]))
	suite.T().Logf("address: %p\n", &lines[0])
	suite.T().Logf("uintptr: 0x%x\n", address)

	length := len(lines)
	suite.T().Logf("length: %d\n", length)

	buffer := make([]byte, length)

	proc, err := os.FindProcess(pid)
	assert.NoError(suite.T(), err)
	suite.T().Logf("find: %+v\n", proc)
	assert.NoError(suite.T(), ReadAddress(pid, address, buffer))

	data, err := CopyAddress(pid, address, uint64(length))
	assert.NoError(suite.T(), err)
	assert.EqualValues(suite.T(), lines, data)
	assert.EqualValues(suite.T(), length, len(data))

	assert.NoError(suite.T(), suite.cmd.Wait())
	suite.T().Logf("test cmd exit: %t\n", suite.cmd.ProcessState.Exited())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(PmTestSuite))
}
