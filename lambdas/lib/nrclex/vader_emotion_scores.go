package nrclex

type VaderEmotionScore struct {
	Compound float64 `json:"compound"`
	Neg      float64 `json:"neg"`
	Neu      float64 `json:"neu"`
	Pos      float64 `json:"pos"`
}
