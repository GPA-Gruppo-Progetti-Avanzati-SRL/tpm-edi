package iso20022CbiConv_test

import (
	_ "embed"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi/iso20022CbiConv"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"testing"
)

//go:embed example-pain.001.001.03.xml
var example []byte

const example_stip_mo_001_00_04 = "example-document-cbi-stip-mo-001-00-04.xml"

func TestPain_001_001_03_To_Stip_Mo_001_00_04_00_Conv(t *testing.T) {
	stipData, err := iso20022CbiConv.Pain_001_001_03_To_Stip_Mo_001_00_04_00_XMLDataConv(example)
	require.NoError(t, err)

	err = os.WriteFile(example_stip_mo_001_00_04, stipData, fs.ModePerm)
	require.NoError(t, err)

	defer os.Remove(example_stip_mo_001_00_04)
}
