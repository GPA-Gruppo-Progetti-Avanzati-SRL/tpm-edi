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

type ConvOptions struct {
	camt_053_001_02_Adapter camt_053_001_02.DocumentAdapter
}

type ConvOption func(opts *ConvOptions)

func WithConvAdapter(adapter camt_053_001_02.DocumentAdapter) ConvOption {
	return func(opts *ConvOptions) {
		opts.camt_053_001_02_Adapter = adapter
	}
}

func Rnd_001_06_05_Rh_To_Camt_053_001_02_DataConv(rhData []byte, opts ...ConvOption) ([]camt_053_001_02.Document, error) {

	options := ConvOptions{}
	for _, o := range opts {
		o(&options)
	}

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

			setOptions := []camt_053_001_02.SetOpOption{camt_053_001_02.SetOpWithLog(true), camt_053_001_02.SetOpWithSkipIfEmpty(true)}
			doc.BkToCstmrStmt.Stmt = append(doc.BkToCstmrStmt.Stmt, camt_053_001_02.AccountStatement2{})
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Id,
				strings.Join([]string{
					rhBofRec.Get(rnd_rh.RecordType),
					rhBofRec.Get(rnd_rh.Sender),
					rhBofRec.Get(rnd_rh.Recipient),
					rhBofRec.Get(rnd_rh.SupportName)}, "/"), setOptions...)

			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_ElctrncSeqNb, rhStmtRec.OpeningBalance.Get(rnd_rh.ProgrNumber), setOptions...)
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_CreDtTm, rhBofRec.Get(rnd_rh.CreationDate), setOptions...)
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_CpyDplctInd,
				rhStmtRec.OpeningBalance.Get(rnd_rh.Reason,
					reader.WithDefaultValue("COPY"),
					reader.WithValueMappings([]reader.KeyValue{{"93001", "COPY"}, {reader.GetPropertyOtherwiseMappingKey, "COPY"}}),
				), setOptions...)
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Acct_Id_IBAN,
				strings.Join([]string{
					rhStmtRec.OpeningBalance.Get(rnd_rh.CountryCode),
					rhStmtRec.OpeningBalance.Get(rnd_rh.Cin),
					rhStmtRec.OpeningBalance.Get(rnd_rh.BankAbi),
					rhStmtRec.OpeningBalance.Get(rnd_rh.BankCab),
					rhStmtRec.OpeningBalance.Get(rnd_rh.CurrentAccountCode),
					rhStmtRec.OpeningBalance.Get(rnd_rh.CheckDigit)}, ""),
				setOptions...)

			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Acct_Tp_Cd, "CACC", setOptions...)
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Acct_Tp_Prtry, rhStmtRec.OpeningBalance.Get(rnd_rh.AccountType), setOptions...)
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Acct_Ccy, rhStmtRec.OpeningBalance.Get(rnd_rh.CurrencyCode), setOptions...)
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Acct_Nm, rhStmtRec.OpeningBalance.Get(rnd_rh.Descr), setOptions...)
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Acct_Svcr_FinInstnId_ClrSysMmbId_MmbId, rhStmtRec.OpeningBalance.Get(rnd_rh.BankAbi), setOptions...)

			doc.BkToCstmrStmt.Stmt[len(doc.BkToCstmrStmt.Stmt)-1].Bal = append(doc.BkToCstmrStmt.Stmt[len(doc.BkToCstmrStmt.Stmt)-1].Bal, camt_053_001_02.CashBalance3{})
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Bal_Tp_CdOrPrtry_Cd, "OPBD", setOptions...) // 61
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Bal_Amt_Value, rhStmtRec.OpeningBalance.Get(rnd_rh.OpeningBalance), setOptions...)
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Bal_CdtDbtInd,
				rhStmtRec.OpeningBalance.Get(rnd_rh.Sign, reader.WithValueMappings([]reader.KeyValue{{"D", "DBIT"}, {"C", "CRDT"}, {reader.GetPropertyOtherwiseMappingKey, "ERR"}})), setOptions...)
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Bal_Dt_Dt, rhStmtRec.OpeningBalance.Get(rnd_rh.AccountingDate), setOptions...)

			doc.BkToCstmrStmt.Stmt[len(doc.BkToCstmrStmt.Stmt)-1].Bal = append(doc.BkToCstmrStmt.Stmt[len(doc.BkToCstmrStmt.Stmt)-1].Bal, camt_053_001_02.CashBalance3{})
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Bal_Tp_CdOrPrtry_Cd, "CLBD", setOptions...) // 64
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Bal_Amt_Value, rhStmtRec.ClosingBalance.Get(rnd_rh.AccountsBalance), setOptions...)
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Bal_Amt_Ccy, rhStmtRec.ClosingBalance.Get(rnd_rh.CurrencyCode), setOptions...)
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Bal_CdtDbtInd,
				rhStmtRec.ClosingBalance.Get(rnd_rh.AccountsBalanceSign, reader.WithValueMappings([]reader.KeyValue{{"D", "DBIT"}, {"C", "CRDT"}, {reader.GetPropertyOtherwiseMappingKey, "ERR"}})), setOptions...)
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Bal_Dt_Dt, rhStmtRec.ClosingBalance.Get(rnd_rh.AccountingDate), setOptions...)

			doc.BkToCstmrStmt.Stmt[len(doc.BkToCstmrStmt.Stmt)-1].Bal = append(doc.BkToCstmrStmt.Stmt[len(doc.BkToCstmrStmt.Stmt)-1].Bal, camt_053_001_02.CashBalance3{})
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Bal_Tp_CdOrPrtry_Cd, "CLAV", setOptions...) // 64
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Bal_Amt_Value, rhStmtRec.ClosingBalance.Get(rnd_rh.CashBalance), setOptions...)
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Bal_Amt_Ccy, rhStmtRec.ClosingBalance.Get(rnd_rh.CurrencyCode), setOptions...)
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Bal_CdtDbtInd,
				rhStmtRec.ClosingBalance.Get(rnd_rh.CashBalanceSign, reader.WithValueMappings([]reader.KeyValue{{"D", "DBIT"}, {"C", "CRDT"}, {reader.GetPropertyOtherwiseMappingKey, "ERR"}})), setOptions...)
			_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Bal_Dt_Dt, rhStmtRec.ClosingBalance.Get(rnd_rh.AccountingDate), setOptions...)

			if !rhStmtRec.ExpCashOnHand.IsEmpty() {
				doc.BkToCstmrStmt.Stmt[len(doc.BkToCstmrStmt.Stmt)-1].Bal = append(doc.BkToCstmrStmt.Stmt[len(doc.BkToCstmrStmt.Stmt)-1].Bal, camt_053_001_02.CashBalance3{})
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Bal_Tp_CdOrPrtry_Cd, "FWAV", setOptions...) // 65
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Bal_Amt_Value, rhStmtRec.ExpCashOnHand.Get(rnd_rh.FirstCashBalance), setOptions...)
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Bal_CdtDbtInd,
					rhStmtRec.ExpCashOnHand.Get(rnd_rh.FirstCashSign, reader.WithValueMappings([]reader.KeyValue{{"D", "DBIT"}, {"C", "CRDT"}, {reader.GetPropertyOtherwiseMappingKey, "ERR"}})), setOptions...)
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Bal_Dt_Dt, rhStmtRec.ExpCashOnHand.Get(rnd_rh.FifthCashOnHandDate), setOptions...)
			}

			for _, m := range rhStmtRec.Movements {
				doc.BkToCstmrStmt.Stmt[len(doc.BkToCstmrStmt.Stmt)-1].Ntry = append(doc.BkToCstmrStmt.Stmt[len(doc.BkToCstmrStmt.Stmt)-1].Ntry, camt_053_001_02.ReportEntry2{})
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_NtryRef, m.Get(rnd_rh.MovmntProgrNumber), setOptions...)
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_Amt_Value, m.Get(rnd_rh.Movmntamount), setOptions...)
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_CdtDbtInd,
					m.Get(rnd_rh.MovmntSign, reader.WithValueMappings([]reader.KeyValue{{"D", "DBIT"}, {"C", "CRDT"}, {reader.GetPropertyOtherwiseMappingKey, "ERR"}})), setOptions...)
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_Sts, "BOOK")
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_BookgDt_Dt, m.Get(rnd_rh.AccountingDate), setOptions...)
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_ValDt_Dt, m.Get(rnd_rh.ValueDate), setOptions...)
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_AcctSvcrRef, m.Get(rnd_rh.BankRef), setOptions...)
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_BkTxCd_Prtry_Cd, strings.Join([]string{m.Get(rnd_rh.CbiReason), m.Get(rnd_rh.InternalReason)}, "/"), setOptions...)

				zz1 := m.FirstByStructureFlag("ZZ1")
				// Note the value is sort of extra field because it is a chardata value
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_AmtDtls_InstdAmt_Amt_Value, zz1.Get(rnd_rh.OrigAmnt), setOptions...)
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_AmtDtls_InstdAmt_Amt_Ccy, zz1.Get(rnd_rh.OrigAmntCurrencyCode), setOptions...)
				// Missing Exchange rate... where to put?
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_AmtDtls_TxAmt_Amt_Value, zz1.Get(rnd_rh.PaidAmnt), setOptions...)
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_AmtDtls_TxAmt_Amt_Ccy, zz1.Get(rnd_rh.PaidAmntCurrencyCode), setOptions...)
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_AmtDtls_CntrValAmt_Amt_Value, zz1.Get(rnd_rh.TrxAmnt), setOptions...)
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_AmtDtls_CntrValAmt_Amt_Value, zz1.Get(rnd_rh.TrxAmntCurrencyCode), setOptions...)

				_ = doc.Set(camt_053_001_02.MustSetArrayItemPathModifiers(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_Chrgs_Amt_Value, []string{"", "", "+"}), zz1.Get(rnd_rh.CommissionAmnt), setOptions...)
				_ = doc.Set(camt_053_001_02.MustSetArrayItemPathModifiers(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_Chrgs_Amt_Value, []string{"", "", "+"}), zz1.Get(rnd_rh.CommissionFeesAmnt), setOptions...)

				id1 := m.FirstByStructureFlag("ID1")
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_NtryDtls_Btch_MsgId, id1.Get(rnd_rh.MsgId), setOptions...)
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_NtryDtls_Btch_PmtInfId, m.Get(rnd_rh.CustRefMovmntDescr), setOptions...)
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_NtryDtls_TxDtls_Refs_MsgId, id1.Get(rnd_rh.MsgId), setOptions...)
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_NtryDtls_TxDtls_Refs_PmtInfId, id1.Get(rnd_rh.CustRefMovmntDescr), setOptions...)
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_NtryDtls_TxDtls_Refs_EndToEndId, id1.Get(rnd_rh.End2EndId), setOptions...)
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_NtryDtls_TxDtls_Refs_ChqNb, m.Get(rnd_rh.ChequeNumber), setOptions...)

				yyy := m.FirstByStructureFlag("YYY")
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_NtryDtls_TxDtls_RltdPties_InitgPty_Id_OrgId_Othr_Id, yyy.Get(rnd_rh.OrderingPrtyTaxpayerCode), setOptions...)

				// Dipende dalla causale cbi row. 176
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_NtryDtls_TxDtls_RltdPties_Dbtr_Id_OrgId_Othr_Id, strings.TrimSuffix(strings.Join([]string{zz1.Get(rnd_rh.OrderingPrtyDescr), zz1.Get(rnd_rh.Country)}, "/"), "/"), setOptions...)
				// Dipende dalla causale cbi row. 176
				zz2 := m.FirstByStructureFlag("ZZ2")
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_NtryDtls_TxDtls_RltdPties_Dbtr_Id_OrgId_Othr_Id, zz2.Get(rnd_rh.OrderingPrty), setOptions...)

				zz3 := m.FirstByStructureFlag("ZZ3")
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_NtryDtls_TxDtls_RltdPties_Cdtr_Id_OrgId_Othr_Id, zz3.Get(rnd_rh.Payee), setOptions...)

				ri1 := m.FirstByStructureFlag("RI1")
				ri2 := m.FirstByStructureFlag("RI2")
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_NtryDtls_TxDtls_RmtInf_Ustrd, strings.TrimSuffix(strings.Join([]string{ri1.Get(rnd_rh.ReconcData), ri2.Get(rnd_rh.ReconcData)}, "/"), "/"), setOptions...)

				zz4 := m.FirstByStructureFlag("/ZZ4/")
				_ = doc.Set(camt_053_001_02.Path_BkToCstmrStmt_Stmt_Ntry_NtryDtls_TxDtls_AddtlTxInf, strings.Join([]string{zz3.Get(rnd_rh.PaymentReason), strings.TrimSuffix(zz4.Get(rnd_rh.PaymentReason), "/ZZ4/"), m.AdditionalInfo()}, " "), setOptions...)
			}

			fmt.Println(rhStmtRec.String())
		}

		rhEofRec, err = r.ReadBatchOfStatementsEORecord()
		if err != nil && err != io.EOF {
			return nil, err
		}

		fmt.Println(rhEofRec.String())

		if options.camt_053_001_02_Adapter != nil {
			_, err = options.camt_053_001_02_Adapter(&doc)
			if err != nil {
				return nil, err
			}
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

func logDocumentSetError(err error) {
	if err != nil {
		log.Error().Err(err).Msg("error in setting property of iso20022 message")
	}
}
