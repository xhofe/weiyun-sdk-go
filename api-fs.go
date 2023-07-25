package weiyunsdkgo

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/go-resty/resty/v2"
)

type File struct {
	FileID    string    `json:"file_id"`
	FileName  string    `json:"filename"`
	FileSize  int64     `json:"file_size"`
	FileSha   string    `json:"file_sha"`
	FileCtime TimeStamp `json:"file_ctime"`
	FileMtime TimeStamp `json:"file_mtime"`

	ExtInfo struct {
		ThumbURL string `json:"thumb_url"`
	} `json:"ext_info"`
}

type Folder struct {
	DirKey   string    `json:"dir_key"`
	DirName  string    `json:"dir_name"`
	DirCtime TimeStamp `json:"dir_ctime"`
	DirMtime TimeStamp `json:"dir_mtime"`
}

type DiskDirBatchListData struct {
	DirList []DiskListData `json:"dir_list"`
}

type DiskListData struct {
	DirList        []Folder `json:"dir_list"`
	FileList       []File   `json:"file_list"`
	PdirKey        string   `json:"pdir_key"`
	FinishFlag     bool     `json:"finish_flag"`
	TotalDirCount  int      `json:"total_dir_count"`
	TotalFileCount int      `json:"total_file_count"`
	TotalSpace     int      `json:"total_space"`
	HideDirCount   int      `json:"hide_dir_count"`
	HideFileCount  int      `json:"hide_file_count"`
}

