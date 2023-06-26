package hgapi

import (
	"bufio"
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"

	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	extractorLocation = "./tools/extractor"

	maxWorkers       *int = nil
	maxTasksPerChild *int = nil
)

// TODO: Change to a direct call of Python API.
func unpack(ctx context.Context, src string, dst string) error {
	err := preCheck()
	if err != nil {
		return err
	}

	srcAbs, err := filepath.Abs(src)
	if err != nil {
		return errors.WithStack(err)
	}
	dstAbs, err := filepath.Abs(dst)
	if err != nil {
		return errors.WithStack(err)
	}

	args := []string{"-u", "main.py", "unpack", srcAbs, dstAbs}
	if maxWorkers != nil {
		args = append(args, "-w", strconv.Itoa(*maxWorkers))
	}
	if maxTasksPerChild != nil {
		args = append(args, "-t", strconv.Itoa(*maxTasksPerChild))
	}
	cmd := exec.CommandContext(ctx, "python", args...)
	cmd.Dir = extractorLocation

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdout)
	go func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			t := scanner.Text()
			// TODO: Info -> Debug
			log.Info().
				Str("output", t).
				Msg("Output from stdout of the extractor... ")
			split := strings.SplitN(t, "=>", 2)
			if len(split) != 2 {
				continue
			}
			srcFile := filepath.ToSlash(filepath.Clean(split[0]))
			dstFile := filepath.ToSlash(filepath.Clean(split[1]))
			log.Info().
				Str("src", srcFile).
				Str("dst", dstFile).
				Msg("Unpacked the resource src to dst. ")
		}
	}(scanner)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	errScanner := bufio.NewScanner(stderr)
	go func(errScanner *bufio.Scanner) {
		for errScanner.Scan() {
			t := errScanner.Text()
			log.Warn().
				Str("output", t).
				Msg("Output from stderr of the extractor... ")
		}
	}(errScanner)

	err = cmd.Start()
	if err != nil {
		return errors.WithStack(err)
	}

	err = cmd.Wait()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func unpackPreCheck() error {
	// check extractor existence
	_, err := os.Stat(extractorLocation)
	if err != nil {
		return errors.Wrap(err, "extractor doesn't exist")
	}
	return nil
}
