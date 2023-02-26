package table

import "bytes"

type TableCollection []*Table

func (tc TableCollection) Len() int {
	return len(tc)
}

func (tc TableCollection) Less(i, j int) bool {
	return tc[i].Name < tc[j].Name
}

func (tc TableCollection) Swap(i, j int) {
	tc[i], tc[j] = tc[j], tc[i]
}

func (tc TableCollection) String() string {
	var buf = new(bytes.Buffer)
	for _, v := range tc {
		buf.WriteString(v.String())
	}
	return buf.String()
}
