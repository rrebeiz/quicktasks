package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var (
	ErrInvalidIDParam = errors.New("invalid ID parameter")
)

type envelope = map[string]any

func (app *application) getParamID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParamFromCtx(r.Context(), "id"), 10, 64)
	if err != nil {
		return 0, ErrInvalidIDParam
	}
	if id < 1 {
		return 0, ErrInvalidIDParam
	}
	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, data envelope, status int, headers http.Header) error {
	var output []byte

	if app.config.environment == "dev" {
		js, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			return err
		}
		output = js

	} else {
		js, err := json.Marshal(data)
		if err != nil {
			return err
		}
		output = js
	}

	output = append(output, '\n')

	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err := w.Write(output)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1024 * 1024 * 1
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&data)
	if err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var syntaxError *json.SyntaxError
		var maxBytesError *http.MaxBytesError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("badly contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("unmarshal type error on field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("unmarshal type error on location %d", unmarshalTypeError.Offset)

		case errors.As(err, &invalidUnmarshalError):
			panic("cannot unmarshal, developer error perhaps you didn't pass it by reference?")

		case errors.As(err, &maxBytesError):
			app.errorLog.Printf("body must not be larger than %d bytes", maxBytesError.Limit)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}
