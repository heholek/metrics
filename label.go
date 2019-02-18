package metrics

import (
	"strings"
	"sync"
)

// FNV-1a
const (
	hashOffset = 14695981039346656037
	hashPrime  = 1099511628211
)

type map1Label struct {
	mutex       sync.Mutex
	name        string
	labelName   string
	labelHashes []uint64
}

type map2Label struct {
	mutex       sync.Mutex
	name        string
	labelNames  [2]string
	labelHashes []uint64
}

type map3Label struct {
	mutex       sync.Mutex
	name        string
	labelNames  [3]string
	labelHashes []uint64
}

// Map1LabelCounter is a Counter composition with a fixed label.
// Multiple goroutines may invoke methods on a Map1LabelCounter simultaneously.
type Map1LabelCounter struct {
	map1Label
	counters []*Counter
}

// Map2LabelCounter is a Counter composition with 2 fixed labels.
// Multiple goroutines may invoke methods on a Map2LabelCounter simultaneously.
type Map2LabelCounter struct {
	map2Label
	counters []*Counter
}

// Map3LabelCounter is a Counter composition with 3 fixed labels.
// Multiple goroutines may invoke methods on a Map3LabelCounter simultaneously.
type Map3LabelCounter struct {
	map3Label
	counters []*Counter
}

// Map1LabelGauge is a Gauge composition with a fixed label.
// Multiple goroutines may invoke methods on a Map1LabelGauge simultaneously.
type Map1LabelGauge struct {
	map1Label
	gauges []*Gauge
}

// Map2LabelGauge is a Gauge composition with 2 fixed labels.
// Multiple goroutines may invoke methods on a Map2LabelGauge simultaneously.
type Map2LabelGauge struct {
	map2Label
	gauges []*Gauge
}

// Map3LabelGauge is a Gauge composition with 3 fixed labels.
// Multiple goroutines may invoke methods on a Map3LabelGauge simultaneously.
type Map3LabelGauge struct {
	map3Label
	gauges []*Gauge
}

// Map1LabelHistogram is a Histogram composition with a fixed label.
// Multiple goroutines may invoke methods on a Map1LabelHistogram simultaneously.
type Map1LabelHistogram struct {
	map1Label
	buckets    []float64
	histograms []*Histogram
}

// Map2LabelHistogram is a Histogram composition with 2 fixed labels.
// Multiple goroutines may invoke methods on a Map2LabelHistogram simultaneously.
type Map2LabelHistogram struct {
	map2Label
	buckets    []float64
	histograms []*Histogram
}

// Map3LabelHistogram is a Histogram composition with 3 fixed labels.
// Multiple goroutines may invoke methods on a Map3LabelHistogram simultaneously.
type Map3LabelHistogram struct {
	map3Label
	histograms []*Histogram
}

// Map1LabelSample is a Sample composition with a fixed label.
// Multiple goroutines may invoke methods on a Map1LabelSample simultaneously.
type Map1LabelSample struct {
	map1Label
	samples []*Sample
}

// Map2LabelSample is a Sample composition with 2 fixed labels.
// Multiple goroutines may invoke methods on a Map2LabelSample simultaneously.
type Map2LabelSample struct {
	map2Label
	samples []*Sample
}

// Map3LabelSample is a Sample composition with 3 fixed labels.
// Multiple goroutines may invoke methods on a Map3LabelSample simultaneously.
type Map3LabelSample struct {
	map3Label
	samples []*Sample
}

// With returns a dedicated Counter for a label. The value
// maps to the name as defined at Must2LabelCounter. With
// registers a new Counter if the label hasn't been used before.
// Remember that each label represents a new time series,
// which can dramatically increase the amount of data stored.
func (l1 *Map1LabelCounter) With(value string) *Counter {
	hash := uint64(hashOffset)
	for i := 0; i < len(value); i++ {
		hash ^= uint64(value[i])
		hash *= hashPrime
	}

	l1.mutex.Lock()

	for i, h := range l1.labelHashes {
		if h == hash {
			hit := l1.counters[i]

			l1.mutex.Unlock()
			return hit
		}
	}

	l1.labelHashes = append(l1.labelHashes, hash)
	c := &Counter{prefix: format1LabelPrefix(l1.name, l1.labelName, value)}
	l1.counters = append(l1.counters, c)

	l1.mutex.Unlock()
	return c
}

