package strn_mo_001_00_01_02

import (
	"fmt"
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

func Strn_St_001_CBIBkToCstmrStmtReqLogMsg_To_Camt_053_001_02_Conv(in *strn_st_001_cbibktocstmrstmt_reqmsg.Document, opts ...ConvOption) ([]*camt_053_001_02.Document, error) {

	const semLogContext = "strn-st-001-cbibktocstmrstmtreqlogmsg-to-camt-053-001-02::conv"

	options := ConvOptions{}
	for _, o := range opts {
		o(&options)
	}

	var camts []*camt_053_001_02.Document
	for _, env := range in.CBIEnvelBkToCstmrStmtReqLogMsg {
		pain := camt_053_001_02.Document{
			BkToCstmrStmt: camt_053_001_02.BankToCustomerStatementV02{
				GrpHdr: env.CBIBkToCstmrStmtReqLogMsg.CBIDlyStmtReqLogMsg.GrpHdr,
				Stmt:   env.CBIBkToCstmrStmtReqLogMsg.CBIDlyStmtReqLogMsg.Stmt,
			},
		}

		var err error
		if options.camt_053_001_02_Adapter != nil {
			_, err = options.camt_053_001_02_Adapter(&pain)
			if err != nil {
				return nil, err
			}
		}

		camts = append(camts, &pain)
	}

	return camts, nil
}

func Strn_St_001_CBIBkToCstmrStmtReqLogMsg_To_Camt_053_001_02_XMLDataConv(strnData []byte, opts ...ConvOption) ([][]byte, error) {

	const semLogContext = "strn-st-001-cbibktocstmrstmtreqlogmsg-to-camt-053-001-02::xml-data-conv"

	strn, err := strn_st_001_cbibktocstmrstmt_reqmsg.NewDocumentFromXML(strnData)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	camts, err := Strn_St_001_CBIBkToCstmrStmtReqLogMsg_To_Camt_053_001_02_Conv(strn, opts...)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	var camtsData [][]byte
	for _, camt := range camts {
		camtData, err := camt.ToXML()
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}

		camtsData = append(camtsData, camtData)
	}

	return camtsData, nil
}

func Strn_St_001_CBIBkToCstmrStmtReqLogMsg_To_Camt_053_001_02_XMLFileConv(inFn string, outFn string, opts ...ConvOption) ([]string, error) {

	const semLogContext = "strn-st-001-cbibktocstmrstmtreqlogmsg-to-camt-053-001-02::xml-file-conv"

	strnData, err := os.ReadFile(inFn)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	camtsData, err := Strn_St_001_CBIBkToCstmrStmtReqLogMsg_To_Camt_053_001_02_XMLDataConv(strnData, opts...)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	var outfiles []string
	for i, camtData := range camtsData {
		outf := fmt.Sprintf(outFn, i)
		err = os.WriteFile(fmt.Sprintf(outFn, i), camtData, fs.ModePerm)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}

		outfiles = append(outfiles, outf)
	}

	return outfiles, nil
}
