package dnsrewrite_test

import (
	"bytes"
	"testing"

	"example.com/dnsrewrite"
)

func TestBug01_CachedAnswerRDataAlias(t *testing.T) {
	e := dnsrewrite.New()
	defer e.Close()
	data := []byte{10, 0, 0, 1}
	ans := &dnsrewrite.Answer{
		Question: dnsrewrite.Question{Name: "a.local", Type: dnsrewrite.TypeA, Class: 1},
		RCode:    dnsrewrite.RCodeNoError,
		Answers:  []dnsrewrite.RR{{Name: "a.local", Type: dnsrewrite.TypeA, Class: 1, TTL: 30, Data: data}},
	}
	cached, err := e.CacheAnswer("k1", ans)
	if err != nil {
		t.Fatal(err)
	}
	data[0] = 99
	if cached.Answers[0].Data[0] == 99 {
		t.Fatal("returned cache RDATA shared with caller")
	}
	got, ok := e.CachedAnswer("k1")
	if !ok {
		t.Fatal("missing cache")
	}
	if !bytes.Equal(got.Answers[0].Data, []byte{10, 0, 0, 1}) {
		t.Fatalf("cache RDATA aliased: %v", got.Answers[0].Data)
	}
}