// With returns a dedicated Counter for a label combination. The values
// map to the names (in order) as defined at Must2LabelCounter. With
// registers a new Counter if the combination hasn't been used before.
// Remember that each label combination represents a new time series,
// which can dramatically increase the amount of data stored.
func (l2 *Map2LabelCounter) With(value1, value2 string) *Counter {
	hash := uint64(hashOffset)
	hash ^= uint64(len(value1))
	hash *= hashPrime
	for i := 0; i < len(value1); i++ {
		hash ^= uint64(value1[i])
		hash *= hashPrime
	}
	for i := 0; i < len(value2); i++ {
		hash ^= uint64(value2[i])
		hash *= hashPrime
	}

	l2.mutex.Lock()

	for i, h := range l2.labelHashes {
		if h == hash {
			hit := l2.counters[i]

			l2.mutex.Unlock()
			return hit
		}
	}

	l2.labelHashes = append(l2.labelHashes, hash)
	c := &Counter{prefix: format2LabelPrefix(l2.name, &l2.labelNames, value1, value2)}
	l2.counters = append(l2.counters, c)

	l2.mutex.Unlock()
	return c
}

// With returns a dedicated Counter for a label combination. The values
// map to the names (in order) as defined at Must3LabelCounter. With
// registers a new Counter if the combination hasn't been used before.
// Remember that each label combination represents a new time series,
// which can dramatically increase the amount of data stored.
func (l3 *Map3LabelCounter) With(value1, value2, value3 string) *Counter {
	hash := uint64(hashOffset)
	hash ^= uint64(len(value1))
	hash *= hashPrime
	for i := 0; i < len(value1); i++ {
		hash ^= uint64(value1[i])
		hash *= hashPrime
	}
	hash ^= uint64(len(value2))
	hash *= hashPrime
	for i := 0; i < len(value2); i++ {
		hash ^= uint64(value2[i])
		hash *= hashPrime
	}
	for i := 0; i < len(value3); i++ {
		hash ^= uint64(value3[i])
		hash *= hashPrime
	}

	l3.mutex.Lock()

	for i, h := range l3.labelHashes {
		if h == hash {
			hit := l3.counters[i]

			l3.mutex.Unlock()
			return hit
		}
	}

	l3.labelHashes = append(l3.labelHashes, hash)
	c := &Counter{prefix: format3LabelPrefix(l3.name, &l3.labelNames, value1, value2, value3)}
	l3.counters = append(l3.counters, c)

	l3.mutex.Unlock()
	return c
}

// With returns a dedicated Gauge for a label. The value
// maps to the name as defined at Must2LabelGauge. With
// registers a new Gauge if the label hasn't been used before.
// Remember that each label represents a new time series,
// which can dramatically increase the amount of data stored.
func (l1 *Map1LabelGauge) With(value string) *Gauge {
	hash := uint64(hashOffset)
	for i := 0; i < len(value); i++ {
		hash ^= uint64(value[i])
		hash *= hashPrime
	}

	l1.mutex.Lock()

	for i, h := range l1.labelHashes {
		if h == hash {
			hit := l1.gauges[i]

			l1.mutex.Unlock()
			return hit
		}
	}

	l1.labelHashes = append(l1.labelHashes, hash)
	g := &Gauge{prefix: format1LabelPrefix(l1.name, l1.labelName, value)}
	l1.gauges = append(l1.gauges, g)

	l1.mutex.Unlock()
	return g
}

// With returns a dedicated Gauge for a label combination. The values
// map to the names (in order) as defined at Must2LabelGauge. With
// registers a new Gauge if the combination hasn't been used before.
// Remember that each label combination represents a new time series,
// which can dramatically increase the amount of data stored.
func (l2 *Map2LabelGauge) With(value1, value2 string) *Gauge {
	hash := uint64(hashOffset)
	hash ^= uint64(len(value1))
	hash *= hashPrime
	for i := 0; i < len(value1); i++ {
		hash ^= uint64(value1[i])
		hash *= hashPrime
	}
	for i := 0; i < len(value2); i++ {
		hash ^= uint64(value2[i])
		hash *= hashPrime
	}

	l2.mutex.Lock()

	for i, h := range l2.labelHashes {
		if h == hash {
			hit := l2.gauges[i]

			l2.mutex.Unlock()
			return hit
		}
	}

	l2.labelHashes = append(l2.labelHashes, hash)
	g := &Gauge{prefix: format2LabelPrefix(l2.name, &l2.labelNames, value1, value2)}
	l2.gauges = append(l2.gauges, g)

	l2.mutex.Unlock()
	return g
}

