package transcoder

import (
	"io/ioutil"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/utils"
)

var _timer *time.Ticker

func StopThumbnailGenerator() {
	if _timer != nil {
		_timer.Stop()
	}
}

// StartThumbnailGenerator starts generating thumbnails.
func StartThumbnailGenerator(chunkPath string, variantIndex int) {
	// Every 20 seconds create a thumbnail from the most
	// recent video segment.
	_timer = time.NewTicker(20 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-_timer.C:
				if err := fireThumbnailGenerator(chunkPath, variantIndex); err != nil {
					log.Errorln("Unable to generate thumbnail:", err)
				}
			case <-quit:
				log.Debug("thumbnail generator has stopped")
				_timer.Stop()
				return
			}
		}
	}()
}

func fireThumbnailGenerator(segmentPath string, variantIndex int) error {
	// JPG takes less time to encode than PNG
	outputFile := path.Join(config.WebRoot, "thumbnail.jpg")
	previewGifFile := path.Join(config.WebRoot, "preview.gif")

	framePath := path.Join(segmentPath, strconv.Itoa(variantIndex))
	files, err := ioutil.ReadDir(framePath)
	if err != nil {
		return err
	}

	var modTime time.Time
	var names []string
	for _, fi := range files {
		if path.Ext(fi.Name()) != ".ts" {
			continue
		}

		if fi.Mode().IsRegular() {
			if !fi.ModTime().Before(modTime) {
				if fi.ModTime().After(modTime) {
					modTime = fi.ModTime()
					names = names[:0]
				}
				names = append(names, fi.Name())
			}
		}
	}

	if len(names) == 0 {
		return nil
	}

	mostRecentFile := path.Join(framePath, names[0])
	ffmpegPath := utils.ValidatedFfmpegPath(data.GetFfMpegPath())

	thumbnailCmdFlags := []string{
		ffmpegPath,
		"-y",                 // Overwrite file
		"-threads 1",         // Low priority processing
		"-t 1",               // Pull from frame 1
		"-i", mostRecentFile, // Input
		"-f image2",  // format
		"-vframes 1", // Single frame
		outputFile,
	}

	ffmpegCmd := strings.Join(thumbnailCmdFlags, " ")
	if _, err := exec.Command("sh", "-c", ffmpegCmd).Output(); err != nil {
		return err
	}

	// If YP support is enabled also create an animated GIF preview
	if data.GetDirectoryEnabled() {
		makeAnimatedGifPreview(mostRecentFile, previewGifFile)
	}

	return nil
}

func makeAnimatedGifPreview(sourceFile string, outputFile string) {
	ffmpegPath := utils.ValidatedFfmpegPath(data.GetFfMpegPath())

	// Filter is pulled from https://engineering.giphy.com/how-to-make-gifs-with-ffmpeg/
	animatedGifFlags := []string{
		ffmpegPath,
		"-y",             // Overwrite file
		"-threads 1",     // Low priority processing
		"-i", sourceFile, // Input
		"-t 1", // Output is one second in length
		"-filter_complex", "\"[0:v] fps=8,scale=w=480:h=-1:flags=lanczos,split [a][b];[a] palettegen=stats_mode=full [p];[b][p] paletteuse=new=1\"",
		outputFile,
	}

	ffmpegCmd := strings.Join(animatedGifFlags, " ")
	if _, err := exec.Command("sh", "-c", ffmpegCmd).Output(); err != nil {
		log.Errorln(err)
	}
}
