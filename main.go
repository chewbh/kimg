package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	manifest := []string{}

	if fi.Mode()&os.ModeNamedPipe != 0 {
		// input passed by pipe
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			manifest = append(manifest, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	} else {
		fmt.Println("no pipe :(")
	}

	for _, l := range extractImagesFromManifest(manifest) {
		fmt.Printf("%s\n", l)
	}

	os.Exit(0)
}

func extractImagesFromManifest(manifest []string) []string {
	imageLines := filter(manifest, func(rawline string) bool {
		// detect line in manifest that potentially contains image tag
		line := strings.TrimSpace(rawline)
		return strings.HasPrefix(line, "image:") || strings.Contains(line, ".image:")
	})

	images := mapTo(imageLines, func(rawLine string) string {
		return strings.TrimSpace(after((rawLine), "image:"))
	})
	return images
}

func filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func mapTo(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func after(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}
