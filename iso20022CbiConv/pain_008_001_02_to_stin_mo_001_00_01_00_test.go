package iso20022CbiConv_test

import (
	_ "embed"
	pain_008_001_02_common "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/008.001.02/common"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/008.001.02/pain_008_001_02"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/xsdt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi/iso20022CbiConv"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"testing"
)

//go:embed example-pain.008.001.02.xml
var examplePain00800102 []byte

const example_stin_mo_001_00_01 = "example-document-cbi-stin-mo-001-00-01.xml"

func TestPain_008_001_02_To_Stin_Mo_001_00_01_00_XMLConv(t *testing.T) {
	stinData, err := iso20022CbiConv.Pain_008_001_02_To_Stin_Mo_001_00_01_00_XMLDataConv(examplePain00800102)
	require.NoError(t, err)

	err = os.WriteFile(example_stin_mo_001_00_01, stinData, fs.ModePerm)
	require.NoError(t, err)

	defer os.Remove(example_stin_mo_001_00_01)
}

func TestPain_008_001_02_To_Stin_Mo_001_00_01_00_Conv(t *testing.T) {

	pain, err := pain_008_001_02.NewDocumentFromXML(examplePain00800102)
	require.NoError(t, err)

	stinObj, err := iso20022CbiConv.Pain_008_001_02_To_Stin_Mo_001_00_01_00_Conv(pain)
	require.NoError(t, err)

	stinObj.GrpHdr.CtrlSum = xsdt.MustToDecimal(float64(220.0))
	stinObj.GrpHdr.InitgPty.Id = &pain_008_001_02.Party6Choice{
		OrgId: &pain_008_001_02.OrganisationIdentification4{
			Othr: []pain_008_001_02.GenericOrganisationIdentification1{{
				Id:   pain_008_001_02_common.MustToMax35Text(pain_008_001_02_common.Max35TextSample),
				Issr: pain_008_001_02_common.MustToMax35Text(pain_008_001_02_common.Max35TextSample)},
			},
		},
	}

	stinObj.PmtInf[0].CdtrAgt.FinInstnId.ClrSysMmbId =
		&pain_008_001_02.ClearingSystemMemberIdentification2{
			MmbId: pain_008_001_02_common.MustToMax35Text(pain_008_001_02_common.Max35TextSample),
		}

	stinData, err := stinObj.ToXML()
	require.NoError(t, err)

	err = os.WriteFile(example_stin_mo_001_00_01, stinData, fs.ModePerm)
	require.NoError(t, err)

	defer os.Remove(example_stin_mo_001_00_01)
}
