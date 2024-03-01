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
	Api *api.ApiClient
}

type MontanaProviderModel struct {
	Token types.String `tfsdk:"token"`
}

func (p *MontanaProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "montana"
}

func (p *MontanaProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Montana Provider",
		MarkdownDescription: "Montana Provider",
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Description:         "Joke Api token",
				MarkdownDescription: "Joke Api token",
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

	providerClient := api.ApiClient{}
	resp.DataSourceData = &providerClient
	resp.ResourceData = &providerClient
}

func (p *MontanaProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewJokeResource,
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
