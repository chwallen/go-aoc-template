# Go template for Advent of Code

<!--toc:start-->
- [Go template for Advent of Code](#go-template-for-advent-of-code)
  - [Instructions](#instructions)
    - [Scaffold the first day](#scaffold-the-first-day)
    - [Testing](#testing)
    - [Running](#running)
<!--toc:end-->

This is a template for writing Go solutions to [Advent of Code](https://adventofcode.com).
Each year is kept in a separate folder, where days are grouped into folders.
Things that are shared between days are placed inside the _internal_ directory.

## Instructions

Start by replacing all occurrences of _chwallen/advent-of-code_ with your own
module name.

The Makefile includes commands to scaffold a day. Each command requires the
variable `AOC_COOKIE` to be set either as an environment variable (recommended)
or passed when running a command. The `YEAR` and `DAY` variables are optional
and defaults to the current year and day.

- `make new` scaffolds a new day, with the current year and day used as defaults.
- `make input` download the input file for a day.
- `make description` download the description for a day as Markdown.

### Scaffold the first day

The following section gives an example of how you would create the first day of
the advent of code challenge.

Begin by retrieving your session cookie used for the AOC website. You can do so
by following these steps:

1. Log in to Advent of Code in your browser with your preferred method
2. Open the developer panel
3. Go to the _Network_ tab and (optionally) filter by _Doc_
4. Navigate to, or refresh, the start page
5. Click the document load request
6. Copy everything after after _session=_ from the _Cookie_ header
7. Set the value as an environment variable via `export AOC_COOKIE=<value>`
  (bash/zsh/etc.), or `set -x AOC_COOKIE <value>` (fish)

Now, scaffold day 1 by running `make new DAY=1`. This will create the folder
`<year>/day01` with the files _day01.go_, _day01_test.go_, _input.txt_, and
_description.md_ (change the year by setting `YEAR=<other-year>`). The file
_day01.go_ will receive the input as `lines []string`.

### Testing

To test your solution against the example input, the `exampleInput` in
_dayXX_test.go_ has to be filled along with the expected value in the tests
slice. Run `make test` to test _all_ your solutions. Go caches the results for
subsequent runs if the code is unchanged. In case you want to test a
specific year (or even day), set the `YEAR` and/or `DAY` variables to limit the
test pattern.

Each day part accepts extra arguments as vararg, and the `Test` struct has a
field for them. This is useful if an example uses a reduced problem space as it
allows your solution to dynamically change its parameters.

### Running

Use `make run` whenever you want to run the day you are currently working on.
Once you have solved the first part, run `make description` to get the
description for the second part.

You can run a specific day part by setting `YEAR`, `DAY`, and `PART` when
running `make run`. It's also possible to run all days for a specific year by
running `make runall YEAR=<year>`.
