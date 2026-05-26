package cotahist

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	recordTypeHeader = "00"
	recordTypeData   = "01"
	recordTypeFooter = "99"

	lineLength = 245
)

// Parse loads every TIPREG=01 line from r into memory.
// Use Walk for large files (annual COTAHIST has millions of records).
func Parse(r io.Reader) ([]Record, error) {
	var records []Record
	err := Walk(r, func(rec Record) error {
		records = append(records, rec)
		return nil
	})
	return records, err
}

// Walk streams the COTAHIST text in r and invokes fn for each data record.
// Returning a non-nil error from fn stops parsing and surfaces that error.
func Walk(r io.Reader, fn func(Record) error) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if len(line) < 2 {
			continue
		}

		switch line[0:2] {
		case recordTypeHeader, recordTypeFooter:
			continue
		case recordTypeData:
			rec, err := ParseLine(line)
			if err != nil {
				return fmt.Errorf("line %d: %w", lineNum, err)
			}
			if err := fn(rec); err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan: %w", err)
	}
	return nil
}

// ParseLine parses a single 245-char COTAHIST data line (TIPREG=01).
func ParseLine(line string) (Record, error) {
	if len(line) < lineLength {
		return Record{}, fmt.Errorf("line too short: got %d, want %d", len(line), lineLength)
	}

	date, err := time.Parse("20060102", line[2:10])
	if err != nil {
		return Record{}, fmt.Errorf("parse date %q: %w", line[2:10], err)
	}

	rec := Record{
		Date:          date,
		BDICode:       strings.TrimSpace(line[10:12]),
		Ticker:        strings.TrimSpace(line[12:24]),
		MarketType:    strings.TrimSpace(line[24:27]),
		ShortName:     strings.TrimSpace(line[27:39]),
		Specification: strings.TrimSpace(line[39:49]),
		Currency:      strings.TrimSpace(line[52:56]),
		OpenPrice:     parsePrice(line[56:69]),
		HighPrice:     parsePrice(line[69:82]),
		LowPrice:      parsePrice(line[82:95]),
		AveragePrice:  parsePrice(line[95:108]),
		ClosePrice:    parsePrice(line[108:121]),
		BestBid:       parsePrice(line[121:134]),
		BestAsk:       parsePrice(line[134:147]),
		Trades:        parseInt(line[147:152]),
		Quantity:      parseInt(line[152:170]),
		Volume:        parsePrice(line[170:188]),
		StrikePrice:   parsePrice(line[188:201]),
		QuoteFactor:   parseInt(line[210:217]),
		ISIN:          strings.TrimSpace(line[230:242]),
		DistNumber:    strings.TrimSpace(line[242:245]),
	}

	if exp := parseOptionalDate(line[202:210]); exp != nil {
		rec.ExpiryDate = exp
	}

	return rec, nil
}

// parsePrice converts a fixed-width COTAHIST price (last 2 digits implicit decimals) to float64.
func parsePrice(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return float64(n) / 100.0
}

func parseInt(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

func parseOptionalDate(s string) *time.Time {
	s = strings.TrimSpace(s)
	// Non-derivative rows use sentinel dates that aren't real expiries.
	if s == "" || s == "00000000" || s == "99991231" {
		return nil
	}
	t, err := time.Parse("20060102", s)
	if err != nil {
		return nil
	}
	return &t
}
