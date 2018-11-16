package ffmpeg

import (
	"github.com/rylio/ytdl"
	"os/exec"
	"log"
)

type VideoInfo struct {
	Info *ytdl.VideoInfo
	Formats []ytdl.Format
}

func GetVideoInfo(link string)(VideoInfo,error){
	vidInfo:=VideoInfo{}
	vid,err:=ytdl.GetVideoInfo(link)
	if err!=nil{
		return vidInfo,err
	}
	vidInfo.Formats=vid.Formats
	vidInfo.Info=vid
	return vidInfo,err
}

func (info VideoInfo)GetDownloadLinkAndConvert(audioName string)error{
	//count:=0
	for _,format:=range info.Formats {
		url, err := info.Info.GetDownloadURL(format)
		if err != nil {
			log.Println("Error getting download URL: ",err)
			return err
		}
		downloadLink:=url.String()
		//log.Printf("[%v]",count)
		cmd := exec.Command("ffmpeg", "-i", downloadLink, audioName)
		errConvert := cmd.Run()
		if errConvert == nil {
			break
		}
		//log.Println("inappropriate format")
		//count++
	}
	return nil
}