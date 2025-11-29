package main

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
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

	yearDirPaths, e := filepath.Glob(filepath.Join(rootFolder, "20[1-9][0-9]"))
	if e != nil {
		log.Fatalf("Failed to glob year directories in %s: %v", rootFolder, e)
	}

	td := TemplateData{RootFolder: filepath.ToSlash(rootFolder)}
	latestEdit := time.Time{}

	for _, yearDirPath := range yearDirPaths {
		year, _ := strconv.Atoi(yearDirPath[len(yearDirPath)-4:])

		dayDirPaths, _ := filepath.Glob(filepath.Join(yearDirPath, "day[0-2][0-9]"))
		for _, dayDirPath := range dayDirPaths {
			dayDir := dayDirPath[len(dayDirPath)-5:]
			dayFilePath := filepath.Join(dayDirPath, fmt.Sprintf("%s.go", dayDir))
			dayFileInfo, err := os.Stat(dayFilePath)
			if err != nil {
				log.Printf("Could not stat file %s, ignoring %d/%s", dayFilePath, year, dayDir)
				continue
			}

			dayNumber, _ := strconv.Atoi(dayDir[3:])
			day := Day{
				Year:  year,
				Day:   dayNumber,
				Parts: []int{1, 2},
			}
			if (year <= 2024 && dayNumber == 25) || (year >= 2025 && dayNumber == 12) {
				day.Parts = day.Parts[:1]
			}

			td.Days = append(td.Days, day)

			modTime := dayFileInfo.ModTime()
			if modTime.After(latestEdit) {
				latestEdit = modTime
				td.CurrentDay = day
			}
		}
	}

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
