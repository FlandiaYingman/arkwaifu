package agdapi

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres/arkconsts"
	"github.com/google/go-github/v52/github"
	"github.com/pkg/errors"
	"regexp"
)

const (
	RepoOwner = "Kengxxiao"
	RepoName  = "ArknightsGameData"
)

var (
	githubClient = github.NewClient(nil)

	// RepoCommitMessageRegex is for matching that a commit is a resource-updating commit.
	//
	// Examples:
	// - [CN UPDATE] Client:2.0.01 Data:23-05-11-16-35-19-8a6fe7 [BOT TEST] - Match: CN, 2.0.01, 23-05-11-16-35-19-8a6fe7
	// - [EN UPDATE] Client:15.9.01 Data:23-04-25-10-10-55-972129 - Match: EN, 15.9.01, 23-04-25-10-10-55-972129
	RepoCommitMessageRegex = regexp.MustCompile("\\[(CN|EN|JP|KR) UPDATE\\] Client:([.\\d]*) Data:([-\\w]*)")

	// RepoFilePathRegex is for matching that a file in the repo is a part of the resources of Arknights.
	// Examples:
	// - .github/workflows/push.yml - Not Match
	// - README.md - Not Match
	// - en_US/gamedata/[uc]lua/GlobalConfig.lua - Match: GameServer=en_US, Name=[uc]lua/GlobalConfig.lua
	// - zh_CN/gamedata/[uc]lua/GlobalConfig.lua - Match: GameServer=zh_CN, Name=[uc]lua/GlobalConfig.lua
	RepoFilePathRegex = regexp.MustCompile("^(zh_CN|en_US|ja_JP|ko_KR)\\/(.*)$")
)

// GetLatestResourceVersion returns the latest resource version from the repository.
//
// Only the first 30 commits will be checked. If those commits do not match RepoCommitMessageRegex, an error will be returned.
func GetLatestResourceVersion(ctx context.Context) (*ResourceVersion, error) {
	commits, _, err := githubClient.Repositories.ListCommits(ctx, RepoOwner, RepoName, nil)
	if err != nil {
		return nil, err
	}

	for _, commit := range commits {
		result := RepoCommitMessageRegex.FindStringSubmatch(commit.GetCommit().GetMessage())
		if result == nil {
			continue
		}
		return &ResourceVersion{
			GameServer:      arkconsts.MustParseServer(result[1]),
			ClientVersion:   result[2],
			ResourceVersion: result[3],
			CommitSHA:       commit.GetSHA(),
		}, nil
	}

	return nil, errors.New("Cannot find a commit in the first 30 commits of the repo that matches the regex.")
}

// GetResourceVersion returns the resource version from the repository specified by resVersion.
//
// Only the first 30 commits will be checked. If those commits do not match RepoCommitMessageRegex and their resource version does not match, an error will be returned.
func GetResourceVersion(ctx context.Context, resVersion string) (ResourceVersion, error) {
	commits, _, err := githubClient.Repositories.ListCommits(ctx, RepoOwner, RepoName, nil)
	if err != nil {
		return ResourceVersion{}, err
	}

	for _, commit := range commits {
		result := RepoCommitMessageRegex.FindStringSubmatch(commit.GetCommit().GetMessage())
		if result == nil {
			continue
		}
		if result[3] != resVersion {
			continue
		}
		return ResourceVersion{
			GameServer:      arkconsts.MustParseServer(result[1]),
			ClientVersion:   result[2],
			ResourceVersion: result[3],
			CommitSHA:       commit.GetSHA(),
		}, nil
	}

	return ResourceVersion{}, errors.New("Cannot find a commit in the first 30 commits of the repo that matches the regex and matches the specified resVersion.")
}

// GetResourceInfoList returns a list of resource info for a specific resource version.
//
// Up to 100,000 entries of resource info are supported. If there are more than 100,000 entries, an error will be returned.
func GetResourceInfoList(ctx context.Context, resourceVersion ResourceVersion) ([]ResourceInfo, error) {
	tree, _, err := githubClient.Git.GetTree(ctx, RepoOwner, RepoName, resourceVersion.CommitSHA, true)
	if err != nil {
		return nil, err
	}
	if tree.GetTruncated() { // -> tree.IsTruncated()
		return nil, errors.New("Cannot fetch all resource info. There are more than 100,000 entries in the repo.")
	}

	infoList := make([]ResourceInfo, 0)
	for _, treeEntry := range tree.Entries {
		// Only files are considered.
		if treeEntry.GetType() != "blob" {
			continue
		}
		regexResult := RepoFilePathRegex.FindStringSubmatch(treeEntry.GetPath())
		// Only legal resources are considered.
		if regexResult == nil {
			continue
		}
		infoList = append(infoList, ResourceInfo{
			Name: regexResult[2],
		})
	}

	return infoList, nil
}
