package video

import (
	"errors"
	"fmt"
	"os/exec"
)

func Concatenate(destPath string, srcPaths ...string) error {
	if len(destPath) == 0 {
		return errors.New("destPath cannot be empty")
	}
	if len(srcPaths) == 0 {
		return errors.New("no source path provided")
	}

	// construct command in the form:
	// MP4Box -add src_1.mp4 -cat src_2.mp4 ... -cat src_n.mp4 dest.mp4
	cmdStr := fmt.Sprintf("MP4Box -add %s", srcPaths[0])
	for _, srcPath := range srcPaths[1:] {
		cmdStr += fmt.Sprintf(" -cat %s", srcPath)
	}
	cmdStr += fmt.Sprintf(" %s", destPath)

	err := exec.Command(cmdStr).Run()
	return err
}
