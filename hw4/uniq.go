package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	modePlain  = "plain"
	modeCount  = "count"
	modeDups   = "dups"
	modeUnique = "unique"
)

type Config struct {
	mode       string
	ignoreCase bool
	skipFields int
	skipChars  int
	inputPath  string
	outputPath string
}

func main() {
	cfg := parseConfig()

	in := os.Stdin
	var inFile *os.File
	var err error
	if cfg.inputPath != "" {
		inFile, err = os.Open(cfg.inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка открытия входного файла: %v\n", err)
			os.Exit(1)
		}
		defer inFile.Close()
		in = inFile
	}

	out := os.Stdout
	var outFile *os.File
	if cfg.outputPath != "" {
		outFile, err = os.Create(cfg.outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка создания выходного файла: %v\n", err)
			os.Exit(1)
		}
		defer outFile.Close()
		out = outFile
	}

	if err := run(in, out, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
}

func parseConfig() Config {
	var cfg Config

	cFlag := flag.Bool("c", false, "подсчитать количество повторов и вывести перед строкой")
	dFlag := flag.Bool("d", false, "вывести только повторяющиеся строки")
	uFlag := flag.Bool("u", false, "вывести только уникальные строки")
	iFlag := flag.Bool("i", false, "игнорировать регистр")
	fFlag := flag.Int("f", 0, "не учитывать первые num полей")
	sFlag := flag.Int("s", 0, "не учитывать первые chars символов")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Использование: uniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Проверка взаимной исключаемости -c, -d, -u
	modes := 0
	if *cFlag {
		modes++
	}
	if *dFlag {
		modes++
	}
	if *uFlag {
		modes++
	}
	if modes > 1 {
		fmt.Fprintln(os.Stderr, "Параметры -c, -d и -u нельзя использовать одновременно.")
		flag.Usage()
		os.Exit(2)
	}

	cfg.mode = modePlain
	if *cFlag {
		cfg.mode = modeCount
	} else if *dFlag {
		cfg.mode = modeDups
	} else if *uFlag {
		cfg.mode = modeUnique
	}

	cfg.ignoreCase = *iFlag
	if *fFlag > 0 {
		cfg.skipFields = *fFlag
	}
	if *sFlag > 0 {
		cfg.skipChars = *sFlag
	}

	args := flag.Args()
	switch len(args) {
	case 0:
	// stdin/stdout
	case 1:
		cfg.inputPath = args[0]
	case 2:
		cfg.inputPath = args[0]
		cfg.outputPath = args[1]
	default:
		fmt.Fprintln(os.Stderr, "Слишком много позиционных аргументов.")
		flag.Usage()
		os.Exit(2)
	}

	return cfg
}

func run(in io.Reader, out io.Writer, cfg Config) error {
	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	var prevLine string
	var prevKey string
	count := 0

	for scanner.Scan() {
		line := scanner.Text()
		key := buildKey(line, cfg)

		if count == 0 {
			prevLine = line
			prevKey = key
			count = 1
			continue
		}

		if key == prevKey {
			count++
		} else {
			if shouldPrint(count, cfg.mode) {
				printGroup(writer, prevLine, count, cfg)
			}
			prevLine = line
			prevKey = key
			count = 1
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if count > 0 && shouldPrint(count, cfg.mode) {
		printGroup(writer, prevLine, count, cfg)
	}

	return nil
}

func shouldPrint(count int, mode string) bool {
	switch mode {
	case modeDups:
		return count > 1
	case modeUnique:
		return count == 1
	default: // plain или count
		return true
	}
}

func printGroup(w io.Writer, line string, count int, cfg Config) {
	if cfg.mode == modeCount {
		fmt.Fprintf(w, "%d %s\n", count, line)
	} else {
		fmt.Fprintln(w, line)
	}
}

func buildKey(line string, cfg Config) string {
	start := 0
	if cfg.skipFields > 0 {
		start = skipFieldsIndex(line, cfg.skipFields)
	}
	if cfg.skipChars > 0 {
		start = skipCharsIndex(line, start, cfg.skipChars)
	}
	if start > len(line) {
		start = len(line)
	}
	key := line[start:]
	if cfg.ignoreCase {
		key = strings.ToLower(key)
	}
	return key
}

func skipFieldsIndex(s string, num int) int {
	if num <= 0 {
		return 0
	}

	i := 0
	n := len(s)
	for f := 0; f < num && i < n; f++ {
		for i < n && s[i] == ' ' {
			i++
		}
		for i < n && s[i] != ' ' {
			i++
		}
	}
	for i < n && s[i] == ' ' {
		i++
	}

	return i
}

func skipCharsIndex(s string, start int, num int) int {
	if num <= 0 {
		return start
	}
	i := start + num
	if i > len(s) {
		i = len(s)
	}
	return i
}
