package redshift

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/redshiftdata"
)

// dataAPIClient abstracts the Redshift Data API for testing.
// The real implementation is *redshiftdata.Client.
type dataAPIClient interface {
	ExecuteStatement(ctx context.Context, params *redshiftdata.ExecuteStatementInput, optFns ...func(*redshiftdata.Options)) (*redshiftdata.ExecuteStatementOutput, error)
	DescribeStatement(ctx context.Context, params *redshiftdata.DescribeStatementInput, optFns ...func(*redshiftdata.Options)) (*redshiftdata.DescribeStatementOutput, error)
	GetStatementResult(ctx context.Context, params *redshiftdata.GetStatementResultInput, optFns ...func(*redshiftdata.Options)) (*redshiftdata.GetStatementResultOutput, error)
	ListTables(ctx context.Context, params *redshiftdata.ListTablesInput, optFns ...func(*redshiftdata.Options)) (*redshiftdata.ListTablesOutput, error)
	DescribeTable(ctx context.Context, params *redshiftdata.DescribeTableInput, optFns ...func(*redshiftdata.Options)) (*redshiftdata.DescribeTableOutput, error)
}

// Compile-time check that the real client satisfies the interface.
var _ dataAPIClient = (*redshiftdata.Client)(nil)
