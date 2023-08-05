package arkassets

import (
	"bufio"
	"context"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/fileutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	extractorLocation = "./tools/extractor"
)

func unpack(ctx context.Context, src string) (string, error) {
	tempDir, err := os.MkdirTemp("", "arkassets_unpack-*")
	if err != nil {
		return "", err
	}

	err = findExtractor()
	if err != nil {
		return "", err
	}

	srcAbs, err := filepath.Abs(src)
	if err != nil {
		return "", errors.WithStack(err)
	}
	dstAbs, err := filepath.Abs(tempDir)
	if err != nil {
		return "", errors.WithStack(err)
	}

	args := []string{"-u", "main.py", "unpack", srcAbs, dstAbs}
	cmd := exec.CommandContext(ctx, "python", args...)
	cmd.Dir = extractorLocation

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(stdout)
	go func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			t := scanner.Text()
			log.Debug().
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
		return "", err
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
		return "", errors.WithStack(err)
	}

	err = cmd.Wait()
	if err != nil {
		return "", errors.WithStack(err)
	}

	return tempDir, nil
}

func findExtractor() error {
	exists, err := fileutil.Exists(extractorLocation)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("cannot find extractor")
	}
	return nil
}
