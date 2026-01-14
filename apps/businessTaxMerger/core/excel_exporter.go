package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// ExportToExcel 匯出營業稅資料到 Excel
func ExportToExcel(allocation [][]*TxtFileInfo, outputFolder string, maxRowsPerExcel int) error {
	timestamp := time.Now().Format("20060102_150405")

	for i, fileGroup := range allocation {
		fileName := fmt.Sprintf("營業人進銷項資料_%d_%s.xlsx", i+1, timestamp)
		fullPath := filepath.Join(outputFolder, fileName)

		fmt.Printf("正在產出第 %d 個 Excel 檔案...\n", i+1)

		// 讀取並解析所有分配到這個 Excel 的檔案
		allRecords := make([]*TaxRecord, 0)

		fmt.Printf("  讀取並解析 %d 個檔案...\n", len(fileGroup))
		for _, fileInfo := range fileGroup {
			records, err := parseFile(fileInfo.FilePath)
			if err != nil {
				fmt.Printf("  警告: 解析檔案 %s 失敗: %v\n", fileInfo.FileName, err)
				continue
			}
			allRecords = append(allRecords, records...)
		}

		fmt.Printf("  ✓ 已讀取 %d 筆資料\n", len(allRecords))

		// 產生 Excel
		fmt.Println("  產生 Excel 檔案...")
		if err := createExcelFile(fullPath, allRecords, i+1); err != nil {
			return fmt.Errorf("產生 Excel 檔案失敗: %v", err)
		}

		fmt.Printf("  ✓ 已產出: %s\n", fileName)
	}

	return nil
}

// parseFile 解析檔案並返回記錄列表
func parseFile(filePath string) ([]*TaxRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	records := make([]*TaxRecord, 0)
	fileName := filepath.Base(filePath)

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		record, err := ParseLine(line, lineNumber, fileName)
		if err != nil {
			fmt.Printf("    警告: 檔案 %s 第 %d 行解析失敗: %v\n", fileName, lineNumber, err)
			continue
		}

		records = append(records, record)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

// createExcelFile 建立 Excel 檔案
func createExcelFile(filePath string, records []*TaxRecord, fileNumber int) error {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "營業人進銷項資料"

	// 創建工作表
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// 設定標題
	if err := setupHeaders(f, sheetName); err != nil {
		return err
	}

	// 寫入資料
	if err := writeData(f, sheetName, records); err != nil {
		return err
	}

	// 設定為預設工作表
	f.SetActiveSheet(index)

	// 刪除預設的 Sheet1
	if err := f.DeleteSheet("Sheet1"); err != nil {
		// 忽略錯誤，可能 Sheet1 不存在
	}

	// 儲存檔案
	if err := f.SaveAs(filePath); err != nil {
		return err
	}

	return nil
}

// setupHeaders 設定 Excel 標題列
func setupHeaders(f *excelize.File, sheetName string) error {
	// 定義標題（橘色底的欄位）
	headers := []string{
		"格式代號",
		"申報營業人稅籍編號",
		"流水號",
		"資料所屬年度",
		"資料所屬月份",
		"買受人統一編號",
		// "發票訖號",
		"銷售人統一編號",
		"發票(起)號碼", // 發票字軌 + 發票(起)號碼 合併
		"銷售金額",
		// "營業稅稅基",
		"課稅別",
		"營業稅額",
		"扣抵代號",
		"彙加註記",
		// "分攤註記",
	}

	// 創建標題樣式 (橘色背景 + 粗體 + 置中 + 邊框)
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFC000"}, // 橘色
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return err
	}

	// 寫入標題
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return err
		}
		if err := f.SetCellStyle(sheetName, cell, cell, headerStyle); err != nil {
			return err
		}
	}

	// 凍結第一列
	if err := f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	}); err != nil {
		return err
	}

	return nil
}

