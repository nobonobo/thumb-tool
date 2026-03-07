package main

type UploadMeta struct {
	Title                  string   `json:"title"`
	Description            string   `json:"description,omitempty"`
	CategoryID             string   `json:"categoryId,omitempty"`
	PrivacyStatus          string   `json:"privacyStatus,omitempty"`
	Embeddable             bool     `json:"embeddable,omitempty"`
	License                string   `json:"license,omitempty"`
	MadeForKids            bool     `json:"madeForKids,omitempty"`
	ContainsSyntheticMedia bool     `json:"containsSyntheticMedia,omitempty"`
	PlaylistIds            []string `json:"playlistIds,omitempty"`
	Language               string   `json:"language,omitempty"`
}

type Template struct {
	SVG        string            `json:"svg"`
	Defaults   map[string]string `json:"defaults"`
	UploadMeta UploadMeta        `json:"upload"`
}

type ProbeInfo struct {
	Streams []Streams `json:"streams"`
	Frames  []Frame   `json:"frames"`
}

type Disposition struct {
	Default         int `json:"default"`
	Dub             int `json:"dub"`
	Original        int `json:"original"`
	Comment         int `json:"comment"`
	Lyrics          int `json:"lyrics"`
	Karaoke         int `json:"karaoke"`
	Forced          int `json:"forced"`
	HearingImpaired int `json:"hearing_impaired"`
	VisualImpaired  int `json:"visual_impaired"`
	CleanEffects    int `json:"clean_effects"`
	AttachedPic     int `json:"attached_pic"`
	TimedThumbnails int `json:"timed_thumbnails"`
	NonDiegetic     int `json:"non_diegetic"`
	Captions        int `json:"captions"`
	Descriptions    int `json:"descriptions"`
	Metadata        int `json:"metadata"`
	Dependent       int `json:"dependent"`
	StillImage      int `json:"still_image"`
	Multilayer      int `json:"multilayer"`
}

type Streams struct {
	Index          int               `json:"index"`
	CodecName      string            `json:"codec_name"`
	CodecLongName  string            `json:"codec_long_name"`
	Profile        string            `json:"profile"`
	CodecType      string            `json:"codec_type"`
	CodecTagString string            `json:"codec_tag_string"`
	CodecTag       string            `json:"codec_tag"`
	Width          int               `json:"width"`
	Height         int               `json:"height"`
	CodedWidth     int               `json:"coded_width"`
	CodedHeight    int               `json:"coded_height"`
	HasBFrames     int               `json:"has_b_frames"`
	PixFmt         string            `json:"pix_fmt"`
	Level          int               `json:"level"`
	ColorRange     string            `json:"color_range"`
	ColorSpace     string            `json:"color_space"`
	ColorTransfer  string            `json:"color_transfer"`
	ColorPrimaries string            `json:"color_primaries"`
	ChromaLocation string            `json:"chroma_location"`
	Refs           int               `json:"refs"`
	RFrameRate     string            `json:"r_frame_rate"`
	AvgFrameRate   string            `json:"avg_frame_rate"`
	TimeBase       string            `json:"time_base"`
	StartPts       int               `json:"start_pts"`
	StartTime      string            `json:"start_time"`
	ExtradataSize  int               `json:"extradata_size"`
	Disposition    Disposition       `json:"disposition"`
	Tags           map[string]string `json:"tags"`
}

type Frame struct {
	MediaType               string `json:"media_type"`
	StreamIndex             int    `json:"stream_index"`
	KeyFrame                int    `json:"key_frame"`
	Pts                     int    `json:"pts"`
	PtsTime                 string `json:"pts_time"`
	PktDts                  int    `json:"pkt_dts"`
	PktDtsTime              string `json:"pkt_dts_time"`
	BestEffortTimestamp     int    `json:"best_effort_timestamp"`
	BestEffortTimestampTime string `json:"best_effort_timestamp_time"`
	Duration                int    `json:"duration"`
	DurationTime            string `json:"duration_time"`
	PktPos                  string `json:"pkt_pos"`
	PktSize                 string `json:"pkt_size"`
	Width                   int    `json:"width"`
	Height                  int    `json:"height"`
	CropTop                 int    `json:"crop_top"`
	CropBottom              int    `json:"crop_bottom"`
	CropLeft                int    `json:"crop_left"`
	CropRight               int    `json:"crop_right"`
	PixFmt                  string `json:"pix_fmt"`
	SampleAspectRatio       string `json:"sample_aspect_ratio"`
	PictType                string `json:"pict_type"`
	InterlacedFrame         int    `json:"interlaced_frame"`
	TopFieldFirst           int    `json:"top_field_first"`
	Lossless                int    `json:"lossless"`
	RepeatPict              int    `json:"repeat_pict"`
	ColorRange              string `json:"color_range"`
	ColorSpace              string `json:"color_space"`
	ColorPrimaries          string `json:"color_primaries"`
	ColorTransfer           string `json:"color_transfer"`
	ChromaLocation          string `json:"chroma_location"`
}
