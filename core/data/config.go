package data

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

const extraContentKey = "extra_page_content"
const streamTitleKey = "stream_title"
const streamKeyKey = "stream_key"
const logoPathKey = "logo_path"
const serverSummaryKey = "server_summary"
const serverNameKey = "server_name"
const serverURLKey = "server_url"
const httpPortNumberKey = "http_port_number"
const rtmpPortNumberKey = "rtmp_port_number"
const serverMetadataTagsKey = "server_metadata_tags"
const directoryEnabledKey = "directory_enabled"
const directoryRegistrationKeyKey = "directory_registration_key"
const socialHandlesKey = "social_handles"
const peakViewersSessionKey = "peak_viewers_session"
const peakViewersOverallKey = "peak_viewers_overall"
const lastDisconnectTimeKey = "last_disconnect_time"
const ffmpegPathKey = "ffmpeg_path"
const nsfwKey = "nsfw"
const s3StorageEnabledKey = "s3_storage_enabled"
const s3StorageConfigKey = "s3_storage_config"
const videoLatencyLevel = "video_latency_level"
const videoStreamOutputVariantsKey = "video_stream_output_variants"

// GetExtraPageBodyContent will return the user-supplied body content.
func GetExtraPageBodyContent() string {
	content, err := _datastore.GetString(extraContentKey)
	if err != nil {
		log.Errorln(extraContentKey, err)
		return config.GetDefaults().PageBodyContent
	}

	return content
}

// SetExtraPageBodyContent will set the user-supplied body content.
func SetExtraPageBodyContent(content string) error {
	return _datastore.SetString(extraContentKey, content)
}

// GetStreamTitle will return the name of the current stream.
func GetStreamTitle() string {
	title, err := _datastore.GetString(streamTitleKey)
	if err != nil {
		return ""
	}

	return title
}

// SetStreamTitle will set the name of the current stream.
func SetStreamTitle(title string) error {
	return _datastore.SetString(streamTitleKey, title)
}

// GetStreamKey will return the inbound streaming password.
func GetStreamKey() string {
	key, err := _datastore.GetString(streamKeyKey)
	if err != nil {
		log.Errorln(streamKeyKey, err)
		return ""
	}

	return key
}

// SetStreamKey will set the inbound streaming password.
func SetStreamKey(key string) error {
	return _datastore.SetString(streamKeyKey, key)
}

// GetLogoPath will return the path for the logo, relative to webroot.
func GetLogoPath() string {
	logo, err := _datastore.GetString(logoPathKey)
	if err != nil {
		log.Errorln(logoPathKey, err)
		return config.GetDefaults().Logo
	}

	if logo == "" {
		return config.GetDefaults().Logo
	}

	return logo
}

// SetLogoPath will set the path for the logo, relative to webroot.
func SetLogoPath(logo string) error {
	return _datastore.SetString(logoPathKey, logo)
}

// GetServerSummary will return the server summary text.
func GetServerSummary() string {
	summary, err := _datastore.GetString(serverSummaryKey)
	if err != nil {
		log.Errorln(serverSummaryKey, err)
		return ""
	}

	return summary
}

// SetServerSummary will set the server summary text.
func SetServerSummary(summary string) error {
	return _datastore.SetString(serverSummaryKey, summary)
}

// GetServerName will return the server name text.
func GetServerName() string {
	name, err := _datastore.GetString(serverNameKey)
	if err != nil {
		log.Errorln(serverNameKey, err)
		return ""
	}

	return name
}

// SetServerName will set the server name text.
func SetServerName(name string) error {
	return _datastore.SetString(serverNameKey, name)
}

// GetServerURL will return the server URL.
func GetServerURL() string {
	url, err := _datastore.GetString(serverURLKey)
	if err != nil {
		return ""
	}

	return url
}

// SetServerURL will set the server URL.
func SetServerURL(url string) error {
	return _datastore.SetString(serverURLKey, url)
}

// GetHTTPPortNumber will return the server HTTP port.
func GetHTTPPortNumber() int {
	port, err := _datastore.GetNumber(httpPortNumberKey)
	if err != nil {
		log.Errorln(httpPortNumberKey, err)
		return config.GetDefaults().WebServerPort
	}

	if port == 0 {
		return config.GetDefaults().WebServerPort
	}
	return int(port)
}

// SetHTTPPortNumber will set the server HTTP port.
func SetHTTPPortNumber(port float64) error {
	return _datastore.SetNumber(httpPortNumberKey, port)
}

