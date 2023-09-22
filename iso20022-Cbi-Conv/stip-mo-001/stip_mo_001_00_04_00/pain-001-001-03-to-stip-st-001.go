package stip_mo_001_00_04_00

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-cbi/cbi/stip_mo_001/stip_mo_001_00_04_00/stip_st_001"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/001.001.03/common"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/001.001.03/pain_001_001_03"
	"github.com/rs/zerolog/log"
	"io/fs"
	"os"
)

type ConvOptions struct {
	pain_001_001_03_Adapter pain_001_001_03.DocumentAdapter
}

type ConvOption func(opts *ConvOptions)

func WithConvAdapter(adapter pain_001_001_03.DocumentAdapter) ConvOption {
	return func(opts *ConvOptions) {
		opts.pain_001_001_03_Adapter = adapter
	}
}

func Pain_001_001_03_To_Stip_St_001_Conv(in *pain_001_001_03.Document, opts ...ConvOption) ([]*stip_st_001.Document, error) {

	const semLogContext = "pain-001-001-03-to-stip-st-001::conv"

	options := ConvOptions{}
	for _, o := range opts {
		o(&options)
	}

	var err error
	if options.pain_001_001_03_Adapter != nil {
		in, err = options.pain_001_001_03_Adapter(in)
		if err != nil {
			return nil, err
		}
	}

	var stips []*stip_st_001.Document
	if len(in.CstmrCdtTrfInitn.PmtInf) == 1 {
		stip := stip_st_001.Document{
			GrpHdr: in.CstmrCdtTrfInitn.GrpHdr,
			PmtInf: in.CstmrCdtTrfInitn.PmtInf[0],
		}

		stips = append(stips, &stip)
	} else {
		for i, pmtinf := range in.CstmrCdtTrfInitn.PmtInf {

			hdr := *in.CstmrCdtTrfInitn.GrpHdr
			hdr.NbOfTxs = "1"

			var msgId string
			if len(hdr.MsgId) <= 32 {
				msgId = fmt.Sprintf("%s_%02d", hdr.MsgId, i)
			} else {
				msgId = fmt.Sprintf("%s_%02d", hdr.MsgId[:32], i)
			}
			hdr.MsgId, err = common.ToMax35Text(msgId)
			if err != nil {
				log.Error().Err(err).Msg(semLogContext)
			}

			stip := stip_st_001.Document{
				GrpHdr: &hdr,
				PmtInf: pmtinf,
			}

			stips = append(stips, &stip)
		}
	}

	return stips, nil
}

func Pain_001_001_03_To_Stip_St_001_XMLDataConv(painData []byte, opts ...ConvOption) ([][]byte, error) {

	const semLogContext = "pain-001-001-03-to-stip-st-001::xml-data-conv"

	pain, err := pain_001_001_03.NewDocumentFromXML(painData)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	stips, err := Pain_001_001_03_To_Stip_St_001_Conv(pain, opts...)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	var stipsData [][]byte
	for _, stip := range stips {
		stipData, err := stip.ToXML()
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}

		stipsData = append(stipsData, stipData)
	}

	return stipsData, nil
}

func Pain_001_001_03_To_Stip_St_001_XMLFileConv(inFn string, outFn string, opts ...ConvOption) ([]string, error) {

	const semLogContext = "pain-001-001-03-to-stip-st-001::xml-file-conv"

	painData, err := os.ReadFile(inFn)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	stipsData, err := Pain_001_001_03_To_Stip_St_001_XMLDataConv(painData, opts...)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	var outfiles []string
	for i, stipData := range stipsData {
		outf := fmt.Sprintf(outFn, i)
		err = os.WriteFile(fmt.Sprintf(outFn, i), stipData, fs.ModePerm)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}

		outfiles = append(outfiles, outf)
	}

	return outfiles, nil
}
