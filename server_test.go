package quark

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestServer(t *testing.T) (*Server, *httptest.Server) {
	t.Helper()
	node, err := NewNode()
	assert.NoError(t, err)
	srv := NewServer(node)
	ts := httptest.NewServer(srv.Handler())
	t.Cleanup(ts.Close)
	return srv, ts
}

func httpJSON(t *testing.T, method, url string, body any) *http.Response {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		assert.NoError(t, json.NewEncoder(&buf).Encode(body))
	}
	req, err := http.NewRequest(method, url, &buf)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	return resp
}

func TestServerMineAndBalance(t *testing.T) {
	srv, ts := newTestServer(t)

	resp := httpJSON(t, http.MethodPost, ts.URL+"/mine", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	resp = httpJSON(t, http.MethodGet, ts.URL+"/balance?address="+srv.Node.Address(), nil)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var out map[string]int64
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&out))
	assert.Equal(t, MiningReward, out["balance"])
}

func TestServerSubmitTx(t *testing.T) {
	srv, ts := newTestServer(t)
	_, err := srv.Node.Mine()
	assert.NoError(t, err)

	tx := NewTransaction(srv.Node.Address(), "recipient", 10)
	tx.Nonce = 1
	assert.NoError(t, tx.Sign(srv.Node.Miner.Wallet))

	resp := httpJSON(t, http.MethodPost, ts.URL+"/tx", tx)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
	assert.Equal(t, 1, srv.Node.Mempool.Len())
}

func TestServerSendAndMine(t *testing.T) {
	srv, ts := newTestServer(t)
	_, err := srv.Node.Mine()
	assert.NoError(t, err)

	body := sendRequest{Recipient: "recipient", Amount: 20, Nonce: 1}
	resp := httpJSON(t, http.MethodPost, ts.URL+"/send", body)
	resp.Body.Close()
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
	assert.Equal(t, 1, srv.Node.Mempool.Len())

	resp = httpJSON(t, http.MethodPost, ts.URL+"/mine", nil)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int64(20), srv.Node.Balance("recipient"))
}

func TestServerGetChain(t *testing.T) {
	srv, ts := newTestServer(t)
	_, err := srv.Node.Mine()
	assert.NoError(t, err)

	resp := httpJSON(t, http.MethodGet, ts.URL+"/chain", nil)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var chain BlockChain
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&chain))
	assert.Equal(t, srv.Node.Chain.Length(), chain.Length())
}

func TestServerAddPeerAndList(t *testing.T) {
	_, ts := newTestServer(t)

	resp := httpJSON(t, http.MethodPost, ts.URL+"/peers", peerRequest{URL: "http://other"})
	resp.Body.Close()
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)

	resp = httpJSON(t, http.MethodGet, ts.URL+"/peers", nil)
	defer resp.Body.Close()
	var out map[string][]string
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&out))
	assert.Equal(t, []string{"http://other"}, out["peers"])
}

func TestServerRejectsInvalidTx(t *testing.T) {
	_, ts := newTestServer(t)

	bad := NewTransaction("nobody", "recipient", 10)
	resp := httpJSON(t, http.MethodPost, ts.URL+"/tx", bad)
	resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
