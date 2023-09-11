package engine

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"

	"github.com/v1tbrah/kvdb/model"
	"github.com/v1tbrah/kvdb/parse"
)

type storage interface {
	Set(key, value string)
	Get(key string) string
	Delete(key string)
}

type Engine struct {
	host, port string

	storage storage
}

func NewEngine(host, port string, storage storage) (*Engine, error) {
	if storage == nil {
		return nil, errors.New("storage is nil")
	}

	return &Engine{
		host:    host,
		port:    port,
		storage: storage,
	}, nil
}

func (e *Engine) Launch(ctx context.Context) error {
	listen, err := net.Listen("tcp", net.JoinHostPort(e.host, e.port))
	if err != nil {
		return err
	}
	slog.Info("tcp server started", slog.String("host", e.host), slog.String("port", e.port))

	var listenClosedFlag bool
	go func() {

		for {
			conn, err := listen.Accept()
			if err != nil {
				if listenClosedFlag {
					return
				}

				slog.Error("listen.Accept", slog.String("error", err.Error()))
				continue
			}
			go e.handleIncomingRequests(conn)
		}
	}()

	<-ctx.Done()

	listenClosedFlag = true
	if errClose := listen.Close(); errClose != nil {
		slog.Error("listen.Close", slog.String("error", errClose.Error()))
	}

	return ctx.Err()
}

func (e *Engine) handleIncomingRequests(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			slog.Error("conn.Close", slog.String("error", err.Error()))
		}
	}()

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			slog.Error("reader.ReadString", slog.String("error", err.Error()))
			return
		}

		// Deletes \r\n from end of line
		cleanedLine := strings.TrimRight(line, "\r\n")
		slog.Debug("received line", slog.String("line", cleanedLine))

		out, err := e.processInputData(cleanedLine)
		if err != nil {
			_, errWr := io.WriteString(conn, err.Error())
			if errWr != nil {
				slog.Error("io.WriteString", slog.String("errWr", errWr.Error()), slog.String("err", err.Error()))
			}
			continue
		}

		_, errWr := io.WriteString(conn, out+"\n")
		if errWr != nil {
			slog.Error("io.WriteString", slog.String("errWr", errWr.Error()), slog.String("out", out))
			continue
		}
	}
}

func (e *Engine) processInputData(inputData string) (out string, err error) {
	parsedData, err := parse.ParseData(inputData)
	if err != nil {
		return "", err
	}

	switch parsedData.OpType {
	case model.OpTypeSet:
		e.storage.Set(parsedData.Key, parsedData.Value)
		return "OK", nil
	case model.OpTypeGet:
		out = e.storage.Get(parsedData.Key)
		return out, nil
	case model.OpTypeDelete:
		e.storage.Delete(parsedData.Key)
		return "OK", nil
	default:
		return "", fmt.Errorf("unsupported operation: %s", parsedData.OpType)
	}
}