// GetRTMPPortNumber will return the server RTMP port.
func GetRTMPPortNumber() int {
	port, err := _datastore.GetNumber(rtmpPortNumberKey)
	if err != nil {
		log.Errorln(rtmpPortNumberKey, err)
		return config.GetDefaults().RTMPServerPort
	}

	if port == 0 {
		return config.GetDefaults().RTMPServerPort
	}

	return int(port)
}

// SetRTMPPortNumber will set the server RTMP port.
func SetRTMPPortNumber(port float64) error {
	return _datastore.SetNumber(rtmpPortNumberKey, port)
}

// GetServerMetadataTags will return the metadata tags.
func GetServerMetadataTags() []string {
	tagsString, err := _datastore.GetString(serverMetadataTagsKey)
	if err != nil {
		log.Errorln(serverMetadataTagsKey, err)
		return []string{}
	}

	return strings.Split(tagsString, ",")
}

// SetServerMetadataTags will return the metadata tags.
func SetServerMetadataTags(tags []string) error {
	tagString := strings.Join(tags, ",")
	return _datastore.SetString(serverMetadataTagsKey, tagString)
}

// GetDirectoryEnabled will return if this server should register to YP.
func GetDirectoryEnabled() bool {
	enabled, err := _datastore.GetBool(directoryEnabledKey)
	if err != nil {
		return config.GetDefaults().YPEnabled
	}

	return enabled
}

// SetDirectoryEnabled will set if this server should register to YP.
func SetDirectoryEnabled(enabled bool) error {
	return _datastore.SetBool(directoryEnabledKey, enabled)
}

// SetDirectoryRegistrationKey will set the YP protocol registration key.
func SetDirectoryRegistrationKey(key string) error {
	return _datastore.SetString(directoryRegistrationKeyKey, key)
}

// GetDirectoryRegistrationKey will return the YP protocol registration key.
func GetDirectoryRegistrationKey() string {
	key, _ := _datastore.GetString(directoryRegistrationKeyKey)
	return key
}

// GetSocialHandles will return the external social links.
func GetSocialHandles() []models.SocialHandle {
	var socialHandles []models.SocialHandle

	configEntry, err := _datastore.Get(socialHandlesKey)
	if err != nil {
		log.Errorln(socialHandlesKey, err)
		return socialHandles
	}

	if err := configEntry.getObject(&socialHandles); err != nil {
		log.Errorln(err)
		return socialHandles
	}

	return socialHandles
}

// SetSocialHandles will set the external social links.
func SetSocialHandles(socialHandles []models.SocialHandle) error {
	var configEntry = ConfigEntry{Key: socialHandlesKey, Value: socialHandles}
	return _datastore.Save(configEntry)
}

// GetPeakSessionViewerCount will return the max number of viewers for this stream.
func GetPeakSessionViewerCount() int {
	count, err := _datastore.GetNumber(peakViewersSessionKey)
	if err != nil {
		return 0
	}
	return int(count)
}

// SetPeakSessionViewerCount will set the max number of viewers for this stream.
func SetPeakSessionViewerCount(count int) error {
	return _datastore.SetNumber(peakViewersSessionKey, float64(count))
}

// GetPeakOverallViewerCount will return the overall max number of viewers.
func GetPeakOverallViewerCount() int {
	count, err := _datastore.GetNumber(peakViewersOverallKey)
	if err != nil {
		return 0
	}
	return int(count)
}

// SetPeakOverallViewerCount will set the overall max number of viewers.
func SetPeakOverallViewerCount(count int) error {
	return _datastore.SetNumber(peakViewersOverallKey, float64(count))
}

// GetLastDisconnectTime will return the time the last stream ended.
func GetLastDisconnectTime() (time.Time, error) {
	var disconnectTime time.Time
	configEntry, err := _datastore.Get(lastDisconnectTimeKey)
	if err != nil {
		return disconnectTime, err
	}

	if err := configEntry.getObject(disconnectTime); err != nil {
		return disconnectTime, err
	}

	return disconnectTime, nil
}

// SetLastDisconnectTime will set the time the last stream ended.
func SetLastDisconnectTime(disconnectTime time.Time) error {
	var configEntry = ConfigEntry{Key: lastDisconnectTimeKey, Value: disconnectTime}
	return _datastore.Save(configEntry)
}

