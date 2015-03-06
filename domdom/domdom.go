//Package domdom handles communication and XML parsing with the DomDomSoft Anime Downloader server.
package domdom

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
)

const service string = "http://anime.domdomsoft.com/Services/MainService.asmx"
const postContentType = "text/xml"

//The Episode Type represents an episode file
type Episode struct {
	SeriesName, FileName, FileSize string
}

//ListEpisodes gives a slice of Episodes given a Series title and a domdom key. The key can be left as an empty string for non-premium users.
func ListEpisodes(title, key string) ([]Episode, error) {
	type getListEpisodeData struct {
		Title, Serial string
	}

	var empty []Episode

	getListEpisode, err := template.New("getListEpisode").Parse(`<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="http://tempuri.org/">
  <SOAP-ENV:Body>
    <ns1:GetListEpisode>
      <ns1:animeTitle>{{.Title}}</ns1:animeTitle>
      <ns1:serial>{{.Serial}}</ns1:serial>
    </ns1:GetListEpisode>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>
`)
	if err != nil {
		return empty, err
	}

	currentData := getListEpisodeData{
		Title:  title,
		Serial: key,
	}

	var buf bytes.Buffer
	err = getListEpisode.Execute(&buf, currentData)
	if err != nil {
		return empty, err
	}

	resp, err := http.Post(service, postContentType, &buf)
	if err != nil {
		return empty, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return empty, err
	}

	type XmlEpisodeFile struct {
		XMLName xml.Name `xml:"EpisodeFile"`
		Name    string   `xml:"Name"`
		Size    string   `xml:"FileSize"`
	}

	type XmlResult struct {
		XMLName     xml.Name `xml:"GetListEpisodeResult"`
		EpisodeFile []XmlEpisodeFile
	}

	type XmlResponse struct {
		XMLName              xml.Name `xml:"GetListEpisodeResponse"`
		GetListEpisodeResult XmlResult
	}

	type XmlBody struct {
		XMLName                xml.Name `xml:"Body"`
		GetListEpisodeResponse XmlResponse
	}

	type XmlEnvelope struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    XmlBody
	}

	list := XmlEnvelope{}
	err = xml.Unmarshal(body, &list)
	if err != nil {
		return empty, err
	}

	out := make([]Episode, len(list.Body.GetListEpisodeResponse.GetListEpisodeResult.EpisodeFile))

	for index, elem := range list.Body.GetListEpisodeResponse.GetListEpisodeResult.EpisodeFile {
		out[index].SeriesName = title
		out[index].FileName = elem.Name
		out[index].FileSize = elem.Size
	}
	return out, nil
}

//UpdateAnimeList downloads the xml list of Anime Series titles from the server and saves it to a file.
func UpdateAnimeList(filePath string) error {
	getAnimeList := `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
<soap:Body>
<GetAnimeList xmlns="http://tempuri.org/" />
</soap:Body>
</soap:Envelope>
`
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Post(service, postContentType, strings.NewReader(getAnimeList))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

//The Anime struct represents a series
type Anime struct {
	Title    string
	NumFiles string
}

//GetAnimeList downloads the list of Anime from the server and returns it as a slice of Animes.
func GetAnimeList() ([]Anime, error) {
	getAnimeList := `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
<soap:Body>
<GetAnimeList xmlns="http://tempuri.org/" />
</soap:Body>
</soap:Envelope>
`
	var empty []Anime

	resp, err := http.Post(service, postContentType, strings.NewReader(getAnimeList))
	if err != nil {
		return empty, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return empty, err
	}

	return parseXmlAnimeList(body)
}

func parseXmlAnimeList(xmlBody []byte) ([]Anime, error) {

	var empty []Anime

	type XmlAnime struct {
		XMLName xml.Name `xml:"Anime"`
		Id      string   `xml:"Id"`
		Title   string   `xml:"Title"`
		NumFile string   `xml:"NumFile"`
	}

	type XmlResult struct {
		XMLName xml.Name `xml:"GetAnimeListResult"`
		Anime   []XmlAnime
	}

	type XmlResponse struct {
		XMLName            xml.Name `xml:"GetAnimeListResponse"`
		GetAnimeListResult XmlResult
	}

	type XmlBody struct {
		XMLName              xml.Name `xml:"Body"`
		GetAnimeListResponse XmlResponse
	}

	type XmlEnvelope struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    XmlBody
	}

	envelope := XmlEnvelope{}
	err := xml.Unmarshal(xmlBody, &envelope)
	if err != nil {
		return empty, err
	}
	targetlength := len(envelope.Body.GetAnimeListResponse.GetAnimeListResult.Anime)
	if targetlength < 1 {
		return empty, errors.New("No Anime were parsed")
	}

	out := make([]Anime, targetlength)

	for index, elem := range envelope.Body.GetAnimeListResponse.GetAnimeListResult.Anime {
		out[index].Title = elem.Title
		out[index].NumFiles = elem.NumFile
	}
	return out, nil
}

