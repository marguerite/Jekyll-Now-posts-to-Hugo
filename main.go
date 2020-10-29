package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func main() {
	var input, output string
	flag.StringVar(&input, "input", "posts", "input posts dir")
	flag.StringVar(&output, "output", "../hugo/content/posts", "output posts dir")
	flag.Parse()

	if _, err := os.Stat(output); os.IsNotExist(err) {
		os.MkdirAll(output, 0755)
	}

	f, err := os.Open(input)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	files, err := f.Readdir(-1)
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile("^(\\d+-\\d+-\\d+)-(.*)\\.md$")
	for _, v := range files {
		if !strings.HasSuffix(v.Name(), ".md") {
			continue
		}
		name := v.Name()
		if !re.MatchString(name) {
			continue
		}
		m := re.FindStringSubmatch(name)
		title := strings.Replace(m[2], "-", " ", -1)
		t, err := time.Parse("2006-1-2", m[1])
		if err != nil {
			panic(err)
		}
		day := t.Format("2006-01-02T15:04:05+08:00")

		str := "---\n"
		str += "title: \"" + title + "\"\n"
		str += "date: " + day + "\n"
		str += "draft: false\n"
		str += "---\n"

		rd, err := os.Open(filepath.Join(input, name))
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(rd)
		for scanner.Scan() {
			str += scanner.Text() + "\n"
		}
		rd.Close()
		ioutil.WriteFile(filepath.Join(output, strings.ReplaceAll(title, " ", "_")+".md"), []byte(str), 0644)

	}
}
