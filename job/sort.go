package job

// ByCreatedDescending is a type for sorting an array of jobs by created date, descending
type ByCreatedDescending []*Job

func (l ByCreatedDescending) Len() int           { return len(l) }
func (l ByCreatedDescending) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l ByCreatedDescending) Less(i, j int) bool { return l[i].Created.Unix() > l[j].Created.Unix() }
