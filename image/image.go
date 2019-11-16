package image

import (
	"fmt"
	"net/url"
)

const screenshotBaseURL string = "https://files.opskins.media/file/opskins-patternindex/"

// BuildURL builds the screenshot URL
func BuildURL(defIndex, paintIndex, paintSeed int, inspectURL string) (string, error) {
	screenshotURL, err := url.Parse(fmt.Sprintf("%s", screenshotBaseURL))
	if err != nil {
		return "", err
	}

	screenshotURL.Path += fmt.Sprintf("%d_%d_%d.jpg", defIndex, paintIndex, paintSeed)

	return screenshotURL.String(), nil
}
