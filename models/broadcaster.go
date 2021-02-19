package models

import "time"

// Broadcaster represents the details around the inbound broadcasting connection.
type Broadcaster struct {
	RemoteAddr    string               `json:"remoteAddr"`
	StreamDetails InboundStreamDetails `json:"streamDetails"`
	Time          time.Time            `json:"time"`
}

// InboundStreamDetails represents an inbound broadcast stream.
type InboundStreamDetails struct {
	Width          int     `json:"width"`
	Height         int     `json:"height"`
	VideoFramerate float32 `json:"framerate"`
	VideoBitrate   int     `json:"videoBitrate"`
	VideoCodec     string  `json:"videoCodec"`
	AudioBitrate   int     `json:"audioBitrate"`
	AudioCodec     string  `json:"audioCodec"`
	Encoder        string  `json:"encoder"`
	VideoOnly      bool    `json:"-"`
}

// RTMPStreamMetadata is the raw metadata that comes in with a RTMP connection.
type RTMPStreamMetadata struct {
	Width          int         `json:"width"`
	Height         int         `json:"height"`
	VideoBitrate   float32     `json:"videodatarate"`
	VideoCodec     interface{} `json:"videocodecid"`
	VideoFramerate float32     `json:"framerate"`
	AudioBitrate   float32     `json:"audiodatarate"`
	AudioCodec     interface{} `json:"audiocodecid"`
	Encoder        string      `json:"encoder"`
}