// 查询文件、文件夹
// 数量限制 500
func (c *WeiYunClient) DiskDirFileList(dirKey string, paramOption []ParamOption, opts ...RestyOption) (*DiskListData, error) {
	param := Json{
		//"pdir_key": pdirKey,
		"dir_key": dirKey,

		"start": 0,
		"count": 500,

		"sort_field":    2,
		"reverse_order": false,

		"get_type": 0,

		"get_abstract_url":    false,
		"get_dir_detail_info": false,
	}
	ApplyParamOption(param, paramOption...)

	var resp DiskListData
	_, err := c.WeiyunQdiskRequest("DiskDirList", 2208, param, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type DiskDirFileBatchListParam struct {
	DirKey      string
	ParamOption []ParamOption
}

// 批量查询文件、文件夹
func (c *WeiYunClient) DiskDirFileBatchList(batchParam []DiskDirFileBatchListParam, commonParamOption []ParamOption, opts ...RestyOption) ([]DiskListData, error) {
	param := Json{
		//"pdir_key": pdirKey,
		"dir_list": MustSliceConvert(batchParam, func(b DiskDirFileBatchListParam) Json {
			dParam := Json{
				"dir_key": b.DirKey,

				"start": 0,
				"count": 500,

				"sort_field":    2,
				"reverse_order": false,

				"get_type": 0,

				"get_abstract_url":    false,
				"get_dir_detail_info": false,
			}
			ApplyParamOption(dParam, commonParamOption...)
			ApplyParamOption(dParam, b.ParamOption...)
			return dParam
		}),
	}

	var resp DiskDirBatchListData
	_, err := c.WeiyunQdiskRequest("DiskDirBatchList", 2209, param, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return resp.DirList, nil
}

type FolderParam struct {
	PPdirKey string `json:"ppdir_key,omitempty"` // 父父目录ID(打包下载忽略)
	PdirKey  string `json:"pdir_key,omitempty"`  // 父目录ID(打包下载忽略)
	DirKey   string `json:"dir_key,omitempty"`   // 目录ID(创建忽略)
	DirName  string `json:"dir_name,omitempty"`  // 目录名称
}

// 文件夹重命名
func (c *WeiYunClient) DiskDirAttrModify(dParam FolderParam, newDirName string, opts ...RestyOption) error {
	param := Json{
		"ppdir_key":    dParam.PPdirKey,
		"pdir_key":     dParam.PdirKey,
		"dir_key":      dParam.DirKey,
		"src_dir_name": dParam.DirName,
		"dst_dir_name": newDirName,
	}
	_, err := c.WeiyunQdiskClientRequest("DiskDirAttrModify", 2605, param, nil, opts...)
	return err
}

// 文件夹删除
func (c *WeiYunClient) DiskDirDelete(dParam FolderParam, opts ...RestyOption) error {
	param := Json{
		"dir_list": []FolderParam{dParam},
	}
	_, err := c.WeiyunQdiskClientRequest("DiskDirFileBatchDeleteEx", 2509, param, nil, opts...)
	return err
}

// 文件夹移动
func (c *WeiYunClient) DiskDirMove(srcParam FolderParam, dstParam FolderParam, opts ...RestyOption) error {
	param := Json{
		"src_ppdir_key": srcParam.PPdirKey,
		"src_pdir_key":  srcParam.PdirKey,
		"dir_list":      []FolderParam{srcParam},
		"dst_ppdir_key": dstParam.PPdirKey,
		"dst_pdir_key":  dstParam.PdirKey,
	}
	_, err := c.WeiyunQdiskClientRequest("DiskDirFileBatchMove", 2618, param, nil, opts...)
	return err
}

// 文件夹创建
func (c *WeiYunClient) DiskDirCreate(dParam FolderParam, opts ...RestyOption) (*Folder, error) {
	param := Json{
		"ppdir_key":         dParam.PPdirKey,
		"pdir_key":          dParam.PdirKey,
		"dir_name":          dParam.DirName,
		"file_exist_option": 2,
		"create_type":       1,
	}
	var folder Folder
	_, err := c.WeiyunQdiskClientRequest("DiskDirCreate", 2614, param, &folder, opts...)
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

type FileParam struct {
	PPdirKey string `json:"ppdir_key,omitempty"` // 父父目录ID(打包下载忽略)
	PdirKey  string `json:"pdir_key,omitempty"`  // 父目录ID
	FileID   string `json:"file_id,omitempty"`   // 文件ID
	FileName string `json:"filename,omitempty"`  // 文件名称(打包下载忽略)
}

// 文件重命名
func (c *WeiYunClient) DiskFileRename(fParam FileParam, newFileName string, opts ...RestyOption) error {
	param := Json{
		"ppdir_key":    fParam.PPdirKey,
		"pdir_key":     fParam.PdirKey,
		"file_id":      fParam.FileID,
		"src_filename": fParam.FileName,
		"filename":     newFileName,
	}
	_, err := c.WeiyunQdiskClientRequest("DiskFileRename", 2605, param, nil, opts...)
	return err
}

// 文件删除
func (c *WeiYunClient) DiskFileDelete(fParam FileParam, opts ...RestyOption) error {
	param := Json{
		"file_list": []FileParam{fParam},
	}

	_, err := c.WeiyunQdiskClientRequest("DiskDirFileBatchDeleteEx", 2509, param, nil, opts...)
	return err
}

// 文件移动
func (c *WeiYunClient) DiskFileMove(srcParam FileParam, dstParam FolderParam, opts ...RestyOption) error {
	param := Json{
		"src_ppdir_key": srcParam.PPdirKey,
		"src_pdir_key":  srcParam.PdirKey,
		"file_list":     []FileParam{srcParam},
		"dst_ppdir_key": dstParam.PPdirKey,
		"dst_pdir_key":  dstParam.PdirKey,
	}
	_, err := c.WeiyunQdiskClientRequest("DiskDirFileBatchMove", 2618, param, nil, opts...)
	return err
}

type DiskFileDownloadData struct {
	CookieName  string `json:"cookie_name"`
	CookieValue string `json:"cookie_value"`

	DownloadUrl string `json:"download_url"`
}

// 文件下载
func (c *WeiYunClient) DiskFileDownload(fParam FileParam, opts ...RestyOption) (*DiskFileDownloadData, error) {
	param := Json{
		"file_list":     []FileParam{fParam},
		"download_type": 0,
	}
	var resp DiskFileDownloadData
	_, err := c.WeiyunQdiskClientRequest("DiskFileBatchDownload", 2402, param, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type DiskFilePackageDownloadParam struct {
	PdirKey  string `json:"pdir_key"`
	PdirList any    `json:"pdir_list"`
}

// 文件打包下载
func (c *WeiYunClient) DiskFilePackageDownload(param []DiskFilePackageDownloadParam, zipFilename string, opts ...RestyOption) (*DiskFileDownloadData, error) {
	param_ := Json{
		"pdir_list": MustSliceConvert(param, func(p DiskFilePackageDownloadParam) Json {
			list := batchParamConvert(p.PdirList)
			list["pdir_key"] = p.PdirKey
			return list
		}),
		"zip_filename": zipFilename,
	}

	var resp DiskFileDownloadData
	_, err := c.WeiyunQdiskRequest("DiskFilePackageDownload", 2403, param_, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func batchParamConvert(param any) Json {
	param_ := Json{}
	switch v := param.(type) {
	case FileParam, *FileParam:
		param_["file_list"] = []any{param}
	case FolderParam, *FolderParam:
		param_["dir_list"] = []any{param}
	case []FileParam, []*FileParam:
		param_["file_list"] = param
	case []FolderParam, []*FolderParam:
		param_["dir_list"] = param
	case []any:
		var fileList []any
		var dirList []any
		for _, vv := range v {
			switch vv.(type) {
			case FileParam, *FileParam:
				fileList = append(fileList, vv)
			case []FolderParam, []*FolderParam:
				dirList = append(dirList, vv)
			}
		}
		param_["file_list"] = fileList
		param_["dir_list"] = dirList
	case Json:
		param_ = v
	}
	return param_
}

type UploadAuth struct {
	UploadKey string `json:"upload_key"`
	Ex        string `json:"ex"`
}

type UploadChannelData struct {
	ID     int `json:"id"`
	Offset int `json:"offset"`
	Len    int `json:"len"`
}

type UpdloadFileParam struct {
	PPdirKey string // 父父目录ID
	PdirKey  string // 父目录ID

	FileName string
	FileSize int64
	File     io.ReadSeeker

	ChannelCount    int // 上传通道数量
	FileExistOption int // 文件存在时操作 6，4
}

type PreUploadData struct {
	CommonUploadRsp File `json:"common_upload_rsp"`

	FileExist       bool `json:"file_exist"`        // 文件是否存在
	ExtChannelCount int  `json:"ext_channel_count"` // 存在通道数

	UploadScr int `json:"upload_scr"` // 未知

	// 上传授权
	UploadAuth

	ChannelList []UploadChannelData `json:"channel_list"` // 上传通道

	Speedlimit int `json:"speedlimit"` // 上传速度限制
	FlowState  int `json:"flow_state"`

	UploadState     int `json:"upload_state"`      // 上传状态
	UploadedDataLen int `json:"uploaded_data_len"` // 已经上传的长度
}

func (c *WeiYunClient) PreUpload(ctx context.Context, param UpdloadFileParam, opts ...RestyOption) (*PreUploadData, error) {
	const blockSize = 1024 * 1024
	lastBlockSize := (param.FileSize % blockSize) // 最后一块大小
	if lastBlockSize == 0 {
		lastBlockSize = blockSize
	}
	beforeBlockSize := param.FileSize - lastBlockSize // 除去最后一刻的大小

	type BlockInfo struct {
		Sha1   string `json:"sha"`
		Offset int64  `json:"offset"`
		Size   int64  `json:"size"`
	}

	// before
	hash := sha1.New()
	blockInfoList := make([]BlockInfo, 0, (param.FileSize/blockSize)+1)
	for offset := int64(0); offset < beforeBlockSize; offset += blockSize {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		if _, err := io.CopyN(hash, param.File, blockSize); err != nil {
			return nil, err
		}
		blockInfoList = append(blockInfoList, BlockInfo{
			Sha1:   hex.EncodeToString(GetSha1State(hash)),
			Offset: offset,
			Size:   blockSize,
		})
	}

	// between
	var checkData string
	var checkSha1 string
	var buf [128]byte
	if lastBlockSize > 128 {
		if n, err := io.CopyN(hash, param.File, lastBlockSize-128); err != nil {
			fmt.Println(n, err)
			return nil, err
		}
		checkSha1 = hex.EncodeToString(GetSha1State(hash))
	}

	// after
	n, err := io.ReadFull(io.TeeReader(param.File, hash), buf[:])
	if err != nil {
		return nil, err
	}

	fileHash := hex.EncodeToString(hash.Sum(nil))
	checkData = base64.StdEncoding.EncodeToString(buf[:n])
	if checkSha1 == "" {
		checkSha1 = fileHash
	}

	blockInfoList = append(blockInfoList, BlockInfo{
		Sha1:   fileHash,
		Offset: beforeBlockSize,
		Size:   lastBlockSize,
	})

	paramJson := Json{
		"common_upload_req": Json{
			"ppdir_key":         param.PPdirKey,
			"pdir_key":          param.PdirKey,
			"file_size":         param.FileSize,
			"filename":          param.FileName,
			"file_exist_option": param.FileExistOption,
			"use_mutil_channel": true,
		},
		"upload_scr":      0,
		"channel_count":   param.ChannelCount, //
		"check_sha":       checkSha1,
		"check_data":      checkData,
		"block_size":      blockSize,
		"block_info_list": blockInfoList,
	}

	var resp struct {
		Body PreUploadData `json:"weiyun.PreUploadMsgRsp_body"`
	}
	_, err = c.UploadRequest("PreUpload", 247120, paramJson, &resp, append([]RestyOption{func(request *resty.Request) { request.SetContext(ctx) }}, opts...)...)
	if err != nil {
		return nil, err
	}
	resp.Body.CommonUploadRsp.FileSha = fileHash
	resp.Body.CommonUploadRsp.FileSize = param.FileSize
	return &resp.Body, nil
}

type AddChannelData struct {
	// 源通道数量
	OrigChannelCount int `json:"orig_channel_count"`
	// 当前通道数量
	FinalChannelCount int `json:"final_channel_count"`

	// 源通道信息
	OrigChannels []struct {
		ID     int `json:"id"`
		Offset int `json:"offset"`
		Len    int `json:"len"`
	} `json:"orig_channels"`
	// 增加通道信息
	AddChannels []struct {
		ID     int `json:"id"`
		Offset int `json:"offset"`
		Len    int `json:"len"`
	} `json:"channels"`
}

// 增加上传通道
func (c *WeiYunClient) AddUploadChannel(origChannelCount, destChannelCount int, auth UploadAuth, opts ...RestyOption) (*AddChannelData, error) {
	param := Json{
		"upload_key": auth.UploadKey,
		"ex":         auth.Ex,

		"orig_channel_count": origChannelCount,
		"dest_channel_count": destChannelCount,

		"speed": 4303,
	}

	var resp struct {
		Body AddChannelData `json:"weiyun.AddChannelMsgRsp_body"`
	}
	_, err := c.UploadRequest("AddChannel", 247122, param, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp.Body, nil
}

type UploadPieceData struct {
	Channel UploadChannelData `json:"channel"` // 下一个上传通道
	Ex      string            `json:"ex"`

	UploadState int `json:"upload_state"` // 上传状态
	FlowState   int `json:"flow_state"`
}

func (c *WeiYunClient) UploadFile(ctx context.Context, channel UploadChannelData, auth UploadAuth, r io.Reader, opts ...RestyOption) (*UploadPieceData, error) {
	param := Json{
		"upload_key": auth.UploadKey,
		"ex":         auth.Ex,
		"channel":    channel,
	}
	var resp struct {
		Body UploadPieceData `json:"weiyun.AddChannelMsgRsp_body"`
	}
	_, err := c.UploadRequest("UploadPiece", 247121, param, &resp, append([]RestyOption{func(request *resty.Request) {
		request.SetFileReader("upload", "blob", r)
	}}, opts...)...)
	if err != nil {
		return nil, err
	}
	return &resp.Body, nil
}
