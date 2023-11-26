package wal

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type fInfo struct {
	fName string
	fNum  int
}

func (w *WAL) GetNamesWALFiles(withOrderByNumberASC bool) ([]string, error) {
	dirWithWalFiles, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	allWALFiles := make([]fInfo, 0)
	for _, f := range dirWithWalFiles {
		if f.IsDir() {
			continue
		}

		var n int
		if n, err = getWALFileNumber(f.Name()); err == nil {
			allWALFiles = append(allWALFiles, fInfo{
				fName: f.Name(),
				fNum:  n,
			})
		}
	}

	if withOrderByNumberASC {
		sort.Slice(allWALFiles, func(i, j int) bool {
			return allWALFiles[i].fNum < allWALFiles[j].fNum
		})
	}

	result := make([]string, 0, len(allWALFiles))
	for _, f := range allWALFiles {
		result = append(result, f.fName)
	}

	return result, nil
}

func findLastWALFile() (*os.File, error) {
	dirWithWalFiles, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var (
		maxFileNumber     int
		fileWithMaxNumber os.DirEntry
	)
	for _, f := range dirWithWalFiles {
		if f.IsDir() {
			continue
		}

		var n int
		if n, err = getWALFileNumber(f.Name()); err == nil {
			if fileWithMaxNumber == nil || n > maxFileNumber {
				fileWithMaxNumber = f
				maxFileNumber = n
			}
		}
	}

	var lastWALFile *os.File
	if fileWithMaxNumber == nil {
		lastWALFile, err = prepareNewWALFile(generateNewWALFileName(0))
	} else {
		lastWALFile, err = os.OpenFile(fileWithMaxNumber.Name(), os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_SYNC, os.ModePerm)
	}

	if err != nil {
		return nil, err
	}

	return lastWALFile, nil
}

func prepareNewWALFile(currWALFileName string) (*os.File, error) {
	number, err := getWALFileNumber(currWALFileName)
	if err != nil {
		return nil, fmt.Errorf("getting current wal file suffix: %w", err)
	}

	newWALFileNumber := number + 1

	newWALFileName := generateNewWALFileName(newWALFileNumber)

	newWALFile, err := os.OpenFile(newWALFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_SYNC, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("preparing new wal file with name '%s': %w", newWALFileName, err)
	}

	return newWALFile, nil
}

func getWALFileNumber(walFileName string) (int, error) {
	nameWithoutExtension, ok := strings.CutSuffix(walFileName, ".wal")
	if !ok {
		return -1, errors.New("it's not .wal file")
	}

	parts := strings.Split(nameWithoutExtension, symbolBetweenStartAndSuffixNamesWALFiles)
	if len(parts) != 2 {
		return -1, fmt.Errorf("invalid file name %s", walFileName)
	}

	number, err := strconv.Atoi(parts[1])
	if err != nil {
		return -1, fmt.Errorf("converting suffix to number: %w", err)
	}

	return number, nil
}

func generateNewWALFileName(newNumber int) string {
	return fmt.Sprintf("wal%s%d.wal", symbolBetweenStartAndSuffixNamesWALFiles, newNumber)
}
