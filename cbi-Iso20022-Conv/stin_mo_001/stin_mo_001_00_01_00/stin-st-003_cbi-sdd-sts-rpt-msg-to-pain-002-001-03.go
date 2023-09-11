package stin_mo_001_00_01_00

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-cbi/cbi/stin_mo_001/stin_mo_001_00_01_00/stin_st_003_cbisdd_stsrptmsg"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/002.001.03/pain_002_001_03"
	"github.com/rs/zerolog/log"
	"io/fs"
	"os"
)

func Stin_St_003_CbiSdd_StsRptMsg_To_Pain_002_001_03_Conv(in *stin_st_003_cbisdd_stsrptmsg.Document, opts ...ConvOption) ([]*pain_002_001_03.Document, error) {

	const semLogContext = "stin_st_003_cbi-sdd_stsrptmsg_to_pain_002_001_03::conv"

	options := ConvOptions{}
	for _, o := range opts {
		o(&options)
	}

	var pains []*pain_002_001_03.Document
	for _, env := range in.CBIEnvelSDDStsRptLogMsg {
		pain := pain_002_001_03.Document{
			CstmrPmtStsRpt: &pain_002_001_03.CustomerPaymentStatusReportV03{
				GrpHdr:            env.CBISDDStsRptLogMsg.GrpHdr,
				OrgnlGrpInfAndSts: env.CBISDDStsRptLogMsg.OrgnlGrpInfAndSts,
				OrgnlPmtInfAndSts: env.CBISDDStsRptLogMsg.OrgnlPmtInfAndSts,
			},
		}

		var err error
		if options.pain_002_001_03_Adapter != nil {
			_, err = options.pain_002_001_03_Adapter(&pain)
			if err != nil {
				return nil, err
			}
		}

		pains = append(pains, &pain)
	}

	return pains, nil
}

func Stin_St_003_CbiSdd_StsRptMsg_To_Pain_002_001_03_XMLDataConv(stinData []byte, opts ...ConvOption) ([][]byte, error) {

	const semLogContext = "stin_st_003_cbi-sdd_stsrptmsg_to_pain_002_001_03::xml-data-conv"

	stin, err := stin_st_003_cbisdd_stsrptmsg.NewDocumentFromXML(stinData)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	pains, err := Stin_St_003_CbiSdd_StsRptMsg_To_Pain_002_001_03_Conv(stin, opts...)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	var painsData [][]byte
	for _, pain := range pains {
		painData, err := pain.ToXML()
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}

		painsData = append(painsData, painData)
	}

	return painsData, nil
}

func Stin_St_003_CbiSdd_StsRptMsg_To_Pain_002_001_03_XMLFileConv(inFn string, outFn string, opts ...ConvOption) ([]string, error) {

	const semLogContext = "stin_st_003_cbi-sdd_stsrptmsg_to_pain_002_001_03::xml-file-conv"

	stinData, err := os.ReadFile(inFn)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	painsData, err := Stin_St_003_CbiSdd_StsRptMsg_To_Pain_002_001_03_XMLDataConv(stinData, opts...)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	var outfiles []string
	for i, painData := range painsData {
		outf := fmt.Sprintf(outFn, i)
		err = os.WriteFile(fmt.Sprintf(outFn, i), painData, fs.ModePerm)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}

		outfiles = append(outfiles, outf)
	}

	return outfiles, nil
}
