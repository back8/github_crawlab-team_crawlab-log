package log

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func setup(driver *SeaweedFSLogDriver) {
	_ = driver.m.DeleteDir("/logs/" + driver.prefix)
}

func cleanup(driver *SeaweedFSLogDriver) {
	_ = driver.m.DeleteDir("/logs/" + driver.prefix)
}

func TestNewSeaweedFSDriver(t *testing.T) {
	_, err := NewSeaweedFSLogDriver("test")
	require.Nil(t, err)
}

func TestSeaweedFSLogDriver_Write(t *testing.T) {
	driver, err := NewSeaweedFSLogDriver("test")
	require.Nil(t, err)

	setup(driver)

	content0 := ""
	content1 := ""
	for i := 0; i < 10001; i++ {
		line := fmt.Sprintf("line: %d", i+1)
		err = driver.Write(line)
		require.Nil(t, err)
		if i < 1000 {
			content0 += line + "\n"
		} else if 1000 <= i && i < 2000 {
			content1 += line + "\n"
		}
	}

	time.Sleep(4 * time.Second)

	ok, err := driver.m.Exists("/logs/test/00000000")
	require.Nil(t, err)
	require.True(t, ok)

	data, err := driver.m.GetFile("/logs/test/00000000")
	require.Nil(t, err)
	require.Equal(t, content0, string(data))

	ok, err = driver.m.Exists("/logs/test/00000001")
	require.Nil(t, err)
	require.True(t, ok)

	data, err = driver.m.GetFile("/logs/test/00000001")
	require.Nil(t, err)
	require.Equal(t, content1, string(data))

	files, err := driver.m.ListDir("/logs/test")
	require.Nil(t, err)
	require.Equal(t, 11, len(files))

	for i := 10001; i < 20001; i++ {
		line := fmt.Sprintf("line: %d", i+1)
		err = driver.Write(line)
		require.Nil(t, err)
	}

	time.Sleep(4 * time.Second)

	files, err = driver.m.ListDir("/logs/test")
	require.Nil(t, err)
	require.Equal(t, 21, len(files))

	cleanup(driver)
}

func TestSeaweedFSLogDriver_WriteLines(t *testing.T) {
	driver, err := NewSeaweedFSLogDriver("test")
	require.Nil(t, err)

	setup(driver)

	batch := 500
	var lines []string
	for i := 0; i < 10; i++ {
		for j := 0; j < batch; j++ {
			line := fmt.Sprintf("line: %d", i*batch+j+1)
			lines = append(lines, line)
		}
		err = driver.WriteLines(lines)
		require.Nil(t, err)
		lines = []string{}
		time.Sleep(1 * time.Second)
	}

	time.Sleep(3 * time.Second)

	files, err := driver.m.ListDir("/logs/test")
	require.Nil(t, err)
	require.Equal(t, 10*batch/1000, len(files))

	cleanup(driver)
}