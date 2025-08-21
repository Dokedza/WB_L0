module github.com/dokedza/WB_L0

go 1.24.2

replace go1f => ./

require go1f v0.0.0

require (
	github.com/lib/pq v1.10.9
	github.com/segmentio/kafka-go v0.4.48
)

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
)
