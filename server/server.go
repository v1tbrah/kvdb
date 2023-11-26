package server

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"strings"

	"github.com/v1tbrah/kvdb/txctx"
)

type dbengine interface {
	Process(ctx context.Context, in string, withSaveToWAL bool) (out string, err error)
}

type Server struct {
	host, port string

	dbengine dbengine
}

func New(host, port string, dbengine dbengine) (*Server, error) {
	if dbengine == nil {
		return nil, errors.New("dbengine is nil")
	}

	return &Server{
		host:     host,
		port:     port,
		dbengine: dbengine,
	}, nil
}

func (e *Server) Launch(ctx context.Context) error {
	listen, err := net.Listen("tcp", net.JoinHostPort(e.host, e.port))
	if err != nil {
		return err
	}
	slog.Info("tcp server started", slog.String("host", e.host), slog.String("port", e.port))

	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				if ctx.Err() != nil {
					return
				}

				slog.Error("listen.Accept", slog.String("error", err.Error()))
				continue
			}

			go e.handleIncomingRequests(ctx, conn)
		}
	}()

	<-ctx.Done()

	if errClose := listen.Close(); errClose != nil {
		slog.Error("listen.Close", slog.String("error", errClose.Error()))
	}

	return ctx.Err()
}

func (e *Server) handleIncomingRequests(ctx context.Context, conn net.Conn) {
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

		ctxTx, ctxTxCancel := context.WithCancel(ctx)
		ctxTx = txctx.CtxWithTx(ctxTx)
		out, err := e.dbengine.Process(ctxTx, cleanedLine, true)
		if err != nil {
			if _, errWr := io.WriteString(conn, err.Error()+"\n"); errWr != nil {
				slog.Error("io.WriteString", slog.String("errWr", errWr.Error()), slog.String("err", err.Error()))
			}
			ctxTxCancel()
			continue
		}

		if _, errWr := io.WriteString(conn, out+"\n"); errWr != nil {
			slog.Error("io.WriteString", slog.String("errWr", errWr.Error()), slog.String("out", out))
		}
		ctxTxCancel()
	}
}
