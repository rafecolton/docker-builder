package job

// List is a type for sorting an array of jobs
type List []*Job

func (l List) Len() int           { return len(l) }
func (l List) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l List) Less(i, j int) bool { return l[i].Created.Unix() > l[j].Created.Unix() }
