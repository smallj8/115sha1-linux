package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"strconv"
	"strings"

	"os"
)

var path_zong string

// 判断 文件夹 or 文件
func readdir(pathname string) bool {
	s, _ := os.Stat(pathname)
	if s.IsDir() {
		return true
	} else {
		return false
	}
}

// 计算文件的sha1并返回
func SHA1File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	h := sha1.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return "", err
	}
	//return fmt.Sprintf("%x",h.Sum(nil)), nil
	return hex.EncodeToString(h.Sum(nil)), nil
}

// 计算文件的头部sha1
func SHA1File2(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	h := sha1.New()
	if err != nil {
		return "", err
	}
	b := make([]byte, 131072)
	file.Read(b)
	h.Write(b)

	return hex.EncodeToString(h.Sum(nil)), nil
}

// 获取文件大小
func getsize(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		return 1
	}
	return fi.Size()
}

// 处理文件
func handlefile(pathname string) string {
	sha1_all, _ := SHA1File(pathname)
	sha1_all = strings.ToUpper(sha1_all)
	sha1_128, _ := SHA1File2(pathname)
	sha1_128 = strings.ToUpper(sha1_128)
	filesize := strconv.FormatInt(getsize(pathname), 10)
	rsp := strings.Split(pathname, "/")
	if getsize(pathname) <= 131072 {
		sha1_128 = sha1_all
	}
	a := "115://" + rsp[len(rsp)-1] + "|" + filesize + "|" + sha1_all + "|" + sha1_128
	fmt.Println(a)
	return a
}

// 处理文件夹
func handledir(pathname, shatxt string) {
	fmt.Println(pathname)
	GetAllFile(pathname, shatxt)
}

// 判断是否是文件夹
func GetAllFile(pathname, shatxt string) error {
	rd, err := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			a := pathname + fi.Name() + "/"
			GetAllFile(a, shatxt)
		} else {
			filename := pathname + fi.Name()
			sha_one := handlefile(filename)
			dir_path := strings.Replace(pathname, path_zong, "", -1)
			dir_path = strings.Replace(dir_path, "/", "|", -1)
			if dir_path == "" {
				file, _ := os.OpenFile(shatxt, os.O_WRONLY|os.O_APPEND, 0666)
				file.WriteString(sha_one + "\n")
				defer file.Close()
			} else {
				sha_two := sha_one + "|" + dir_path
				sha_two = strings.TrimSuffix(sha_two, "|")
				file, _ := os.OpenFile(shatxt, os.O_WRONLY|os.O_APPEND, 0666)
				file.WriteString(sha_two + "\n")
				defer file.Close()
			}

		}
	}
	return err
}

func main() {
	// 获取参数 pathname
	path := os.Args[1]

	sha_out := "/root/alist/getsha1/"

	// 判断是否是绝对路径如果不是把相对路径转为绝对路径
	if !filepath.IsAbs(path) {
		path, _ = filepath.Abs(path)
	}

	path_zong = path
	// 判断根目录是文件还是文件夹
	isdir := readdir(path)
	if isdir {
		if !strings.HasSuffix(path, "/") {
			path = path + "/"
			path_zong = path
		}
		pathsplist := strings.Split(path, "/")
		filename := pathsplist[len(pathsplist)-2] + "_本地计算(带目录).txt"
		sha1_txt := sha_out + filename
		os.OpenFile(sha1_txt, os.O_CREATE, 0666)
		handledir(path, sha1_txt)
	} else {
		handlefile(path)
	}
}
