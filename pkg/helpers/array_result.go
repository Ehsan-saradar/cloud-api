package helpers

// ArrayResult will used in methods with array return and pagination input
// and contains array items and count of total items.
type ArrayResult struct {
	TotalCount int         `json:"total_count"`
	Items      interface{} `json:"items"`
}
