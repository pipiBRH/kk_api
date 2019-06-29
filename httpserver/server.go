package httpserver

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Addr   string
	Port   int
	Engine *gin.Engine
}

func NewServer(addr string, port int, engine *gin.Engine) *Server {
	return &Server{
		Addr:   addr,
		Port:   port,
		Engine: engine,
	}
}

func (s *Server) InitServer() {
	s.RouterSetting()
	s.Engine.Run(fmt.Sprintf("%s:%d", s.Addr, s.Port))
}
