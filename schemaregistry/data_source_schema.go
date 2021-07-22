package schemaregistry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/riferrei/srclient"
)

func dataSourceSchema() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSubjectRead,
		Schema: map[string]*schema.Schema{
			"subject": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The subject related to the schema",
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version of the schema",
			},
			"schema_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The schema ID",
			},
			"schema": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The schema string",
			},
			"references": {
				Type:        schema.TypeList,
				Computed: 	 true,
				Description: "The referenced schema names list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The referenced schema name",
						},
						"subject": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The subject related to the schema",
						},
						"version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The version of the schema",
						},
					},
				},
			},
		},
	}
}

func dataSourceSubjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	subject := d.Get("subject").(string)

	client := m.(*srclient.SchemaRegistryClient)

	schema, err := client.GetLatestSchemaWithArbitrarySubject(subject)
	if err != nil {
		return diag.FromErr(err)
		// return diag.FromErr(fmt.Errorf("unknown schema for subject '%s'", subject))
	}

	d.Set("schema", schema.Schema())
	d.Set("schema_id", schema.ID())
	d.Set("version", schema.Version())

	if err = d.Set("references", FromRegistryReferences(schema.References())); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(formatSchemaVersionID(subject))

	return diags
}

func FromRegistryReferences(references []srclient.Reference) []interface{} {
	if len(references) == 0 {
		return make([]interface{}, 0)
	}

	refs := make([]interface{}, 0, len(references))
	for _, reference := range references {
		refs = append(refs, map[string]interface{}{
			"name": reference.Name,
			"subject": reference.Subject,
			"version": reference.Version,
		})
	}

	return refs
}