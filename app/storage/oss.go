package storage

import (
	"bytes"
	"io"
	"mime"
	"path/filepath"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
	"oss/app/conf"
)

// OSSStorage 阿里云oss
type OSSStorage struct {
	Bucket     *oss.Bucket
	c          *conf.AliYun
	bucketName string
}

func NewOSS() (OSSStorage, error) {
	os := OSSStorage{}
	client, err := oss.New(conf.Settings.Storage.EndPoint, conf.Settings.Storage.AccessKeyID, conf.Settings.Storage.AccessKeySecret)
	if err != nil {
		return os, err
	}
	os.Bucket, err = client.Bucket(conf.Settings.Storage.DefaultBucket)
	if err != nil {
		return os, err
	}
	os.c = &conf.Settings.Storage
	os.bucketName = conf.Settings.Storage.DefaultBucket
	return os, nil
}

func (o OSSStorage) Get(objectKey string) ([]byte, error) {
	content, err := o.Bucket.GetObject(objectKey)
	if err != nil {
		return nil, errors.Wrap(err, "GetObject failed")
	}
	return io.ReadAll(content)
}

func (o OSSStorage) Put(savePathName string, data []byte) error {
	contentType, err := getContentTypeByPath(savePathName)
	if err != nil {
		return errors.Wrap(err, "oss GetContentType failed")
	}
	options := []oss.Option{oss.ContentType(contentType)}
	if err = o.Bucket.PutObject(savePathName, bytes.NewReader(data), options...); err != nil {
		return errors.Wrap(err, "oss PutObject failed")
	}
	return nil
}

func (o OSSStorage) IsExist(objectKey string) (ok bool, err error) {
	return o.Bucket.IsObjectExist(objectKey)
}

func (o OSSStorage) PutFromFile(savePathName string, localPath string) error {
	options := []oss.Option{
		oss.SetHeader("Cache-Control", "max-age=1000"),
	}
	if err := o.Bucket.PutObjectFromFile(savePathName, localPath, options...); err != nil {
		return errors.Wrap(err, "oss PutFromFile")
	}
	return nil
}

func (o OSSStorage) Delete(objectKeys ...string) ([]string, error) {
	result, err := o.Bucket.DeleteObjects(objectKeys)
	if err != nil {
		return nil, errors.Wrap(err, "oss delete")
	}
	return result.DeletedObjects, nil
}

func (o OSSStorage) List(path string) (name []string) {
	ls, _ := o.Bucket.ListObjects(oss.Prefix(path))
	for _, object := range ls.Objects {
		name = append(name, object.Key)
	}
	return
}

func getContentTypeByPath(filePath string) (contentType string, err error) {
	ext := getExtByFilePath(filePath)
	if ext == "" {
		err = errors.New("file ext is required")
		return
	}
	contentType = mime.TypeByExtension(ext)
	if contentType == "" {
		err = errors.New("invalid file ext")
		return
	}
	return
}

func getExtByFilePath(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext
}
