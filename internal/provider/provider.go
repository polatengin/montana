package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/polatengin/montana/internal/api"
)

var _ provider.Provider = &MontanaProvider{}
var _ provider.ProviderWithFunctions = &MontanaProvider{}

type MontanaProvider struct {
	Api *api.ProviderClient
}

type MontanaProviderModel struct {
	UseCli types.Bool `tfsdk:"use_cli"`

	TenantId     types.String `tfsdk:"tenant_id"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

func (p *MontanaProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "montana"
}

func (p *MontanaProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Montana Provider",
		MarkdownDescription: "Montana Provider",
		Attributes: map[string]schema.Attribute{
			"use_cli": schema.BoolAttribute{
				Description:         "Flag to indicate whether to use the CLI for authentication",
				MarkdownDescription: "Flag to indicate whether to use the CLI for authentication. ",
				Optional:            true,
			},
			"tenant_id": schema.StringAttribute{
				Description:         "The id of the AAD tenant that Montana uses to authenticate with",
				MarkdownDescription: "The id of the AAD tenant that Montana uses to authenticate with",
				Optional:            true,
			},
			"client_id": schema.StringAttribute{
				Description:         "The client id of the Montana app registration",
				MarkdownDescription: "The client id of the Montana app registration",
				Optional:            true,
			},
			"client_secret": schema.StringAttribute{
				Description:         "The secret of the Montana app registration",
				MarkdownDescription: "The secret of the Montana app registration",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *MontanaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data MontanaProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	providerClient := api.ProviderClient{}
	resp.DataSourceData = &providerClient
	resp.ResourceData = &providerClient
}

func (p *MontanaProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPalindromeResource,
	}
}

func (p *MontanaProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPalindromeDataSource,
	}
}

func (p *MontanaProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New() func() provider.Provider {
	return func() provider.Provider {
		return &MontanaProvider{}
	}
}