// With returns a dedicated Gauge for a label combination. The values
// map to the names (in order) as defined at Must3LabelGauge. With
// registers a new Gauge if the combination hasn't been used before.
// Remember that each label combination represents a new time series,
// which can dramatically increase the amount of data stored.
func (l3 *Map3LabelGauge) With(value1, value2, value3 string) *Gauge {
	hash := uint64(hashOffset)
	hash ^= uint64(len(value1))
	hash *= hashPrime
	for i := 0; i < len(value1); i++ {
		hash ^= uint64(value1[i])
		hash *= hashPrime
	}
	hash ^= uint64(len(value2))
	hash *= hashPrime
	for i := 0; i < len(value2); i++ {
		hash ^= uint64(value2[i])
		hash *= hashPrime
	}
	for i := 0; i < len(value3); i++ {
		hash ^= uint64(value3[i])
		hash *= hashPrime
	}

	l3.mutex.Lock()

	for i, h := range l3.labelHashes {
		if h == hash {
			hit := l3.gauges[i]

			l3.mutex.Unlock()
			return hit
		}
	}

	l3.labelHashes = append(l3.labelHashes, hash)
	g := &Gauge{prefix: format3LabelPrefix(l3.name, &l3.labelNames, value1, value2, value3)}
	l3.gauges = append(l3.gauges, g)

	l3.mutex.Unlock()
	return g
}

// With returns a dedicated Histogram for a label. The value
// maps to the name as defined at Must2LabelHistogram. With
// registers a new Histogram if the label hasn't been used before.
// Remember that each label represents a new time series,
// which can dramatically increase the amount of data stored.
func (l1 *Map1LabelHistogram) With(value string, buckets ...float64) *Histogram {
	hash := uint64(hashOffset)
	for i := 0; i < len(value); i++ {
		hash ^= uint64(value[i])
		hash *= hashPrime
	}

	l1.mutex.Lock()

	for i, h := range l1.labelHashes {
		if h == hash {
			hit := l1.histograms[i]

			l1.mutex.Unlock()
			return hit
		}
	}

	l1.labelHashes = append(l1.labelHashes, hash)
	h := newHistogram(l1.name, buckets)
	l1.histograms = append(l1.histograms, h)

	h.sumPrefix = merge1LabelPrefix(h.sumPrefix, l1.labelName, value)
	h.countPrefix = merge1LabelPrefix(h.countPrefix, l1.labelName, value)
	for i, s := range h.bucketPrefixes {
		h.bucketPrefixes[i] = merge1LabelPrefix(s, l1.labelName, value)
	}

	l1.mutex.Unlock()
	return h
}

// With returns a dedicated Histogram for a label combination. The values
// map to the names (in order) as defined at Must2LabelHistogram. With
// registers a new Histogram if the combination hasn't been used before.
// Remember that each label combination represents a new time series,
// which can dramatically increase the amount of data stored.
func (l2 *Map2LabelHistogram) With(value1, value2 string, buckets ...float64) *Histogram {
	hash := uint64(hashOffset)
	hash ^= uint64(len(value1))
	hash *= hashPrime
	for i := 0; i < len(value1); i++ {
		hash ^= uint64(value1[i])
		hash *= hashPrime
	}
	for i := 0; i < len(value2); i++ {
		hash ^= uint64(value2[i])
		hash *= hashPrime
	}

	l2.mutex.Lock()

	for i, h := range l2.labelHashes {
		if h == hash {
			hit := l2.histograms[i]

			l2.mutex.Unlock()
			return hit
		}
	}

	l2.labelHashes = append(l2.labelHashes, hash)
	h := newHistogram(l2.name, buckets)
	l2.histograms = append(l2.histograms, h)

	h.sumPrefix = merge2LabelPrefix(h.sumPrefix, &l2.labelNames, value1, value2)
	h.countPrefix = merge2LabelPrefix(h.countPrefix, &l2.labelNames, value1, value2)
	for i, s := range h.bucketPrefixes {
		h.bucketPrefixes[i] = merge2LabelPrefix(s, &l2.labelNames, value1, value2)
	}

	l2.mutex.Unlock()
	return h
}

