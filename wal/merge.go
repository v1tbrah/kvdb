package wal

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/v1tbrah/kvdb/dbengine"
)

const mergeCheckInterval = time.Hour

func startProcessMergeWALFiles() {
	if err := mergeWALFiles(); err != nil {
		slog.Error("mergeWALFiles", err)
	}

	for range time.Tick(mergeCheckInterval) {
		if err := mergeWALFiles(); err != nil {
			slog.Error("mergeWALFiles", err)
		}
	}
}

func mergeWALFiles() error {
	dirWithWalFiles, err := os.ReadDir(".")
	if err != nil {
		return err
	}

	// max heap by fNum (for contain 2 files with min numbers)
	twoFirstFiles := make([]fInfo, 0)
	var hasMoreThenTwoWALFiles bool

	for _, f := range dirWithWalFiles {
		if f.IsDir() {
			continue
		}

		var n int
		if n, err = getWALFileNumber(f.Name()); err == nil {
			fI := fInfo{fName: f.Name(), fNum: n}
			if len(twoFirstFiles) < 2 {
				twoFirstFiles = append(twoFirstFiles, fI)
				heapify(twoFirstFiles, 0)
				continue
			}

			hasMoreThenTwoWALFiles = true

			if n < twoFirstFiles[0].fNum {
				twoFirstFiles[0] = fI
				heapify(twoFirstFiles, 0)
			}
		}
	}

	if hasMoreThenTwoWALFiles {
		return mergeTwoWALFiles(twoFirstFiles[1], twoFirstFiles[0])
	}

	return nil
}

func mergeTwoWALFiles(file1, file2 fInfo) error {
	f1, err := os.Open(file1.fName)
	if err != nil {
		return fmt.Errorf("open first file for merge '%s': %w", file1.fName, err)
	}
	defer f1.Close()

	f2, err := os.Open(file2.fName)
	if err != nil {
		return fmt.Errorf("open second file for merge '%s': %w", file2.fName, err)
	}
	defer f2.Close()

	type operationAndValue struct {
		operation string
		value     string
	}

	keyToOperationAndValue := make(map[string]operationAndValue)

	sf1 := bufio.NewScanner(f1)
	for sf1.Scan() {
		parts := strings.SplitN(sf1.Text(), " ", 3)

		if len(parts) == 2 {
			keyToOperationAndValue[parts[1]] = operationAndValue{operation: parts[0]}
		} else if len(parts) == 3 {
			keyToOperationAndValue[parts[1]] = operationAndValue{operation: parts[0], value: parts[2]}
		}
	}

	sf2 := bufio.NewScanner(f2)
	for sf2.Scan() {
		parts := strings.SplitN(sf2.Text(), " ", 3)

		if len(parts) == 2 {
			keyToOperationAndValue[parts[1]] = operationAndValue{operation: parts[0]}
		} else if len(parts) == 3 {
			keyToOperationAndValue[parts[1]] = operationAndValue{operation: parts[0], value: parts[2]}
		}
	}

	finalOperations := make([]string, 0)
	for key, oAndV := range keyToOperationAndValue {
		if strings.ToUpper(oAndV.operation) != dbengine.OpTypeDelete.String() {
			var s strings.Builder

			s.Grow(len(oAndV.operation) + 1 + len(key) + 1 + len(oAndV.value))

			s.WriteString(oAndV.operation)
			s.WriteString(" ")
			s.WriteString(key)
			s.WriteString(" ")
			s.WriteString(oAndV.value)

			finalOperations = append(finalOperations, s.String())
		}
	}

	f1.Close()
	f2.Close()

	// TODO что если:
	// 1: после удаления первого файла не удастся удалить второй?
	// 2: получится удалить оба, но не удастся записать результат мерджа?
	if err = os.Remove(file1.fName); err != nil {
		return fmt.Errorf("remove first WAL file before save merge results '%s': %w", file1.fName, err)
	}
	if err = os.Remove(file2.fName); err != nil {
		return fmt.Errorf("remove second WAL file before save merge results '%s': %w", file2.fName, err)
	}

	mergedF, err := os.OpenFile(file2.fName, os.O_CREATE|os.O_WRONLY|os.O_SYNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open second file to save merge results '%s': %w", file2.fName, err)
	}

	for _, o := range finalOperations {
		_, err = mergedF.WriteString(o + "\n")
		if err != nil {
			return err
		}
	}

	return mergedF.Sync()
}

// max heapify by fNum
func heapify(arr []fInfo, idx int) {
	largestValueIdx := idx
	leftChildNode := 2*idx + 1
	rightChildNode := 2*idx + 2

	if leftChildNode < len(arr) && arr[leftChildNode].fNum > arr[largestValueIdx].fNum {
		largestValueIdx = leftChildNode
	}

	if rightChildNode < len(arr) && arr[rightChildNode].fNum > arr[largestValueIdx].fNum {
		largestValueIdx = rightChildNode
	}

	if largestValueIdx != idx {
		arr[idx], arr[largestValueIdx] = arr[largestValueIdx], arr[idx]
		heapify(arr, largestValueIdx)
	}
}
