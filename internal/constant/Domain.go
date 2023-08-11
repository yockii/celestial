package constant

type TimestampRangeCondition struct {
	Start int64 `json:"start,omitempty" query:"start"`
	End   int64 `json:"end,omitempty" query:"end"`
}
