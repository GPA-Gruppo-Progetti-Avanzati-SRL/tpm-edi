package stin_mo_001_00_01_00

import (
	stin_st_002_cbisdd_techvalstsmsg "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-cbi/cbi/stin_mo_001/stin_mo_001_00_01_00/stin_st_002_cbisdd_techvalstsmsg"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/002.001.03/pain_002_001_03"
	"github.com/rs/zerolog/log"
	"io/fs"
	"os"
)

type ConvOptions struct {
	pain_002_001_03_Adapter pain_002_001_03.DocumentAdapter
}

type ConvOption func(opts *ConvOptions)

func WithConvAdapter(adapter pain_002_001_03.DocumentAdapter) ConvOption {
	return func(opts *ConvOptions) {
		opts.pain_002_001_03_Adapter = adapter
	}
}

func Stin_St_002_CbiSdd_TechValStsMsg_To_Pain_002_001_03_Conv(in *stin_st_002_cbisdd_techvalstsmsg.Document, opts ...ConvOption) (*pain_002_001_03.Document, error) {

	const semLogContext = "stin-st-002_cbi-sdd-tech-val-sts-msg-to-pain-002-001-03::conv"

	options := ConvOptions{}
	for _, o := range opts {
		o(&options)
	}

	pain := pain_002_001_03.Document{
		CstmrPmtStsRpt: pain_002_001_03.CustomerPaymentStatusReportV03{
			GrpHdr:            in.CBIEnvelSDDTechValStsLogMsg.CBISDDTechValStsLogMsg[0].GrpHdr,
			OrgnlGrpInfAndSts: in.CBIEnvelSDDTechValStsLogMsg.CBISDDTechValStsLogMsg[0].OrgnlGrpInfAndSts,
			OrgnlPmtInfAndSts: in.CBIEnvelSDDTechValStsLogMsg.CBISDDTechValStsLogMsg[0].OrgnlPmtInfAndSts,
		},
	}

	var err error
	if options.pain_002_001_03_Adapter != nil {
		_, err = options.pain_002_001_03_Adapter(&pain)
		if err != nil {
			return nil, err
		}
	}

	return &pain, nil
}

func Stin_St_002_CbiSdd_TechValStsMsg_To_Pain_002_001_03_XMLDataConv(stinData []byte, opts ...ConvOption) ([]byte, error) {

	const semLogContext = "stin-st-002_cbi-sdd-tech-val-sts-msg-to-pain-002-001-03::xml-data-conv"

	stin, err := stin_st_002_cbisdd_techvalstsmsg.NewDocumentFromXML(stinData)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	pain, err := Stin_St_002_CbiSdd_TechValStsMsg_To_Pain_002_001_03_Conv(stin, opts...)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	painData, err := pain.ToXML()
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	return painData, nil
}

func Stin_St_002_CbiSdd_TechValStsMsg_To_Pain_002_001_03_XMLFileConv(inFn string, outFn string, opts ...ConvOption) error {

	const semLogContext = "stin-st-002_cbi-sdd-tech-val-sts-msg-to-pain-002-001-03::xml-file-conv"

	stinData, err := os.ReadFile(inFn)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	painData, err := Stin_St_002_CbiSdd_TechValStsMsg_To_Pain_002_001_03_XMLDataConv(stinData, opts...)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	err = os.WriteFile(outFn, painData, fs.ModePerm)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	return nil
}
