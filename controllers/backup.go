package controllers

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"blog/helpers"
	"blog/system"
)

//func BackupPost(c *gin.Context) {
//	var (
//		err error
//		res = gin.H{}
//	)
//	defer writeJSON(c, res)
//	err = Backup()
//	if err != nil {
//		res["message"] = err.Error()
//		return
//	}
//	res["succeed"] = true
//}

func RestorePost(c *gin.Context) {
	var (
		fileName  string
		fileUrl   string
		err       error
		res       = gin.H{}
		resp      *http.Response
		bodyBytes []byte
	)
	defer writeJSON(c, res)
	fileName = c.PostForm("fileName")
	if fileName == "" {
		res["message"] = "fileName cannot be empty."
		return
	}
	fileUrl = system.GetConfiguration().QiniuFileServer + fileName
	resp, err = http.Get(fileUrl)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	defer resp.Body.Close()

	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	bodyBytes, err = helpers.Decrypt(bodyBytes, system.GetConfiguration().BackupKey)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	err = ioutil.WriteFile(fileName, bodyBytes, os.ModePerm)
	if err == nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

//func Backup() (err error) {
//	var (
//		u           *url.URL
//		exist       bool
//		ret         PutRet
//		bodyBytes   []byte
//		encryptData []byte
//	)
//	u, err = url.Parse(system.GetConfiguration().DSN)
//	if err != nil {
//		seelog.Debug("parse dsn error:%v", err)
//		return
//	}
//	exist, _ = helpers.PathExists(u.Path)
//	if !exist {
//		err = errors.New("database file doesn't exists.")
//		seelog.Debug("database file doesn't exists.")
//		return
//	}
//	seelog.Debug("start backup...")
//	bodyBytes, err = ioutil.ReadFile(u.Path)
//	if err != nil {
//		seelog.Error(err)
//		return
//	}
//	encryptData, err = helpers.Encrypt(bodyBytes, system.GetConfiguration().BackupKey)
//	if err != nil {
//		seelog.Error(err)
//		return
//	}
//
//	putPolicy := storage.PutPolicy{
//		Scope: system.GetConfiguration().QiniuBucket,
//	}
//	mac := qbox.NewMac(system.GetConfiguration().QiniuAccessKey, system.GetConfiguration().QiniuSecretKey)
//	token := putPolicy.UploadToken(mac)
//	cfg := storage.Config{}
//	uploader := storage.NewFormUploader(&cfg)
//	putExtra := storage.PutExtra{}
//
//	fileName := fmt.Sprintf("blog_%s.db", helpers.GetCurrentTime().Format("20060102150405"))
//	err = uploader.Put(context.Background(), &ret, token, fileName, bytes.NewReader(encryptData), int64(len(encryptData)), &putExtra)
//	if err != nil {
//		seelog.Debugf("backup error:%v", err)
//		return
//	}
//	seelog.Debug("backup succeefully.")
//	return err
//}
