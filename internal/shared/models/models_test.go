package models

import "testing"

func TestLancamentoComandaTableName(t *testing.T) {
	if got := (LancamentoComanda{}).TableName(); got != "lancamentocomanda" {
		t.Fatalf("TableName() = %q", got)
	}
}

func TestLancamentoComandaItemTableName(t *testing.T) {
	if got := (LancamentoComandaItem{}).TableName(); got != "lancamentocomandaitem" {
		t.Fatalf("TableName() = %q", got)
	}
}
