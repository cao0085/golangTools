package core

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

// AnalyzeFiles 分析所有 TXT 檔案的行數（並行處理）
func AnalyzeFiles(txtFiles []string) ([]*TxtFileInfo, error) {
	fileInfoList := make([]*TxtFileInfo, len(txtFiles))

	// 使用 WaitGroup 進行並行處理
	var wg sync.WaitGroup
	var mu sync.Mutex
	errors := make([]error, 0)

	fmt.Println("正在分析檔案...")
	fmt.Println()

	for i, file := range txtFiles {
		wg.Add(1)
		// 並行處理
		go func(index int, filePath string) {
			defer wg.Done() // 函式結束時自動執行

			lineCount, err := countFileLines(filePath)
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("檔案 %s 讀取失敗: %v", filePath, err))
				mu.Unlock()
				return
			}

			mu.Lock()
			fileInfoList[index] = NewTxtFileInfo(filePath, lineCount)
			mu.Unlock()
		}(i, file)
	}

	wg.Wait()

	if len(errors) > 0 {
		fmt.Println("警告: 部分檔案讀取失敗")
		for _, err := range errors {
			fmt.Printf("  - %v\n", err)
		}
	}

	// 過濾掉失敗的檔案
	validFiles := make([]*TxtFileInfo, 0)
	for _, info := range fileInfoList {
		if info != nil {
			validFiles = append(validFiles, info)
		}
	}

	totalLines := 0
	for _, f := range validFiles {
		totalLines += f.LineCount
	}

	fmt.Printf("總計：%d 個檔案，共 %d 行\n", len(validFiles), totalLines)

	return validFiles, nil
}

// countFileLines 計算檔案行數
func countFileLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

// ValidateAndAllocateFiles 驗證並分配檔案到 Excel 檔案中
func ValidateAndAllocateFiles(
	fileInfoList []*TxtFileInfo,
	maxRowsPerExcel int,
	desiredExcelCount int,
) ([][]*TxtFileInfo, error) {
	// 檢查是否有任何單一檔案超過最大列數限制
	for _, fileInfo := range fileInfoList {
		if fileInfo.LineCount > maxRowsPerExcel {
			return nil, fmt.Errorf(
				"檔案 '%s' 有 %d 行，超過單一 Excel 最大列數限制 %d 行，此檔案無法處理",
				fileInfo.FileName,
				fileInfo.LineCount,
				maxRowsPerExcel,
			)
		}
	}

	// 模擬分配檔案到 Excel
	allocation := simulateFileAllocation(fileInfoList, maxRowsPerExcel)

	// 檢查是否超過使用者期望的檔案個數
	if len(allocation) > desiredExcelCount {
		return nil, fmt.Errorf(
			"根據檔案大小和不可分割規則，需要至少 %d 個 Excel 檔案才能容納所有資料，但使用者只設定了 %d 個檔案。\n建議：增加檔案個數至 %d 個或更多，或增加每個 Excel 的最大列數限制",
			len(allocation),
			desiredExcelCount,
			len(allocation),
		)
	}

	return allocation, nil
}

// simulateFileAllocation 模擬檔案分配到 Excel
// 規則：一個檔案不可分割，如果加上該檔案會超出最大限制，就分配到下一個 Excel
func simulateFileAllocation(fileInfoList []*TxtFileInfo, maxRowsPerExcel int) [][]*TxtFileInfo {
	allocation := make([][]*TxtFileInfo, 0)
	currentExcel := make([]*TxtFileInfo, 0)
	currentRowCount := 0

	for _, fileInfo := range fileInfoList {
		// 如果加入當前檔案會超過限制，則開啟新的 Excel
		if currentRowCount+fileInfo.LineCount > maxRowsPerExcel && len(currentExcel) > 0 {
			allocation = append(allocation, currentExcel)
			currentExcel = make([]*TxtFileInfo, 0)
			currentRowCount = 0
		}

		// 將檔案加入當前 Excel
		currentExcel = append(currentExcel, fileInfo)
		currentRowCount += fileInfo.LineCount
	}

	// 加入最後一個 Excel（如果有內容）
	if len(currentExcel) > 0 {
		allocation = append(allocation, currentExcel)
	}

	return allocation
}

// DisplayAllocation 顯示分配結果
func DisplayAllocation(allocation [][]*TxtFileInfo, maxRowsPerExcel int) {
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("檔案分配結果：")
	fmt.Println("═══════════════════════════════════════════════════")

	for i, excelFiles := range allocation {
		totalRows := 0
		for _, f := range excelFiles {
			totalRows += f.LineCount
		}
		percentage := float64(totalRows) / float64(maxRowsPerExcel) * 100

		fmt.Println()
		fmt.Printf("Excel 檔案 #%d：\n", i+1)
		fmt.Printf("  包含 %d 個 TXT 檔案，共 %d 行 (%.2f%%)\n", len(excelFiles), totalRows, percentage)

		for _, file := range excelFiles {
			fmt.Printf("    - %s: %d 行\n", file.FileName, file.LineCount)
		}
	}

	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════")
}
