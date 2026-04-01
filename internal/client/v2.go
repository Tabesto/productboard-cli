package client

import (
	"net/url"
	"regexp"
	"strings"
)

// v2EntityTypes maps V1 list paths to their V2 entity type names.
var v2EntityTypes = map[string]string{
	"/features":       "feature",
	"/products":       "product",
	"/components":     "component",
	"/initiatives":    "initiative",
	"/objectives":     "objective",
	"/key-results":    "keyResult",
	"/releases":       "release",
	"/release-groups": "releaseGroup",
	"/companies":      "company",
	"/users":          "user",
}

// v2DirectPrefixPaths are V1 paths that just need a /v2 prefix in V2 mode.
var v2DirectPrefixPaths = []string{
	"/notes",
	"/jira-integrations",
	"/plugin-integrations",
	"/webhooks",
	"/members",
	"/teams",
}

// v2ParamMapping maps V1 query parameter names to V2 equivalents.
var v2ParamMapping = map[string]string{
	"statusId":   "status[id]",
	"statusName": "status[name]",
	"parentId":   "parent[id]",
	"ownerEmail": "owner[email]",
}

// entitySinglePathRe matches paths like /features/{id}, /products/{id}, etc.
var entitySinglePathRe = regexp.MustCompile(`^/(features|products|components|initiatives|objectives|key-results|releases|release-groups|companies|users)/([^/]+)$`)

// linksPathRe matches paths like /features/{id}/links/initiatives.
var linksPathRe = regexp.MustCompile(`^/(features|objectives|initiatives)/([^/]+)/links/(features|objectives|initiatives)$`)

// noteLinksPathRe matches paths like /notes/{id}/links.
var noteLinksPathRe = regexp.MustCompile(`^/notes/([^/]+)/links$`)

// noteTagsPathRe matches paths like /notes/{id}/tags.
var noteTagsPathRe = regexp.MustCompile(`^/notes/([^/]+)/tags$`)

// translateV2Path converts a V1-style path to a V2 path and returns
// any additional query parameters that should be added.
func translateV2Path(path string) (v2Path string, extraParams map[string]string) {
	extraParams = make(map[string]string)

	// Check entity single-resource paths: /features/{id} → /v2/entities/{id}
	if m := entitySinglePathRe.FindStringSubmatch(path); m != nil {
		return "/v2/entities/" + m[2], extraParams
	}

	// Check entity-to-entity link paths: /features/{id}/links/initiatives → /v2/entities/{id}/relationships
	if m := linksPathRe.FindStringSubmatch(path); m != nil {
		targetType := v2EntityTypes["/"+m[3]]
		extraParams["target[type]"] = targetType
		return "/v2/entities/" + m[2] + "/relationships", extraParams
	}

	// Check note links path: /notes/{id}/links → /v2/notes/{id}/relationships
	if m := noteLinksPathRe.FindStringSubmatch(path); m != nil {
		return "/v2/notes/" + m[1] + "/relationships", extraParams
	}

	// Check note tags path: /notes/{id}/tags — no V2 equivalent, keep as /v2/notes/{id}/tags
	if m := noteTagsPathRe.FindStringSubmatch(path); m != nil {
		return "/v2/notes/" + m[1] + "/tags", extraParams
	}

	// Check entity list paths: /features → /v2/entities?type[]=feature
	if entityType, ok := v2EntityTypes[path]; ok {
		extraParams["type[]"] = entityType
		return "/v2/entities", extraParams
	}

	// Check direct prefix paths: /notes → /v2/notes, /notes/{id} → /v2/notes/{id}
	for _, prefix := range v2DirectPrefixPaths {
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return "/v2" + path, extraParams
		}
	}

	// Feature-release assignments — V2 uses entity relationships
	if path == "/feature-release-assignments" || path == "/feature-release-assignments/assignment" {
		return "/v2" + path, extraParams
	}

	// Special cases
	if path == "/feature-statuses" {
		return "/v2/entities/configurations/feature", extraParams
	}
	if path == "/hierarchy-entities/custom-fields" || strings.HasPrefix(path, "/hierarchy-entities/custom-fields/") {
		suffix := strings.TrimPrefix(path, "/hierarchy-entities/custom-fields")
		return "/v2/entities/configurations" + suffix, extraParams
	}
	if path == "/hierarchy-entities/custom-fields-values" {
		return "/v2/entities/fields/values", extraParams
	}
	if strings.HasPrefix(path, "/hierarchy-entities/custom-fields-values/") {
		suffix := strings.TrimPrefix(path, "/hierarchy-entities/custom-fields-values")
		return "/v2/entities/fields" + suffix + "/values", extraParams
	}

	// Fallback: just prefix with /v2
	return "/v2" + path, extraParams
}

