package pain_008_001_02_to_stin_st_001_cbisdd_reqmsg

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-cbi/cbi/stin_mo_001/stin_mo_001_00_01_00/stin_st_001_cbisdd_reqmsg"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/008.001.02/pain_008_001_02"
	"github.com/rs/zerolog/log"
	"io/fs"
	"os"
)

type ConvOptions struct {
	inputAdapter  pain_008_001_02.DocumentAdapter
	outputAdapter stin_st_001_cbisdd_reqmsg.DocumentAdapter
}

type ConvOption func(opts *ConvOptions)

func WithInputAdapter(adapter pain_008_001_02.DocumentAdapter) ConvOption {
	return func(opts *ConvOptions) {
		opts.inputAdapter = adapter
	}
}

func WithOutputAdapter(adapter stin_st_001_cbisdd_reqmsg.DocumentAdapter) ConvOption {
	return func(opts *ConvOptions) {
		opts.outputAdapter = adapter
	}
}

func Conv(in *pain_008_001_02.Document, opts ...ConvOption) (*stin_st_001_cbisdd_reqmsg.Document, error) {

	const semLogContext = "pain-008-001-02-to-stin-st-001_cbi-sdd-req-msg::conv"

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

	stin := &stin_st_001_cbisdd_reqmsg.Document{
		GrpHdr: in.CstmrDrctDbtInitn.GrpHdr,
		PmtInf: in.CstmrDrctDbtInitn.PmtInf,
	}

	if options.outputAdapter != nil {
		stin, err = options.outputAdapter(stin)
		if err != nil {
			return nil, err
		}
	}

	return stin, nil
}

func XMLDataConv(painData []byte, opts ...ConvOption) ([]byte, error) {

	const semLogContext = "pain-008-001-02-to-stin-st-001_cbi-sdd-req-msg::xml-data-conv"

	pain, err := pain_008_001_02.NewDocumentFromXML(painData)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	stip, err := Conv(pain, opts...)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	stipData, err := stip.ToXML()
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	return stipData, nil
}

func XMLFileConv(inFn string, outFn string, opts ...ConvOption) error {

	const semLogContext = "pain-008-001-02-to-stin-st-001_cbi-sdd-req-msg::xml-file-conv"

	painData, err := os.ReadFile(inFn)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	stipData, err := XMLDataConv(painData, opts...)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	err = os.WriteFile(outFn, stipData, fs.ModePerm)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	return nil
}
