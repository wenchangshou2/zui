package form
type ActiveWindowByPidRequestBody struct{
	Data struct {
		Pid int
		X int
		Y int
		Width int
		Height int
	}
}
