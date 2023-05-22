package gater

import (
	"encoding/xml"
	"net/http"
	"strings"

	"github.com/bnb-chain/greenfield-storage-provider/base/types/gfsperrors"
	"github.com/bnb-chain/greenfield-storage-provider/base/types/gfsptask"
	coremodule "github.com/bnb-chain/greenfield-storage-provider/core/module"
	"github.com/bnb-chain/greenfield-storage-provider/model"
	"github.com/bnb-chain/greenfield-storage-provider/model/job"
	"github.com/bnb-chain/greenfield-storage-provider/pkg/log"
	servicetypes "github.com/bnb-chain/greenfield-storage-provider/service/types"
	"github.com/bnb-chain/greenfield-storage-provider/util"
	storagetypes "github.com/bnb-chain/greenfield/x/storage/types"
)

// putObjectHandler handles the upload object request.
func (g *GateModular) putObjectHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		reqCtx     *RequestContext
		authorized bool
		objectInfo *storagetypes.ObjectInfo
		params     *storagetypes.Params
	)

	defer func() {
		reqCtx.Cancel()
		if err != nil {
			reqCtx.SetError(gfsperrors.MakeGfSpError(err))
			reqCtx.SetHttpCode(int(gfsperrors.MakeGfSpError(err).GetHttpStatusCode()))
			MakeErrorResponse(w, gfsperrors.MakeGfSpError(err))
		} else {
			reqCtx.SetHttpCode(http.StatusOK)
		}
		log.CtxDebugw(reqCtx.Context(), reqCtx.String())
	}()

	reqCtx, err = NewRequestContext(r)
	if err != nil {
		return
	}
	if reqCtx.NeedVerifyAuthorizer() {
		authorized, err = g.baseApp.GfSpClient().VerifyAuthorize(reqCtx.Context(),
			coremodule.AuthOpTypePutObject, reqCtx.Account(), reqCtx.bucketName, reqCtx.objectName)
		if err != nil {
			log.CtxErrorw(reqCtx.Context(), "failed to verify authorize", "error", err)
			return
		}
		if !authorized {
			log.CtxErrorw(reqCtx.Context(), "no permission to operator")
			err = ErrNoPermission
			return
		}
	}

	objectInfo, err = g.baseApp.Consensus().QueryObjectInfo(reqCtx.Context(), reqCtx.bucketName, reqCtx.objectName)
	if err != nil {
		log.CtxErrorw(reqCtx.Context(), "failed to get object info from consensus", "error", err)
		err = ErrConsensus
		return
	}
	if objectInfo.GetPayloadSize() == 0 || objectInfo.GetPayloadSize() > g.maxPayloadSize {
		log.CtxErrorw(reqCtx.Context(), "failed to put object payload size is zero")
		err = ErrInvalidPayloadSize
		return
	}
	params, err = g.baseApp.Consensus().QueryStorageParams(reqCtx.Context())
	if err != nil {
		log.CtxErrorw(reqCtx.Context(), "failed to get storage params from consensus", "error", err)
		err = ErrConsensus
		return
	}
	task := &gfsptask.GfSpUploadObjectTask{}
	task.InitUploadObjectTask(objectInfo, params, g.baseApp.TaskTimeout(task, objectInfo.GetPayloadSize()))
	ctx := log.WithValue(reqCtx.Context(), log.CtxKeyTask, task.Key().String())
	err = g.baseApp.GfSpClient().UploadObject(ctx, task, r.Body)
	if err != nil {
		log.CtxErrorw(ctx, "failed to upload payload data", "error", err)
	}
	log.CtxDebugw(ctx, "succeed to upload payload data")
}

