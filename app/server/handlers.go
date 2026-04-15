package server

import (
	"net/http"
	"os"
	"winresizer/core"
	"winresizer/utils"

	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine) {
	// 정적 파일 — embed 패키지로 바이너리에 내장
	r.Static("/static", "./ui/static")
	r.LoadHTMLGlob("./ui/templates/*")

	r.GET("/", handleIndex)

	api := r.Group("/api")
	{
		api.GET("/status", handleGetStatus)
		api.GET("/config", handleGetConfig)
		api.POST("/config", handlePostConfig)
		api.POST("/config/reset", handleResetConfig)
		api.GET("/execute", handleExecuteGet)
		api.POST("/execute", handleExecutePost)
	}
}

// GET /
func handleIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// GET /api/status
func handleGetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"accessibility_granted":   core.CheckAccessibilityPermission(),
		"input_monitoring_granted": core.CheckInputMonitoringPermission(),
		"pid":                     os.Getpid(),
	})
}

// GET /api/config
func handleGetConfig(c *gin.Context) {
	cfg, err := core.GetConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cfg)
}

// POST /api/config
func handlePostConfig(c *gin.Context) {
	var cfg core.Config
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다."})
		return
	}
	if err := core.SaveConfig(&cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 설정 변경 후 단축키 리스너 재시작
	go core.RestartHotkeyManager()
	utils.Log.Infof("설정 저장 완료")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// POST /api/config/reset
func handleResetConfig(c *gin.Context) {
	cfg, err := core.LoadDefaultConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Log.Infof("기본 설정값 요청됨 (저장 미수행)")
	c.JSON(http.StatusOK, gin.H{"status": "ok", "config": cfg})
}

// GET /api/execute?mode=left_half
func handleExecuteGet(c *gin.Context) {
	mode := c.Query("mode")
	if mode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mode 파라미터가 필요합니다."})
		return
	}
	if err := core.ExecuteWindowCommand(mode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "mode": mode})
}

// POST /api/execute
func handleExecutePost(c *gin.Context) {
	var body struct {
		Mode string `json:"mode"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Mode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mode 필드가 필요합니다."})
		return
	}
	if err := core.ExecuteWindowCommand(body.Mode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "mode": body.Mode})
}
