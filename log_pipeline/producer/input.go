package producer

import (
	"net"

	"go.uber.org/zap"
)

func HandleProducer(conn net.Conn, logger *zap.Logger) {
	defer conn.Close()
	logger.Info("New connection", zap.String("remote_addr", conn.RemoteAddr().String()))

}
