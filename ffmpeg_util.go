package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// CodecMapping コーデックごとの優先順位リスト
var CodecMapping = map[string][]string{
	"h264": {
		"h264_nvenc", "h264_amf", "h264_qsv", "h264_vaapi", "libx264",
	},
	"hevc": {
		"hevc_nvenc", "hevc_amf", "hevc_qsv", "hevc_vaapi", "libx264",
	},
	"av1": {
		"av1_nvenc", "av1_amf", "av1_qsv", "libsvtav1", "libaom-av1",
	},
}

// DetectHardwareEncoder 指定コーデックに対応する最適HWエンコーダを検出
func DetectHardwareEncoder(codec string) (string, error) {
	priorityList, exists := CodecMapping[strings.ToLower(codec)]
	if !exists {
		return "", fmt.Errorf("未サポートのコーデック: %s", codec)
	}

	for _, encoder := range priorityList {
		if err := testEncoder(encoder); err == nil {
			return encoder, nil
		}
	}

	// HWなしの場合、ソフトウェアエンコーダを返す
	for _, encoder := range priorityList {
		if strings.HasPrefix(encoder, "lib") {
			return encoder, nil
		}
	}

	return "", fmt.Errorf("コーデック %s のエンコーダが見つかりませんでした", codec)
}

// testEncoder 個別エンコーダ動作テスト
func testEncoder(encoder string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-y",
		"-f", "lavfi",
		"-i", "testsrc=duration=1:size=320x240:rate=30,format=yuv420p",
		"-t", "0.1",
		"-c:v", encoder,
		"-f", "null", "-",
	)

	var devNull bytes.Buffer
	cmd.Stdout = &devNull
	cmd.Stderr = &devNull

	return cmd.Run()
}

// GetProbeInfo executes ffprobe to extract keyframes info from a video and parses the JSON output.
func GetProbeInfo(ffprobePath, filePath string, sec int) ([]Frame, error) {
	if ffprobePath == "" {
		ffprobePath = "ffprobe"
	}
	args := []string{
		"-loglevel", "error",
		"-select_streams", "v:0",
		"-skip_frame", "nokey",
		"-show_frames",
		"-of", "json",
	}
	if sec > 0 {
		args = append(args, "-read_intervals", fmt.Sprintf("%%+%d", sec))
	}
	args = append(args, filePath)
	cmd := exec.Command(ffprobePath, args...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	// capture stderr in case of error? (optional, skip for now to keep simple)

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var info ProbeInfo
	if err := json.Unmarshal(stdout.Bytes(), &info); err != nil {
		return nil, err
	}

	return info.Frames, nil
}
