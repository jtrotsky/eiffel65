package image

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jtrotsky/eiffel65/util"
)

const (
	screenshotBaseURL string = "https://metjm.net/shared/screenshots-v5.php"
)

// ScreenshotPayload contains the result of the screenshot job from metjm.net.
type ScreenshotPayload struct {
	Result Screenshot `json:"result,omitempty"`
}

// Screenshot contains the screenshot summary from metjm.net.
type Screenshot struct {
	ID       int64  `json:"screen_id,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

// GetScreenshot returns an image URL of a screenshot for a given CSGO asset.
func GetScreenshot(inspectURL string) (string, error) {
	// https://metjm.net/shared/screenshots-v5.php?cmd=request_new_link&inspect_link=steam%3A%2F%2Frungame%2F730%2F76561202255233023%2F%2Bcsgo_econ_action_preview%2520M1929109707517347231A14516089565D14420669795467175295&user_uuid=2b28928a-d665-ded6-3a24-cab3d42f7f28&user_client=1&custom_rotation_id=0&use_logo=0&mode=7&resolution=4&forceOpskins=0
	screenshotURL, err := url.Parse(fmt.Sprintf("%s", screenshotBaseURL))
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("cmd", "request_new_link")
	params.Add("user_uuid", util.GenerateUUID())
	params.Add("user_client", "1")
	params.Add("use_logo", "0")
	params.Add("mode", "7")
	params.Add("resolution", "4")
	params.Add("forceOpskins", "0")
	screenshotURL.RawQuery = params.Encode()

	screenshotURL.RawQuery += fmt.Sprintf("&inspect_link=%s", inspectURL)

	log.Println(screenshotURL)

	response, err := http.DefaultClient.Get(screenshotURL.String())
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	screenshotPayload := ScreenshotPayload{}
	err = json.NewDecoder(response.Body).Decode(&screenshotPayload)
	if err != nil {
		return "", err
	}

	screenshotLink, err := lookupScreenshotJob(screenshotPayload.Result.ID)
	if err != nil {
		return "", err
	}

	return screenshotLink, nil
}

// lookupScreenshotJob looks up a screenshot based on the ID that was returned
// when starting the job of creating it.
func lookupScreenshotJob(screenshotID int64) (string, error) {
	// 	https://metjm.net/shared/screenshots-v5.php?cmd=request_screenshot_status&id=2654287279347
	screenshotJobURL, err := url.Parse(fmt.Sprintf("%s", screenshotBaseURL))
	if err != nil {
		return "", err
	}

	screenshotIDStr := strconv.FormatInt(screenshotID, 10)

	params := url.Values{}
	params.Add("cmd", "request_screenshot_status")
	params.Add("id", screenshotIDStr)
	screenshotJobURL.RawQuery = params.Encode()

	log.Println(screenshotJobURL)

	response, err := http.DefaultClient.Get(screenshotJobURL.String())
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	screenshotJobPayload := ScreenshotPayload{}
	err = json.NewDecoder(response.Body).Decode(&screenshotJobPayload)
	if err != nil {
		return "", err
	}

	return screenshotJobPayload.Result.ImageURL, err
}
