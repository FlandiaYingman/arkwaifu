package data

import (
	"context"
	"github.com/cavaliergopher/grab/v3"
	"github.com/google/go-github/v42/github"
	"github.com/mholt/archiver/v4"
	"github.com/pkg/errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func getZipballLink(owner string, repo string, ref string) (url.URL, error) {
	ctx := context.Background()
	client := github.NewClient(nil)

	opts := github.RepositoryContentGetOptions{Ref: ref}
	link, _, err := client.Repositories.GetArchiveLink(ctx, owner, repo, github.Zipball, &opts, true)
	return *link, err
}

func downloadZipball(link url.URL, dest string) (string, error) {
	err := os.MkdirAll(dest, 0755)
	if err != nil {
		return "", err
	}
	resp, err := grab.Get(dest, link.String())
	return resp.Filename, err
}

func extractZipball(zipball string, root string, files []string, dest string) error {
	ctx := context.Background()
	zip := archiver.Zip{}

	file, err := os.Open(zipball)
	if err != nil {
		return err
	}
	defer file.Close()

	return zip.Extract(ctx, file, files, func(ctx context.Context, f archiver.File) error {
		diskFilePath := filepath.Join(dest, strings.TrimPrefix(f.NameInArchive, root))
		if f.IsDir() {
			return os.MkdirAll(diskFilePath, os.ModePerm)
		}

		diskFile, err := os.Create(diskFilePath)
		if err != nil {
			return err
		}
		defer diskFile.Close()

		zipEntryFile, err := f.Open()
		if err != nil {
			return err
		}
		defer zipEntryFile.Close()

		_, err = io.Copy(diskFile, zipEntryFile)
		return err
	})
}

func findCommitByResVersion(owner string, repo string, resVersion string) (string, error) {
	client := github.NewClient(nil)
	page := 0
	perPage := 100
	until := time.Now()
	for {
		page++
		commits, _, err := client.Repositories.ListCommits(
			context.Background(),
			owner,
			repo,
			&github.CommitsListOptions{
				Path:  "/zh_CN",
				Until: until,
				ListOptions: github.ListOptions{
					Page:    page,
					PerPage: perPage,
				},
			},
		)
		if err != nil {
			return "", err
		}
		if len(commits) == 0 {
			break
		}
		for _, c := range commits {
			message := c.GetCommit().GetMessage()
			if strings.Contains(message, "CN UPDATE") && strings.Contains(message, resVersion) {
				return c.GetSHA(), nil
			}
		}
	}
	return "", nil
}

func findLatestCommit(owner string, repo string) (string, string, error) {
	client := github.NewClient(nil)
	commits, _, err := client.Repositories.ListCommits(
		context.Background(),
		owner,
		repo,
		&github.CommitsListOptions{
			Path: "/zh_CN",
			ListOptions: github.ListOptions{
				Page:    1,
				PerPage: 1,
			},
		},
	)
	if err != nil {
		return "", "", err
	}
	if len(commits) != 1 {
		return "", "", errors.Errorf("unexpected len(commits) != 1; commits: %v", commits)
	}

	commit := commits[0]
	return commit.GetCommit().GetMessage(), commit.GetSHA(), nil
}
