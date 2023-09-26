package pain_008_001_02_to_stin_st_001_cbisdd_reqmsg_test

import (
	_ "embed"
	pain_008_001_02_common "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/008.001.02/common"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/008.001.02/pain_008_001_02"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/xsdt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi/iso20022-Cbi-Conv/stin-mo-001/stin_mo_001_00_01_00/pain_008_001_02_to_stin_st_001_cbisdd_reqmsg"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"testing"
)

//go:embed example-pain.008.001.02.xml
var examplePain00800102 []byte

const example_stin_st_001_cbisdd_reqmsg = "example-stin_st_001_cbisdd_reqmsg.xml"

func Pain008_001_02_Adapter(pain *pain_008_001_02.Document) (*pain_008_001_02.Document, error) {
	pain.CstmrDrctDbtInitn.GrpHdr.CtrlSum = xsdt.MustToDecimal(float64(220.0))
	pain.CstmrDrctDbtInitn.GrpHdr.InitgPty.Id = &pain_008_001_02.Party6Choice{
		OrgId: &pain_008_001_02.OrganisationIdentification4{
			Othr: []pain_008_001_02.GenericOrganisationIdentification1{{
				Id:   pain_008_001_02_common.MustToMax35Text(pain_008_001_02_common.Max35TextSample),
				Issr: pain_008_001_02_common.MustToMax35Text(pain_008_001_02_common.Max35TextSample)},
			},
		},
	}

	pain.CstmrDrctDbtInitn.PmtInf[0].CdtrAgt.FinInstnId.ClrSysMmbId =
		&pain_008_001_02.ClearingSystemMemberIdentification2{
			MmbId: pain_008_001_02_common.MustToMax35Text(pain_008_001_02_common.Max35TextSample),
		}

	pain.CstmrDrctDbtInitn.PmtInf[0].DrctDbtTxInf[0].PmtId.InstrId = "1"
	pain.CstmrDrctDbtInitn.PmtInf[0].CdtrAgt.FinInstnId.BIC = ""
	pain.CstmrDrctDbtInitn.PmtInf[0].DrctDbtTxInf[0].DbtrAgt.FinInstnId.BIC = "UNCRITMM"
	pain.CstmrDrctDbtInitn.PmtInf[0].DrctDbtTxInf[0].DbtrAgt.FinInstnId.Othr = nil
	return pain, nil
}

func Test_XMLConv(t *testing.T) {

	adapter := pain_008_001_02.DocumentAdapter(Pain008_001_02_Adapter)
	stinData, err := pain_008_001_02_to_stin_st_001_cbisdd_reqmsg.XMLDataConv(examplePain00800102, pain_008_001_02_to_stin_st_001_cbisdd_reqmsg.WithInputAdapter(adapter))
	require.NoError(t, err)

	err = os.WriteFile(example_stin_st_001_cbisdd_reqmsg, stinData, fs.ModePerm)
	require.NoError(t, err)

	defer os.Remove(example_stin_st_001_cbisdd_reqmsg)
}

func Test_Conv(t *testing.T) {

	adapter := pain_008_001_02.DocumentAdapter(Pain008_001_02_Adapter)

	pain, err := pain_008_001_02.NewDocumentFromXML(examplePain00800102)
	require.NoError(t, err)

	stinObj, err := pain_008_001_02_to_stin_st_001_cbisdd_reqmsg.Conv(pain, pain_008_001_02_to_stin_st_001_cbisdd_reqmsg.WithInputAdapter(adapter))
	require.NoError(t, err)

	/*
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

		stinObj.PmtInf[0].DrctDbtTxInf[0].PmtId.InstrId = "1"
		stinObj.PmtInf[0].CdtrAgt.FinInstnId.BIC = ""
		stinObj.PmtInf[0].DrctDbtTxInf[0].DbtrAgt.FinInstnId.BIC = "UNCRITMM"
		stinObj.PmtInf[0].DrctDbtTxInf[0].DbtrAgt.FinInstnId.Othr = nil
	*/

	stinData, err := stinObj.ToXML()
	require.NoError(t, err)

	err = os.WriteFile(example_stin_st_001_cbisdd_reqmsg, stinData, fs.ModePerm)
	require.NoError(t, err)

	defer os.Remove(example_stin_st_001_cbisdd_reqmsg)
}