func parseRange(rangeStr string) (bool, int64, int64) {
	if rangeStr == "" {
		return false, -1, -1
	}
	rangeStr = strings.ToLower(rangeStr)
	rangeStr = strings.ReplaceAll(rangeStr, " ", "")
	if !strings.HasPrefix(rangeStr, "bytes=") {
		return false, -1, -1
	}
	rangeStr = rangeStr[len("bytes="):]
	if strings.HasSuffix(rangeStr, "-") {
		rangeStr = rangeStr[:len(rangeStr)-1]
		rangeStart, err := util.StringToUint64(rangeStr)
		if err != nil {
			return false, -1, -1
		}
		return true, int64(rangeStart), -1
	}
	pair := strings.Split(rangeStr, "-")
	if len(pair) == 2 {
		rangeStart, err := util.StringToUint64(pair[0])
		if err != nil {
			return false, -1, -1
		}
		rangeEnd, err := util.StringToUint64(pair[1])
		if err != nil {
			return false, -1, -1
		}
		return true, int64(rangeStart), int64(rangeEnd)
	}
	return false, -1, -1
}

// getObjectHandler handles the download object request.
func (g *GateModular) getObjectHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		reqCtx     *RequestContext
		authorized bool
		objectInfo *storagetypes.ObjectInfo
		bucketInfo *storagetypes.BucketInfo
		params     *storagetypes.Params
	)
	defer func() {
		reqCtx.Cancel()
		if err != nil {
			reqCtx.SetError(gfsperrors.MakeGfSpError(err))
			reqCtx.SetHttpCode(int(gfsperrors.MakeGfSpError(err).GetHttpStatusCode()))
			MakeErrorResponse(w, gfsperrors.MakeGfSpError(err))
		} else {
			reqCtx.SetHttpCode(http.StatusOK)
		}
		log.CtxDebugw(reqCtx.Context(), reqCtx.String())
	}()

	reqCtx, err = NewRequestContext(r)
	if err != nil {
		return
	}
	if reqCtx.NeedVerifyAuthorizer() {
		authorized, err = g.baseApp.GfSpClient().VerifyAuthorize(reqCtx.Context(),
			coremodule.AuthOpTypeGetObject, reqCtx.Account(), reqCtx.bucketName, reqCtx.objectName)
		if err != nil {
			log.CtxErrorw(reqCtx.Context(), "failed to verify authorize", "error", err)
			return
		}
		if !authorized {
			log.CtxErrorw(reqCtx.Context(), "no permission to operator")
			return
		}
	}

	objectInfo, err = g.baseApp.Consensus().QueryObjectInfo(reqCtx.Context(), reqCtx.bucketName, reqCtx.objectName)
	if err != nil {
		log.CtxErrorw(reqCtx.Context(), "failed to get object info from consensus", "error", err)
		err = ErrConsensus
		return
	}
	bucketInfo, err = g.baseApp.Consensus().QueryBucketInfo(reqCtx.Context(), objectInfo.GetBucketName())
	if err != nil {
		log.CtxErrorw(reqCtx.Context(), "failed to get bucket info from consensus", "error", err)
		err = ErrConsensus
		return
	}
	params, err = g.baseApp.Consensus().QueryStorageParams(reqCtx.Context())
	if err != nil {
		log.CtxErrorw(reqCtx.Context(), "failed to get storage params from consensus", "error", err)
		err = ErrConsensus
		return
	}

	isRange, rangeStart, rangeEnd := parseRange(reqCtx.request.Header.Get(model.RangeHeader))
	if isRange && (rangeEnd < 0 || rangeEnd >= int64(objectInfo.GetPayloadSize())) {
		rangeEnd = int64(objectInfo.GetPayloadSize()) - 1
	}
	if isRange && (rangeStart < 0 || rangeEnd < 0 || rangeStart > rangeEnd) {
		err = ErrInvalidRange
		return
	}
	var low int64
	var high int64
	if isRange {
		low = rangeStart
		high = rangeEnd
	} else {
		low = 0
		high = int64(objectInfo.GetPayloadSize()) - 1
	}

	task := &gfsptask.GfSpDownloadObjectTask{}
	task.InitDownloadObjectTask(objectInfo, bucketInfo, params, g.baseApp.TaskPriority(task), reqCtx.Account(),
		low, high, g.baseApp.TaskTimeout(task, uint64(high-low+1)), g.baseApp.TaskMaxRetry(task))
	data, err := g.baseApp.GfSpClient().GetObject(reqCtx.Context(), task)
	if err != nil {
		log.CtxErrorw(reqCtx.Context(), "failed to download object", "error", err)
		return
	}
	w.Write(data)
	w.Header().Set(model.ContentTypeHeader, objectInfo.GetContentType())
	if isRange {
		w.Header().Set(model.ContentRangeHeader, "bytes "+util.Uint64ToString(uint64(low))+
			"-"+util.Uint64ToString(uint64(high)))
	} else {
		w.Header().Set(model.ContentLengthHeader, util.Uint64ToString(objectInfo.GetPayloadSize()))
	}
	log.CtxDebugw(reqCtx.Context(), "succeed to download object")
}

