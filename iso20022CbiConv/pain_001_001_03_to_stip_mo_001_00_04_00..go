package iso20022CbiConv

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-cbi/cbi/stip_mo_001/stip_mo_001_00_04_00"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/pain/001.001.03/pain_001_001_03"
	"github.com/rs/zerolog/log"
	"io/fs"
	"os"
)

func Pain_001_001_03_To_Stip_Mo_001_00_04_00_Conv(in *pain_001_001_03.Document) (*stip_mo_001_00_04_00.Document, error) {

	const semLogContext = "pain_001_001_03_to_stip_mo_001_00_04_00::conv"

	stip := stip_mo_001_00_04_00.Document{
		GrpHdr: in.CstmrCdtTrfInitn.GrpHdr,
		PmtInf: in.CstmrCdtTrfInitn.PmtInf,
	}

	return &stip, nil
}

func Pain_001_001_03_To_Stip_Mo_001_00_04_00_XMLDataConv(painData []byte) ([]byte, error) {

	const semLogContext = "pain_001_001_03_to_stip_mo_001_00_04_00::xml-data-conv"

	pain, err := pain_001_001_03.NewDocumentFromXML(painData)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	stip, err := Pain_001_001_03_To_Stip_Mo_001_00_04_00_Conv(pain)
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

func Pain_001_001_03_To_Stip_Mo_001_00_04_00_XMLFileConv(inFn string, outFn string) error {

	const semLogContext = "pain_001_001_03_to_stip_mo_001_00_04_00::xml-file-conv"

	painData, err := os.ReadFile(inFn)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	stipData, err := Pain_001_001_03_To_Stip_Mo_001_00_04_00_XMLDataConv(painData)
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
