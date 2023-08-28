package iso20022CbiConv

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-cbi/cbi/stin_mo_001/stin_mo_001_00_01_00"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/008.001.02/pain_008_001_02"
	"github.com/rs/zerolog/log"
	"io/fs"
	"os"
)

func Pain_008_001_02_To_Stin_Mo_001_00_01_00_Conv(in *pain_008_001_02.Document) (*stin_mo_001_00_01_00.Document, error) {

	const semLogContext = "pain_008_001_02_to_stin_mo_001_00_01_00::conv"

	stip := stin_mo_001_00_01_00.Document{
		GrpHdr: in.CstmrDrctDbtInitn.GrpHdr,
		PmtInf: in.CstmrDrctDbtInitn.PmtInf,
	}

	return &stip, nil
}

func Pain_008_001_02_To_Stin_Mo_001_00_01_00_XMLDataConv(painData []byte) ([]byte, error) {

	const semLogContext = "pain_008_001_02_to_stin_mo_001_00_01_00::xml-data-conv"

	pain, err := pain_008_001_02.NewDocumentFromXML(painData)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	stip, err := Pain_008_001_02_To_Stin_Mo_001_00_01_00_Conv(pain)
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

func Pain_008_001_02_To_Stin_Mo_001_00_01_00_XMLFileConv(inFn string, outFn string) error {

	const semLogContext = "pain_008_001_02_to_stin_mo_001_00_01_00::xml-file-conv"

	painData, err := os.ReadFile(inFn)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	stipData, err := Pain_008_001_02_To_Stin_Mo_001_00_01_00_XMLDataConv(painData)
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
