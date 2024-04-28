package tax

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func HandlePersonalCalculationsCSV(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("taxFile")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		defer src.Close()

		requests, err := parseCsv(src, file.Filename)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		taxes, err := calculateTaxes(requests, db)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return respondWithTaxes(c, taxes)
	}

}

func parseCsv(file io.Reader, fileName string) ([]Request, error) {
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext != ".csv" {
		return nil, errors.New("file extension must be .csv")
	}

	baseName := strings.ToLower(filepath.Base(fileName))
	if baseName != "taxes.csv" {
		return nil, errors.New("file name must be taxes.csv")
	}

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 1 {
		return nil, errors.New("empty file")
	}
	header := records[0]
	if len(header) != 3 || header[0] != "totalIncome" || header[1] != "wht" || header[2] != "donation" {
		return nil, errors.New("invalid header row")
	}

	var requests []Request
	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) != 3 {
			return nil, errors.New("invalid row data")
		}

		totalIncome, err := strconv.ParseFloat(record[0], 64)
		if err != nil || totalIncome < 0 {
			return nil, errors.New("invalid totalIncome value")
		}

		wht, err := strconv.ParseFloat(record[1], 64)
		if err != nil || wht < 0 {
			return nil, errors.New("invalid wht value")
		}

		donation, err := strconv.ParseFloat(record[2], 64)
		if err != nil || donation < 0 {
			return nil, errors.New("invalid donation value")
		}

		requests = append(requests, Request{
			TotalIncome: totalIncome,
			WHT:         wht,
			Allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        donation,
				},
			},
		})
	}

	return requests, nil
}

func calculateTaxes(requests []Request, db *sql.DB) ([]TaxResponse, error) {
	var responses []TaxResponse
	for _, req := range requests {
		resp, err := calculateTax(req, db)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

func calculateTax(req Request, db *sql.DB) (TaxResponse, error) {
	tax, _, _, err := CalculateTax(db, req.TotalIncome, req.WHT, req.Allowances)
	if err != nil {
		return TaxResponse{}, err
	}

	return TaxResponse{
		TotalIncome: req.TotalIncome,
		Tax:         tax,
	}, nil
}

func respondWithTaxes(c echo.Context, taxes []TaxResponse) error {
	jsonResp, err := json.Marshal(struct {
		Taxes []TaxResponse `json:"taxes"`
	}{
		Taxes: taxes,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSONBlob(http.StatusOK, jsonResp)
}
