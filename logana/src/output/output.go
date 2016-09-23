package output

import (
	"strconv"
)

type Output_record struct {
	Name            interface{}
	Fmt_Percent     string
	Number, Percent float64
}

func (record *Output_record) Format_percent() {

	record.Fmt_Percent = strconv.FormatFloat(record.Percent, 'f', 2, 64) + "%"
}

type Output_slice struct {
	Records
	Sum float64
}

type Records []Output_record

func (s *Output_slice) Output_slice_gen(counter map[interface{}]float64, stat map[interface{}]float64, sum float64) {
	s.Sum = sum

	record := Output_record{}

	for key, value := range counter {
		record.Name = key
		record.Number = value
		record.Percent = stat[key]
		record.Format_percent()
		s.Records = append(s.Records, record)
	}
}

func (s Records) Less(i, j int) bool {
	return s[j].Percent < s[i].Percent
}

func (a Records) Len() int {
	return len(a)
}
func (a Records) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
