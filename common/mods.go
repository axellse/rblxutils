package common

type RcmfFile struct {
	Spec  string     `json:"spec"`
	Rules []RcmfRule `json:"rules"`
}

type Sources struct {
	Expressions []string `json:"expressions"`
	Ids []int `json:"ids"`
	Types []int `json:"types"`
	Files []string `json:"files"`
}

type RcmfRule struct {
	Sources  Sources `json:"sources"`
	Data RcmfData `json:"data"`
}

type RcmfData struct {
	Blob  []byte `json:"blob"`
	Key   string `json:"key"` //dots can be used for nesting (eg. Settings.ContentFolder)
	Value any `json:"value"`
}