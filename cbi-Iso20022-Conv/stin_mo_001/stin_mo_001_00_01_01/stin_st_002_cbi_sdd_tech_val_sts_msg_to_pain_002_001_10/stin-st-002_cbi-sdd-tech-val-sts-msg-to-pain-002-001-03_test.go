package stin_st_002_cbi_sdd_tech_val_sts_msg_to_pain_002_001_10_test

import (
	_ "embed"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-cbi/cbi/stin_mo_001/stin_mo_001_00_01_01/stin_st_002_cbisdd_techvalstsmsg"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi/cbi-Iso20022-Conv/stin_mo_001/stin_mo_001_00_01_01/stin_st_002_cbi_sdd_tech_val_sts_msg_to_pain_002_001_10"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"testing"
)

//go:embed example-stin-st-002-cbisdd-techvalstsmsg.xml
var exampleStinSt002 []byte

const example_pain_002_001_03 = "example-pain.002.001.03_%d.xml"

func Test_Conv(t *testing.T) {

	stin, err := stin_st_002_cbisdd_techvalstsmsg.NewDocumentFromXML(exampleStinSt002)
	require.NoError(t, err)

	pains, err := stin_st_002_cbi_sdd_tech_val_sts_msg_to_pain_002_001_10.Conv(stin)
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
