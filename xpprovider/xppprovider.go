package xpprovider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"

	"github.com/hashicorp/terraform-provider-aws/internal/provider"
)

type AWSConfig conns.Config
type AWSClient conns.AWSClient

func Provider() *schema.Provider {
	p, _ := provider.New(context.TODO())
	return p
}

func (ac *AWSConfig) GetClient(ctx context.Context, client *AWSClient) (*conns.AWSClient, diag.Diagnostics) {
	return (*conns.Config)(ac).ConfigureProvider(ctx, (*conns.AWSClient)(client))
}
