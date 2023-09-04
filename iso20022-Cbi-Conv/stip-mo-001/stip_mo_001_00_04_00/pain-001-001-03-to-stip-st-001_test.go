package stip_mo_001_00_04_00_test

import (
	_ "embed"
	pain_001_001_03_common "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/001.001.03/common"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/001.001.03/pain_001_001_03"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/xsdt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi/iso20022-Cbi-Conv/stip-mo-001/stip_mo_001_00_04_00"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"testing"
)

//go:embed example-pain.001.001.03.xml
var example []byte

const example_stip_st_001 = "example-stip-st-001.xml"

func Pain001_001_03_Adapter(pain *pain_001_001_03.Document) (*pain_001_001_03.Document, error) {
	pain.CstmrCdtTrfInitn.GrpHdr.CtrlSum = xsdt.MustToDecimal(float64(220.0))
	pain.CstmrCdtTrfInitn.GrpHdr.InitgPty.Id = &pain_001_001_03.Party6Choice{
		OrgId: &pain_001_001_03.OrganisationIdentification4{
			Othr: []pain_001_001_03.GenericOrganisationIdentification1{{
				Id:   pain_001_001_03_common.MustToMax35Text(pain_001_001_03_common.Max35TextSample),
				Issr: pain_001_001_03_common.MustToMax35Text(pain_001_001_03_common.Max35TextSample)},
			},
		},
	}

	pain.CstmrCdtTrfInitn.PmtInf[0].DbtrAgt.FinInstnId.ClrSysMmbId =
		&pain_001_001_03.ClearingSystemMemberIdentification2{
			MmbId: pain_001_001_03_common.MustToMax35Text(pain_001_001_03_common.Max35TextSample),
		}

	pain.CstmrCdtTrfInitn.PmtInf[0].CdtTrfTxInf[0].PmtId.InstrId = "1"
	return pain, nil
}

func TestPain_001_001_03_To_Stip_St_001_XMLConv(t *testing.T) {
	adapter := pain_001_001_03.DocumentAdapter(Pain001_001_03_Adapter)
	stipData, err := stip_mo_001_00_04_00.Pain_001_001_03_To_Stip_St_001_XMLDataConv(example, stip_mo_001_00_04_00.WithConvAdapter(adapter))
	require.NoError(t, err)

	err = os.WriteFile(example_stip_st_001, stipData, fs.ModePerm)
	require.NoError(t, err)

	defer os.Remove(example_stip_st_001)
}

func TestPain_001_001_03_To_Stip_St_001_Conv(t *testing.T) {
	adapter := pain_001_001_03.DocumentAdapter(Pain001_001_03_Adapter)
	pain, err := pain_001_001_03.NewDocumentFromXML(example)
	require.NoError(t, err)

	stipObj, err := stip_mo_001_00_04_00.Pain_001_001_03_To_Stip_St_001_Conv(pain, stip_mo_001_00_04_00.WithConvAdapter(adapter))
	require.NoError(t, err)

	stipData, err := stipObj.ToXML()
	require.NoError(t, err)

	err = os.WriteFile(example_stip_st_001, stipData, fs.ModePerm)
	require.NoError(t, err)

	defer os.Remove(example_stip_st_001)
}
