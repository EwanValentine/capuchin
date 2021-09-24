package source

// S3 source
type S3 struct {
	Bucket string
}

// NewS3 -
func NewS3(bucket string) *S3 {
	return &S3{Bucket: bucket}
}

// Load - raw data is written into a single S3 bucket,
// the files are split into chunks, each file's name consists
// of the date of the first item in the batch, and the last.
// E.g. 20190101_20190201.csv
func (s *S3) Load(start, end string) error {
	return nil
}
