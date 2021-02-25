package zeroapi

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/zerogo-hub/zero-helper/file"
)

// File 文件相关
type File interface {
	// File 获取上传文件信息
	File(key string) (multipart.File, *multipart.FileHeader, error)

	// Files 获取上传文件信息，可能有多个文件
	// cbs 在存盘前修改 multipart.FileHeader
	Files(destDirectory string, cbs ...func(Context, *multipart.FileHeader)) (int64, error)

	// DownloadFile 下载文件
	// path 文件路径
	// filename 文件名称
	DownloadFile(path string, filename ...string)
}

// upload 从临时文件夹或者内存中写入到指定位置的文件夹中
func upload(dest string, header *multipart.FileHeader) (int64, error) {
	// 打开临时文件或者内存中的文件内容
	src, err := header.Open()
	if err != nil {
		return 0, err
	}
	defer src.Close()

	file, err := os.OpenFile(filepath.Join(dest, header.Filename), os.O_WRONLY|os.O_CREATE, os.FileMode(0666))
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return io.Copy(file, src)
}

func (ctx *context) File(key string) (multipart.File, *multipart.FileHeader, error) {
	if err := ctx.req.ParseMultipartForm(ctx.app.FileMaxMemory()); err != nil {
		return nil, nil, err
	}

	return ctx.req.FormFile(key)
}

func (ctx *context) Files(destDirectory string, cbs ...func(Context, *multipart.FileHeader)) (int64, error) {
	if err := ctx.req.ParseMultipartForm(ctx.app.FileMaxMemory()); err != nil {
		return 0, err
	}

	// MultipartForm: 需要先调用 ParseMultipartForm，
	// including file uploads

	if ctx.req.MultipartForm != nil {
		if f := ctx.req.MultipartForm.File; f != nil {
			var l int64
			for _, files := range f {
				for _, file := range files {
					for _, cb := range cbs {
						cb(ctx, file)
					}
					length, err := upload(destDirectory, file)
					if err != nil {
						return 0, err
					}
					l += length
				}
			}
			return l, nil
		}
	}

	return 0, http.ErrMissingFile
}

func (ctx *context) DownloadFile(path string, filename ...string) {
	if !file.IsExist(path) {
		http.ServeFile(ctx.res, ctx.req, path)
		return
	}

	fname := ""
	if len(filename) > 0 && filename[0] != "" {
		fname = filename[0]
	} else {
		fname = file.NameRand(path, 8)
	}

	ctx.AddHeader("Content-Disposition", "attachment; filename="+fname)
	ctx.AddHeader("Content-Description", "File Transfer")
	ctx.AddHeader("Content-Type", "app/octet-stream")
	ctx.AddHeader("Content-Transfer-Encoding", "binary")
	ctx.AddHeader("Expires", "0")
	ctx.AddHeader("Cache-Control", "must-revalidate")
	ctx.AddHeader("Pragma", "public")
	http.ServeFile(ctx.res, ctx.req, path)
}
