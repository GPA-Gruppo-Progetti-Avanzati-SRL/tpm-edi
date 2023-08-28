package iso20022CbiConv_test

import (
	_ "embed"
	pain_001_001_03_common "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/001.001.03/common"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/001.001.03/pain_001_001_03"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/xsdt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi/iso20022CbiConv"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"testing"
)

//go:embed example-pain.001.001.03.xml
var example []byte

const example_stip_mo_001_00_04 = "example-document-cbi-stip-mo-001-00-04.xml"

func TestPain_001_001_03_To_Stip_Mo_001_00_04_00_XMLConv(t *testing.T) {
	stipData, err := iso20022CbiConv.Pain_001_001_03_To_Stip_Mo_001_00_04_00_XMLDataConv(example)
	require.NoError(t, err)

	err = os.WriteFile(example_stip_mo_001_00_04, stipData, fs.ModePerm)
	require.NoError(t, err)

	defer os.Remove(example_stip_mo_001_00_04)
}

func TestPain_001_001_03_To_Stip_Mo_001_00_04_00_Conv(t *testing.T) {

	pain, err := pain_001_001_03.NewDocumentFromXML(example)
	require.NoError(t, err)

	stipObj, err := iso20022CbiConv.Pain_001_001_03_To_Stip_Mo_001_00_04_00_Conv(pain)
	require.NoError(t, err)

	stipObj.GrpHdr.CtrlSum = xsdt.MustToDecimal(float64(220.0))
	stipObj.GrpHdr.InitgPty.Id = &pain_001_001_03.Party6Choice{
		OrgId: &pain_001_001_03.OrganisationIdentification4{
			Othr: []pain_001_001_03.GenericOrganisationIdentification1{{
				Id:   pain_001_001_03_common.MustToMax35Text(pain_001_001_03_common.Max35TextSample),
				Issr: pain_001_001_03_common.MustToMax35Text(pain_001_001_03_common.Max35TextSample)},
			},
		},
	}

	stipObj.PmtInf[0].DbtrAgt.FinInstnId.ClrSysMmbId =
		&pain_001_001_03.ClearingSystemMemberIdentification2{
			MmbId: pain_001_001_03_common.MustToMax35Text(pain_001_001_03_common.Max35TextSample),
		}

	stipObj.PmtInf[0].CdtTrfTxInf[0].PmtId.InstrId = "1"

	stipData, err := stipObj.ToXML()
	require.NoError(t, err)

	err = os.WriteFile(example_stip_mo_001_00_04, stipData, fs.ModePerm)
	require.NoError(t, err)

	defer os.Remove(example_stip_mo_001_00_04)
}
