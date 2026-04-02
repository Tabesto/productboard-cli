package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/tabesto/productboard-cli/internal/client"
	"github.com/tabesto/productboard-cli/internal/config"
	"github.com/tabesto/productboard-cli/internal/health"
)

func getClient() (*client.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return client.New(cfg)
}

func toJSON(data any) (string, error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func errorResult(err error) *mcp.CallToolResult {
	if apiErr, ok := err.(*client.APIError); ok {
		return mcp.NewToolResultError(apiErr.Message)
	}
	return mcp.NewToolResultError(err.Error())
}

func getLimit(request mcp.CallToolRequest) int {
	limit := request.GetInt("limit", defaultLimit)
	if limit <= 0 {
		return defaultLimit
	}
	return limit
}

// --- Features ---

func handleListFeatures(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	params := map[string]string{
		"statusId":   request.GetString("status_id", ""),
		"statusName": request.GetString("status_name", ""),
		"parentId":   request.GetString("parent_id", ""),
		"archived":   request.GetString("archived", ""),
		"ownerEmail": request.GetString("owner_email", ""),
		"noteId":     request.GetString("note_id", ""),
	}
	data, err := c.GetList("/features", params, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetFeature(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetSingle(fmt.Sprintf("/features/%s", id))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleListFeatureInitiatives(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetLinkedResources(fmt.Sprintf("/features/%s/links/initiatives", id), getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleListFeatureObjectives(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetLinkedResources(fmt.Sprintf("/features/%s/links/objectives", id), getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

// --- Notes ---

func handleListNotes(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	params := map[string]string{
		"term":        request.GetString("term", ""),
		"dateFrom":    request.GetString("date_from", ""),
		"dateTo":      request.GetString("date_to", ""),
		"createdFrom": request.GetString("created_from", ""),
		"createdTo":   request.GetString("created_to", ""),
		"updatedFrom": request.GetString("updated_from", ""),
		"updatedTo":   request.GetString("updated_to", ""),
		"featureId":   request.GetString("feature_id", ""),
		"companyId":   request.GetString("company_id", ""),
		"ownerEmail":  request.GetString("owner_email", ""),
		"source":      request.GetString("source", ""),
	}
	if tags := request.GetString("tags", ""); tags != "" {
		params["anyTag"] = tags
	}
	data, err := c.GetList("/notes", params, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetNote(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetSingle("/notes/" + id)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleListNoteTags(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetLinkedResources("/notes/"+id+"/tags", getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleListNoteLinks(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetLinkedResources("/notes/"+id+"/links", getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleListFeedbackForms(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	data, err := c.GetList("/feedback-form-configurations", nil, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetFeedbackForm(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetSingle("/feedback-form-configurations/" + id)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

// --- Products ---

func handleListProducts(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	data, err := c.GetList("/products", nil, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetProduct(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetSingle("/products/" + id)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

// --- Components ---

func handleListComponents(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	data, err := c.GetList("/components", nil, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetComponent(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetSingle("/components/" + id)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

// --- Releases ---

func handleListReleases(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	params := map[string]string{
		"releaseGroupId": request.GetString("release_group_id", ""),
	}
	data, err := c.GetList("/releases", params, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetRelease(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetSingle("/releases/" + id)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleListReleaseGroups(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	data, err := c.GetList("/release-groups", nil, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetReleaseGroup(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetSingle("/release-groups/" + id)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleListFeatureReleaseAssignments(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	params := map[string]string{
		"featureId":    request.GetString("feature_id", ""),
		"releaseId":    request.GetString("release_id", ""),
		"releaseState": request.GetString("release_state", ""),
		"endDateFrom":  request.GetString("end_date_from", ""),
		"endDateTo":    request.GetString("end_date_to", ""),
	}
	data, err := c.GetList("/feature-release-assignments", params, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetFeatureReleaseAssignment(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	params := map[string]string{
		"featureId": request.GetString("feature_id", ""),
		"releaseId": request.GetString("release_id", ""),
	}
	body, err := c.Get("/feature-release-assignments/assignment", params)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(string(body)), nil
}

// --- Objectives ---

func handleListObjectives(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	data, err := c.GetList("/objectives", nil, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetObjective(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetSingle("/objectives/" + id)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleListObjectiveFeatures(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetLinkedResources(fmt.Sprintf("/objectives/%s/links/features", id), getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleListObjectiveInitiatives(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetLinkedResources(fmt.Sprintf("/objectives/%s/links/initiatives", id), getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

// --- Key Results ---

func handleListKeyResults(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	data, err := c.GetList("/key-results", nil, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetKeyResult(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetSingle("/key-results/" + id)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

// --- Initiatives ---

func handleListInitiatives(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	data, err := c.GetList("/initiatives", nil, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetInitiative(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetSingle("/initiatives/" + id)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleListInitiativeObjectives(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetLinkedResources(fmt.Sprintf("/initiatives/%s/links/objectives", id), getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleListInitiativeFeatures(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetLinkedResources(fmt.Sprintf("/initiatives/%s/links/features", id), getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

// --- Companies ---

func handleListCompanies(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	params := map[string]string{
		"term":      request.GetString("term", ""),
		"hasNotes":  request.GetString("has_notes", ""),
		"featureId": request.GetString("feature_id", ""),
	}
	data, err := c.GetList("/companies", params, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetCompany(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetSingle("/companies/" + id)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleListCompanyCustomFields(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	data, err := c.GetList("/companies/custom-fields", nil, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetCompanyCustomField(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetSingle("/companies/custom-fields/" + id)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetCompanyCustomFieldValue(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	companyID, err := request.RequireString("company_id")
	if err != nil {
		return mcp.NewToolResultError("company_id is required"), nil
	}
	fieldID, err := request.RequireString("field_id")
	if err != nil {
		return mcp.NewToolResultError("field_id is required"), nil
	}
	body, err := c.Get(fmt.Sprintf("/companies/%s/custom-fields/%s/value", companyID, fieldID), nil)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(string(body)), nil
}

// --- Users ---

func handleListUsers(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	data, err := c.GetList("/users", nil, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetUser(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetSingle("/users/" + id)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

// --- Custom Fields ---

func handleListCustomFields(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	fieldType, err := request.RequireString("type")
	if err != nil {
		return mcp.NewToolResultError("type is required"), nil
	}
	params := map[string]string{"type": fieldType}
	data, err := c.GetList("/hierarchy-entities/custom-fields", params, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetCustomField(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetSingle("/hierarchy-entities/custom-fields/" + id)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleListCustomFieldValues(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	params := map[string]string{
		"type":              request.GetString("type", ""),
		"customFieldId":     request.GetString("custom_field_id", ""),
		"hierarchyEntityId": request.GetString("hierarchy_entity_id", ""),
	}
	data, err := c.GetList("/hierarchy-entities/custom-fields-values", params, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetCustomFieldValue(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	params := map[string]string{
		"customFieldId":     request.GetString("custom_field_id", ""),
		"hierarchyEntityId": request.GetString("hierarchy_entity_id", ""),
	}
	body, err := c.Get("/hierarchy-entities/custom-fields-values/value", params)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(string(body)), nil
}

// --- Feature Statuses ---

func handleListFeatureStatuses(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	data, err := c.GetList("/feature-statuses", nil, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

// --- Plugin Integrations ---

func handleListPluginIntegrations(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	data, err := c.GetList("/plugin-integrations", nil, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetPluginIntegration(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetSingle("/plugin-integrations/" + id)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleListPluginIntegrationConnections(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetList(fmt.Sprintf("/plugin-integrations/%s/connections", id), nil, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetPluginIntegrationConnection(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	featureID, err := request.RequireString("feature_id")
	if err != nil {
		return mcp.NewToolResultError("feature_id is required"), nil
	}
	data, err := c.GetSingle(fmt.Sprintf("/plugin-integrations/%s/connections/%s", id, featureID))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

// --- Jira Integrations ---

func handleListJiraIntegrations(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	data, err := c.GetList("/jira-integrations", nil, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetJiraIntegration(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetSingle("/jira-integrations/" + id)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleListJiraIntegrationConnections(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	params := map[string]string{
		"issueKey": request.GetString("issue_key", ""),
		"issueId":  request.GetString("issue_id", ""),
	}
	data, err := c.GetList(fmt.Sprintf("/jira-integrations/%s/connections", id), params, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetJiraIntegrationConnection(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	featureID, err := request.RequireString("feature_id")
	if err != nil {
		return mcp.NewToolResultError("feature_id is required"), nil
	}
	data, err := c.GetSingle(fmt.Sprintf("/jira-integrations/%s/connections/%s", id, featureID))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

// --- Feature Health ---

func handleFeaturesHealthList(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}

	// Build server-side params (only archived is server-side)
	params := map[string]string{}
	if !request.GetBool("include_archived", false) {
		params["archived"] = "false"
	}
	if c.IsV2() {
		params["fields[]"] = "health,name,status,owner"
	}

	// Fetch ALL features for client-side filtering
	data, err := c.GetList("/features", params, 0)
	if err != nil {
		return errorResult(err), nil
	}

	// Parse date filters
	opts := health.FilterOpts{
		IncludeNoHealth: request.GetBool("include_no_health", false),
		IncludeArchived: request.GetBool("include_archived", false),
		HealthStatus:    request.GetString("health_status", ""),
		StatusName:      request.GetString("status", ""),
		OwnerEmail:      request.GetString("owner", ""),
	}
	if since := request.GetString("updated_since", ""); since != "" {
		t, err := time.Parse("2006-01-02", since)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid updated_since date format %q, expected YYYY-MM-DD", since)), nil
		}
		opts.UpdatedSince = &t
	}
	if before := request.GetString("updated_before", ""); before != "" {
		t, err := time.Parse("2006-01-02", before)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid updated_before date format %q, expected YYYY-MM-DD", before)), nil
		}
		opts.UpdatedBefore = &t
	}

	filtered := health.FilterAndSort(data, opts)
	filtered = health.ApplyLimit(filtered, getLimit(request))

	s, err := toJSON(filtered)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleFeaturesHealthGet(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	var params map[string]string
	if c.IsV2() {
		params = map[string]string{"fields[]": "health,name,status,owner"}
	}
	data, err := c.GetSingleWithParams("/features/"+id, params)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

// --- Webhooks ---

func handleListWebhooks(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	data, err := c.GetList("/webhooks", nil, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetWebhook(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetSingle("/webhooks/" + id)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

// --- Members (V2 only) ---

func handleListMembers(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	if !c.IsV2() {
		return mcp.NewToolResultError("The 'list_members' tool requires API V2."), nil
	}
	params := map[string]string{
		"query":   request.GetString("query", ""),
		"roles[]": request.GetString("role", ""),
	}
	data, err := c.GetList("/members", params, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetMember(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	if !c.IsV2() {
		return mcp.NewToolResultError("The 'get_member' tool requires API V2."), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetSingle("/members/" + id)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

// --- Teams (V2 only) ---

func handleListTeams(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	if !c.IsV2() {
		return mcp.NewToolResultError("The 'list_teams' tool requires API V2."), nil
	}
	params := map[string]string{
		"query": request.GetString("query", ""),
	}
	data, err := c.GetList("/teams", params, getLimit(request))
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}

func handleGetTeam(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	if !c.IsV2() {
		return mcp.NewToolResultError("The 'get_team' tool requires API V2."), nil
	}
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}
	data, err := c.GetSingle("/teams/" + id)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(data)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}
