package pain_001_001_03_to_stip_st_001_test

import (
	_ "embed"
	"fmt"
	pain_001_001_03_common "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/001.001.03/common"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/001.001.03/pain_001_001_03"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/xsdt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi/iso20022-Cbi-Conv/stip-mo-001/stip_mo_001_00_04_00/pain_001_001_03_to_stip_st_001"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"testing"
)

//go:embed example-pain.001.001.03.xml
var example []byte

const example_stip_st_001 = "example-stip-st-001_%d.xml"

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

	for i := range pain.CstmrCdtTrfInitn.PmtInf {
		pain.CstmrCdtTrfInitn.PmtInf[i].DbtrAgt.FinInstnId.ClrSysMmbId =
			&pain_001_001_03.ClearingSystemMemberIdentification2{
				MmbId: pain_001_001_03_common.MustToMax35Text(pain_001_001_03_common.Max35TextSample),
			}

		for j := range pain.CstmrCdtTrfInitn.PmtInf[i].CdtTrfTxInf {
			pain.CstmrCdtTrfInitn.PmtInf[i].CdtTrfTxInf[j].PmtId.InstrId = "1"
		}
	}

	return pain, nil
}

func TestPain_001_001_03_To_Stip_St_001_XMLConv(t *testing.T) {
	adapter := pain_001_001_03.DocumentAdapter(Pain001_001_03_Adapter)
	stipsData, err := pain_001_001_03_to_stip_st_001.XMLDataConv(example,
		pain_001_001_03_to_stip_st_001.WithInputAdapter(adapter),
		pain_001_001_03_to_stip_st_001.WithOutputAdapter(pain_001_001_03_to_stip_st_001.DefaultOutputAdapter))
	require.NoError(t, err)

	for i, stipData := range stipsData {
		outF := fmt.Sprintf(example_stip_st_001, i)
		err = os.WriteFile(outF, stipData, fs.ModePerm)
		require.NoError(t, err)

		defer os.Remove(outF)
	}
}

func TestPain_001_001_03_To_Stip_St_001_Conv(t *testing.T) {
	adapter := pain_001_001_03.DocumentAdapter(Pain001_001_03_Adapter)
	pain, err := pain_001_001_03.NewDocumentFromXML(example)
	require.NoError(t, err)

	stipsObj, err := pain_001_001_03_to_stip_st_001.Conv(pain, pain_001_001_03_to_stip_st_001.WithInputAdapter(adapter))
	require.NoError(t, err)

	for i, stipObj := range stipsObj {
		stipData, err := stipObj.ToXML()
		require.NoError(t, err)

		outF := fmt.Sprintf(example_stip_st_001, i)
		err = os.WriteFile(outF, stipData, fs.ModePerm)
		require.NoError(t, err)

		defer os.Remove(outF)
	}
}
