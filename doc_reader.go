package dom

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"

	"github.com/goutlz/errz"

	"github.com/PuerkitoBio/goquery"
)

type DocumentReader func(request *http.Request) (*Document, error)

func NewDocumentReader(client *http.Client) (DocumentReader, error) {
	if client == nil {
		return nil, errz.New("Empty http client")
	}

	return func(request *http.Request) (*Document, error) {
		resp, err := client.Do(request)
		if err != nil {
			return nil, errz.Wrapf(err, "Failed to do request: %v", request)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, errz.Newf("Request failed: %v. Status code: %v. Status: %s", request, resp.StatusCode, resp.Status)
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errz.Wrap(err, "Failed to read response body")
		}

		bodyReader, err := convertEncodingToUtf8(b)
		if err != nil {
			return nil, errz.Wrap(err, "Failed to convert body encoding to utf-8")
		}

		doc, err := goquery.NewDocumentFromReader(bodyReader)
		if err != nil {
			return nil, errz.Wrap(err, "Failed to create document from reader")
		}

		return &Document{
			docEl: doc,
		}, nil
	}, nil
}

func convertEncodingToUtf8(b []byte) (io.Reader, error) {
	enc, _, _, err := determineEncoding(bytes.NewReader(b))
	if err != nil {
		return nil, errz.Wrap(err, "Failed to determine body encoding")
	}

	bodyReader := transform.NewReader(bytes.NewReader(b), enc.NewDecoder())
	return bodyReader, nil
}

func determineEncoding(r io.Reader) (encoding.Encoding, string, bool, error) {
	b, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		return nil, "", false, errz.Wrap(err, "Failed to peek bytes from reader")
	}

	enc, name, certain := charset.DetermineEncoding(b, "")
	return enc, name, certain, nil
}
