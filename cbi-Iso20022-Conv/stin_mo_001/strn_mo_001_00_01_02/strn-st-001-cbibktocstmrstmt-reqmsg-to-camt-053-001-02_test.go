package strn_mo_001_00_01_02_test

import (
	_ "embed"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-cbi/cbi/strn_mo_001/strn_mo_001_00_01_02/strn_st_001_cbibktocstmrstmt_reqmsg"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi/cbi-Iso20022-Conv/stin_mo_001/strn_mo_001_00_01_02"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"testing"
)

//go:embed example-strn-st-001-cbibktocstmrstmt-reqmsg.xml
var exampleStrnSt001 []byte

const example_pain_002_001_03 = "example-pain.002.001.03.xml"

func TestStrn_St_001_CBIBkToCstmrStmtReqLogMsg_To_Camt_053_001_02_Conv(t *testing.T) {

	strn, err := strn_st_001_cbibktocstmrstmt_reqmsg.NewDocumentFromXML(exampleStrnSt001)
	require.NoError(t, err)

	camt, err := strn_mo_001_00_01_02.Strn_St_001_CBIBkToCstmrStmtReqLogMsg_To_Camt_053_001_02_Conv(strn)
	require.NoError(t, err)

	camtData, err := camt.ToXML()
	require.NoError(t, err)

	err = os.WriteFile(example_pain_002_001_03, camtData, fs.ModePerm)
	require.NoError(t, err)

	defer os.Remove(example_pain_002_001_03)
}
