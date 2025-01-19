package controllers

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	cons "project_sem/internal/constants"
	"strings"

	"net/http"
	"project_sem/internal/logger"
	"project_sem/internal/services"
)

var (
	log *logger.Logger
)

type PriceHandler struct {
	ps services.PriceService
}

func NewPriceHandler(logger *logger.Logger, ps *services.PriceService) PriceHandler {
	log = logger
	return PriceHandler{ps: *ps}
}

func (ph PriceHandler) SavePrice(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	log.WithField("reqID", reqID).Info("Save the prices")

	body, err := io.ReadAll(r.Body)
	defer func() {
		err = r.Body.Close()
		if err != nil {
			log.
				WithField("reqID", reqID).
				WithError(err).Error("Close body")
		}
	}()
	if err != nil {
		log.
			WithField("reqID", reqID).
			WithError(err).Error("Error reading request body")
		http.Error(w, "Error reading request body", http.StatusBadRequest)
	}

	buf, err := unzipBody(body)
	if err != nil {
		log.WithField("reqID", reqID).
			WithError(err).Error("The attempt to unzip archive failed")
		http.Error(w, "The attempt to unzip archive failed", http.StatusBadRequest)
	}

	reader := csv.NewReader(buf)
	items := make([][]string, 0)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.WithField("reqID", reqID).WithError(err).Error("Read the record from CSV")
		}
		items = append(items, record)
	}

	price, err := ph.ps.SaveItem(r.Context(), items)
	if err != nil {
		log.WithField("reqID", reqID).
			WithError(err).Error("The attempt to save the price failed")
		http.Error(w, "The attempt to save the price failed", http.StatusBadRequest)
	}

	log.WithField("reqID", reqID).
		WithField("TotalPrice", price).
		Info("The prices have been saved")

	body, err = json.Marshal(price)
	if err != nil {
		log.WithField("reqID", reqID).
			WithError(err).Error("The attempt to marshal the price failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(body)
	if err != nil {
		log.WithField("reqID", reqID).
			WithError(err).Error("The attempt to write body failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Add(cons.KeyContentType, cons.TypeJSONContent)
	w.WriteHeader(http.StatusCreated)
}

func (ph PriceHandler) GetPrice(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	log.WithField("reqID", reqID).Info("Get all prices")

	item, err := ph.ps.AllItem(r.Context())
	if err != nil {
		log.WithField("reqID", reqID).
			WithError(err).Error("The attempt to got the prices failed")
		http.Error(w, "The attempt to got the prices failed", http.StatusInternalServerError)
	}

	var b bytes.Buffer

	writer := csv.NewWriter(&b)
	err = writer.WriteAll(item)
	if err != nil {
		log.WithField("reqID", reqID).
			WithError(err).Error("Error while writing to CSV")
		http.Error(w, "Error while writing to CSV", http.StatusInternalServerError)
	}

	writer.Flush()
	err = writer.Error()
	if err != nil {
		log.WithField("reqID", reqID).
			WithError(err).Error("Error while writing to CSV")
		http.Error(w, "Error while writing to CSV", http.StatusInternalServerError)
	}

	body, err := zipBody(&b)
	if err != nil {
		log.WithField("reqID", reqID).
			WithError(err).Error("The attempt to zip the prices failed")
		http.Error(w, "The attempt to zip the prices failed", http.StatusInternalServerError)
	}

	_, err = w.Write(body.Bytes())
	if err != nil {
		log.WithField("reqID", reqID).
			WithError(err).Error("The attempt to write body failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Add(cons.KeyContentType, cons.TypeZIPContent)
	w.WriteHeader(http.StatusCreated)
}

func unzipBody(body []byte) (*bytes.Buffer, error) {
	reader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	for _, file := range reader.File {
		if !strings.HasSuffix(file.Name, ".csv") {
			continue
		}

		f, err := file.Open()
		if err != nil {
			return nil, err
		}

		if _, err := io.Copy(&buf, f); err != nil {
			return nil, err
		}

		err = f.Close()
		if err != nil {
			return nil, err
		}
	}
	return &buf, nil
}

func zipBody(data *bytes.Buffer) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	file, err := zipWriter.Create("data.csv")
	if err != nil {
		return nil, err
	}

	_, err = file.Write(data.Bytes())
	if err != nil {
		return nil, err
	}

	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf, nil
}
