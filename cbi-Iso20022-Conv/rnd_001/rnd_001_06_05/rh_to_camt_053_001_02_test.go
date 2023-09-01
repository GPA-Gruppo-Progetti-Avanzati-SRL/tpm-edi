package rnd_001_06_05_test

import (
	_ "embed"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi/cbi-Iso20022-Conv/rnd_001/rnd_001_06_05"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

//go:embed example-cbi-rnd-001-06-05.txt
var example []byte

func TestRnd_001_06_05_Rh_To_Camt_053_001_02_DataConv(t *testing.T) {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	camts, err := rnd_001_06_05.Rnd_001_06_05_Rh_To_Camt_053_001_02_DataConv(example)
	require.NoError(t, err)

	for _, camt := range camts {
		xmlCamt, err := camt.ToXML()
		require.NoError(t, err)

		t.Log(string(xmlCamt))
	}

}
