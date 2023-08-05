package arkdata

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/fileutil"
	"github.com/google/go-github/v52/github"
	"github.com/pkg/errors"
	"os"
	"regexp"
)

const (
	RepoOwner          = "Kengxxiao"
	RepoName           = "ArknightsGameData"
	CompositeRepoOwner = "FlandiaYingman"
	CompositeRepoName  = "ArknightsGameDataComposite"
	DefaultPrefix      = "assets/torappu/dynamicassets"
)

var (
	githubClient = github.NewClient(nil)

	// RepoCommitMessageRegex is for matching that a commit is a resource-updating commit.
	//
	// Examples:
	//   - [CN UPDATE] Client:2.0.01 Data:23-05-11-16-35-19-8a6fe7 [BOT TEST] - Match: CN, 2.0.01, 23-05-11-16-35-19-8a6fe7
	//   - [EN UPDATE] Client:15.9.01 Data:23-04-25-10-10-55-972129 - Match: EN, 15.9.01, 23-04-25-10-10-55-972129
	RepoCommitMessageRegex = regexp.MustCompile("\\[(CN|EN|JP|KR) UPDATE] Client:([.\\d]*) Data:([-\\w]*)")
)

type DataInfo struct {
	Name string
}

type DataVersion struct {
	GameServer      ark.Server
	ClientVersion   string
	ResourceVersion string
	CommitSHA       string
}

func GetGameData(ctx context.Context, server ark.Server, version ark.Version, patterns []string, dst string) error {
	return getGameData(ctx, server, version, RepoOwner, RepoName, patterns, dst)
}
func GetLatestDataVersion(ctx context.Context, server ark.Server) (*DataVersion, error) {
	return getLatestDataVersion(ctx, server, RepoOwner, RepoName)
}
func GetDataVersion(ctx context.Context, server ark.Server, resVer string) (*DataVersion, error) {
	return getDataVersion(ctx, server, resVer, RepoOwner, RepoName)
}

func GetCompositeGameData(ctx context.Context, server ark.Server, version ark.Version, patterns []string, dst string) error {
	return getGameData(ctx, server, version, CompositeRepoOwner, CompositeRepoName, patterns, dst)
}
func GetLatestCompositeDataVersion(ctx context.Context, server ark.Server) (*DataVersion, error) {
	return getLatestDataVersion(ctx, server, CompositeRepoOwner, CompositeRepoName)
}
func GetCompositeDataVersion(ctx context.Context, server ark.Server, resVer string) (*DataVersion, error) {
	return getDataVersion(ctx, server, resVer, CompositeRepoOwner, CompositeRepoName)
}

func getGameData(ctx context.Context, server ark.Server, version ark.Version, repoOwner string, repoName string, patterns []string, dst string) error {
	var sha string
	if version != "" {
		dataVersion, err := getDataVersion(ctx, server, version, repoOwner, repoName)
		if err != nil {
			return errors.WithStack(err)
		}
		sha = dataVersion.CommitSHA
	} else {
		dataVersion, err := getLatestDataVersion(ctx, server, repoOwner, repoName)
		if err != nil {
			return errors.WithStack(err)
		}
		sha = dataVersion.CommitSHA
	}

	zipball, err := download(ctx, repoOwner, repoName, sha)
	if err != nil {
		return errors.WithStack(err)
	}
	defer os.RemoveAll(zipball)

	data, err := unzip(ctx, zipball, patterns, server)
	if err != nil {
		return errors.WithStack(err)
	}
	defer os.RemoveAll(data)

	err = fileutil.MoveAllContent(data, dst)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
func getLatestDataVersion(ctx context.Context, server ark.Server, repoOwner string, repoName string) (*DataVersion, error) {
	commit, err := findCommit(ctx, repoOwner, repoName, func(commit *github.RepositoryCommit) bool {
		version := parseCommit(commit)
		return version != nil && version.GameServer == server
	})
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find latest DataVersion")
	}
	return parseCommit(commit), nil
}
func getDataVersion(ctx context.Context, server ark.Server, resVer string, repoOwner string, repoName string) (*DataVersion, error) {
	commit, err := findCommit(ctx, repoOwner, repoName, func(commit *github.RepositoryCommit) bool {
		version := parseCommit(commit)
		return version != nil && version.GameServer == server && version.ResourceVersion == resVer
	})
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find DataVersion by %s", resVer)
	}
	if commit == nil {
		return nil, errors.Errorf("no DataVersion can be found by %s", resVer)
	}
	return parseCommit(commit), nil
}

func parseCommit(commit *github.RepositoryCommit) *DataVersion {
	matches := RepoCommitMessageRegex.FindStringSubmatch(commit.GetCommit().GetMessage())
	if matches == nil {
		return nil
	} else {
		return &DataVersion{
			GameServer:      ark.MustParseServer(matches[1]),
			ClientVersion:   matches[2],
			ResourceVersion: matches[3],
			CommitSHA:       commit.GetSHA(),
		}
	}
}
func findCommit(ctx context.Context, owner, repo string, predicate func(*github.RepositoryCommit) bool) (*github.RepositoryCommit, error) {
	const perPage = 100
	var page = 1
	var opts = &github.CommitsListOptions{}
	for {
		opts.PerPage = perPage
		opts.Page = page
		page++

		commits, _, err := githubClient.Repositories.ListCommits(ctx, owner, repo, opts)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		for _, commit := range commits {
			if predicate(commit) {
				return commit, nil
			}
		}

		if len(commits) < perPage {
			break
		}
	}
	return nil, nil
}
