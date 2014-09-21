package job

type JobList []*Job

func (l JobList) Len() int           { return len(l) }
func (l JobList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l JobList) Less(i, j int) bool { return l[i].Created.Unix() > l[j].Created.Unix() }
