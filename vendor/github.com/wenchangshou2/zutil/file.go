package zutil

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

)

func IsExist(path string)bool{
	_,err:=os.Stat(path)
	return err==nil||os.IsExist(err)
}

func GetResourceAbsolutePath(rPath string,id string)string{
	return path.Join(rPath,id)
}
func IsAbsolutePath(p string)bool{
	return path.IsAbs(p)
}
func GetSize(f multipart.File)(int ,error){
	content,err:=ioutil.ReadAll(f)
	return len(content),err
}
func GetExt(fileName string)string{
	return path.Ext(fileName)
}
func CheckExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

func IsNotExistMkDir(src string) error {
	if exist := CheckExist(src); exist == true {
		if err := os.MkdirAll(src,0755); err != nil {
			return err
		}
	}

	return nil
}
func ReMkdir(src string)error{
	if err:=IsExistDelete(src);err!=nil{
		return err
	}

	if err:=IsNotExistMkDir(src);err!=nil{
		return  err
	}
	return nil
}
//如果目录存在就直接删除
func IsExistDelete(src string) error{
	if exist:=CheckExist(src);!exist{
		os.RemoveAll(src)
	}
	return  nil
}

func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	perm := CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}
func substr(s string,pos,length int)string{
	runes:=[]rune(s)
	l:=pos+length

	if l>len(runes){
		l=len(runes)
	}
	return string(runes[pos:l])
}
func getParentDirector(director string) string{
	director=strings.ReplaceAll(director,"\\","/")
	return  substr(director,0,strings.LastIndex(director,"/"))
}
//解析zip文件
//sourcePath:源文件
//targetPath:目标路径
func Unzip(sourcePath,targetPath string)(err error){
	zipReader,err:=zip.OpenReader(sourcePath)
	defer zipReader.Close()
	if err!=nil{
		return
	}
	for _,file:=range zipReader.Reader.File{
		zippedFile,err:=file.Open()
		defer zippedFile.Close()
		if err!=nil{
			log.Fatal(err)
		}

		utf8Reader := transform.NewReader(strings.NewReader(file.Name),
			simplifiedchinese.GBK.NewDecoder())
		fileName, err := ioutil.ReadAll(utf8Reader)

		extractedFilePath:=filepath.Join(
			targetPath,string(fileName),
		)
		if file.FileInfo().IsDir(){ //目录处理
			os.MkdirAll(extractedFilePath,file.Mode())
		}else{ //文件处理
			parentPath:=getParentDirector(extractedFilePath) //获取上一级目录
			IsNotExistMkDir(parentPath)//如果上一级目录不存在就自动创建
			outputFile,err:=os.OpenFile(
				extractedFilePath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				file.Mode())//创建输出文件
			if err!=nil{
				log.Fatal(err)
			}
			defer outputFile.Close()
			_,err=io.Copy(outputFile,zippedFile)//将压缩包内的数据写入到目标文件
			if err!=nil{
				log.Fatal(err)
			}
		}
	}
	return
}
//获取文件名
func GetFileName(path string)string{
	idx:=strings.LastIndex(path,"/")
	if idx==-1{
		return ""
	}
	return path[idx+1:]
}
func GetFileExt(p string)string{
	return filepath.Ext(p)
}
func CopyFile(source,target string)(error){
	sfi,err:=os.Stat(source)
	if err!=nil{
		return err
	}
	if !sfi.Mode().IsRegular(){
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi,err:=os.Stat(target)
	if err!=nil{
		if !os.IsNotExist(err){
			return  err
		}
	}else{
		if !(dfi.Mode().IsRegular()){
			return fmt.Errorf("CopyFile:no-regular destination file%s (%q)",dfi.Name(),dfi.Mode().String())
		}
		if os.SameFile(sfi,dfi){
			return nil
		}
	}
	if err=os.Link(source,target);err==nil{
		return nil
	}
	err=copyFileContents(source,target)

	return err
}
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
//获取文件的md5
func GetFileMd5(filePath string)(string,error){
	var returnMD5String string
	file,err:=os.Open(filePath)
	if err!=nil{
		return returnMD5String,err
	}
	defer file.Close()
	hash:=md5.New()
	if _,err:=io.Copy(hash,file);err!=nil{
		return returnMD5String,nil
	}
	hashInBytes:=hash.Sum(nil)[:16]
	returnMD5String=hex.EncodeToString(hashInBytes)
	return returnMD5String,nil
}


func GetWorkDirFileListRecursive(workDir string) ([]string,error){
	var files[]string
	err:=filepath.Walk(workDir, func(path string, info os.FileInfo, err error) error {
		files=append(files,path)
		return nil
	})
	return files,err
}
func GetWorkDirFileList(workDir string)(fileList []string,err error){
	files,err:=ioutil.ReadDir(workDir)
	if err!=nil{
		return
	}
	for _,f:=range files{
		fileName:=path.Join(workDir,f.Name())
		fileList=append(fileList,fileName)
	}
	return
}