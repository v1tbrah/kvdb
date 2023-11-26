package wal

import (
	"context"
	"fmt"
	"os"
	"sync"
)

const (
	symbolBetweenStartAndSuffixNamesWALFiles = "_"
)

type WAL struct {
	currWALFile     *os.File
	currWALFileSize int

	isSyncWrite bool

	operationsBuff             []string
	operationsBuffMu           *sync.Mutex
	operationsBuffSize         int
	operationsBuffFlushTrigger chan struct{}
}

func New(isSyncWrite bool) (*WAL, error) {
	lastWALFile, err := findLastWALFile()
	if err != nil {
		return nil, fmt.Errorf("find last wal file: %w", err)
	}

	stat, err := lastWALFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("check last wal file stat: %w", err)
	}

	go startProcessMergeWALFiles()

	w := &WAL{
		currWALFile:     lastWALFile,
		currWALFileSize: int(stat.Size()),

		isSyncWrite: isSyncWrite,

		operationsBuff:             make([]string, 0),
		operationsBuffMu:           new(sync.Mutex),
		operationsBuffSize:         0,
		operationsBuffFlushTrigger: make(chan struct{}),
	}

	if !isSyncWrite {
		go w.startProcessingBuff()
	}

	return w, nil
}

func (w *WAL) Save(ctx context.Context, op string) error {
	if w.isSyncWrite {
		return w.saveSync(ctx, op)
	}

	w.addOpToBuff(ctx, op)

	return nil
}

func (w *WAL) saveSync(_ context.Context, op string) error {
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