// translateV2Params converts V1-style query parameters to V2 equivalents.
func translateV2Params(params map[string]string) map[string]string {
	if params == nil {
		return nil
	}
	translated := make(map[string]string, len(params))
	for k, v := range params {
		if v2Key, ok := v2ParamMapping[k]; ok {
			translated[v2Key] = v
		} else {
			translated[k] = v
		}
	}
	return translated
}

// flattenV2Entity merges the "fields" sub-object to the top level of a V2 entity.
// This preserves backward compatibility with display code that reads from the top level.
func flattenV2Entity(entity map[string]interface{}) {
	if entity == nil {
		return
	}
	fields, ok := entity["fields"].(map[string]interface{})
	if !ok {
		return
	}
	for k, v := range fields {
		entity[k] = v
	}
	delete(entity, "fields")

	// Normalize V2 "health" field to V1 "lastHealthUpdate" format
	normalizeV2Health(entity)
}

// v2HealthStatusMap maps V2 camelCase health statuses to V1 kebab-case.
var v2HealthStatusMap = map[string]string{
	"onTrack": "on-track",
	"atRisk":  "at-risk",
	"offTrack": "off-track",
}

// normalizeV2Health converts V2 "health" field to V1 "lastHealthUpdate" format.
func normalizeV2Health(entity map[string]interface{}) {
	health, ok := entity["health"].(map[string]interface{})
	if !ok || health == nil {
		return
	}

	// Map V2 field names to V1 equivalents
	normalized := make(map[string]interface{})
	for k, v := range health {
		normalized[k] = v
	}

	// V2 "comment" → V1 "message"
	if comment, ok := normalized["comment"]; ok {
		normalized["message"] = comment
		delete(normalized, "comment")
	}

	// V2 "lastUpdatedAt" → V1 "createdAt" (used for date display)
	if updatedAt, ok := normalized["lastUpdatedAt"]; ok {
		normalized["createdAt"] = updatedAt
	}

	// V2 status camelCase → V1 kebab-case
	if status, ok := normalized["status"].(string); ok {
		if mapped, exists := v2HealthStatusMap[status]; exists {
			normalized["status"] = mapped
		}
	}

	entity["lastHealthUpdate"] = normalized
	delete(entity, "health")
}

// flattenV2Relationship extracts the target entity from a V2 relationship object
// so existing display code can read id/name/etc from the top level.
// V2 relationship: {type: "link", source: {id, type}, target: {id, type}}
// After flattening: {id: target.id, type: target.type, name: target.name, ...}
func flattenV2Relationship(rel map[string]interface{}) map[string]interface{} {
	target, ok := rel["target"].(map[string]interface{})
	if !ok {
		return rel
	}
	// Return the target entity as a flat map
	result := make(map[string]interface{})
	for k, v := range target {
		result[k] = v
	}
	// Also flatten fields if present (target might have fields sub-object)
	flattenV2Entity(result)
	return result
}

// isRelationshipPath returns true if the V2 path is a relationships endpoint.
func isRelationshipPath(v2Path string) bool {
	return strings.HasSuffix(v2Path, "/relationships") || strings.Contains(v2Path, "/relationships?")
}

// extractPageCursorFromURL extracts the pageCursor query parameter from a full URL.
func extractPageCursorFromURL(nextURL string) string {
	u, err := url.Parse(nextURL)
	if err != nil {
		return ""
	}
	return u.Query().Get("pageCursor")
}
