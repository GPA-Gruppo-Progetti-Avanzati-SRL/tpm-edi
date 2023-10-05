package pain_001_001_09_to_stip_st_001

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-cbi/cbi/stip_mo_001/stip_mo_001_00_04_01/stip_st_001"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/001.001.09/common"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/001.001.09/pain_001_001_09"
	"github.com/rs/zerolog/log"
	"io/fs"
	"os"
	"strings"
)

type ConvOptions struct {
	inputAdapter      pain_001_001_09.DocumentAdapter
	outputAdapter     stip_st_001.DocumentsAdapter
	allowPain00100103 bool
}

type ConvOption func(opts *ConvOptions)

func WithInputAdapter(adapter pain_001_001_09.DocumentAdapter) ConvOption {
	return func(opts *ConvOptions) {
		opts.inputAdapter = adapter
	}
}

func WithOutputAdapter(adapter stip_st_001.DocumentsAdapter) ConvOption {
	return func(opts *ConvOptions) {
		opts.outputAdapter = adapter
	}
}

// WithAllowPain00100103 inspected by XMLDataConv to make the pain.001.001.03 look like the pain.001.001.09
func WithAllowPain00100103(b bool) ConvOption {
	return func(opts *ConvOptions) {
		opts.allowPain00100103 = b
	}
}

func DefaultOutputAdapter(in []*stip_st_001.Document) ([]*stip_st_001.Document, error) {

	const semLogContext = "stip_mo_001_00_04_01::stips-st-001-default-adapter"

	var err error
	if len(in) < 2 {
		return in, nil
	}

	for i, stip := range in {
		stip.GrpHdr.NbOfTxs = "1"

		var msgId string
		if len(stip.GrpHdr.MsgId) <= 32 {
			msgId = fmt.Sprintf("%s_%02d", stip.GrpHdr.MsgId, i)
		} else {
			msgId = fmt.Sprintf("%s_%02d", stip.GrpHdr.MsgId[:32], i)
		}
		stip.GrpHdr.MsgId, err = common.ToMax35Text(msgId)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}
	}

	return in, nil
}

func Conv(in *pain_001_001_09.Document, opts ...ConvOption) ([]*stip_st_001.Document, error) {

	const semLogContext = "pain-001-001-09-to-stip-st-001::conv"

	options := ConvOptions{}
	for _, o := range opts {
		o(&options)
	}

	var err error
	if options.inputAdapter != nil {
		in, err = options.inputAdapter(in)
		if err != nil {
			return nil, err
		}
	}

	var stips []*stip_st_001.Document
	if len(in.CstmrCdtTrfInitn.PmtInf) == 1 {
		stip := &stip_st_001.Document{
			GrpHdr: in.CstmrCdtTrfInitn.GrpHdr,
			PmtInf: in.CstmrCdtTrfInitn.PmtInf[0],
		}

		stips = append(stips, stip)

	} else {
		for _, pmtinf := range in.CstmrCdtTrfInitn.PmtInf {
			hdr := *in.CstmrCdtTrfInitn.GrpHdr
			stip := &stip_st_001.Document{
				GrpHdr: &hdr,
				PmtInf: pmtinf,
			}

			stips = append(stips, stip)
		}
	}

	if options.outputAdapter != nil {
		stips, err = options.outputAdapter(stips)
		if err != nil {
			return nil, err
		}
	}

	return stips, nil
}

func XMLDataConv(painData []byte, opts ...ConvOption) ([][]byte, error) {

	const semLogContext = "pain-001-001-09-to-stip-st-001::xml-data-conv"

	options := ConvOptions{}
	for _, o := range opts {
		o(&options)
	}

	if options.allowPain00100103 {
		ndx := strings.Index(string(painData), "urn:iso:std:iso:20022:tech:xsd:pain.001.001.03")
		if ndx >= 0 {
			painData = ([]byte)(strings.Replace(string(painData), "urn:iso:std:iso:20022:tech:xsd:pain.001.001.03", "urn:iso:std:iso:20022:tech:xsd:pain.001.001.09", -1))
		}
	}

	pain, err := pain_001_001_09.NewDocumentFromXML(painData)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	stips, err := Conv(pain, opts...)
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

func XMLFileConv(inFn string, outFn string, opts ...ConvOption) ([]string, error) {

	const semLogContext = "pain-001-001-09-to-stip-st-001::xml-file-conv"

	painData, err := os.ReadFile(inFn)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	stipsData, err := XMLDataConv(painData, opts...)
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
