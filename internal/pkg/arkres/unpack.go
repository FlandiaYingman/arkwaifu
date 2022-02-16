package arkres

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	extractorLocation = "./arkwaifu-extractor"
)

func unpackResources(src string, dst string) error {
	//TODO: Change to a direct call of Python API.

	// check extractor existence
	_, err := os.Stat(extractorLocation)
	if err != nil {
		return err
	}

	install := exec.Command("pipenv", "install")
	install.Dir = "./arkwaifu-extractor"
	output, err := install.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w (%v), output:%q", err, install, output)
	}

	srcAbs, err := filepath.Abs(src)
	if err != nil {
		return err
	}
	dstAbs, err := filepath.Abs(dst)
	if err != nil {
		return err
	}

	cmd := exec.Command("pipenv", "run", "python", "main.py", "unpack", srcAbs, dstAbs)
	cmd.Dir = "./arkwaifu-extractor"

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		t := scanner.Text()
		split := strings.SplitN(t, "=>", 2)
		if len(split) != 2 {
			continue
		}
		srcFile := filepath.ToSlash(filepath.Clean(split[0]))
		dstFile := filepath.ToSlash(filepath.Clean(split[1]))
		log.WithFields(log.Fields{
			"src": srcFile,
			"dst": dstFile,
		}).Infof("Resource unpacked.")
	}

	return nil
}
