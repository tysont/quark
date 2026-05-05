package quark

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func waitFor(t *testing.T, cond func() bool, msg string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timeout waiting for: %s", msg)
}

func setupTwoNodes(t *testing.T) (*Server, *httptest.Server, *Server, *httptest.Server) {
	t.Helper()
	srvA, tsA := newTestServer(t)
	srvB, tsB := newTestServer(t)
	srvA.AddPeer(tsB.URL)
	srvB.AddPeer(tsA.URL)
	return srvA, tsA, srvB, tsB
}

func TestE2EMinePropagatesToPeer(t *testing.T) {
	srvA, tsA, srvB, _ := setupTwoNodes(t)

	resp := httpJSON(t, http.MethodPost, tsA.URL+"/mine", nil)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	waitFor(t, func() bool { return srvB.Node.Chain.Length() == 2 }, "B receives A's block")
	assert.Equal(t, MiningReward, srvB.Node.Balance(srvA.Node.Address()))
}

func TestE2ETxPropagatesToPeer(t *testing.T) {
	srvA, tsA, srvB, _ := setupTwoNodes(t)

	resp := httpJSON(t, http.MethodPost, tsA.URL+"/mine", nil)
	resp.Body.Close()
	waitFor(t, func() bool { return srvB.Node.Chain.Length() == 2 }, "B sees A's first block")

	tx := NewTransaction(srvA.Node.Address(), "recipient", 10)
	tx.Nonce = 1
	assert.NoError(t, tx.Sign(srvA.Node.Miner.Wallet))

	resp = httpJSON(t, http.MethodPost, tsA.URL+"/tx", tx)
	resp.Body.Close()
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)

	waitFor(t, func() bool { return srvB.Node.Mempool.Len() == 1 }, "B receives A's tx")
}

func TestE2ESyncReplacesShorterChain(t *testing.T) {
	srvA, tsA := newTestServer(t)
	srvB, _ := newTestServer(t)

	for i := 0; i < 3; i++ {
		_, err := srvA.Node.Mine()
		assert.NoError(t, err)
	}
	_, err := srvB.Node.Mine()
	assert.NoError(t, err)

	srvB.AddPeer(tsA.URL)
	replaced, err := srvB.Sync()
	assert.NoError(t, err)
	assert.True(t, replaced)
	assert.Equal(t, srvA.Node.Chain.Length(), srvB.Node.Chain.Length())
	assert.Equal(t, MiningReward*3, srvB.Node.Balance(srvA.Node.Address()))
	assert.True(t, srvB.Node.Chain.IsValid())
}

func TestE2ESyncIgnoresShorterPeer(t *testing.T) {
	srvA, tsA := newTestServer(t)
	srvB, _ := newTestServer(t)

	_, err := srvA.Node.Mine()
	assert.NoError(t, err)
	for i := 0; i < 3; i++ {
		_, err := srvB.Node.Mine()
		assert.NoError(t, err)
	}

	srvB.AddPeer(tsA.URL)
	replaced, err := srvB.Sync()
	assert.NoError(t, err)
	assert.False(t, replaced)
	assert.Equal(t, 4, srvB.Node.Chain.Length())
}

func TestE2EPersistRestartReplay(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "node.json")

	srvA, tsA := newTestServer(t)
	for i := 0; i < 2; i++ {
		_, err := srvA.Node.Mine()
		assert.NoError(t, err)
	}
	assert.NoError(t, srvA.Node.Save(path))

	loaded, err := LoadNode(path)
	assert.NoError(t, err)
	srvB := NewServer(loaded)
	tsB := httptest.NewServer(srvB.Handler())
	defer tsB.Close()

	resp := httpJSON(t, http.MethodGet, tsB.URL+"/chain", nil)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, srvA.Node.Chain.Length(), loaded.Chain.Length())

	// A has not been touched
	assert.Equal(t, srvA.Node.Address(), loaded.Address())
	assert.Equal(t, srvA.Node.Balance(srvA.Node.Address()),
		loaded.Balance(srvA.Node.Address()))
	_ = tsA
}

func TestE2ETwoNodesMineAndConverge(t *testing.T) {
	srvA, tsA, srvB, tsB := setupTwoNodes(t)

	resp := httpJSON(t, http.MethodPost, tsA.URL+"/mine", nil)
	resp.Body.Close()
	waitFor(t, func() bool { return srvB.Node.Chain.Length() == 2 }, "B sees A block 1")

	resp = httpJSON(t, http.MethodPost, tsB.URL+"/mine", nil)
	resp.Body.Close()
	waitFor(t, func() bool { return srvA.Node.Chain.Length() == 3 }, "A sees B block 2")

	resp = httpJSON(t, http.MethodPost, tsA.URL+"/mine", nil)
	resp.Body.Close()
	waitFor(t, func() bool { return srvB.Node.Chain.Length() == 4 }, "B sees A block 3")

	assert.Equal(t, srvA.Node.Chain.Length(), srvB.Node.Chain.Length())
	assert.Equal(t, MiningReward*2, srvA.Node.Balance(srvA.Node.Address()))
	assert.Equal(t, MiningReward, srvA.Node.Balance(srvB.Node.Address()))
}
