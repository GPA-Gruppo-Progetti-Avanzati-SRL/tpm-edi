package rnd_001_06_05

import (
	_ "embed"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fixedlengthfile/reader"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-cbi/cbi/rnd_001/rnd_001_06_05/rnd_rh"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-edi-iso20022/iso-20022/messages/camt/053.001.02/camt_053_001_02"
	"github.com/rs/zerolog/log"
	"io"
	"strings"
)

func Rnd_001_06_05_Rh_To_Camt_053_001_02_DataConv(rhData []byte) ([]camt_053_001_02.Document, error) {

	r, err := rnd_rh.NewReader(rhData)
	if err != nil {
		return nil, err
	}

	var docs []camt_053_001_02.Document

	var rhBofRec reader.Record
	var rhEofRec reader.Record
	for err == nil {
		rhBofRec, err = r.ReadBatchOfStatementsBORecord()
		if err != nil && err != io.EOF {
			return nil, err
		}

		if rhBofRec.IsEmpty() {
			return docs, nil
		}

		doc := camt_053_001_02.NewDocument()
		_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_GrpHdr_MsgId, rhBofRec.Get(rnd_rh.SupportName), camt_053_001_02.SetOpWithLog(true))
		_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_GrpHdr_CreDtTm, rhBofRec.Get(rnd_rh.CreationDate), camt_053_001_02.SetOpWithLog(true))
		_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_GrpHdr_MsgRcpt_Id_OrgId_Othr_Id, rhBofRec.Get(rnd_rh.Recipient), camt_053_001_02.SetOpWithLog(true))
		_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_GrpHdr_AddtlInf, rhBofRec.Get(rnd_rh.Sender), camt_053_001_02.SetOpWithLog(true))
		var rhStmtRec rnd_rh.Statement
		for err == nil {
			rhStmtRec, err = r.ReadStatement()
			if err != nil && err != io.EOF {
				return nil, err
			}

			if rhStmtRec.IsEmpty() {
				break
			}

			doc.BkToCstmrStmt.Stmt = append(doc.BkToCstmrStmt.Stmt, camt_053_001_02.AccountStatement2{})
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Id,
				strings.Join([]string{
					rhBofRec.Get(rnd_rh.RecordType),
					rhBofRec.Get(rnd_rh.Sender),
					rhBofRec.Get(rnd_rh.Recipient),
					rhBofRec.Get(rnd_rh.SupportName)}, "/"), camt_053_001_02.SetOpWithLog(true))

			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_ElctrncSeqNb, rhStmtRec.OpeningBalance.Get(rnd_rh.ProgrNumber), camt_053_001_02.SetOpWithLog(true))
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_CreDtTm, rhBofRec.Get(rnd_rh.CreationDate), camt_053_001_02.SetOpWithLog(true))
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_CpyDplctInd,
				rhStmtRec.OpeningBalance.Get(rnd_rh.Reason,
					reader.WithDefaultValue("COPY"),
					reader.WithValueMappings([]reader.KeyValue{{"93001", "COPY"}, {reader.GetPropertyOtherwiseMappingKey, "COPY"}}),
				), camt_053_001_02.SetOpWithLog(true))
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Acct_Id_IBAN,
				strings.Join([]string{
					rhStmtRec.OpeningBalance.Get(rnd_rh.CountryCode),
					rhStmtRec.OpeningBalance.Get(rnd_rh.Cin),
					rhStmtRec.OpeningBalance.Get(rnd_rh.BankAbi),
					rhStmtRec.OpeningBalance.Get(rnd_rh.BankCab),
					rhStmtRec.OpeningBalance.Get(rnd_rh.CurrentAccountCode),
					rhStmtRec.OpeningBalance.Get(rnd_rh.CheckDigit)}, ""),
				camt_053_001_02.SetOpWithLog(true))

			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Acct_Tp_Prtry, rhStmtRec.OpeningBalance.Get(rnd_rh.AccountType, reader.WithDefaultValue("NOTPROVIDED")), camt_053_001_02.SetOpWithLog(true))

			fmt.Println(rhStmtRec.String())
		}

		rhEofRec, err = r.ReadBatchOfStatementsEORecord()
		if err != nil && err != io.EOF {
			return nil, err
		}

		fmt.Println(rhEofRec.String())
		docs = append(docs, doc)
	}

	return docs, nil
}

func logDocumentSetError(err error) {
	if err != nil {
		log.Error().Err(err).Msg("error in setting property of iso20022 message")
	}
}
