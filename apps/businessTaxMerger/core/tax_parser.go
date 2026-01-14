package core

import (
	"fmt"
	"strings"
)

// ParseLine 解析單行資料
// line: 原始行資料
// lineNumber: 行號
// sourceFileName: 來源檔案名稱
func ParseLine(line string, lineNumber int, sourceFileName string) (*TaxRecord, error) {
	if len(strings.TrimSpace(line)) == 0 {
		return nil, fmt.Errorf("資料行不可為空")
	}

	record := &TaxRecord{
		RawData:        line,
		LineNumber:     lineNumber,
		SourceFileName: sourceFileName,
	}

	// 基本資訊
	record.FormatCode = safeSubstring(line, 0, 2)              // 1-2
	record.DeclarantTaxId = safeSubstring(line, 2, 9)          // 3-11
	record.SequenceNumber = safeSubstring(line, 11, 7)         // 12-18

	// 資料所屬年月
	record.DataYear = safeSubstring(line, 18, 3)               // 19-21
	record.DataMonth = safeSubstring(line, 21, 2)              // 22-23

	// 買受人/營業/銷售人統編號 (欄位共用 24-31)
	record.BuyerTaxId = safeSubstring(line, 23, 8)             // 24-31
	record.BusinessNumber = safeSubstring(line, 23, 8)         // 24-31 (共用)

	record.SellerTaxId = safeSubstring(line, 31, 8)            // 32-39

	// 統一發票
	record.InvoicePrefix = safeSubstring(line, 39, 2)          // 40-41
	record.InvoiceStartNumber = safeSubstring(line, 41, 8)     // 42-49

	// 彙總張數 (欄位共用 32-35)
	record.TotalSheets = safeSubstring(line, 31, 4)            // 32-35
	record.Blank1 = safeSubstring(line, 35, 4)                 // 36-39

	// 其他憑證號碼/公用事業流水號 (欄位共用 40-49)
	record.OtherVoucherNumber = safeSubstring(line, 39, 10)    // 40-49
	record.UtilitySequenceNumber = safeSubstring(line, 39, 10) // 40-49 (共用)

	// 海關代徵
	record.Blank2 = safeSubstring(line, 31, 4)                 // 32-35
	record.CustomsTaxPaymentNumber = safeSubstring(line, 35, 14) // 36-49

	// 金額 (欄位共用 50-61)
	record.SalesAmount = safeSubstring(line, 49, 12)           // 50-61
	record.TaxBase = safeSubstring(line, 49, 12)               // 50-61 (共用)

	// 課稅別
	record.TaxType = safeSubstring(line, 61, 1)                // 62-62

	// 營業稅額
	record.TaxAmount = safeSubstring(line, 62, 10)             // 63-72

	// 扣抵代號
	record.DeductionCode = safeSubstring(line, 72, 1)          // 73-73
	record.Blank3 = safeSubstring(line, 73, 5)                 // 74-78

	// 特種稅額類稅率
	record.SpecialTaxRate = safeSubstring(line, 78, 1)         // 79-79

	// 彙加性記/分攤性記 (欄位共用 80-80)
	record.AggregationMark = safeSubstring(line, 79, 1)        // 80-80
	record.AllocationMark = safeSubstring(line, 79, 1)         // 80-80 (共用)

	// 通關方式註記
	record.CustomsClearanceMark = safeSubstring(line, 80, 1)   // 81-81

	return record, nil
}

// safeSubstring 安全的字串截取（處理超出範圍的情況）
func safeSubstring(source string, startIndex, length int) string {
	if len(source) == 0 {
		return ""
	}

	// 將 string 轉為 rune 陣列以正確處理 UTF-8 字元
	runes := []rune(source)

	if startIndex >= len(runes) {
		return ""
	}

	endIndex := startIndex + length
	if endIndex > len(runes) {
		endIndex = len(runes)
	}

	return strings.TrimSpace(string(runes[startIndex:endIndex]))
}
