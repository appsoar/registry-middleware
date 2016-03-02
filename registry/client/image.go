package client

import (
//	""
)

//注意`json:"*"中的*必须要和json数据完全一致,否则会出现某些域无数据的情况`
type Jwk struct {
	Crv string `json:"crv"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

type Header struct {
	Jwk Jwk    `json:"jwk"`
	Alg string `json:"alg"`
}
type Signatures struct {
	Header    Header `json:"header"`
	Signature string `json:"signature"`
	Protected string `json:"protected"`
}

type History struct {
	V1Compatibility string `json:"v1Compatibility"`
}

type FsLayers struct {
	BlobSum string `json:"blobSum"`
}

/*镜像manifest数据*/
type Manifests struct {
	SchemaVersion int          `json:"schemaVersion"`
	Name          string       `json:"name"`
	Tag           string       `json:"tag"`
	Architecture  string       `json:"architecture"`
	FsLayers      []FsLayers   `json:"fsLayers"`
	History       []History    `json:"history"`
	Signatures    []Signatures `json:"signatures"`
}
