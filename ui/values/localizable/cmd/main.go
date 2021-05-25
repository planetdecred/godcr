package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/planetdecred/godcr/ui/values/localizable"
)

var rex = regexp.MustCompile(`(?m)("(?:\\.|[^"\\])*")\s*=\s*("(?:\\.|[^"\\])*")`) // "key"="value"
const commentPrefix = "/"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Invalid arguments")
		return
	}

	readIntoMap := func(m map[string]string, localizableStrings string) {
		scanner := bufio.NewScanner(strings.NewReader(localizableStrings))
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, commentPrefix) {
				continue
			}

			matches := rex.FindAllStringSubmatch(line, -1)
			if len(matches) == 0 {
				continue
			}

			kv := matches[0]
			key := trimQuotes(kv[1])
			value := trimQuotes(kv[2])

			m[key] = value
		}
	}

	en := make(map[string]string)
	readIntoMap(en, localizable.EN)

	fileBuff, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	translation := make(map[string]string)
	readIntoMap(translation, string(fileBuff))

	var sb strings.Builder
	for k := range en {
		translationValue, ok := translation[k]
		if ok {
			sb.WriteString("\"" + k + "\" = \"")
			sb.WriteString(translationValue)
			sb.WriteString("\";\n")
		}
	}

	err = ioutil.WriteFile("translated.txt", []byte(sb.String()), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
	}
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}