// SetNSFW will set if this stream has NSFW content.
func SetNSFW(isNSFW bool) error {
	return _datastore.SetBool(nsfwKey, isNSFW)
}

// GetNSFW will return if this stream has NSFW content.
func GetNSFW() bool {
	nsfw, err := _datastore.GetBool(nsfwKey)
	if err != nil {
		return false
	}
	return nsfw
}

// SetFfmpegPath will set the custom ffmpeg path.
func SetFfmpegPath(path string) error {
	return _datastore.SetString(ffmpegPathKey, path)
}

// GetFfMpegPath will return the ffmpeg path.
func GetFfMpegPath() string {
	path, err := _datastore.GetString(ffmpegPathKey)
	if err != nil {
		return ""
	}
	return path
}

// GetS3Config will return the external storage configuration.
func GetS3Config() models.S3 {
	configEntry, err := _datastore.Get(s3StorageConfigKey)
	if err != nil {
		return models.S3{Enabled: false}
	}

	var s3Config models.S3
	if err := configEntry.getObject(&s3Config); err != nil {
		return models.S3{Enabled: false}
	}

	return s3Config
}

// SetS3Config will set the external storage configuration.
func SetS3Config(config models.S3) error {
	var configEntry = ConfigEntry{Key: s3StorageConfigKey, Value: config}
	return _datastore.Save(configEntry)
}

// GetS3StorageEnabled will return if external storage is enabled.
func GetS3StorageEnabled() bool {
	enabled, err := _datastore.GetBool(s3StorageEnabledKey)
	if err != nil {
		log.Errorln(err)
		return false
	}

	return enabled
}

// SetS3StorageEnabled will enable or disable external storage.
func SetS3StorageEnabled(enabled bool) error {
	return _datastore.SetBool(s3StorageEnabledKey, enabled)
}

// GetStreamLatencyLevel will return the stream latency level.
func GetStreamLatencyLevel() models.LatencyLevel {
	level, err := _datastore.GetNumber(videoLatencyLevel)
	if err != nil || level == 0 {
		level = 4
	}

	return models.GetLatencyLevel(int(level))
}

// SetStreamLatencyLevel will set the stream latency level.
func SetStreamLatencyLevel(level float64) error {
	return _datastore.SetNumber(videoLatencyLevel, level)
}

// GetStreamOutputVariants will return all of the stream output variants.
func GetStreamOutputVariants() []models.StreamOutputVariant {
	configEntry, err := _datastore.Get(videoStreamOutputVariantsKey)
	if err != nil {
		return config.GetDefaults().StreamVariants
	}

	var streamOutputVariants []models.StreamOutputVariant
	if err := configEntry.getObject(&streamOutputVariants); err != nil {
		return config.GetDefaults().StreamVariants
	}

	if len(streamOutputVariants) == 0 {
		return config.GetDefaults().StreamVariants
	}

	return streamOutputVariants
}

// SetStreamOutputVariants will set the stream output variants.
func SetStreamOutputVariants(variants []models.StreamOutputVariant) error {
	var configEntry = ConfigEntry{Key: videoStreamOutputVariantsKey, Value: variants}
	return _datastore.Save(configEntry)
}

// VerifySettings will perform a sanity check for specific settings values.
func VerifySettings() error {
	if GetStreamKey() == "" {
		return errors.New("no stream key set. Please set one in your config file")
	}

	return nil
}

// FindHighestVideoQualityIndex will return the highest quality from a slice of variants.
func FindHighestVideoQualityIndex(qualities []models.StreamOutputVariant) int {
	type IndexedQuality struct {
		index   int
		quality models.StreamOutputVariant
	}

	if len(qualities) < 2 {
		return 0
	}

	indexedQualities := make([]IndexedQuality, 0)
	for index, quality := range qualities {
		indexedQuality := IndexedQuality{index, quality}
		indexedQualities = append(indexedQualities, indexedQuality)
	}

	sort.Slice(indexedQualities, func(a, b int) bool {
		if indexedQualities[a].quality.IsVideoPassthrough && !indexedQualities[b].quality.IsVideoPassthrough {
			return true
		}

		if !indexedQualities[a].quality.IsVideoPassthrough && indexedQualities[b].quality.IsVideoPassthrough {
			return false
		}

		return indexedQualities[a].quality.VideoBitrate > indexedQualities[b].quality.VideoBitrate
	})

	return indexedQualities[0].index
}
