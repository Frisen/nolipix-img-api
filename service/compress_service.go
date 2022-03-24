package service

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"nolipix-img-api/internal/aliyun"
	"nolipix-img-api/util"
	"path"
	"strings"
	"time"

	cv "gocv.io/x/gocv"
)

// 压缩
func Compress(url string, rows, cols int) (fileName string, err error) {
	client := aliyun.GetSHClient()
	t2 := time.Now()
	//从oss获取图片
	r, err := client.GetObject(url)
	fmt.Println("get oss --->", time.Since(t2))
	if err != nil {
		return
	}
	defer r.Close()
	t3 := time.Now()
	srcBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	//构造Mat对象
	mat, err := cv.IMDecode(srcBytes, cv.IMReadUnchanged)
	if err != nil {
		return
	}
	defer mat.Close()
	fmt.Println("decode --->", time.Since(t3))

	t1 := time.Now()
	//压缩图片
	dst, err := compress(&mat, mat.Rows(), mat.Cols(), rows, cols)
	fmt.Println("compress --->", time.Since(t1))
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	t4 := time.Now()
	err = png.Encode(buf, dst)
	if err != nil {
		return
	}
	fmt.Println("encode --->", time.Since(t4))
	//命名压缩图片
	fullName := path.Base(url)
	suffix := path.Ext(fullName)
	dir := path.Dir(url)
	fileName = strings.TrimSuffix(fullName, suffix)
	fileName = dir + fileName + "_thumb" + suffix
	t := time.Now()
	//图片上传oss
	err = client.PutObject(fileName, buf)
	fmt.Println("put oss --->", time.Since(t))
	if err != nil {
		return
	}
	return fileName, nil
}

func compress(src *cv.Mat, srcRows, srcCols, dstRows, dstCols int) (image.Image, error) {
	defer src.Close()
	if src.Empty() {
		return nil, errors.New("Mat is empty")
	}
	w, h := util.CalXY(float64(srcCols), float64(srcRows), float64(dstCols), float64(dstRows))
	dst := cv.NewMat()
	cv.Resize(*src, &dst, image.Pt(w, h), 0, 0, cv.InterpolationArea)
	buffer, err := cv.IMEncodeWithParams(cv.PNGFileExt, dst, []int{cv.IMWritePngCompression, 9, cv.IMWritePngStrategy, cv.IMWritePngStrategyHuffmanOnly})
	if err != nil {
		return nil, err
	}
	mat, err := cv.IMDecode(buffer.GetBytes(), cv.IMReadUnchanged)
	if err != nil {
		return nil, err
	}
	img, err := mat.ToImage()
	if err != nil {
		return nil, err
	}
	return img, nil
}
