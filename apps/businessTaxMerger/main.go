package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"accountingTools/apps/businessTaxMerger/core"
)

func main() {
	continueProgram := true

	for continueProgram {
		// clearScreen()
		fmt.Print("\033[H\033[2J")
		// printHeader()
		fmt.Println("╔══════════════════════════════════════════════════════════╗")
		fmt.Println("║            營業人進銷項資料檔合併工具 v2.0 (Go)          ║")
		fmt.Println("║                                                          ║")
		fmt.Println("║    流程：選擇資料夾 → 分析 TXT → 合併 → 匯出 Excel       ║")
		fmt.Println("╚══════════════════════════════════════════════════════════╝")
		fmt.Println()

		// Step 1: 選擇資料夾
		folderPath, txtFiles, err := selectFolderAndFiles()
		if err != nil {
			fmt.Printf("錯誤: %v\n", err)
			fmt.Println("按 Enter 重新選擇資料夾...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			continue
		}

		if len(txtFiles) == 0 {
			fmt.Println("資料夾中沒有找到 TXT 檔案。")
			fmt.Println("按 Enter 重新選擇資料夾...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			continue
		}

		fmt.Printf("\n找到 %d 個 TXT 檔案：\n", len(txtFiles))
		displayLimit := 10
		if len(txtFiles) < displayLimit {
			displayLimit = len(txtFiles)
		}
		for i := 0; i < displayLimit; i++ {
			fmt.Printf("  %d. %s\n", i+1, filepath.Base(txtFiles[i]))
		}
		if len(txtFiles) > 10 {
			fmt.Printf("  ... 以及其他 %d 個檔案\n", len(txtFiles)-10)
		}

		fmt.Print("\n是否開始處理這些檔案？(y/n): ")
		if !confirmYes() {
			continue
		}

		// Step 2: 分析 TXT 檔案
		fmt.Println()
		fileInfoList, err := core.AnalyzeFiles(txtFiles)
		if err != nil {
			fmt.Printf("分析檔案時發生錯誤: %v\n", err)
			continue
		}

		if len(fileInfoList) == 0 {
			fmt.Println("沒有成功分析到任何檔案！")
			fmt.Print("是否要重新嘗試？(y/n): ")
			if !confirmYes() {
				continueProgram = false
			}
			continue
		}

		// Step 3: 取得使用者參數
		maxRowsPerExcel, desiredExcelCount := getUserParameters()

		// Step 4: 驗證並分配檔案
		fmt.Println()
		fmt.Println("正在驗證檔案分配...")

		allocation, err := core.ValidateAndAllocateFiles(fileInfoList, maxRowsPerExcel, desiredExcelCount)
		if err != nil {
			fmt.Println()
			fmt.Printf("❌ 驗證失敗：%v\n", err)
			fmt.Print("\n是否要重新設定參數？(y/n): ")
			if !confirmYes() {
				continueProgram = false
			}
			continue
		}

		// 顯示分配結果
		core.DisplayAllocation(allocation, maxRowsPerExcel)

		// 確認是否繼續
		fmt.Println()
		fmt.Print("是否確認以上分配並開始產出 Excel？(y/n): ")
		if !confirmYes() {
			continue
		}

		// Step 5: 產出 Excel 檔案
		fmt.Println()
		fmt.Println("═══════════════════════════════════════════════════")
		fmt.Println("開始產出 Excel 檔案...")
		fmt.Println("═══════════════════════════════════════════════════")
		fmt.Println()

		if err := core.ExportToExcel(allocation, folderPath, maxRowsPerExcel); err != nil {
			fmt.Println()
			fmt.Printf("❌ Excel 匯出失敗：%v\n", err)
		} else {
			fmt.Println()
			fmt.Println("✓ 所有 Excel 檔案產出成功！")
			fmt.Printf("   輸出位置: %s\n", folderPath)
		}

		// 詢問是否繼續
		fmt.Println()
		fmt.Println()
		fmt.Print("====== 是否要繼續處理其他資料夾？(y/n): ")
		if !confirmYes() {
			continueProgram = false
		}
	}

	fmt.Println("\n感謝使用，再見！")
	fmt.Println("按 Enter 離開...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// selectFolderAndFiles 選擇資料夾並取得 TXT 檔案
func selectFolderAndFiles() (string, []string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("請輸入要處理的資料夾路徑（或直接拖曳資料夾）：")
	fmt.Print("> ")

	folderPath, err := reader.ReadString('\n')
	if err != nil {
		return "", nil, err
	}

	folderPath = strings.TrimSpace(folderPath)
	folderPath = strings.Trim(folderPath, "\"") // 移除拖曳時的引號

	// 檢查資料夾是否存在
	info, err := os.Stat(folderPath)
	if err != nil {
		return "", nil, fmt.Errorf("資料夾不存在或無法存取")
	}

	if !info.IsDir() {
		return "", nil, fmt.Errorf("路徑不是資料夾")
	}

	// 取得所有 TXT 檔案
	txtFiles, err := filepath.Glob(filepath.Join(folderPath, "*.txt"))
	if err != nil {
		return "", nil, err
	}

	return folderPath, txtFiles, nil
}

// getUserParameters 取得使用者參數
func getUserParameters() (int, int) {
	reader := bufio.NewReader(os.Stdin)

	// 取得每個 Excel 最大列數
	var maxRowsPerExcel int
	for {
		fmt.Println()
		fmt.Print("請輸入每個 Excel 檔案的最大列數限制 (預設: 1048576): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			maxRowsPerExcel = 1048576
			break
		}

		val, err := strconv.Atoi(input)
		if err != nil || val <= 0 {
			fmt.Println("請輸入有效的正整數")
			continue
		}

		maxRowsPerExcel = val
		break
	}

	// 取得期望的 Excel 檔案個數
	var desiredExcelCount int
	for {
		fmt.Print("請輸入期望產出的 Excel 檔案個數 (預設: 10): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			desiredExcelCount = 10
			break
		}

		val, err := strconv.Atoi(input)
		if err != nil || val <= 0 {
			fmt.Println("請輸入有效的正整數")
			continue
		}

		desiredExcelCount = val
		break
	}

	return maxRowsPerExcel, desiredExcelCount
}

// confirmYes 確認是否為 yes
func confirmYes() bool {
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
