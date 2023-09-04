package strn_mo_001_00_01_02

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-cbi/cbi/strn_mo_001/strn_mo_001_00_01_02/strn_st_001_cbibktocstmrstmt_reqmsg"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/camt/053.001.02/camt_053_001_02"
	"github.com/rs/zerolog/log"
	"io/fs"
	"os"
)

type ConvOptions struct {
	camt_053_001_02_Adapter camt_053_001_02.DocumentAdapter
}

type ConvOption func(opts *ConvOptions)

func WithConvAdapter(adapter camt_053_001_02.DocumentAdapter) ConvOption {
	return func(opts *ConvOptions) {
		opts.camt_053_001_02_Adapter = adapter
	}
}

func Strn_St_001_CBIBkToCstmrStmtReqLogMsg_To_Camt_053_001_02_Conv(in *strn_st_001_cbibktocstmrstmt_reqmsg.Document, opts ...ConvOption) (*camt_053_001_02.Document, error) {

	const semLogContext = "strn-st-001-cbibktocstmrstmtreqlogmsg-to-camt-053-001-02::conv"

	options := ConvOptions{}
	for _, o := range opts {
		o(&options)
	}

	pain := camt_053_001_02.Document{
		BkToCstmrStmt: camt_053_001_02.BankToCustomerStatementV02{
			GrpHdr: in.GrpHdr,
			Stmt:   in.Stmt,
		},
	}

	var err error
	if options.camt_053_001_02_Adapter != nil {
		_, err = options.camt_053_001_02_Adapter(&pain)
		if err != nil {
			return nil, err
		}
	}

	return &pain, nil
}

func Strn_St_001_CBIBkToCstmrStmtReqLogMsg_To_Camt_053_001_02_XMLDataConv(strnData []byte, opts ...ConvOption) ([]byte, error) {

	const semLogContext = "strn-st-001-cbibktocstmrstmtreqlogmsg-to-camt-053-001-02::xml-data-conv"

	strn, err := strn_st_001_cbibktocstmrstmt_reqmsg.NewDocumentFromXML(strnData)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	camt, err := Strn_St_001_CBIBkToCstmrStmtReqLogMsg_To_Camt_053_001_02_Conv(strn, opts...)
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

func Strn_St_001_CBIBkToCstmrStmtReqLogMsg_To_Camt_053_001_02_XMLFileConv(inFn string, outFn string, opts ...ConvOption) error {

	const semLogContext = "strn-st-001-cbibktocstmrstmtreqlogmsg-to-camt-053-001-02::xml-file-conv"

	strnData, err := os.ReadFile(inFn)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	camtData, err := Strn_St_001_CBIBkToCstmrStmtReqLogMsg_To_Camt_053_001_02_XMLDataConv(strnData, opts...)
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
