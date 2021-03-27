// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/27

package e

var (
	ParsePayloadFailed = payload{}
	Forbidden          = forbidden{}
	Unknown            = unknown{}
	GetTokenFail       = getTokenFail{}
	GetAbsPathFail     = absPath{}
	MkPathFail         = mkPath{}
	NotFound           = notFound{}
	OperateDbFail      = operateDb{}
	StatisticsFileFail = statisticsFile{}
	ParseUrlFail       = parseUrl{}
	DeleteBucketFail   = deleteBucket{}
	OpenFileFail       = openFile{}
	DeleteFileFail     = deleteFile{}
	GenerateTokenFail  = generateToken{}
	GenerateSignFail   = generateSign{}
	Recovery           = recovery{}
	InvalidEmail       = invalidEmail{}
	MarshalFail        = marshal{}
	UnmarshalFail      = unmarshal{}
)
