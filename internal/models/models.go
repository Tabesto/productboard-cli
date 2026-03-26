package models

// PaginatedResponse wraps paginated API responses.
type PaginatedResponse struct {
	Data         []map[string]interface{} `json:"data"`
	PageCursor   *string                  `json:"pageCursor"`
	TotalResults *int                     `json:"totalResults"`
}

// SingleResponse wraps a single entity API response.
type SingleResponse struct {
	Data map[string]interface{} `json:"data"`
}

// LinksResponse wraps linked resource API responses.
type LinksResponse struct {
	Data       []map[string]interface{} `json:"data"`
	PageCursor *string                  `json:"pageCursor"`
}

// Feature represents a ProductBoard feature.
type Feature struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Status      *Status                `json:"status"`
	Parent      *Reference             `json:"parent"`
	Owner       *OwnerRef              `json:"owner"`
	Archived    bool                   `json:"archived"`
	CreatedAt   string                 `json:"createdAt"`
	UpdatedAt   string                 `json:"updatedAt"`
	Extra       map[string]interface{} `json:"-"`
}

// Status represents a feature status.
type Status struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Reference is a generic ID reference.
type Reference struct {
	ID string `json:"id"`
}

// OwnerRef is a reference by email.
type OwnerRef struct {
	Email string `json:"email"`
}

// Product represents a ProductBoard product.
type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Component represents a ProductBoard component.
type Component struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Parent      *Reference `json:"parent"`
}

// FeatureStatus represents a feature status definition.
type FeatureStatus struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Note represents a customer feedback note.
type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
	Source    string    `json:"source"`
	Owner     *OwnerRef `json:"owner"`
	Company   *Reference `json:"company"`
}

// Tag represents a note tag.
type Tag struct {
	Name string `json:"name"`
}

// Link represents a note link.
type Link struct {
	Type   string     `json:"type"`
	Target *Reference `json:"target"`
}

// Company represents a ProductBoard company.
type Company struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

// User represents a ProductBoard user.
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Release represents a release.
type Release struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Timeframe    *Timeframe `json:"timeframe"`
	ReleaseGroup *Reference `json:"releaseGroup"`
	State        string     `json:"state"`
}

// Timeframe represents a date range.
type Timeframe struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// ReleaseGroup represents a release group.
type ReleaseGroup struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// FeatureReleaseAssignment represents a feature-release assignment.
type FeatureReleaseAssignment struct {
	Feature *Reference `json:"feature"`
	Release *Reference `json:"release"`
}

// Objective represents a strategic objective.
type Objective struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	State       string    `json:"state"`
	Owner       *OwnerRef `json:"owner"`
}

// KeyResult represents an objective key result.
type KeyResult struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	CurrentValue float64    `json:"currentValue"`
	TargetValue  float64    `json:"targetValue"`
	Objective    *Reference `json:"objective"`
}

// Initiative represents a strategic initiative.
type Initiative struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	State       string    `json:"state"`
	Owner       *OwnerRef `json:"owner"`
}

// CustomField represents a custom field definition.
type CustomField struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Options     []Option `json:"options"`
}

// Option represents a dropdown option for a custom field.
type Option struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// CustomFieldValue represents a custom field value.
type CustomFieldValue struct {
	CustomField     *Reference  `json:"customField"`
	HierarchyEntity *Reference `json:"hierarchyEntity"`
	Value           interface{} `json:"value"`
}

// PluginIntegration represents a plugin integration.
type PluginIntegration struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// JiraIntegration represents a Jira integration.
type JiraIntegration struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Connection represents an integration connection.
type Connection struct {
	FeatureID string `json:"featureId"`
	ExternalID string `json:"externalId"`
}

// Webhook represents a webhook subscription.
type Webhook struct {
	ID     string   `json:"id"`
	URL    string   `json:"url"`
	Events []string `json:"events"`
}

// FeedbackFormConfiguration represents a feedback form config.
type FeedbackFormConfiguration struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Configuration map[string]interface{} `json:"configuration"`
}
