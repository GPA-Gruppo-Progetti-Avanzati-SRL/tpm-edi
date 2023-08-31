package strn_mo_001_00_01_02

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-cbi/cbi/strn_mo_001/strn_mo_001_00_01_02/strn_st_001_cbibktocstmrstmt_reqmsg"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/camt/053.001.02/camt_053_001_02"
	"github.com/rs/zerolog/log"
	"io/fs"
	"os"
)

func Strn_St_001_CBIBkToCstmrStmtReqLogMsg_To_Camt_053_001_02_Conv(in *strn_st_001_cbibktocstmrstmt_reqmsg.Document) (*camt_053_001_02.Document, error) {

	const semLogContext = "strn-st-001-cbibktocstmrstmtreqlogmsg-to-camt-053-001-02::conv"

	pain := camt_053_001_02.Document{
		BkToCstmrStmt: camt_053_001_02.BankToCustomerStatementV02{
			GrpHdr: in.GrpHdr,
			Stmt:   in.Stmt,
		},
	}
	return &pain, nil
}

func Strn_St_001_CBIBkToCstmrStmtReqLogMsg_To_Camt_053_001_02_XMLDataConv(strnData []byte) ([]byte, error) {

	const semLogContext = "strn-st-001-cbibktocstmrstmtreqlogmsg-to-camt-053-001-02::xml-data-conv"

	strn, err := strn_st_001_cbibktocstmrstmt_reqmsg.NewDocumentFromXML(strnData)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	camt, err := Strn_St_001_CBIBkToCstmrStmtReqLogMsg_To_Camt_053_001_02_Conv(strn)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	camtData, err := camt.ToXML()
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	return camtData, nil
}

func Strn_St_001_CBIBkToCstmrStmtReqLogMsg_To_Camt_053_001_02_XMLFileConv(inFn string, outFn string) error {

	const semLogContext = "strn-st-001-cbibktocstmrstmtreqlogmsg-to-camt-053-001-02::xml-file-conv"

	strnData, err := os.ReadFile(inFn)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	camtData, err := Strn_St_001_CBIBkToCstmrStmtReqLogMsg_To_Camt_053_001_02_XMLDataConv(strnData)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	err = os.WriteFile(outFn, camtData, fs.ModePerm)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	return nil
}
