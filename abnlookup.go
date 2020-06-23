package abnlookup

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

var (
	ErrWrongLength     = errors.New("abn should be eleven digits")
	ErrInvalidFormat   = errors.New("abn should only be numbers")
	ErrInvalidChecksum = errors.New("invalid checksum")
	ErrInvalidABN      = errors.New("invalid abn format")
	ErrRecordNotFound  = errors.New("record not found")
	ErrUnknownResponse = errors.New("unknown response from ABR")
)

func Validate(abn string) error {
	if len(abn) != 11 {
		return ErrWrongLength
	}

	for i := 0; i < len(abn); i++ {
		if abn[i] < '0' || abn[i] > '9' {
			return ErrInvalidFormat
		}
	}

	var n int

	n += int(abn[0]-'1') * 10
	n += int(abn[1]-'0') * 1
	n += int(abn[2]-'0') * 3
	n += int(abn[3]-'0') * 5
	n += int(abn[4]-'0') * 7
	n += int(abn[5]-'0') * 9
	n += int(abn[6]-'0') * 11
	n += int(abn[7]-'0') * 13
	n += int(abn[8]-'0') * 15
	n += int(abn[9]-'0') * 17
	n += int(abn[10]-'0') * 19

	if n%89 != 0 {
		return ErrInvalidChecksum
	}

	return nil
}

type ABNData struct {
	ABN  string
	Name string
}

func Fetch(abn string) (*ABNData, error) {
	if err := Validate(abn); err != nil {
		return nil, errors.Wrap(err, "abnlookup.Fetch: preliminary validation failed")
	}

	res, err := http.Get("https://abr.business.gov.au/ABN/View?abn=" + abn)
	if err != nil {
		return nil, errors.Wrap(err, "abnlookup.Fetch: couldn't perform request")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("abnlookup.Fetch: invalid response code; expected 200 OK but got %s", res.Status)
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, errors.Wrap(err, "abnlookup.Fetch: couldn't parse response")
	}
	if doc == nil {
		return nil, errors.Errorf("abnlookup.Fetch: couldn't parse response")
	}

	if strings.Contains(doc.Find("div.process-message").Text(), "No record found matching ABN") {
		return nil, ErrRecordNotFound
	}

	if strings.Contains(doc.Find("div.process-message").Text(), "The number entered is not a valid ABN") {
		return nil, ErrInvalidABN
	}

	if !strings.HasPrefix(doc.Find("title").Text(), "Current details for ABN") {
		return nil, ErrUnknownResponse
	}

	return &ABNData{
		ABN:  abn,
		Name: doc.Find("span[itemprop=legalName]").Text(),
	}, nil
}

func Lookup(abn string) error {
	_, err := Fetch(abn)
	return err
}
