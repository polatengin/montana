package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &MontanaProvider{}
var _ provider.ProviderWithFunctions = &MontanaProvider{}

type MontanaProvider struct {
	version string
}

type MontanaProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *MontanaProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "montana"
}

func (p *MontanaProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
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

	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
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
