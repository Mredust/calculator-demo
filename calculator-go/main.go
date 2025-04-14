package main

import (
	"calculator-go/config"
	calculatorConnect "calculator-go/core/calculator"
	"calculator-go/service"
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net/http"
)

// CORS 中间件
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 允许的源（生产环境应替换为具体域名）
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// 明确允许的方法（必须包含 POST 和 OPTIONS）
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		// 关键：允许 gRPC-Web 和 Connect 协议所需的头部
		w.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Connect-Protocol-Version, X-Grpc-Web, Grpc-Timeout",
		)
		// 允许浏览器暴露自定义头部（如 gRPC 状态）
		w.Header().Set("Access-Control-Expose-Headers",
			"Grpc-Status, Grpc-Message",
		)
		// 预检请求直接返回 204
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 创建计算器服务实例
	calculatorServer := service.NewCalculatorServer()

	// 设置路由
	mux := http.NewServeMux()
	path, handler := calculatorConnect.NewCalculatorServiceHandler(calculatorServer)
	mux.Handle(path, handler)
	corsHandler := enableCORS(mux)
	// 启动服务器
	fmt.Printf("Calculator service running on http://localhost:%s\n", cfg.Server.Port)
	err = http.ListenAndServe(
		":"+cfg.Server.Port,
		h2c.NewHandler(corsHandler, &http2.Server{}),
	)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
