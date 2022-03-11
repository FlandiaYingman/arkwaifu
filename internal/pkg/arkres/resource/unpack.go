package resource

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	extractorLocation = "./tools/extractor"
)

func unpackResources(src string, dst string) error {
	// TODO: Change to a direct call of Python API.

	// check extractor existence
	_, err := os.Stat(extractorLocation)
	if err != nil {
		return err
	}

	srcAbs, err := filepath.Abs(src)
	if err != nil {
		return err
	}
	dstAbs, err := filepath.Abs(dst)
	if err != nil {
		return err
	}

	cmd := exec.Command("python", "-u", "main.py", "unpack", srcAbs, dstAbs)
	cmd.Dir = extractorLocation
	// cmd.Stdout = os.Stdout

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdout)
	go func(scanner *bufio.Scanner) {
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
	}(scanner)

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
