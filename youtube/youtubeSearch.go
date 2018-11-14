package youtube

import (
	"google.golang.org/api/youtube/v3"
	"net/http"
	"google.golang.org/api/googleapi/transport"
)

const MaxSearchCount  = 10

type ClientYoutube struct {
	Client *youtube.Service
}

type Video struct {
	Id string
	Title string
}

func NewYoutubeClient(developerKey string)(ClientYoutube,error){
	client:=&http.Client{
		Transport: &transport.APIKey{Key:developerKey},
	}
	NewCLient:=ClientYoutube{}
	service,err:=youtube.New(client)
	if err!=nil{
		return NewCLient,err
	}
	NewCLient.Client=service
	return NewCLient,err
}

func (client *ClientYoutube)Search(query string)(Video,error){
	vid:=Video{}
	request:=client.Client.Search.List("snippet").Q(query).MaxResults(MaxSearchCount)
	response,err:=request.Do()
	if err!=nil{
		return vid,err
	}
	for _,item:=range response.Items{
		if item.Id.Kind=="youtube#video"{
			vid.Id=item.Id.VideoId
			vid.Title=item.Snippet.Title
		}else {
			continue
		}
	}
	return vid,err
}
