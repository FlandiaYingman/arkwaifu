package arkconsts

import (
	"fmt"
	"github.com/pkg/errors"
)

type GameServer string

const (
	CN GameServer = "CN"
	EN            = "EN"
	JP            = "JP"
	KR            = "KR"
)

func ParseServer(serverString string) (GameServer, error) {
	switch serverString {
	case "CN", "zh_CN":
		return CN, nil
	case "EN", "en_US":
		return EN, nil
	case "JP", "ja_JP":
		return JP, nil
	case "KR", "ko_KR":
		return KR, nil
	}
	return "", errors.New(fmt.Sprintf("Unrecognized server string: %s", serverString))
}

func MustParseServer(serverString string) GameServer {
	server, err := ParseServer(serverString)
	if err != nil {
		panic(err)
	}
	return server
}

func (gs GameServer) GetLanguageCode() string {
	switch gs {
	case CN:
		return "zh_CN"
	case EN:
		return "en_US"
	case JP:
		return "ja_JP"
	case KR:
		return "ko_KR"
	}
	panic(fmt.Sprintf("Unrecognized game server %s", gs))
}
