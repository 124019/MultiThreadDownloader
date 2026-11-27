package netdisk
import (
	"MultiThreadDownloader/downloader/core"
	"fmt"
)
func ChunkRangeInit(intervals []core.CHUNK_Range) []core.CHUNK_Range { // BaiduNetdisk is supported
	var result []core.CHUNK_Range
	for i := range intervals {
		Range := fmt.Sprintf("bytes=%d-%d", intervals[i].St, intervals[i].Ed)
		result = append(result, core.CHUNK_Range{
			St: intervals[i].St,
			Ed: intervals[i].Ed,
			No: intervals[i].No,
			Header: Range,
		})
	}
	// fmt.Println(result)
	return result
}