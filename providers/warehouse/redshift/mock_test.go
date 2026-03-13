package redshift

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata/types"
)

// mockDataAPIClient implements dataAPIClient for testing.
type mockDataAPIClient struct {
	executeErr    error
	describeStatus types.StatusString
	describeErr   error
	resultColumns []types.ColumnMetadata
	resultRecords [][]types.Field
	resultErr     error
	tables        []types.TableMember
	tablesErr     error
	describeCols  []types.ColumnMetadata
	describeTableErr error

	executedSQL   string // captures the SQL from ExecuteStatement
	stmtCounter   int
}

func (m *mockDataAPIClient) ExecuteStatement(ctx context.Context, params *redshiftdata.ExecuteStatementInput, optFns ...func(*redshiftdata.Options)) (*redshiftdata.ExecuteStatementOutput, error) {
	if m.executeErr != nil {
		return nil, m.executeErr
	}
	m.executedSQL = aws.ToString(params.Sql)
	m.stmtCounter++
	return &redshiftdata.ExecuteStatementOutput{
		Id: aws.String(fmt.Sprintf("stmt-%d", m.stmtCounter)),
	}, nil
}

func (m *mockDataAPIClient) DescribeStatement(ctx context.Context, params *redshiftdata.DescribeStatementInput, optFns ...func(*redshiftdata.Options)) (*redshiftdata.DescribeStatementOutput, error) {
	if m.describeErr != nil {
		return nil, m.describeErr
	}
	status := m.describeStatus
	if status == "" {
		status = types.StatusStringFinished
	}
	out := &redshiftdata.DescribeStatementOutput{
		Status: status,
	}
	if status == types.StatusStringFailed {
		out.Error = aws.String("mock query failed")
	}
	return out, nil
}

func (m *mockDataAPIClient) GetStatementResult(ctx context.Context, params *redshiftdata.GetStatementResultInput, optFns ...func(*redshiftdata.Options)) (*redshiftdata.GetStatementResultOutput, error) {
	if m.resultErr != nil {
		return nil, m.resultErr
	}
	return &redshiftdata.GetStatementResultOutput{
		ColumnMetadata: m.resultColumns,
		Records:        m.resultRecords,
	}, nil
}

func (m *mockDataAPIClient) ListTables(ctx context.Context, params *redshiftdata.ListTablesInput, optFns ...func(*redshiftdata.Options)) (*redshiftdata.ListTablesOutput, error) {
	if m.tablesErr != nil {
		return nil, m.tablesErr
	}
	return &redshiftdata.ListTablesOutput{
		Tables: m.tables,
	}, nil
}

func (m *mockDataAPIClient) DescribeTable(ctx context.Context, params *redshiftdata.DescribeTableInput, optFns ...func(*redshiftdata.Options)) (*redshiftdata.DescribeTableOutput, error) {
	if m.describeTableErr != nil {
		return nil, m.describeTableErr
	}
	return &redshiftdata.DescribeTableOutput{
		ColumnList: m.describeCols,
	}, nil
}
