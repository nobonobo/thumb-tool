package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"
)

func main() {
	var tmpl Template
	data, err := os.ReadFile(config.template)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(data, &tmpl); err != nil {
		log.Fatal(err)
	}
	abs, err := filepath.Abs(config.template)
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(abs)
	svgPath := tmpl.SVG
	if !filepath.IsAbs(svgPath) {
		svgPath = filepath.Join(dir, svgPath)
	}
	tempDir, err := os.MkdirTemp("", "thumb-tool-*")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("tempDir:", tempDir)
	defer os.RemoveAll(tempDir)

	log.Println("svg:", svgPath)
	inputPath := config.args[0]
	fileExt := filepath.Ext(inputPath)
	metaInfoPath := filepath.Join(tempDir, "metainfo.json")
	cmd := exec.Command(config.ffprobe,
		"-hide_banner", "-loglevel", "error",
		"-select_streams", "v:0",
		"-show_streams",
		"-of", "json",
		inputPath,
	)
	stdout := bytes.NewBuffer(nil)
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	var probeInfo ProbeInfo
	if err := json.Unmarshal(stdout.Bytes(), &probeInfo); err != nil {
		log.Fatal(err)
	}
	video := probeInfo.Streams[0]
	codec, err := DetectHardwareEncoder(video.CodecName)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s: %s(%s)", inputPath, codec, video.CodecName)
	frames, err := GetProbeInfo(
		config.ffprobe,
		inputPath,
		int(3*config.svgDuration/time.Second/2), // 1.5倍の秒数
	)
	if err != nil {
		log.Fatal(err)
	}
	splitFrame := frames[0]
	for _, frame := range frames[1:] {
		splitFrame = frame
		if frame.Pts >= int(config.svgDuration/time.Millisecond) {
			break
		}
	}
	splitBeforePath := filepath.Join(tempDir, "split_before.mkv")
	splitAfterPath := filepath.Join(tempDir, "split_after.mkv")
	introPath := filepath.Join(tempDir, "intro.mkv")
	cmd = exec.Command(config.ffmpeg,
		"-hide_banner", "-loglevel", "error",
		"-i", inputPath,
		"-t", fmt.Sprintf("%f", float64(splitFrame.Pts)/1000),
		"-avoid_negative_ts", "make_zero",
		"-c", "copy", splitBeforePath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	cmd = exec.Command(config.ffmpeg,
		"-hide_banner", "-loglevel", "error",
		"-ss", fmt.Sprintf("%f", float64(splitFrame.Pts)/1000),
		"-i", inputPath,
		"-avoid_negative_ts", "make_zero",
		"-c", "copy", splitAfterPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	svgOutPath := filepath.Join(tempDir, "template.svg")
	pngPath := filepath.Join(tempDir, "background.png")
	overlayPath := filepath.Join(tempDir, "overlay.png")
	thumbnailPath := filepath.Join(tempDir, "thumbnail.jpg")
	outputPath := filepath.Join(tempDir, "output"+fileExt)
	if err := patchSVG(svgPath, svgOutPath); err != nil {
		log.Fatal(err)
	}
	cmd = exec.Command("ffmpeg",
		"-hide_banner", "-loglevel", "error",
		"-ss", fmt.Sprintf("%.3f", float64(config.offsetSec)/float64(time.Second)), // 高速キーシーク
		"-i", inputPath,
		"-frames:v", "1",
		"-c:v", "png",
		"-q:v", "0", // 最高品質
		"-y", // 上書き許可
		pngPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	cmd = exec.Command(config.inkscape,
		"--without-gui",
		"--export-type=png",
		fmt.Sprintf("--export-width=%d", video.Width),
		fmt.Sprintf("--export-height=%d", video.Height),
		"--export-filename="+overlayPath,
		svgOutPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	width := video.Height * 16 / 9
	cmd = exec.Command("ffmpeg",
		"-hide_banner", "-loglevel", "error",
		"-i", pngPath,
		"-i", overlayPath,
		"-filter_complex", fmt.Sprintf(
			"[0:v][1:v]overlay,format=yuvj420p[ovr];[ovr]crop=%d:%d:%d:%d,format=yuv420p",
			width, video.Height, (video.Width-width)/2, 0,
		),
		"-q:v", "8", // JPEG品質
		"-y", // 上書き許可
		thumbnailPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	cmd = exec.Command("ffmpeg",
		"-hide_banner", "-loglevel", "error",
		"-loop", "1",
		"-t", "10",
		"-i", overlayPath,
		"-i", splitBeforePath,
		"-filter_complex", `[0:v]format=yuva420p,fade=t=out:st=8:d=1.8:alpha=1[ovr];[1:v][ovr]overlay=(W-w)/2:(H-h)/2[outv]`,
		"-map", `[outv]`,
		"-map", "1:a",
		"-c:v", codec,
		"-b:v", "20M",
		"-c:a", "copy",
		"-pix_fmt", "yuv420p",
		"-y",
		introPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	listPath := filepath.Join(tempDir, "list.txt")
	buff := bytes.NewBuffer(nil)
	fmt.Fprintf(buff, "file '%s'\n", introPath)
	fmt.Fprintf(buff, "file '%s'\n", splitAfterPath)
	if err := os.WriteFile(listPath, buff.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
	cmd = exec.Command(config.ffmpeg,
		"-hide_banner", "-loglevel", "error",
		"-f", "concat",
		"-safe", "0",
		"-i", listPath,
		"-c", "copy",
		"-y", outputPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	params := tmpl.Defaults
	for k, v := range config.params {
		params[k] = v
	}
	b, err := json.Marshal(tmpl.UploadMeta)
	if err != nil {
		log.Fatal(err)
	}
	t, err := template.New("upload").Parse(string(b))
	if err != nil {
		log.Fatal(err)
	}
	buff = bytes.NewBuffer(nil)
	if err := t.Execute(buff, params); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(metaInfoPath, buff.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
	cmd = exec.Command(
		config.youtubeuploader,
		"-cache", filepath.Join(config.cacheDir, "request.token"),
		"-secrets", config.secrets,
		"-metaJSON", metaInfoPath,
		"-thumbnail", thumbnailPath,
		"-filename", outputPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Println("meta:", buff.String())
		log.Fatal(err)
	}
}
