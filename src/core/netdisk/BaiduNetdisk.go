package netdisk
import (
	"MultiThreadDownloader/src/core"
	"fmt"
)
func BaiduNetdiskChunkRangeInit(intervals []core.CHUNK_Range) []core.CHUNK_Range {
	var result []core.CHUNK_Range
	for i := range intervals {
		Range := fmt.Sprintf("bytes=%d-%d", intervals[i].St, intervals[i].Ed)
		result = append(result, core.CHUNK_Range{
			St: intervals[i].St,
			Ed: intervals[i].Ed,
			Header: Range,
		})
	}
	// fmt.Println(result)
	return result
}