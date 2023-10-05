package pain_001_001_09_to_stip_st_001_test

import (
	_ "embed"
	"fmt"
	pain_001_001_09_common "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/001.001.09/common"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/001.001.09/pain_001_001_09"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/xsdt"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi/iso20022-Cbi-Conv/stip-mo-001/stip_mo_001_00_04_01/pain_001_001_09_to_stip_st_001"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"testing"
)

//go:embed example-pain.001.001.03.xml
var examplePain03 []byte

const example_stip_st_001 = "example-stip-st-001_%d.xml"

func Pain001_001_09_Adapter(pain *pain_001_001_09.Document) (*pain_001_001_09.Document, error) {
	pain.CstmrCdtTrfInitn.GrpHdr.CtrlSum = xsdt.MustToDecimal(float64(220.0))
	pain.CstmrCdtTrfInitn.GrpHdr.InitgPty.Id = &pain_001_001_09.Party38Choice{
		OrgId: &pain_001_001_09.OrganisationIdentification29{
			Othr: []pain_001_001_09.GenericOrganisationIdentification1{{
				Id:   pain_001_001_09_common.MustToMax35Text(pain_001_001_09_common.Max35TextSample),
				Issr: pain_001_001_09_common.MustToMax35Text(pain_001_001_09_common.Max35TextSample)},
			},
		},
	}

	for i := range pain.CstmrCdtTrfInitn.PmtInf {
		pain.CstmrCdtTrfInitn.PmtInf[i].DbtrAgt.FinInstnId.ClrSysMmbId =
			&pain_001_001_09.ClearingSystemMemberIdentification2{
				MmbId: pain_001_001_09_common.MustToMax35Text(pain_001_001_09_common.Max35TextSample),
			}

		for j := range pain.CstmrCdtTrfInitn.PmtInf[i].CdtTrfTxInf {
			pain.CstmrCdtTrfInitn.PmtInf[i].CdtTrfTxInf[j].PmtId.InstrId = "1"
		}
	}

	return pain, nil
}

func Test_XMLConv(t *testing.T) {
	adapter := pain_001_001_09.DocumentAdapter(Pain001_001_09_Adapter)
	stipsData, err := pain_001_001_09_to_stip_st_001.XMLDataConv(examplePain03,
		pain_001_001_09_to_stip_st_001.WithInputAdapter(adapter),
		pain_001_001_09_to_stip_st_001.WithOutputAdapter(pain_001_001_09_to_stip_st_001.DefaultOutputAdapter),
		pain_001_001_09_to_stip_st_001.WithAllowPain00100103(true))
	require.NoError(t, err)

	for i, stipData := range stipsData {
		outF := fmt.Sprintf(example_stip_st_001, i)
		err = os.WriteFile(outF, stipData, fs.ModePerm)
		require.NoError(t, err)

		defer os.Remove(outF)
	}
}
