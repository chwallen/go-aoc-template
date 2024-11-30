package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/chwallen/advent-of-code/internal"
	"github.com/chwallen/advent-of-code/internal/aoc"
	"github.com/chwallen/advent-of-code/internal/util"
)

type TemplateData struct {
	Year    int
	DayName string
}

func main() {
	year, day, cookie, err := aoc.ParseYearDayCookieFlags()
	if err != nil {
		log.Fatal(err)
	}

	rootFolder := util.GetModuleRootPath()

	dayName := fmt.Sprintf("day%02d", day)
	dirPath := filepath.Join(rootFolder, strconv.Itoa(year), dayName)
	err = os.MkdirAll(dirPath, 0o755)
	if err != nil {
		log.Fatalf("Error creating destination directory %s: %v", dirPath, err)
	}

	templateData := TemplateData{Year: year, DayName: dayName}

	implFileName := filepath.Join(dirPath, fmt.Sprintf("%s.go", dayName))
	err = createFileFromTemplate("day.tmpl", implFileName, templateData)
	if err != nil {
		log.Fatalf("Failed to create %s: %v", implFileName, err)
	}

	testsFileName := filepath.Join(dirPath, fmt.Sprintf("%s_test.go", dayName))
	err = createFileFromTemplate("day_test.tmpl", testsFileName, templateData)
	if err != nil {
		log.Fatalf("Failed to create %s: %v", testsFileName, err)
	}

	err = aoc.DownloadDayInput(year, day, cookie)
	if err != nil {
		log.Fatalf("Failed to download input for day %d (%d): %v", day, year, err)
	}

	err = aoc.DownloadDayDescription(year, day, cookie)
	if err != nil {
		log.Fatalf("Failed to download description for day %d (%d): %v", day, year, err)
	}
}

func createFileFromTemplate(
	templateName string,
	destination string,
	templateData any,
) error {
	f, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			log.Printf("%s already exists, will not overwrite", destination)
			return nil
		}
		return fmt.Errorf("error creating file %s: %w", destination, err)
	}
	defer f.Close()

	err = internal.Templates.ExecuteTemplate(f, templateName, templateData)
	if err != nil {
		return fmt.Errorf("error executing template %s: %w", templateName, err)
	}

	log.Printf("Scaffolded %s\n", destination)
	return nil
}
