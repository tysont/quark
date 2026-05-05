// ABOUTME: HTTP server exposes the Node operations to peers and
// ABOUTME: clients, propagating blocks and syncing chains.
package quark

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type Server struct {
	Node *Node

	mu    sync.Mutex
	peers map[string]bool
}

func NewServer(node *Node) *Server {
	return &Server{
		Node:  node,
		peers: map[string]bool{},
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /tx", s.handleSubmitTx)
	mux.HandleFunc("POST /send", s.handleSend)
	mux.HandleFunc("POST /block", s.handleReceiveBlock)
	mux.HandleFunc("GET /chain", s.handleGetChain)
	mux.HandleFunc("GET /balance", s.handleGetBalance)
	mux.HandleFunc("GET /address", s.handleGetAddress)
	mux.HandleFunc("POST /mine", s.handleMine)
	mux.HandleFunc("POST /peers", s.handleAddPeer)
	mux.HandleFunc("GET /peers", s.handleListPeers)
	mux.HandleFunc("POST /sync", s.handleSync)
	return mux
}

func (s *Server) AddPeer(url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.peers[url] = true
}

func (s *Server) Peers() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]string, 0, len(s.peers))
	for p := range s.peers {
		out = append(out, p)
	}
	return out
}

func (s *Server) handleSubmitTx(w http.ResponseWriter, r *http.Request) {
	var tx Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := s.Node.SubmitTransaction(&tx); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	go s.broadcastTx(&tx)
	writeJSON(w, http.StatusAccepted, map[string]string{"hash": tx.Hash()})
}

type sendRequest struct {
	Recipient string `json:"recipient"`
	Amount    int64  `json:"amount"`
	Nonce     int64  `json:"nonce"`
}

func (s *Server) handleSend(w http.ResponseWriter, r *http.Request) {
	var req sendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if req.Recipient == "" || req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, errors.New("recipient and positive amount required"))
		return
	}
	tx := NewTransaction(s.Node.Address(), req.Recipient, req.Amount)
	tx.Nonce = req.Nonce
	if err := tx.Sign(s.Node.Miner.Wallet); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if err := s.Node.SubmitTransaction(tx); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	go s.broadcastTx(tx)
	writeJSON(w, http.StatusAccepted, map[string]string{"hash": tx.Hash()})
}

func (s *Server) handleReceiveBlock(w http.ResponseWriter, r *http.Request) {
	var block Block
	if err := json.NewDecoder(r.Body).Decode(&block); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := s.Node.ReceiveBlock(&block); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]int{"length": s.Node.Chain.Length()})
}

func (s *Server) handleGetChain(w http.ResponseWriter, r *http.Request) {
	s.Node.mu.Lock()
	chain := s.Node.Chain
	s.Node.mu.Unlock()
	writeJSON(w, http.StatusOK, chain)
}

func (s *Server) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	addr := r.URL.Query().Get("address")
	if addr == "" {
		writeError(w, http.StatusBadRequest, errors.New("address query parameter required"))
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"balance": s.Node.Balance(addr)})
}

func (s *Server) handleGetAddress(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"address": s.Node.Address()})
}

func (s *Server) handleMine(w http.ResponseWriter, r *http.Request) {
	block, err := s.Node.Mine()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	go s.broadcastBlock(block)
	writeJSON(w, http.StatusOK, block)
}

type peerRequest struct {
	URL string `json:"url"`
}

func (s *Server) handleAddPeer(w http.ResponseWriter, r *http.Request) {
	var req peerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if req.URL == "" {
		writeError(w, http.StatusBadRequest, errors.New("url required"))
		return
	}
	s.AddPeer(req.URL)
	writeJSON(w, http.StatusAccepted, map[string]string{"peer": req.URL})
}

func (s *Server) handleListPeers(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string][]string{"peers": s.Peers()})
}

func (s *Server) handleSync(w http.ResponseWriter, r *http.Request) {
	replaced, err := s.Sync()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"replaced": replaced,
		"length":   s.Node.Chain.Length(),
	})
}

func (s *Server) Sync() (bool, error) {
	bestLen := s.Node.Chain.Length()
	var bestChain *BlockChain

	for _, peer := range s.Peers() {
		chain, err := fetchChain(peer)
		if err != nil {
			continue
		}
		if chain.Length() <= bestLen {
			continue
		}
		if chain.Validate() != nil {
			continue
		}
		bestLen = chain.Length()
		bestChain = chain
	}
	if bestChain == nil {
		return false, nil
	}

	s.Node.mu.Lock()
	defer s.Node.mu.Unlock()
	s.Node.Chain = bestChain
	for _, block := range bestChain.Blocks {
		hashes := make([]string, 0, len(block.Data))
		for _, tx := range block.Data {
			hashes = append(hashes, tx.Hash())
		}
		s.Node.Mempool.Remove(hashes...)
	}
	return true, nil
}

func (s *Server) broadcastBlock(block *Block) {
	body, err := json.Marshal(block)
	if err != nil {
		return
	}
	for _, peer := range s.Peers() {
		_, _ = http.Post(peer+"/block", "application/json", bytes.NewReader(body))
	}
}

func (s *Server) broadcastTx(tx *Transaction) {
	body, err := json.Marshal(tx)
	if err != nil {
		return
	}
	for _, peer := range s.Peers() {
		_, _ = http.Post(peer+"/tx", "application/json", bytes.NewReader(body))
	}
}

func fetchChain(peer string) (*BlockChain, error) {
	resp, err := http.Get(peer + "/chain")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("peer returned %d: %s", resp.StatusCode, string(body))
	}
	var chain BlockChain
	if err := json.NewDecoder(resp.Body).Decode(&chain); err != nil {
		return nil, err
	}
	if chain.Config == nil {
		chain.Config = DefaultDifficultyConfig()
	}
	return &chain, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