//LoadAnimeList loads the anime list from a file and returns it as a slice of Animes.
func LoadAnimeList(filePath string) ([]Anime, error) {
	var empty []Anime
	xmlBody, err := ioutil.ReadFile(filePath)
	if err != nil {
		return empty, err
	}
	return parseXmlAnimeList(xmlBody)
}

//GetDownloadLinks retrieves the download links for an Episode from the server given an Episode and a domdom key. The key can be left as an empty string for non-premium users.
func GetDownloadLinks(episode Episode, key string) ([]string, error) {
	title := episode.SeriesName
	filename := episode.FileName

	type requestLinkDownload2Data struct {
		Title, FileName, Serial string
	}

	var empty []string

	requestLinkDownload2, err := template.New("requestLinkDownload2").Parse(`<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="http://tempuri.org/">
  <SOAP-ENV:Body>
    <ns1:RequestLinkDownload2>
      <ns1:animeTitle>{{.Title}}</ns1:animeTitle>
      <ns1:episodeName>{{.FileName}}</ns1:episodeName>
      <ns1:serial>{{.Serial}}</ns1:serial>
    </ns1:RequestLinkDownload2>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>
`)

	currentData := requestLinkDownload2Data{
		Title:    title,
		FileName: filename,
		Serial:   key,
	}

	var buf bytes.Buffer
	err = requestLinkDownload2.Execute(&buf, currentData)
	if err != nil {
		return empty, err
	}

	resp, err := http.Post(service, postContentType, &buf)
	if err != nil {
		return empty, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return empty, err
	}

	type XmlResponse struct {
		XMLName                    xml.Name `xml:"RequestLinkDownload2Response"`
		RequestLinkDownload2Result string
	}

	type XmlBody struct {
		XMLName                      xml.Name `xml:"Body"`
		RequestLinkDownload2Response XmlResponse
	}

	type XmlEnvelope struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    XmlBody
	}

	envelope := XmlEnvelope{}
	err = xml.Unmarshal(body, &envelope)
	if err != nil {
		return empty, err
	}

	downloadLinksString := envelope.Body.RequestLinkDownload2Response.RequestLinkDownload2Result
	out := strings.Split(downloadLinksString, "|||")

	return out, nil
}

//FindEpisodeByName is a convenience function that returns an Episode object given a series name and filename.
func FindEpisodeByName(series, episodeName, key string) (Episode, error) {
	var empty Episode

	episodes, err := ListEpisodes(series, key)
	if err != nil {
		return empty, err
	}

	for _, episode := range episodes {
		if episode.FileName == episodeName {
			return episode, nil
		}
	}
	return empty, errors.New("No such episode")
}

//FindEpisodeByIndex is a convenience function that returns an Episode object given a series name and an episode index (starting at 1).
func FindEpisodeByIndex(series, episodeId, key string) (Episode, error) {
	var empty Episode

	episodes, err := ListEpisodes(series, key)
	if err != nil {
		return empty, err
	}

	for index, episode := range episodes {
		var thisEp string = strconv.Itoa(index + 1)
		if thisEp == episodeId {
			return episode, nil
		}
	}
	return empty, errors.New("No such episode")
}