// With returns a dedicated Histogram for a label combination. The values
// map to the names (in order) as defined at Must3LabelHistogram. With
// registers a new Histogram if the combination hasn't been used before.
// Remember that each label combination represents a new time series,
// which can dramatically increase the amount of data stored.
func (l3 *Map3LabelHistogram) With(value1, value2, value3 string, buckets ...float64) *Histogram {
	hash := uint64(hashOffset)
	hash ^= uint64(len(value1))
	hash *= hashPrime
	for i := 0; i < len(value1); i++ {
		hash ^= uint64(value1[i])
		hash *= hashPrime
	}
	hash ^= uint64(len(value2))
	hash *= hashPrime
	for i := 0; i < len(value2); i++ {
		hash ^= uint64(value2[i])
		hash *= hashPrime
	}
	for i := 0; i < len(value3); i++ {
		hash ^= uint64(value3[i])
		hash *= hashPrime
	}

	l3.mutex.Lock()

	for i, h := range l3.labelHashes {
		if h == hash {
			hit := l3.histograms[i]

			l3.mutex.Unlock()
			return hit
		}
	}

	l3.labelHashes = append(l3.labelHashes, hash)
	h := newHistogram(l3.name, buckets)
	l3.histograms = append(l3.histograms, h)

	h.sumPrefix = merge3LabelPrefix(h.sumPrefix, &l3.labelNames, value1, value2, value3)
	h.countPrefix = merge3LabelPrefix(h.countPrefix, &l3.labelNames, value1, value2, value3)
	for i, s := range h.bucketPrefixes {
		h.bucketPrefixes[i] = merge3LabelPrefix(s, &l3.labelNames, value1, value2, value3)
	}

	l3.mutex.Unlock()
	return h
}

// With returns a dedicated Sample for a label. The value
// maps to the name as defined at Must2LabelSample. With
// registers a new Sample if the label hasn't been used before.
// Remember that each label represents a new time series,
// which can dramatically increase the amount of data stored.
func (l1 *Map1LabelSample) With(value string) *Sample {
	hash := uint64(hashOffset)
	for i := 0; i < len(value); i++ {
		hash ^= uint64(value[i])
		hash *= hashPrime
	}

	l1.mutex.Lock()

	for i, h := range l1.labelHashes {
		if h == hash {
			hit := l1.samples[i]

			l1.mutex.Unlock()
			return hit
		}
	}

	l1.labelHashes = append(l1.labelHashes, hash)
	s := &Sample{prefix: format1LabelPrefix(l1.name, l1.labelName, value)}
	l1.samples = append(l1.samples, s)

	l1.mutex.Unlock()
	return s
}

// With returns a dedicated Sample for a label combination. The values
// map to the names (in order) as defined at Must2LabelSample. With
// registers a new Sample if the combination hasn't been used before.
// Remember that each label combination represents a new time series,
// which can dramatically increase the amount of data stored.
func (l2 *Map2LabelSample) With(value1, value2 string) *Sample {
	hash := uint64(hashOffset)
	hash ^= uint64(len(value1))
	hash *= hashPrime
	for i := 0; i < len(value1); i++ {
		hash ^= uint64(value1[i])
		hash *= hashPrime
	}
	for i := 0; i < len(value2); i++ {
		hash ^= uint64(value2[i])
		hash *= hashPrime
	}

	l2.mutex.Lock()

	for i, h := range l2.labelHashes {
		if h == hash {
			hit := l2.samples[i]

			l2.mutex.Unlock()
			return hit
		}
	}

	l2.labelHashes = append(l2.labelHashes, hash)
	s := &Sample{prefix: format2LabelPrefix(l2.name, &l2.labelNames, value1, value2)}
	l2.samples = append(l2.samples, s)

	l2.mutex.Unlock()
	return s
}

// With returns a dedicated Sample for a label combination. The values
// map to the names (in order) as defined at Must3LabelSample. With
// registers a new Sample if the combination hasn't been used before.
// Remember that each label combination represents a new time series,
// which can dramatically increase the amount of data stored.
func (l3 *Map3LabelSample) With(value1, value2, value3 string) *Sample {
	hash := uint64(hashOffset)
	hash ^= uint64(len(value1))
	hash *= hashPrime
	for i := 0; i < len(value1); i++ {
		hash ^= uint64(value1[i])
		hash *= hashPrime
	}
	hash ^= uint64(len(value2))
	hash *= hashPrime
	for i := 0; i < len(value2); i++ {
		hash ^= uint64(value2[i])
		hash *= hashPrime
	}
	for i := 0; i < len(value3); i++ {
		hash ^= uint64(value3[i])
		hash *= hashPrime
	}

	l3.mutex.Lock()

	for i, h := range l3.labelHashes {
		if h == hash {
			hit := l3.samples[i]

			l3.mutex.Unlock()
			return hit
		}
	}

	l3.labelHashes = append(l3.labelHashes, hash)
	s := &Sample{prefix: format3LabelPrefix(l3.name, &l3.labelNames, value1, value2, value3)}
	l3.samples = append(l3.samples, s)

	l3.mutex.Unlock()
	return s
}

