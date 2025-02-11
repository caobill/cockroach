// Copyright 2021 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/server/serverpb"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/sql"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/catconstants"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondata"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlstats"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlstats/persistedsqlstats/sqlstatsutil"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/cockroachdb/errors"
)

func getTimeFromSeconds(seconds int64) *time.Time {
	if seconds != 0 {
		t := timeutil.Unix(seconds, 0)
		return &t
	}
	return nil
}

func (s *statusServer) CombinedStatementStats(
	ctx context.Context, req *serverpb.CombinedStatementsStatsRequest,
) (*serverpb.StatementsResponse, error) {
	ctx = propagateGatewayMetadata(ctx)
	ctx = s.AnnotateCtx(ctx)

	if err := s.privilegeChecker.requireViewActivityOrViewActivityRedactedPermission(ctx); err != nil {
		return nil, err
	}

	return getCombinedStatementStats(
		ctx,
		req,
		s.sqlServer.pgServer.SQLServer.GetSQLStatsProvider(),
		s.internalExecutor,
		s.st,
		s.sqlServer.execCfg.SQLStatsTestingKnobs)
}

func getCombinedStatementStats(
	ctx context.Context,
	req *serverpb.CombinedStatementsStatsRequest,
	statsProvider sqlstats.Provider,
	ie *sql.InternalExecutor,
	settings *cluster.Settings,
	testingKnobs *sqlstats.TestingKnobs,
) (*serverpb.StatementsResponse, error) {
	startTime := getTimeFromSeconds(req.Start)
	endTime := getTimeFromSeconds(req.End)
	limit := SQLStatsResponseMax.Get(&settings.SV)
	whereClause, orderAndLimit, args := getQueryClausesAndArgs(startTime, endTime, limit, testingKnobs)
	statements, err := collectCombinedStatements(ctx, ie, whereClause, args, orderAndLimit)
	if err != nil {
		return nil, err
	}

	transactions, err := collectCombinedTransactions(ctx, ie, whereClause, args, orderAndLimit)
	if err != nil {
		return nil, err
	}

	response := &serverpb.StatementsResponse{
		Statements:            statements,
		Transactions:          transactions,
		LastReset:             statsProvider.GetLastReset(),
		InternalAppNamePrefix: catconstants.InternalAppNamePrefix,
	}

	return response, nil
}

// getQueryClausesAndArgs returns:
// - where clause (filtering by name and aggregates_ts when defined)
// - order and limit clause
// - args that will replace the clauses above
func getQueryClausesAndArgs(
	start, end *time.Time, limit int64, testingKnobs *sqlstats.TestingKnobs,
) (whereClause string, orderAndLimitClause string, args []interface{}) {
	var buffer strings.Builder
	buffer.WriteString(testingKnobs.GetAOSTClause())

	// Filter out internal statements by app name.
	buffer.WriteString(fmt.Sprintf(" WHERE app_name NOT LIKE '%s%%'", catconstants.InternalAppNamePrefix))

	if start != nil {
		buffer.WriteString(" AND aggregated_ts >= $1")
		args = append(args, *start)
	}

	if end != nil {
		buffer.WriteString(fmt.Sprintf(" AND aggregated_ts <= $%d", len(args)+1))
		args = append(args, *end)
	}

	orderAndLimitClause = fmt.Sprintf(` ORDER BY aggregated_ts DESC LIMIT $%d`, len(args)+1)
	args = append(args, limit)

	return buffer.String(), orderAndLimitClause, args
}

