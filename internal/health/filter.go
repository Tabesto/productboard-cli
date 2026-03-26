package health

import (
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/tabesto/productboard-cli/internal/output"
)

// FilterOpts holds all filter parameters for health feature filtering.
type FilterOpts struct {
	IncludeNoHealth bool
	IncludeArchived bool
	UpdatedSince    *time.Time
	UpdatedBefore   *time.Time
	HealthStatus    string
	StatusName      string
	OwnerEmail      string
}

// FilterAndSort applies client-side filters and sorts features by health update date descending.
func FilterAndSort(features []map[string]interface{}, opts FilterOpts) []map[string]interface{} {
	filtered := make([]map[string]interface{}, 0, len(features))

	for _, f := range features {
		// Skip features without health updates unless explicitly included
		healthUpdate := getHealthUpdate(f)
		if healthUpdate == nil && !opts.IncludeNoHealth {
			continue
		}

		// Client-side status name filter
		if opts.StatusName != "" {
			featureStatus := output.SafeNested(f, "status", "name")
			if !strings.EqualFold(featureStatus, opts.StatusName) {
				continue
			}
		}

		// Client-side owner email filter
		if opts.OwnerEmail != "" {
			ownerEmail := output.SafeNested(f, "owner", "email")
			if !strings.EqualFold(ownerEmail, opts.OwnerEmail) {
				continue
			}
		}

		// Client-side health status filter
		if opts.HealthStatus != "" && healthUpdate != nil {
			hStatus := output.SafeStr(healthUpdate, "status")
			if !strings.EqualFold(hStatus, opts.HealthStatus) {
				continue
			}
		}

		// Client-side date range filter on health update date
		if healthUpdate != nil && (opts.UpdatedSince != nil || opts.UpdatedBefore != nil) {
			updatedAtStr := output.SafeStr(healthUpdate, "createdAt")
			updatedAt, err := time.Parse(time.RFC3339, updatedAtStr)
			if err != nil {
				// If date can't be parsed, skip from date-filtered results
				continue
			}
			if opts.UpdatedSince != nil && updatedAt.Before(*opts.UpdatedSince) {
				continue
			}
			if opts.UpdatedBefore != nil && !updatedAt.Before(*opts.UpdatedBefore) {
				continue
			}
		} else if healthUpdate == nil && (opts.UpdatedSince != nil || opts.UpdatedBefore != nil) {
			// No health update means no date to compare against; skip from date-filtered results
			continue
		}

		filtered = append(filtered, f)
	}

	// Sort by lastHealthUpdate.updatedAt descending; features without health sort to end
	sort.Slice(filtered, func(i, j int) bool {
		hi := getHealthUpdate(filtered[i])
		hj := getHealthUpdate(filtered[j])
		if hi == nil && hj == nil {
			return false
		}
		if hi == nil {
			return false
		}
		if hj == nil {
			return true
		}
		ti, _ := time.Parse(time.RFC3339, output.SafeStr(hi, "createdAt"))
		tj, _ := time.Parse(time.RFC3339, output.SafeStr(hj, "createdAt"))
		return ti.After(tj)
	})

	return filtered
}

// ApplyLimit truncates results to the given limit. A limit of 0 or negative means no limit.
func ApplyLimit(features []map[string]interface{}, limit int) []map[string]interface{} {
	if limit > 0 && len(features) > limit {
		return features[:limit]
	}
	return features
}

// getHealthUpdate extracts the lastHealthUpdate map from a feature, returning nil if absent.
func getHealthUpdate(f map[string]interface{}) map[string]interface{} {
	v, ok := f["lastHealthUpdate"]
	if !ok || v == nil {
		return nil
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return m
}

// GetHealthUpdate is the exported version for use by CLI and MCP packages.
func GetHealthUpdate(f map[string]interface{}) map[string]interface{} {
	return getHealthUpdate(f)
}

var htmlTagRe = regexp.MustCompile(`<[^>]*>`)

// StripHTML removes HTML tags and collapses whitespace for clean table display.
func StripHTML(s string) string {
	s = htmlTagRe.ReplaceAllString(s, " ")
	s = strings.Join(strings.Fields(s), " ")
	return strings.TrimSpace(s)
}

// FormatDate formats an ISO 8601 timestamp to YYYY-MM-DD for table display.
func FormatDate(iso string) string {
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return iso
	}
	return t.Format("2006-01-02")
}