var valueEscapes = strings.NewReplacer("\n", `\n`, `"`, `\"`, `\`, `\\`)

func format1LabelPrefix(name, labelName, labelValue string) string {
	var buf strings.Builder
	buf.Grow(6 + len(name) + len(labelName) + len(labelValue))

	buf.WriteString(name)
	buf.WriteByte('{')
	buf.WriteString(labelName)
	buf.WriteString(`="`)
	valueEscapes.WriteString(&buf, labelValue)
	buf.WriteString(`"} `)

	return buf.String()
}

func format2LabelPrefix(name string, labelNames *[2]string, labelValue1, labelValue2 string) string {
	var buf strings.Builder
	buf.Grow(10 + len(name) + len(labelNames[0]) + len(labelNames[1]) + len(labelValue1) + len(labelValue2))

	buf.WriteString(name)
	buf.WriteByte('{')
	buf.WriteString(labelNames[0])
	buf.WriteString(`="`)
	valueEscapes.WriteString(&buf, labelValue1)
	buf.WriteString(`",`)
	buf.WriteString(labelNames[1])
	buf.WriteString(`="`)
	valueEscapes.WriteString(&buf, labelValue2)
	buf.WriteString(`"} `)

	return buf.String()
}

func format3LabelPrefix(name string, labelNames *[3]string, labelValue1, labelValue2, labelValue3 string) string {
	var buf strings.Builder
	buf.Grow(14 + len(name) + len(labelNames[0]) + len(labelNames[1]) + len(labelNames[2]) + len(labelValue1) + len(labelValue2) + len(labelValue3))

	buf.WriteString(name)
	buf.WriteByte('{')
	buf.WriteString(labelNames[0])
	buf.WriteString(`="`)
	valueEscapes.WriteString(&buf, labelValue1)
	buf.WriteString(`",`)
	buf.WriteString(labelNames[1])
	buf.WriteString(`="`)
	valueEscapes.WriteString(&buf, labelValue2)
	buf.WriteString(`",`)
	buf.WriteString(labelNames[2])
	buf.WriteString(`="`)
	valueEscapes.WriteString(&buf, labelValue3)
	buf.WriteString(`"} `)

	return buf.String()
}

func merge1LabelPrefix(prefix, labelName, labelValue string) string {
	var buf strings.Builder
	buf.Grow(4 + len(prefix) + len(labelName) + len(labelValue))

	buf.WriteString(prefix[:len(prefix)-2])
	buf.WriteByte(',')
	buf.WriteString(labelName)
	buf.WriteString(`="`)
	valueEscapes.WriteString(&buf, labelValue)
	buf.WriteString(`"} `)

	return buf.String()
}

func merge2LabelPrefix(prefix string, labelNames *[2]string, labelValue1, labelValue2 string) string {
	var buf strings.Builder
	buf.Grow(8 + len(prefix) + len(labelNames[0]) + len(labelNames[1]) + len(labelValue1) + len(labelValue2))

	buf.WriteString(prefix[:len(prefix)-2])
	buf.WriteByte(',')
	buf.WriteString(labelNames[0])
	buf.WriteString(`="`)
	valueEscapes.WriteString(&buf, labelValue1)
	buf.WriteString(`",`)
	buf.WriteString(labelNames[1])
	buf.WriteString(`="`)
	valueEscapes.WriteString(&buf, labelValue2)
	buf.WriteString(`"} `)

	return buf.String()
}

func merge3LabelPrefix(prefix string, labelNames *[3]string, labelValue1, labelValue2, labelValue3 string) string {
	var buf strings.Builder
	buf.Grow(12 + len(prefix) + len(labelNames[0]) + len(labelNames[1]) + len(labelNames[2]) + len(labelValue1) + len(labelValue2) + len(labelValue3))

	buf.WriteString(prefix[:len(prefix)-2])
	buf.WriteByte(',')
	buf.WriteString(labelNames[0])
	buf.WriteString(`="`)
	valueEscapes.WriteString(&buf, labelValue1)
	buf.WriteString(`",`)
	buf.WriteString(labelNames[1])
	buf.WriteString(`="`)
	valueEscapes.WriteString(&buf, labelValue2)
	buf.WriteString(`",`)
	buf.WriteString(labelNames[2])
	buf.WriteString(`="`)
	valueEscapes.WriteString(&buf, labelValue3)
	buf.WriteString(`"} `)

	return buf.String()
}