// writeData 寫入資料到 Excel
func writeData(f *excelize.File, sheetName string, records []*TaxRecord) error {
	// 創建資料列樣式 (邊框 + 靠右對齊)
	dataStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
	})
	if err != nil {
		return err
	}

	// 創建數字格式樣式 (邊框 + 靠右對齊 + 小數點兩位)
	numberStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
		NumFmt: 2, // 數字格式: 0.00
	})
	if err != nil {
		return err
	}

	row := 2 // 從第二列開始（第一列是標題）

	for _, record := range records {
		col := 1

		// 寫入各欄位資料

		// 基礎資料 (字串)
		values := []string{
			record.FormatCode,
			record.DeclarantTaxId,
			record.SequenceNumber,
			record.DataYear,
			record.DataMonth,
			record.BuyerTaxId,
			// record.BusinessNumber,
			record.SellerTaxId,
		}

		// 發票(起)號碼 = 發票字軌 + 發票(起)號碼
		invoiceNumber := record.InvoiceStartNumber
		if len(strings.TrimSpace(record.InvoicePrefix)) > 0 {
			invoiceNumber = record.InvoicePrefix + record.InvoiceStartNumber
		}
		values = append(values, invoiceNumber)

		// 寫入字串欄位 (套用 dataStyle 靠右對齊)
		for i, value := range values {
			cell, _ := excelize.CoordinatesToCellName(col+i, row)
			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				return err
			}
			if err := f.SetCellStyle(sheetName, cell, cell, dataStyle); err != nil {
				return err
			}
		}
		col += len(values)

		// 銷售金額 (轉成數字,套用 numberStyle)
		salesAmount := parseAmountToInt(record.SalesAmount)
		salesCell, _ := excelize.CoordinatesToCellName(col, row)
		if err := f.SetCellValue(sheetName, salesCell, salesAmount); err != nil {
			return err
		}
		if err := f.SetCellStyle(sheetName, salesCell, salesCell, numberStyle); err != nil {
			return err
		}
		col++

		// 課稅別 (字串)
		cell, _ := excelize.CoordinatesToCellName(col, row)
		if err := f.SetCellValue(sheetName, cell, record.TaxType); err != nil {
			return err
		}
		if err := f.SetCellStyle(sheetName, cell, cell, dataStyle); err != nil {
			return err
		}
		col++

		// 營業稅額 (轉成數字,套用 numberStyle)
		taxAmount := parseAmountToInt(record.TaxAmount)
		taxCell, _ := excelize.CoordinatesToCellName(col, row)
		if err := f.SetCellValue(sheetName, taxCell, taxAmount); err != nil {
			return err
		}
		if err := f.SetCellStyle(sheetName, taxCell, taxCell, numberStyle); err != nil {
			return err
		}
		col++

		// 扣抵代號 (字串)
		cell, _ = excelize.CoordinatesToCellName(col, row)
		if err := f.SetCellValue(sheetName, cell, record.DeductionCode); err != nil {
			return err
		}
		if err := f.SetCellStyle(sheetName, cell, cell, dataStyle); err != nil {
			return err
		}
		col++

		// 彙加註記 (字串)
		cell, _ = excelize.CoordinatesToCellName(col, row)
		if err := f.SetCellValue(sheetName, cell, record.AggregationMark); err != nil {
			return err
		}
		if err := f.SetCellStyle(sheetName, cell, cell, dataStyle); err != nil {
			return err
		}
		col++

		row++
	}

	// 自動調整欄寬
	for i := 1; i <= 16; i++ {
		col, _ := excelize.ColumnNumberToName(i)
		if err := f.SetColWidth(sheetName, col, col, 15); err != nil {
			return err
		}
	}

	return nil
}

// parseAmountToInt 將金額字串轉換成整數
// 例如: "000000123456" -> 123456
func parseAmountToInt(amountStr string) int64 {
	// 去除空白
	amountStr = strings.TrimSpace(amountStr)

	// 如果是空字串，回傳 0
	if amountStr == "" {
		return 0
	}

	// 轉換成整數
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		// 轉換失敗回傳 0
		return 0
	}

	return amount
}
