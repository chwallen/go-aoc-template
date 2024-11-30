package main

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/chwallen/advent-of-code/internal"
	"github.com/chwallen/advent-of-code/internal/util"
)

type Day struct {
	Year  int
	Day   int
	Parts []int
}

type TemplateData struct {
	RootFolder string
	Days       []Day
	CurrentDay Day
}

func main() {
	rootFolder := util.GetModuleRootPath()

	td := TemplateData{RootFolder: rootFolder}
	latestEdit := time.Time{}
	baseDepth := strings.Count(rootFolder, string(filepath.Separator))
	yearRegex := regexp.MustCompile(`^20\d{2}$`)
	filepath.WalkDir(rootFolder, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		pathParts := strings.Split(path, string(filepath.Separator))
		if path == rootFolder || len(pathParts) == baseDepth+2 {
			return nil
		}

		yearDirectory := pathParts[baseDepth+1]
		if len(pathParts) == baseDepth+3 {
			if yearRegex.MatchString(yearDirectory) {
				return nil
			}
			return filepath.SkipDir
		}

		dayDirectory := pathParts[baseDepth+2]
		if !strings.HasPrefix(dayDirectory, "day") || len(pathParts) == baseDepth+5 {
			return filepath.SkipDir
		}
		if entry.Name() != fmt.Sprintf("%s.go", dayDirectory) {
			return nil
		}

		year, _ := strconv.Atoi(yearDirectory)
		dayNumber, _ := strconv.Atoi(dayDirectory[3:])
		day := Day{
			Year:  year,
			Day:   dayNumber,
			Parts: []int{1},
		}
		if dayNumber < 25 {
			day.Parts = append(day.Parts, 2)
		}

		td.Days = append(td.Days, day)

		t, err := getModTime(entry)
		if err != nil {
			log.Printf("Failed to get latest modification time for %s: %v\n", path, err)
		} else if t.After(latestEdit) {
			latestEdit = t
			td.CurrentDay = day
		}

		return nil
	})

	slices.SortFunc(td.Days, func(a, b Day) int {
		if a.Year-b.Year == 0 {
			return cmp.Compare(a.Day, b.Day)
		}
		return a.Year - b.Year
	})

	templateName := "day_parts.tmpl"
	outputPath := filepath.Join(rootFolder, "cmd", "run", "day_parts.go")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Failed to create %s: %v", outputPath, err)
	}
	defer outputFile.Close()

	err = internal.Templates.ExecuteTemplate(outputFile, templateName, td)
	if err != nil {
		log.Fatalf("Failed to execute template %s: %v", templateName, err)
	}
}

func getModTime(de os.DirEntry) (time.Time, error) {
	fi, err := de.Info()
	if err != nil {
		return time.Time{}, err
	}

	return fi.ModTime(), nil
}
