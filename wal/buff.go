package wal

import (
	"context"
	"log/slog"
	"time"
)

const (
	buffFlushInterval = time.Millisecond * 50
	buffMaxSize       = 16 * 1e3 // 16Kb
)

func (w *WAL) addOpToBuff(_ context.Context, op string) {
	w.operationsBuffMu.Lock()
	w.operationsBuff = append(w.operationsBuff, op)
	w.operationsBuffSize += len(op) + 1 // 1 for 'carriage return' symbol
	w.operationsBuffMu.Unlock()

	if w.operationsBuffSize > buffMaxSize {
		w.operationsBuffFlushTrigger <- struct{}{}
	}
}

func (w *WAL) startProcessingBuff() {
	ticker := time.NewTicker(buffFlushInterval)
	defer ticker.Stop()

	for {
		var (
			buffCopy          []string
			triggeredByTicker bool
		)

		select {
		case <-ticker.C:
			triggeredByTicker = true
		case <-w.operationsBuffFlushTrigger:
		}

		w.operationsBuffMu.Lock()
		buffCopy = make([]string, len(w.operationsBuff))
		copy(buffCopy, w.operationsBuff)
		w.operationsBuff = w.operationsBuff[:0]
		w.operationsBuffSize = 0
		w.operationsBuffMu.Unlock()

		w.flush(buffCopy)

		if !triggeredByTicker {
			ticker.Reset(buffFlushInterval)
		}
	}
}

func (w *WAL) flush(operations []string) {
	for _, op := range operations {
		if err := w.saveSync(context.Background(), op); err != nil {
			slog.Error("w.saveSync", slog.String("op", op), slog.String("err", err.Error()))
		}
	}
}
