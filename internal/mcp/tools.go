package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerTools(s *server.MCPServer) {
	// --- Features ---
	s.AddTool(mcp.NewTool("list_features",
		mcp.WithDescription("List ProductBoard features with optional filters"),
		mcp.WithString("status_id", mcp.Description("Filter by status ID")),
		mcp.WithString("status_name", mcp.Description("Filter by status name")),
		mcp.WithString("parent_id", mcp.Description("Filter by parent feature ID")),
		mcp.WithString("archived", mcp.Description("Filter by archived state (true/false)")),
		mcp.WithString("owner_email", mcp.Description("Filter by owner email")),
		mcp.WithString("note_id", mcp.Description("Filter by note ID")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListFeatures)

	s.AddTool(mcp.NewTool("get_feature",
		mcp.WithDescription("Get a ProductBoard feature by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Feature ID")),
	), handleGetFeature)

	s.AddTool(mcp.NewTool("list_feature_initiatives",
		mcp.WithDescription("List initiatives linked to a feature"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Feature ID")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListFeatureInitiatives)

	s.AddTool(mcp.NewTool("list_feature_objectives",
		mcp.WithDescription("List objectives linked to a feature"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Feature ID")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListFeatureObjectives)

	// --- Notes ---
	s.AddTool(mcp.NewTool("list_notes",
		mcp.WithDescription("List ProductBoard notes with optional filters"),
		mcp.WithString("term", mcp.Description("Search term")),
		mcp.WithString("date_from", mcp.Description("Filter by date from (ISO 8601)")),
		mcp.WithString("date_to", mcp.Description("Filter by date to (ISO 8601)")),
		mcp.WithString("created_from", mcp.Description("Filter by created from (ISO 8601)")),
		mcp.WithString("created_to", mcp.Description("Filter by created to (ISO 8601)")),
		mcp.WithString("updated_from", mcp.Description("Filter by updated from (ISO 8601)")),
		mcp.WithString("updated_to", mcp.Description("Filter by updated to (ISO 8601)")),
		mcp.WithString("feature_id", mcp.Description("Filter by feature ID")),
		mcp.WithString("company_id", mcp.Description("Filter by company ID")),
		mcp.WithString("owner_email", mcp.Description("Filter by owner email")),
		mcp.WithString("source", mcp.Description("Filter by source")),
		mcp.WithString("tags", mcp.Description("Filter by tags (comma-separated)")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListNotes)

	s.AddTool(mcp.NewTool("get_note",
		mcp.WithDescription("Get a ProductBoard note by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Note ID")),
	), handleGetNote)

	s.AddTool(mcp.NewTool("list_note_tags",
		mcp.WithDescription("List tags for a note"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Note ID")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListNoteTags)

	s.AddTool(mcp.NewTool("list_note_links",
		mcp.WithDescription("List links for a note"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Note ID")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListNoteLinks)

	s.AddTool(mcp.NewTool("list_feedback_forms",
		mcp.WithDescription("List ProductBoard feedback form configurations"),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListFeedbackForms)

	s.AddTool(mcp.NewTool("get_feedback_form",
		mcp.WithDescription("Get a feedback form configuration by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Feedback form configuration ID")),
	), handleGetFeedbackForm)

	// --- Products ---
	s.AddTool(mcp.NewTool("list_products",
		mcp.WithDescription("List ProductBoard products"),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListProducts)

	s.AddTool(mcp.NewTool("get_product",
		mcp.WithDescription("Get a ProductBoard product by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Product ID")),
	), handleGetProduct)

	// --- Components ---
	s.AddTool(mcp.NewTool("list_components",
		mcp.WithDescription("List ProductBoard components"),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListComponents)

	s.AddTool(mcp.NewTool("get_component",
		mcp.WithDescription("Get a ProductBoard component by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Component ID")),
	), handleGetComponent)

	// --- Releases ---
	s.AddTool(mcp.NewTool("list_releases",
		mcp.WithDescription("List ProductBoard releases"),
		mcp.WithString("release_group_id", mcp.Description("Filter by release group ID")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListReleases)

	s.AddTool(mcp.NewTool("get_release",
		mcp.WithDescription("Get a ProductBoard release by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Release ID")),
	), handleGetRelease)

	s.AddTool(mcp.NewTool("list_release_groups",
		mcp.WithDescription("List ProductBoard release groups"),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListReleaseGroups)

	s.AddTool(mcp.NewTool("get_release_group",
		mcp.WithDescription("Get a ProductBoard release group by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Release group ID")),
	), handleGetReleaseGroup)

	s.AddTool(mcp.NewTool("list_feature_release_assignments",
		mcp.WithDescription("List feature-release assignments with optional filters"),
		mcp.WithString("feature_id", mcp.Description("Filter by feature ID")),
		mcp.WithString("release_id", mcp.Description("Filter by release ID")),
		mcp.WithString("release_state", mcp.Description("Filter by release state")),
		mcp.WithString("end_date_from", mcp.Description("Filter by end date from (ISO 8601)")),
		mcp.WithString("end_date_to", mcp.Description("Filter by end date to (ISO 8601)")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListFeatureReleaseAssignments)

	s.AddTool(mcp.NewTool("get_feature_release_assignment",
		mcp.WithDescription("Get a specific feature-release assignment by feature ID and release ID"),
		mcp.WithString("feature_id", mcp.Required(), mcp.Description("Feature ID")),
		mcp.WithString("release_id", mcp.Required(), mcp.Description("Release ID")),
	), handleGetFeatureReleaseAssignment)

	// --- Objectives ---
	s.AddTool(mcp.NewTool("list_objectives",
		mcp.WithDescription("List ProductBoard objectives"),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListObjectives)

	s.AddTool(mcp.NewTool("get_objective",
		mcp.WithDescription("Get a ProductBoard objective by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Objective ID")),
	), handleGetObjective)

	s.AddTool(mcp.NewTool("list_objective_features",
		mcp.WithDescription("List features linked to an objective"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Objective ID")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListObjectiveFeatures)

	s.AddTool(mcp.NewTool("list_objective_initiatives",
		mcp.WithDescription("List initiatives linked to an objective"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Objective ID")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListObjectiveInitiatives)

	// --- Key Results ---
	s.AddTool(mcp.NewTool("list_key_results",
		mcp.WithDescription("List ProductBoard key results"),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListKeyResults)

	s.AddTool(mcp.NewTool("get_key_result",
		mcp.WithDescription("Get a ProductBoard key result by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Key result ID")),
	), handleGetKeyResult)

	// --- Initiatives ---
	s.AddTool(mcp.NewTool("list_initiatives",
		mcp.WithDescription("List ProductBoard initiatives"),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListInitiatives)

	s.AddTool(mcp.NewTool("get_initiative",
		mcp.WithDescription("Get a ProductBoard initiative by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Initiative ID")),
	), handleGetInitiative)

	s.AddTool(mcp.NewTool("list_initiative_objectives",
		mcp.WithDescription("List objectives linked to an initiative"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Initiative ID")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListInitiativeObjectives)

	s.AddTool(mcp.NewTool("list_initiative_features",
		mcp.WithDescription("List features linked to an initiative"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Initiative ID")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListInitiativeFeatures)

	// --- Companies ---
	s.AddTool(mcp.NewTool("list_companies",
		mcp.WithDescription("List ProductBoard companies with optional filters"),
		mcp.WithString("term", mcp.Description("Search term")),
		mcp.WithString("has_notes", mcp.Description("Filter by whether company has notes (true/false)")),
		mcp.WithString("feature_id", mcp.Description("Filter by feature ID")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListCompanies)

	s.AddTool(mcp.NewTool("get_company",
		mcp.WithDescription("Get a ProductBoard company by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Company ID")),
	), handleGetCompany)

	s.AddTool(mcp.NewTool("list_company_custom_fields",
		mcp.WithDescription("List custom fields for companies"),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListCompanyCustomFields)

	s.AddTool(mcp.NewTool("get_company_custom_field",
		mcp.WithDescription("Get a company custom field by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Custom field ID")),
	), handleGetCompanyCustomField)

	s.AddTool(mcp.NewTool("get_company_custom_field_value",
		mcp.WithDescription("Get the value of a custom field for a specific company"),
		mcp.WithString("company_id", mcp.Required(), mcp.Description("Company ID")),
		mcp.WithString("field_id", mcp.Required(), mcp.Description("Custom field ID")),
	), handleGetCompanyCustomFieldValue)

	// --- Users ---
	s.AddTool(mcp.NewTool("list_users",
		mcp.WithDescription("List ProductBoard users"),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListUsers)

	s.AddTool(mcp.NewTool("get_user",
		mcp.WithDescription("Get a ProductBoard user by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("User ID")),
	), handleGetUser)

	// --- Custom Fields (hierarchy entities) ---
	s.AddTool(mcp.NewTool("list_custom_fields",
		mcp.WithDescription("List hierarchy entity custom fields"),
		mcp.WithString("type", mcp.Required(), mcp.Description("Entity type (e.g., feature, sub-feature, product, component)")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListCustomFields)

	s.AddTool(mcp.NewTool("get_custom_field",
		mcp.WithDescription("Get a hierarchy entity custom field by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Custom field ID")),
	), handleGetCustomField)

	s.AddTool(mcp.NewTool("list_custom_field_values",
		mcp.WithDescription("List values for hierarchy entity custom fields"),
		mcp.WithString("type", mcp.Description("Entity type")),
		mcp.WithString("custom_field_id", mcp.Description("Filter by custom field ID")),
		mcp.WithString("hierarchy_entity_id", mcp.Description("Filter by hierarchy entity ID")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListCustomFieldValues)

	s.AddTool(mcp.NewTool("get_custom_field_value",
		mcp.WithDescription("Get a specific custom field value for a hierarchy entity"),
		mcp.WithString("custom_field_id", mcp.Required(), mcp.Description("Custom field ID")),
		mcp.WithString("hierarchy_entity_id", mcp.Required(), mcp.Description("Hierarchy entity ID")),
	), handleGetCustomFieldValue)

	// --- Feature Statuses ---
	s.AddTool(mcp.NewTool("list_feature_statuses",
		mcp.WithDescription("List ProductBoard feature statuses"),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListFeatureStatuses)

	// --- Plugin Integrations ---
	s.AddTool(mcp.NewTool("list_plugin_integrations",
		mcp.WithDescription("List ProductBoard plugin integrations"),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListPluginIntegrations)

	s.AddTool(mcp.NewTool("get_plugin_integration",
		mcp.WithDescription("Get a plugin integration by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Plugin integration ID")),
	), handleGetPluginIntegration)

	s.AddTool(mcp.NewTool("list_plugin_integration_connections",
		mcp.WithDescription("List connections for a plugin integration"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Plugin integration ID")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListPluginIntegrationConnections)

	s.AddTool(mcp.NewTool("get_plugin_integration_connection",
		mcp.WithDescription("Get a specific connection for a plugin integration"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Plugin integration ID")),
		mcp.WithString("feature_id", mcp.Required(), mcp.Description("Feature ID")),
	), handleGetPluginIntegrationConnection)

	// --- Jira Integrations ---
	s.AddTool(mcp.NewTool("list_jira_integrations",
		mcp.WithDescription("List ProductBoard Jira integrations"),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListJiraIntegrations)

	s.AddTool(mcp.NewTool("get_jira_integration",
		mcp.WithDescription("Get a Jira integration by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Jira integration ID")),
	), handleGetJiraIntegration)

	s.AddTool(mcp.NewTool("list_jira_integration_connections",
		mcp.WithDescription("List connections for a Jira integration"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Jira integration ID")),
		mcp.WithString("issue_key", mcp.Description("Filter by Jira issue key")),
		mcp.WithString("issue_id", mcp.Description("Filter by Jira issue ID")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListJiraIntegrationConnections)

	s.AddTool(mcp.NewTool("get_jira_integration_connection",
		mcp.WithDescription("Get a specific Jira integration connection"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Jira integration ID")),
		mcp.WithString("feature_id", mcp.Required(), mcp.Description("Feature ID")),
	), handleGetJiraIntegrationConnection)

	// --- Webhooks ---
	s.AddTool(mcp.NewTool("list_webhooks",
		mcp.WithDescription("List ProductBoard webhooks"),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 25)")),
	), handleListWebhooks)

	s.AddTool(mcp.NewTool("get_webhook",
		mcp.WithDescription("Get a ProductBoard webhook by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Webhook ID")),
	), handleGetWebhook)
}
