package strn_st_001_cbibktocstmrstmt_reqmsg_to_camt_053_001_02_test

import (
	_ "embed"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-cbi/cbi/strn_mo_001/strn_mo_001_00_01_02/strn_st_001_cbibktocstmrstmt_reqmsg"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi/cbi-Iso20022-Conv/strn_mo_001/strn_mo_001_00_01_02/strn_st_001_cbibktocstmrstmt_reqmsg_to_camt_053_001_02"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"testing"
)

//go:embed example-in-strn-st-001-cbi-bdy_bk-to-cstmr-stmt-req.xml
var exampleStrnSt001 []byte

const example_camt_053_001_02 = "example-camt_053_001_02_%d.xml"

func Test_Conv(t *testing.T) {

	strn, err := strn_st_001_cbibktocstmrstmt_reqmsg.NewDocumentFromXML(exampleStrnSt001)
	require.NoError(t, err)

	camts, err := strn_st_001_cbibktocstmrstmt_reqmsg_to_camt_053_001_02.Conv(strn)
	require.NoError(t, err)

	for i, camt := range camts {
		camtData, err := camt.ToXML()
		require.NoError(t, err)

		outF := fmt.Sprintf(example_camt_053_001_02, i)
		err = os.WriteFile(outF, camtData, fs.ModePerm)
		require.NoError(t, err)

		defer os.Remove(outF)
	}

}
