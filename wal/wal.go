package wal

import (
	"context"
	"fmt"
	"os"
)

const (
	symbolBetweenStartAndSuffixNamesWALFiles     = "_"
	maxWALFileSize                           int = 16 * 1e6 // 16Mb
)

type WAL struct {
	currWALFile     *os.File
	currWALFileSize int
}

func New() (*WAL, error) {
	lastWALFile, err := findLastWALFile()
	if err != nil {
		return nil, fmt.Errorf("find last wal file: %w", err)
	}

	stat, err := lastWALFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("check last wal file stat: %w", err)
	}

	go startProcessMergeWALFiles()

	return &WAL{
		currWALFile:     lastWALFile,
		currWALFileSize: int(stat.Size()),
	}, nil
}

func (w *WAL) Save(_ context.Context, op string) error {
	// TODO нужно ли здесь блокировку делать???
	if w.currWALFileSize+len(op)+len("\n") > maxWALFileSize {
		newWALFile, err := prepareNewWALFile(w.currWALFile.Name())
		if err != nil {
			return fmt.Errorf("preparing new wal file: %w", err)
		}

		w.currWALFile = newWALFile
		w.currWALFileSize = 0
	}

	_, err := w.currWALFile.WriteString(op + "\n")
	if err != nil {
		return err
	}

	if err = w.currWALFile.Sync(); err != nil {
		return err
	}

	return nil
}
