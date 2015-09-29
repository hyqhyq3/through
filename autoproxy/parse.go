package autoproxy

import (
	"bytes"
	"regexp"
	"strings"
)

func ParseList(l string) (except, rules []Rule, _ error) {
	buf := bytes.NewBuffer([]byte(l))
	except = make([]Rule, 0)
	rules = make([]Rule, 0)
	for {
		var rule Rule
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			break
		}
		if line == "" || line[:1] == "!" {
			continue
		}

		switch {
		case line[:2] == "||":
			rule = &DomainRule{line[2:]}
		case line[:4] == "@@||":
			rule = &DomainRule{line[4:]}
			except = append(except, rule)
			continue
		case line[:1] == "|":
			rule = &PrefixRule{line[1:]}
		case line[:1] == "/" && line[len(line)-1:] == "/":
			expr := line[1 : len(line)-1]
			re, err := regexp.Compile(expr)
			if err != nil {
				return nil, nil, err
			}
			rule = &RegexRule{re}
		default:
			rule = &KeywordRule{line}
		}
		rules = append(rules, rule)
	}

	return
}
