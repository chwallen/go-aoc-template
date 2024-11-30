// Package aoc provides helpers for commands with setup and to interact with the
// website, such as downloading input for a specific day.
package aoc

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	htmlToMarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/chwallen/advent-of-code/internal/util"
	"golang.org/x/net/html"
)

func ParseYearDayCookieFlags() (year, day int, cookie string, err error) {
	now := time.Now()
	flag.IntVar(&year, "year", now.Year(), "The year to target")
	flag.IntVar(&day, "day", now.Day(), "The day to target")
	flag.StringVar(&cookie, "cookie", os.Getenv("AOC_COOKIE"), "The cookie to use for fetching the input and description. Can also be provided as an env-variable.")
	flag.Parse()

	if year < 2015 {
		return 0, 0, "", fmt.Errorf("year must be 2015 or later, got %d", year)
	}
	if day < 1 || day > 25 {
		return 0, 0, "", fmt.Errorf("day must be 1 through 25, got %d", year)
	}
	if cookie == "" {
		flag.Usage()
		os.Exit(2)
	}

	return year, day, cookie, nil
}

func DownloadDayDescription(year, day int, cookie string) error {
	outputPath := createOutputPath(year, day, "description.md")
	bytes, err := downloadAocContentToFile(
		createAdventOfCodeURL(year, day),
		cookie,
		outputPath,
		parseHTMLToMarkdown,
	)

	if err == nil {
		log.Printf("Wrote %d bytes of Markdown to %s\n", bytes, outputPath)
	}
	return err
}

func DownloadDayInput(year, day int, cookie string) error {
	outputPath := createOutputPath(year, day, "input.txt")
	bytes, err := downloadAocContentToFile(
		createAdventOfCodeURL(year, day, "input"),
		cookie,
		outputPath,
		io.ReadAll,
	)

	if err == nil {
		log.Printf("Wrote %d bytes of input to %s\n", bytes, outputPath)
	}
	return err
}

func createAdventOfCodeURL(year, day int, extra ...string) string {
	args := make([]string, 3+len(extra))
	args[0], args[1], args[2] = strconv.Itoa(year), "day", strconv.Itoa(day)
	args = append(args, extra...)
	url, _ := url.JoinPath("https://adventofcode.com", args...)
	return url
}

func createOutputPath(year, day int, outputName string) string {
	return filepath.Join(
		util.GetModuleRootPath(),
		strconv.Itoa(year),
		fmt.Sprintf("day%02d", day),
		outputName,
	)
}

func downloadAocContentToFile(
	url string,
	cookie string,
	outputPath string,
	parser func(r io.Reader) ([]byte, error),
) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create GET request to %s: %w", url, err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: cookie,
	})

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to GET %s: %w", url, err)
	}

	body := res.Body
	defer body.Close()

	data, err := parser(body)
	if err != nil {
		return 0, err
	}

	err = os.WriteFile(outputPath, data, 0o644)
	if err != nil {
		return 0, fmt.Errorf("failed to write data to %s: %w", outputPath, err)
	}

	return len(data), nil
}

func parseHTMLToMarkdown(r io.Reader) ([]byte, error) {
	root, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	dayDescNodes := []*html.Node{}
	depthFirstSearchHTML(root, &dayDescNodes, isDayDescNode)

	var markdown []byte
	for _, node := range dayDescNodes {
		b, err := htmlToMarkdown.ConvertNode(node)
		if err != nil {
			return nil, fmt.Errorf("failed to parse HTML to Markdown: %w", err)
		}
		markdown = append(markdown, b...)
		markdown = append(markdown, '\n', '\n')
	}

	// Remove first byte to turn the ## heading to #, and remove the two trailing newlines.
	return markdown[1 : len(markdown)-2], nil
}

func isDayDescNode(node *html.Node) bool {
	for _, attr := range node.Attr {
		if attr.Key == "class" && attr.Val == "day-desc" {
			return true
		}
	}
	return false
}

func depthFirstSearchHTML(node *html.Node, matches *[]*html.Node, callback func(*html.Node) bool) {
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		if callback(n) {
			*matches = append(*matches, n)
		}
		depthFirstSearchHTML(n, matches, callback)
	}
}
