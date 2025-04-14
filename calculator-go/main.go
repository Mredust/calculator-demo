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
		// 设置CORS头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version, X-Grpc-Web, Grpc-Timeout")
		w.Header().Set("Access-Control-Expose-Headers", "Grpc-Status, Grpc-Message")

		// 预检请求直接返回
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
		//h2c.NewHandler(mux, &http2.Server{}),
	)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
