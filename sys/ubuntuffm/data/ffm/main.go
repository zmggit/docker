package main

import (
	"embed"
	"ffm/tool"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"log"
)

// Filter 水印
func Filter() {
	overlay := ffmpeg.Input("./91cb556b8798ac44c4de15380ad401c5.jpg").Filter("scale", ffmpeg.Args{"64:-1"})
	err := ffmpeg.Filter(
		[]*ffmpeg.Stream{
			ffmpeg.Input("./video.mp4"),
			overlay,
		}, "overlay", ffmpeg.Args{"10:10"}, ffmpeg.KwArgs{"enable": "gte(t,1)"}).
		Output("./out1.mp4").OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		fmt.Println(err, "错误")
	}
}

func Gif() {
	err := ffmpeg.Input("./static/video.mp4", ffmpeg.KwArgs{"ss": "1"}).
		Output("./static/out1.gif", ffmpeg.KwArgs{"s": "320x240", "pix_fmt": "rgb24", "t": "3", "r": "3"}).
		OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		fmt.Println(err, "错误")
	}
}

func Check() {
	// 打开嵌入的文件
	file, err := embeddedImage.Open("static/ffmpeg.exe")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ffmByte, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	if err := tool.WriteFile("test.exe", ffmByte, 0644); err != nil {
		fmt.Println(err, "c错误")
	}

	//installFFM := "Start-Process \"installer.exe\" -ArgumentList \"/S\" -Wait"
}

// 嵌入式
//
//go:embed static/*
var embeddedImage embed.FS

func main() {
	//fmt.Println(os.Getwd())
	//环境校验
	Check()

	//水印
	//Filter()
	//视频变gif
	//Gif()
}
