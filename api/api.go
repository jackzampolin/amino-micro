package api

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/cmd/gaia/app"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/gorilla/mux"
)

type encodeReq struct {
	Tx auth.StdTx `json:"tx"`
}

type encodeResp struct {
	Tx string `json:"tx"`
}

// Server represents the API server
type Server struct {
	Port int `json:"port"`

	Version string
	Commit  string
	Branch  string
}

// Router returns the router
func (s *Server) Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/version", s.VersionHandler)
	router.HandleFunc("/tx/encode", s.EncodeHandler).Methods("POST")

	return router
}

// EncodeHandler takes JSON input and encodes it into amino
func (s Server) EncodeHandler(w http.ResponseWriter, r *http.Request) {
	var m encodeReq

	cdc := app.MakeCodec()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, cdc, http.StatusBadRequest, err.Error())
		return
	}

	err = cdc.UnmarshalJSON(body, m)
	if err != nil {
		writeErrorResponse(w, cdc, http.StatusBadRequest, err.Error())
		return
	}

	// Re-encode it to the wire protocol
	txBytes, err := cdc.MarshalBinaryLengthPrefixed(m.Tx)
	if err != nil {
		writeErrorResponse(w, cdc, http.StatusInternalServerError, err.Error())
		return
	}

	// Encode the bytes to base64
	txBytesBase64 := base64.StdEncoding.EncodeToString(txBytes)

	// Return to client
	output, err := cdc.MarshalJSON(encodeResp{Tx: txBytesBase64})
	if err != nil {
		writeErrorResponse(w, cdc, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

// VersionHandler handles the /version route
func (s *Server) VersionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("{\"version\": \"%s\", \"commit\": \"%s\", \"branch\": \"%s\"}", s.Version, s.Commit, s.Branch)))
}

func writeErrorResponse(w http.ResponseWriter, cdc *codec.Codec, status int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(cdc.MustMarshalJSON(errorResponse{Code: 0, Message: err}))
}

// errorResponse defines the attributes of a JSON error response.
type errorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}
