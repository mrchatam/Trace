package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
)

func decodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	if err := dec.Decode(dst); err != nil {
		return err
	}
	return nil
}

func queryInt(r *http.Request, key string, def int) (int, error) {
	v := strings.TrimSpace(r.URL.Query().Get(key))
	if v == "" {
		return def, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func mapDomainErr(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	var inv *domain.ErrInvalidTransition
	if errors.As(err, &inv) {
		msg := inv.Error()
		code := "CONFLICT"
		status := http.StatusConflict
		if strings.Contains(msg, "missing required capabilities") {
			code = "FORBIDDEN"
			status = http.StatusForbidden
		}
		writeEnvelope(w, status, code, msg, map[string]string{
			"from": inv.From, "to": inv.To, "reason": inv.Reason,
		})
		return true
	}
	var val *domain.ErrValidation
	if errors.As(err, &val) {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", val.Error(), nil)
		return true
	}
	msg := err.Error()
	// loop.requireUUID and similar plain errors ("… must be UUID") → 400, not 500.
	if strings.Contains(msg, "must be UUID") {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", msg, nil)
		return true
	}
	if isNotFoundErr(err) {
		writeEnvelope(w, http.StatusNotFound, "NOT_FOUND", msg, nil)
		return true
	}
	writeEnvelope(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error", nil)
	return true
}

func isNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "sql: no rows") || strings.Contains(msg, "no rows in result set")
}

func (s *Server) handleNotImplemented(w http.ResponseWriter, r *http.Request) {
	writeEnvelope(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "not implemented in this Trace release", nil)
}

func (s *Server) handleRPCForbidden(w http.ResponseWriter, r *http.Request) {
	writeEnvelope(w, http.StatusNotFound, "NOT_FOUND", "MCP /rpc is not available over HTTP; use local stdio MCP", nil)
}
