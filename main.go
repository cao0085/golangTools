package main

import (
	"fmt"
)

func main() {
	// 設定終端機為 UTF-8 編碼 (Windows)
	fmt.Println("\033[?25h") // 顯示游標
	// 保留未來可整合成一個入口
}
