package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
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
		os.Exit(1)
	}

	for _, l := range extractImagesFromManifest(manifest) {
		fmt.Printf("%s\n", l)
	}

	// image/tag:v1.0.0
	// 123.123.123.123:123/image/tag:v1.0.0
	// your-domain.com/image/tag
	// your-domain.com/image/tag:v1.1.1-patch1
	// image/tag
	// image
	// image:v1.1.1-patch
	// ubuntu@sha256:45b23dee08af5e43a7fea6c4cf9c25ccf269ee113168c19722f87876677c5cb2
	// etc...

	tests := []string{
		"image/tag:v1.0.0",
		"123.123.123.123/image/tag:v1.0.0",
		"your-domain.com/image/tag",
		"your-domain.com/image/tag:v1.1.1-patch1",
		"image/tag",
		"image",
		"image:v1.1.1-patch",
		"ubuntu@sha256:45b23dee08af5e43a7fea6c4cf9c25ccf269ee113168c19722f87876677c5cb2",
	}

	for _, t := range tests {
		fmt.Printf("\n%s => %s", t, changeContainerRegistry(t, "some.domain.com"))
	}

	os.Exit(0)
}

func changeContainerRegistry(imageTag string, newContainerRegistry string) string {
	if strings.Contains(imageTag, "/") {
		potentialRegistry := before(imageTag, "/")
		if net.ParseIP(potentialRegistry) != nil {
			return fmt.Sprintf("%s/%s", newContainerRegistry, after(imageTag, "/"))
		}
	}
	return fmt.Sprintf("%s/%s", newContainerRegistry, imageTag)
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

func before(value string, a string) string {
	// Get substring before a string.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	return value[0:pos]
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
