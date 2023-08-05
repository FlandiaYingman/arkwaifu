package arkparser

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

type Directive struct {
	Name   string
	Params map[string]string
}

var DirectiveRegex = regexp.MustCompile(`\[(\w*)(?:\((.*)\))?]|\[name="(.*)"]`)

func ParseDirectives(raw string) []Directive {
	directives := make([]Directive, 0)
	matches := DirectiveRegex.FindAllStringSubmatch(raw, -1)
	for _, match := range matches {
		if match[3] == "" {
			name := strings.ToLower(match[1])
			params := ParseDirectiveParams(match[2])
			directives = append(directives, Directive{
				Name:   name,
				Params: params,
			})
		} else {
			name := ""
			params := make(map[string]string)
			params["name"] = match[3]
			directives = append(directives, Directive{
				Name:   name,
				Params: params,
			})
		}
	}
	return directives
}

func ParseDirectiveParams(rawParams string) map[string]string {
	params := make(map[string]string)
	for _, rawParam := range splitTokens(rawParams) {
		key, value, _ := strings.Cut(rawParam, "=")
		key = strings.ToLower(strings.TrimSpace(key))
		value = strings.Trim(strings.TrimSpace(value), `"`)
		params[key] = value
	}
	return params
}

// splitTokens splits a string separated by comma which is not in double quotes.
func splitTokens(s string) []string {
	if len(s) <= 0 {
		return nil
	}
	tokens := make([]string, 0)
	pos := 0
	inQuotes := false
	for i, char := range s {
		if char == '"' {
			inQuotes = !inQuotes
		} else if char == ',' && !inQuotes {
			tokens = append(tokens, s[pos:i])
			pos = i + 1
		}
	}
	lastToken := s[pos:]
	if lastToken == "," {
		tokens = append(tokens, "")
	} else {
		tokens = append(tokens, lastToken)
	}
	return tokens
}

func (p *Parser) ParseStoryText(storyTextPath string) ([]Directive, error) {
	storyTextPath = fmt.Sprintf("%s.txt", storyTextPath)
	storyTextPath = path.Join(p.Root, p.Prefix, "gamedata/story", storyTextPath)
	storyTextData, err := os.ReadFile(storyTextPath)
	if err != nil {
		return nil, err
	}

	return ParseDirectives(string(storyTextData)), nil
}
