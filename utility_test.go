package quark

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	s := "siamese dreams"
	o := struct {
		Payload string
	}{
		Payload: s,
	}
	p := &o

	b, err := encode(p)
	assert.NoError(t, err)
	d, err := decode(reflect.TypeOf(p), b)
	assert.NoError(t, err)
	assert.Equal(t, p, d)
}

func TestMineHeaderMeetsDifficulty(t *testing.T) {
	d := int32(12)
	prefix := strings.Repeat("0", int(d/4))

	bh := mineHeader("", nil, d)
	assert.True(t, bh.IsValid())
	assert.True(t, strings.HasPrefix(bh.Hash, prefix))
}

func TestMerkleRootDeterministic(t *testing.T) {
	w, err := NewWallet()
	assert.NoError(t, err)

	tx := NewTransaction(w.Address(), "recipient", 10)
	err = tx.Sign(w)
	assert.NoError(t, err)

	r1 := merkleRoot([]*Transaction{tx})
	r2 := merkleRoot([]*Transaction{tx})
	assert.Equal(t, r1, r2)
}

func TestMerkleRootChangesWithTransactions(t *testing.T) {
	w, err := NewWallet()
	assert.NoError(t, err)

	tx1 := NewTransaction(w.Address(), "a", 10)
	err = tx1.Sign(w)
	assert.NoError(t, err)

	tx2 := NewTransaction(w.Address(), "b", 10)
	err = tx2.Sign(w)
	assert.NoError(t, err)

	r1 := merkleRoot([]*Transaction{tx1})
	r2 := merkleRoot([]*Transaction{tx1, tx2})
	assert.NotEqual(t, r1, r2)
}
