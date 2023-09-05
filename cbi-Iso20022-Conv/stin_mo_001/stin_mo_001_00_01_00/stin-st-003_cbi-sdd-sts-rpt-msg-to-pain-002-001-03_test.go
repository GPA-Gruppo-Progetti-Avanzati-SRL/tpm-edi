package stin_mo_001_00_01_00_test

import (
	_ "embed"
	"fmt"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-cbi/cbi/stin_mo_001/stin_mo_001_00_01_00/stin_st_003_cbisdd_stsrptmsg"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi/cbi-Iso20022-Conv/stin_mo_001/stin_mo_001_00_01_00"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"testing"
)

//go:embed example_stin_st_003_cbisdd_stsrptmsg.xml
var exampleStinSt003 []byte

func Test_Stin003CbiSdd_StsRptMsg_To_Pain_002_001_03_Conv(t *testing.T) {

	stin, err := stin_st_003_cbisdd_stsrptmsg.NewDocumentFromXML(exampleStinSt003)
	require.NoError(t, err)

	pains, err := stin_mo_001_00_01_00.Stin_St_003_CbiSdd_StsRptMsg_To_Pain_002_001_03_Conv(stin)
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
