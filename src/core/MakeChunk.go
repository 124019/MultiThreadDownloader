package core

type CHUNK_Range struct {
	No     int
	St     int
	Ed     int
	Header string
}


func InitChunksIndex(total_size int, chunk_size int) []CHUNK_Range {
	var result []CHUNK_Range
	for st_byte := 0; st_byte < total_size; st_byte += chunk_size {
		ed_byte := st_byte + chunk_size - 1
		last_plus := 0
		if ed_byte > total_size {
			ed_byte = total_size - 1 // first is 0, last is total_size - 1
			last_plus = 1
		}
		result = append(result, CHUNK_Range{
			No: st_byte / chunk_size + last_plus, // Because of Golang's "round towards zero" while "Cannot devide evenly", the last chunk's number shoud plus 1
			St: st_byte,
			Ed: ed_byte,
			Header: "",
		})
	}
	// fmt.Println(result)
	return result
}