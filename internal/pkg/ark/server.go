package ark

import (
	"github.com/pkg/errors"
	"strings"
)

// Server represents different server of the game. E.g., CN and EN.
type Server = string

const (
	// CnServer stands for China Server
	CnServer Server = "CN"

	// EnServer stands for English Server
	EnServer Server = "EN"

	// JpServer stands for Japan Server
	JpServer Server = "JP"

	// KrServer stands for Korea Server
	KrServer Server = "KR"

	// TwServer stands for Taiwan Server
	TwServer Server = "TW"
)

var (
	Servers = []Server{CnServer, EnServer, JpServer, KrServer, TwServer}
)

func ParseServer(s string) (Server, error) {
	EqualFolds := func(s string, others ...string) bool {
		for _, other := range others {
			if strings.EqualFold(s, other) {
				return true
			}
		}
		return false
	}
	if EqualFolds(s, "CN", "zh-CN") {
		return CnServer, nil
	}
	if EqualFolds(s, "EN", "en-US", "en-GB") {
		return EnServer, nil
	}
	if EqualFolds(s, "JP", "ja-JP") {
		return JpServer, nil
	}
	if EqualFolds(s, "KR", "ko-KR") {
		return KrServer, nil
	}
	if EqualFolds(s, "TW", "zh-TW") {
		return TwServer, nil
	}
	return "", errors.Errorf("unknown server: %v", s)
}

func MustParseServer(s string) Server {
	server, err := ParseServer(s)
	if err != nil {
		panic(err)
	}
	return server
}

func LanguageCodeUnderscore(s Server) string {
	switch s {
	case CnServer:
		return "zh_CN"
	case EnServer:
		return "en_US"
	case JpServer:
		return "ja_JP"
	case KrServer:
		return "ko_KR"
	case TwServer:
		return "zh_TW"
	}
	return ""
}