func collectCombinedStatements(
	ctx context.Context,
	ie *sql.InternalExecutor,
	whereClause string,
	qargs []interface{},
	orderAndLimit string,
) ([]serverpb.StatementsResponse_CollectedStatementStatistics, error) {

	query := fmt.Sprintf(
		`SELECT
				fingerprint_id,
				transaction_fingerprint_id,
				app_name,
				aggregated_ts,
				max(metadata) AS metadata,
				crdb_internal.merge_statement_stats(array_agg(statistics)) AS statistics,
				max(sampled_plan) AS sampled_plan,
				aggregation_interval
		FROM crdb_internal.statement_statistics %s
		GROUP BY
				fingerprint_id,
				transaction_fingerprint_id,
				app_name,
				aggregated_ts,
				aggregation_interval
		%s`, whereClause, orderAndLimit)

	const expectedNumDatums = 8

	it, err := ie.QueryIteratorEx(ctx, "combined-stmts-by-interval", nil,
		sessiondata.InternalExecutorOverride{
			User: security.NodeUserName(),
		}, query, qargs...)

	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := it.Close()
		if closeErr != nil {
			err = errors.CombineErrors(err, closeErr)
		}
	}()

	var statements []serverpb.StatementsResponse_CollectedStatementStatistics
	var ok bool
	for ok, err = it.Next(ctx); ok; ok, err = it.Next(ctx) {
		var row tree.Datums
		if row = it.Cur(); row == nil {
			return nil, errors.New("unexpected null row")
		}

		if row.Len() != expectedNumDatums {
			return nil, errors.Newf("expected %d columns, receieved %d", expectedNumDatums)
		}

		var statementFingerprintID uint64
		if statementFingerprintID, err = sqlstatsutil.DatumToUint64(row[0]); err != nil {
			return nil, err
		}

		var transactionFingerprintID uint64
		if transactionFingerprintID, err = sqlstatsutil.DatumToUint64(row[1]); err != nil {
			return nil, err
		}

		app := string(tree.MustBeDString(row[2]))
		aggregatedTs := tree.MustBeDTimestampTZ(row[3]).Time

		var metadata roachpb.CollectedStatementStatistics
		metadataJSON := tree.MustBeDJSON(row[4]).JSON
		if err = sqlstatsutil.DecodeStmtStatsMetadataJSON(metadataJSON, &metadata); err != nil {
			return nil, err
		}

		metadata.Key.App = app
		metadata.Key.TransactionFingerprintID =
			roachpb.TransactionFingerprintID(transactionFingerprintID)

		statsJSON := tree.MustBeDJSON(row[5]).JSON
		if err = sqlstatsutil.DecodeStmtStatsStatisticsJSON(statsJSON, &metadata.Stats); err != nil {
			return nil, err
		}

		planJSON := tree.MustBeDJSON(row[6]).JSON
		plan, err := sqlstatsutil.JSONToExplainTreePlanNode(planJSON)
		if err != nil {
			return nil, err
		}
		metadata.Stats.SensitiveInfo.MostRecentPlanDescription = *plan

		aggInterval := tree.MustBeDInterval(row[7]).Duration

		stmt := serverpb.StatementsResponse_CollectedStatementStatistics{
			Key: serverpb.StatementsResponse_ExtendedStatementStatisticsKey{
				KeyData:             metadata.Key,
				AggregatedTs:        aggregatedTs,
				AggregationInterval: time.Duration(aggInterval.Nanos()),
			},
			ID:    roachpb.StmtFingerprintID(statementFingerprintID),
			Stats: metadata.Stats,
		}

		statements = append(statements, stmt)

	}

	if err != nil {
		return nil, err
	}

	return statements, nil
}

func collectCombinedTransactions(
	ctx context.Context,
	ie *sql.InternalExecutor,
	whereClause string,
	qargs []interface{},
	orderAndLimit string,
) ([]serverpb.StatementsResponse_ExtendedCollectedTransactionStatistics, error) {

	query := fmt.Sprintf(
		`SELECT
				app_name,
				aggregated_ts,
				fingerprint_id,
				metadata,
				statistics,
				aggregation_interval
			FROM crdb_internal.transaction_statistics
			%s %s`, whereClause, orderAndLimit)

	const expectedNumDatums = 6

	it, err := ie.QueryIteratorEx(ctx, "combined-txns-by-interval", nil,
		sessiondata.InternalExecutorOverride{
			User: security.NodeUserName(),
		}, query, qargs...)

	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := it.Close()
		if closeErr != nil {
			err = errors.CombineErrors(err, closeErr)
		}
	}()

	var transactions []serverpb.StatementsResponse_ExtendedCollectedTransactionStatistics
	var ok bool
	for ok, err = it.Next(ctx); ok; ok, err = it.Next(ctx) {
		var row tree.Datums
		if row = it.Cur(); row == nil {
			return nil, errors.New("unexpected null row")
		}

		if row.Len() != expectedNumDatums {
			return nil, errors.Newf("expected %d columns, receieved %d", expectedNumDatums, row.Len())
		}

		app := string(tree.MustBeDString(row[0]))
		aggregatedTs := tree.MustBeDTimestampTZ(row[1]).Time
		fingerprintID, err := sqlstatsutil.DatumToUint64(row[2])
		if err != nil {
			return nil, err
		}

		var metadata roachpb.CollectedTransactionStatistics
		metadataJSON := tree.MustBeDJSON(row[3]).JSON
		if err = sqlstatsutil.DecodeTxnStatsMetadataJSON(metadataJSON, &metadata); err != nil {
			return nil, err
		}

		statsJSON := tree.MustBeDJSON(row[4]).JSON
		if err = sqlstatsutil.DecodeTxnStatsStatisticsJSON(statsJSON, &metadata.Stats); err != nil {
			return nil, err
		}

		aggInterval := tree.MustBeDInterval(row[5]).Duration

		txnStats := serverpb.StatementsResponse_ExtendedCollectedTransactionStatistics{
			StatsData: roachpb.CollectedTransactionStatistics{
				StatementFingerprintIDs:  metadata.StatementFingerprintIDs,
				App:                      app,
				Stats:                    metadata.Stats,
				AggregatedTs:             aggregatedTs,
				AggregationInterval:      time.Duration(aggInterval.Nanos()),
				TransactionFingerprintID: roachpb.TransactionFingerprintID(fingerprintID),
			},
		}

		transactions = append(transactions, txnStats)
	}

	if err != nil {
		return nil, err
	}

	return transactions, nil
}
