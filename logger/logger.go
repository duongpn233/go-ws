package logger

import "encoding/json"

type LogRequest struct {
	Lt   string  `json:"lt"`
	Url  string  `json:"url"`
	Urp  string  `json:"urp"`
	Urq  string  `json:"urq"`
	Rt   float32 `json:"rt"`
	St   int     `json:"st"`
	Mt   string  `json:"mt"`
	Rmip string  `json:"rmip"`
	Cip  string  `json:"cip"`
	Bbs  int     `json:"bbs"`
	Cl   int64   `json:"cl"`
	RF   string  `json:"RF"`
	Au   string  `json:"au"`
	Host string  `json:"host"`
	Sn   string  `json:"sn"`
	Tl   string  `json:"tl"`
	Rid  string  `json:"rid"`
	Uid  string  `json:"uid"`
	Usrc string  `json:"usrc"`
}

type LogErr struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewLogError(msg string, err string) LogErr {
	return LogErr{
		Message: msg,
		Error:   err,
	}
}

func (log LogErr) ToJson() string {
	jsonLog, _ := json.Marshal(log)
	return string(jsonLog)
}
