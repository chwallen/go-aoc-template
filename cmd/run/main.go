package main

//go:generate go run ../gen

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"slices"
	"time"

	"github.com/chwallen/advent-of-code/internal/util"
)

type AdventOfCodeResult struct {
	Result   any
	Duration time.Duration
}

type RunnerInput struct {
	Func  func(lines []string, extras ...any) any
	Input string
	Year  int
	Day   int
	Part  int
}

var (
	runAll  bool
	year    int
	day     int
	part    int
	profile string
)

func main() {
	flag.BoolVar(&runAll, "all", false, "Run all days")
	flag.IntVar(&year, "year", 0, "The year to target")
	flag.IntVar(&day, "day", 0, "The day to target")
	flag.IntVar(&part, "part", 0, "The part to target")
	flag.StringVar(&profile, "profile", "", "Write CPU profile to file")
	flag.Parse()

	if profile != "" {
		f, err := os.Create(profile)
		if err != nil {
			log.Fatalf("Could not create CPU profile file %s: %v", profile, err)
		}
		_ = pprof.StartCPUProfile(f)
		defer f.Close()
	}

	if year == 0 && day == 0 && part == 0 {
		runCurrent()
	} else if year < 2015 {
		log.Fatal("Year must be 2015 or later")
	} else if runAll {
		runAllDaysForYear(year)
	} else {
		runDayPart(year, day, part)
	}

	pprof.StopCPUProfile()
}

func runAllDaysForYear(year int) {
	var overallDuration time.Duration

	days := make([][]RunnerInput, 25)

	for _, dayPart := range dayParts {
		if dayPart.Year == year {
			days[dayPart.Day-1] = append(days[dayPart.Day-1], dayPart)
		}
	}

	for i, day := range days {
		if len(day) == 0 {
			continue
		}
		fmt.Printf("(%d) Day %02d\n", year, i+1)
		for _, dayPart := range day {
			r := run(dayPart)
			result, duration := r.Result, r.Duration
			overallDuration += duration
			fmt.Printf(
				"part %d: %v (duration: %d.%03d ms)\n",
				dayPart.Part,
				result,
				duration.Milliseconds(),
				duration.Microseconds()%1000,
			)
		}
		fmt.Println()
	}

	var zeroDuration time.Duration
	if overallDuration == zeroDuration {
		log.Printf("No solutions found for year %d", year)
	} else {
		fmt.Printf(
			"(%d) Overall time elapsed: %d.%03d ms\n",
			year,
			overallDuration.Milliseconds(),
			overallDuration.Microseconds()%1000,
		)
	}
}

func runDayPart(year, day, part int) {
	partIndex, found := slices.BinarySearchFunc(dayParts, 0, func(e RunnerInput, t int) int {
		if e.Year-year != 0 {
			return e.Year - year
		}
		if e.Day-day != 0 {
			return e.Day - day
		}
		return e.Part - part
	})

	if !found {
		log.Fatalf("Did not find a solution for day %d part %d (%d)\n", day, part, year)
	}

	dayPart := dayParts[partIndex]
	fmt.Printf("Day %d part %d (%d)\n", day, part, year)
	r := run(dayPart)
	result, duration := r.Result, r.Duration
	fmt.Println(result)
	fmt.Printf(
		"Time elapsed: %d.%03d ms\n",
		duration.Milliseconds(),
		duration.Microseconds()%1000,
	)
}

func run(dayPart RunnerInput) AdventOfCodeResult {
	lines, err := util.ReadLines(dayPart.Input)
	if err != nil {
		log.Fatalf(
			"Could not run day %d part %d (%d): %v",
			dayPart.Day,
			dayPart.Part,
			dayPart.Year,
			err,
		)
	}

	start := time.Now()
	r := dayPart.Func(lines)
	duration := time.Since(start)

	return AdventOfCodeResult{Result: r, Duration: duration}
}
