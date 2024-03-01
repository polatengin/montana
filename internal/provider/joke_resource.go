package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/polatengin/montana/internal/api"
)

var _ resource.Resource = &JokeResource{}
var _ resource.ResourceWithImportState = &JokeResource{}

func NewJokeResource() resource.Resource {
	return &JokeResource{}
}

type JokeResource struct {
	client *api.ApiClient
}

type JokeResourceModel struct {
	Text     types.String `tfsdk:"text"`
	Category types.String `tfsdk:"category"`
	Id       types.Int64  `tfsdk:"id"`
}

func (r *JokeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_joke"
}

func (r *JokeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Joke resource",

		Attributes: map[string]schema.Attribute{
			"text": schema.StringAttribute{
				MarkdownDescription: "Joke configurable attribute",
				Optional:            true,
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "Joke category",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: loadCategoryValidators(r.client),
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Joke identifier",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func loadCategoryValidators(apiClient *api.ApiClient) []validator.String {
	validators := []validator.String{}

	response, err := apiClient.Execute(context.Background(), "GET", "https://raw.githubusercontent.com/polatengin/montana/main/data/category_list.json", nil, nil, []int{200})

	if err != nil {
		tflog.Error(context.Background(), "Failed to load joke categories")
	}

	var model []struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}

	err = json.Unmarshal(response.BodyAsBytes, &model)

	if err != nil {
		tflog.Error(context.Background(), "Failed to parse joke categories")
	}

	categoryList := make([]string, len(model))

	for i, category := range model {
		categoryList[i] = category.Name
	}

	validators = append(validators, stringvalidator.OneOf(categoryList...))

	return validators
}

func (r *JokeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.ApiClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *JokeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data JokeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.Execute(ctx, "GET", fmt.Sprintf("https://v2.jokeapi.dev/joke/%s?blacklistFlags=nsfw,religious,political,racist,sexist,explicit&type=single", data.Category.ValueString()), nil, nil, []int{200})

	if err != nil {
		resp.Diagnostics.AddError("Failed to fetch joke", err.Error())
		return
	}

	var model struct {
		Id       int64  `json:"id"`
		Safe     bool   `json:"safe"`
		Language string `json:"lang"`
		Error    bool   `json:"error"`
		Category string `json:"category"`
		Type     string `json:"type"`
		Joke     string `json:"joke"`
		Flags    struct {
			Nsfw      bool `json:"nsfw"`
			Religious bool `json:"religious"`
			Political bool `json:"political"`
			Racist    bool `json:"racist"`
			Sexist    bool `json:"sexist"`
			Explicit  bool `json:"explicit"`
		} `json:"flags"`
	}

	err = json.Unmarshal(response.BodyAsBytes, &model)

	if err != nil {
		resp.Diagnostics.AddError("Failed to parse joke", err.Error())
		return
	}

	data.Id = types.Int64Value(model.Id)
	data.Text = types.StringValue(model.Joke)

	tflog.Trace(ctx, "created a resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JokeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data JokeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JokeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data JokeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JokeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data JokeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *JokeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
