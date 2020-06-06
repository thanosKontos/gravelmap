package hgt

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
)

const nasa30mSrtmURL = "http://e4ftl01.cr.usgs.gov/MEASURES/SRTMGL1.003/2000.02.11/%s.SRTMGL1.hgt.zip"

type download struct {
	nasaUsername   string
	nasaPassword   string
	destinationDir string
}

func (d *download) download(dms string) error {
	out, err := os.Create(fmt.Sprintf("/%s/%s.hgt.zip", d.destinationDir, dms))
	defer out.Close()
	if err != nil {
		return err
	}

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:           jar,
		CheckRedirect: d.redirectPolicyFunc,
	}
	req, err := http.NewRequest("GET", fmt.Sprintf(nasa30mSrtmURL, dms), nil)
	req.SetBasicAuth(d.nasaUsername, d.nasaPassword)
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (d *download) redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.SetBasicAuth(d.nasaUsername, d.nasaPassword)
	return nil
}