func (g *GateModular) queryUploadProgressHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		reqCtx     *RequestContext
		authorized bool
		objectInfo *storagetypes.ObjectInfo
		jobState   int32
	)
	defer func() {
		reqCtx.Cancel()
		if err != nil {
			reqCtx.SetError(gfsperrors.MakeGfSpError(err))
			reqCtx.SetHttpCode(int(gfsperrors.MakeGfSpError(err).GetHttpStatusCode()))
			MakeErrorResponse(w, gfsperrors.MakeGfSpError(err))
		} else {
			reqCtx.SetHttpCode(http.StatusOK)
		}
		log.CtxDebugw(reqCtx.Context(), reqCtx.String())
	}()

	reqCtx, err = NewRequestContext(r)
	if err != nil {
		return
	}
	if reqCtx.NeedVerifyAuthorizer() {
		authorized, err = g.baseApp.GfSpClient().VerifyAuthorize(reqCtx.Context(),
			coremodule.AuthOpTypeGetUploadingState, reqCtx.Account(), reqCtx.bucketName, reqCtx.objectName)
		if err != nil {
			log.CtxErrorw(reqCtx.Context(), "failed to verify authorize", "error", err)
			return
		}
		if !authorized {
			log.CtxErrorw(reqCtx.Context(), "no permission to operator")
			err = ErrNoPermission
			return
		}
	}

	objectInfo, err = g.baseApp.Consensus().QueryObjectInfo(reqCtx.Context(), reqCtx.bucketName, reqCtx.objectName)
	if err != nil {
		log.CtxErrorw(reqCtx.Context(), "failed to get object info from consensus", "error", err)
		err = ErrConsensus
		return
	}

	jobState, err = g.baseApp.GfSpClient().GetUploadObjectState(reqCtx.Context(), objectInfo.Id.Uint64())
	if err != nil {
		log.CtxErrorw(reqCtx.Context(), "failed to get uploading job state", "error", err)
		return
	}
	jobStateDescription := job.StateToDescription(servicetypes.JobState(jobState))

	var xmlInfo = struct {
		XMLName             xml.Name `xml:"QueryUploadProgress"`
		Version             string   `xml:"version,attr"`
		ProgressDescription string   `xml:"ProgressDescription"`
	}{
		Version:             model.GnfdResponseXMLVersion,
		ProgressDescription: jobStateDescription,
	}
	xmlBody, err := xml.Marshal(&xmlInfo)
	if err != nil {
		log.Errorw("failed to marshal xml", "error", err)
		err = ErrEncodeResponse
		return
	}
	w.Header().Set(model.ContentTypeHeader, model.ContentTypeXMLHeaderValue)
	if _, err = w.Write(xmlBody); err != nil {
		log.Errorw("failed to write body", "error", err)
		err = ErrEncodeResponse
		return
	}
	log.Debugw("succeed to query upload progress", "xml_info", xmlInfo)
}