package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"winprockill/internal/service"
)

type H struct {
	cmd    *service.WinCommand
	uiTmpl []byte
}

func New(cmd *service.WinCommand, uiTmpl []byte) *H {
	return &H{cmd: cmd, uiTmpl: uiTmpl}
}

func (h *H) Ui(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	_, _ = w.Write(h.uiTmpl)
}

func (h *H) Processes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	p, err := h.cmd.Processes(r.Context())

	sendJson(w, p, err)
}

func (h *H) Control(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	sendJson(w, nil, h.cmd.KillProcesses(r.Context()))
}

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func sendJson(w http.ResponseWriter, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")

	payload := Response{
		Data: data,
	}

	if err != nil {
		payload.Error = err.Error()
	}

	err = json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Println("encode payload: ", err, payload)
	}
}
