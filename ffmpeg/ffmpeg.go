package ffmpeg

import (
	"github.com/rylio/ytdl"
	"log"
	"os/exec"
)

type VideoInfo struct {
	Info *ytdl.VideoInfo
	Format ytdl.Format
}

func GetVideoInfo(link string)(VideoInfo,error){
	vidInfo:=VideoInfo{}
	vid,err:=ytdl.GetVideoInfo(link)
	if err!=nil{
		return vidInfo,err
	}
	if len(vid.Formats)>0{
		vidInfo.Format=vid.Formats[0]
	}else {
		log.Println("Length formats is 0")
	}
	return vidInfo,err
}

func (info *VideoInfo)GetDownloadURL()(string,error){
	url,err:=info.Info.GetDownloadURL(info.Format)
	if err!=nil {
		return "", err
	}
	URL:=url.String()
	return URL,err
}

func ConvertVideoToAudio(downloadLink,audioname string)error{
	cmd:=exec.Command("ffmpeg","-i",downloadLink,audioname)
	err:=cmd.Run()
	if err!=nil{
		return err
	}
	return nil
}