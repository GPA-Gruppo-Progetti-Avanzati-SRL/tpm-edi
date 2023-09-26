package stin_st_003_cbi_sdd_sts_rpt_msg_to_pain_002_001_03_test

import (
	_ "embed"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-cbi/cbi/stin_mo_001/stin_mo_001_00_01_00/stin_st_003_cbisdd_stsrptmsg"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi/cbi-Iso20022-Conv/stin_mo_001/stin_mo_001_00_01_00/stin_st_003_cbi_sdd_sts_rpt_msg_to_pain_002_001_03"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"testing"
)

//go:embed example_stin_st_003_cbisdd_stsrptmsg.xml
var exampleStinSt003 []byte

const example_pain_002_001_03 = "example-pain.002.001.03_%d.xml"

func Test_Conv(t *testing.T) {

	stin, err := stin_st_003_cbisdd_stsrptmsg.NewDocumentFromXML(exampleStinSt003)
	require.NoError(t, err)

	pains, err := stin_st_003_cbi_sdd_sts_rpt_msg_to_pain_002_001_03.Conv(stin)
	require.NoError(t, err)

	for i, pain := range pains {
		painData, err := pain.ToXML()
		require.NoError(t, err)

		outF := fmt.Sprintf(example_pain_002_001_03, i)
		err = os.WriteFile(outF, painData, fs.ModePerm)
		require.NoError(t, err)

		defer os.Remove(outF)
	}
}
