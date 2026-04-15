package server

import (
	"fmt"
	"math/rand"
	"net"
	"winresizer/utils"

	"github.com/gin-gonic/gin"
)

// Server는 gin 웹서버입니다.
type Server struct {
	Port   int
	engine *gin.Engine
}

// New는 웹서버를 초기화하고 라우터를 등록합니다.
func New() *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(noCacheMiddleware())

	registerRoutes(r)

	port := findFreePort(40000, 49999)
	return &Server{Port: port, engine: r}
}

// Start는 웹서버를 백그라운드에서 시작합니다.
func (s *Server) Start() {
	addr := fmt.Sprintf("127.0.0.1:%d", s.Port)
	utils.Log.Infof("웹 서버 시작: http://%s", addr)
	if err := s.engine.Run(addr); err != nil {
		utils.Log.Errorf("웹 서버 오류: %v", err)
	}
}

// findFreePort는 지정된 범위에서 사용 가능한 랜덤 포트를 반환합니다.
func findFreePort(start, end int) int {
	for {
		port := start + rand.Intn(end-start+1)
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		ln, err := net.Listen("tcp", addr)
		if err == nil {
			ln.Close()
			return port
		}
	}
}

// noCacheMiddleware는 브라우저 캐시를 방지하는 헤더를 추가합니다.
func noCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		c.Header("Expires", "0")
		c.Next()
	}
}
