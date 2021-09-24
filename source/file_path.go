package source

import "fmt"

func generateFilePath(start, end int) string {
	return fmt.Sprintf("%d_%d.csv", start, end)
}
