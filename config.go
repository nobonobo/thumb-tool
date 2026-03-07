package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var config = struct {
	secrets         string
	template        string
	inkscape        string
	ffprobe         string
	ffmpeg          string
	youtubeuploader string
	offsetSec       time.Duration
	svgDuration     time.Duration
	params          map[string]string
	args            []string
	configDir       string
	cacheDir        string
}{
	secrets:         "client_secret.json",
	template:        "template.json",
	inkscape:        "inkscape",
	ffprobe:         "ffprobe",
	ffmpeg:          "ffmpeg",
	youtubeuploader: "youtubeuploader",
	offsetSec:       0,
	params:          map[string]string{},
	svgDuration:     10 * time.Second,
}

func init() {
	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	workDir := os.Getenv("WORK")
	if strings.HasPrefix(executable, workDir) {
		executable = "thumb-tool"
	}
	name := strings.TrimSuffix(filepath.Base(executable), filepath.Ext(executable))
	configRoot, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	configDir := filepath.Join(configRoot, name)
	if err := os.MkdirAll(configDir, 0644); err != nil {
		log.Fatal(err)
	}
	cacheRoot, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}
	cacheDir := filepath.Join(cacheRoot, name)
	if err := os.MkdirAll(cacheDir, 0644); err != nil {
		log.Fatal(err)
	}
	config.configDir = configDir
	config.cacheDir = cacheDir
	config.secrets = filepath.Join(configDir, "client_secret.json")
	flag.StringVar(&config.secrets, "secrets", config.secrets, "Google API secrets file")
	flag.StringVar(&config.template, "t", config.template, "Template file")
	flag.StringVar(&config.inkscape, "inkscape", config.inkscape, "Inkscape executable")
	flag.StringVar(&config.ffprobe, "ffprobe", config.ffprobe, "FFprobe executable")
	flag.StringVar(&config.ffmpeg, "ffmpeg", config.ffmpeg, "FFmpeg executable")
	flag.StringVar(&config.youtubeuploader, "youtubeuploader", config.youtubeuploader, "YoutubeUploader executable")
	flag.DurationVar(&config.offsetSec, "offset", config.offsetSec, "Thumbnail offset")
	flag.DurationVar(&config.svgDuration, "duration", config.svgDuration, "Thumbnail duration")
	flag.Parse()
	if _, err := os.Stat(config.secrets); err != nil {
		if os.IsNotExist(err) {
			log.Fatal("secrets not found: you need get client_secret.json from Google API Console and install to ", config.configDir)
		}
		log.Fatal(err)
	}
	if runtime.GOOS == "windows" {
		config.inkscape = "C:\\Program Files\\Inkscape\\bin\\inkscape.exe"
		if !strings.HasSuffix(config.ffprobe, ".exe") {
			config.ffprobe += ".exe"
		}
		if !strings.HasSuffix(config.ffmpeg, ".exe") {
			config.ffmpeg += ".exe"
		}
		if !strings.HasSuffix(config.youtubeuploader, ".exe") {
			config.youtubeuploader += ".exe"
		}
	}
	inkscape, err := exec.LookPath(config.inkscape)
	if err != nil {
		log.Fatal("inkscape not found")
	}
	config.inkscape = inkscape
	ffprobe, err := exec.LookPath(config.ffprobe)
	if err != nil {
		log.Fatal("ffprobe not found")
	}
	ffmpeg, err := exec.LookPath(config.ffmpeg)
	if err != nil {
		log.Fatal("ffmpeg not found")
	}
	config.ffprobe = ffprobe
	config.ffmpeg = ffmpeg
	youtubeuploader, err := exec.LookPath(config.youtubeuploader)
	if err != nil {
		log.Fatal("youtubeuploader not found")
	}
	config.youtubeuploader = youtubeuploader
	for _, arg := range flag.Args() {
		vars := strings.SplitN(arg, "=", 2)
		if len(vars) != 2 {
			config.args = append(config.args, arg)
			continue
		}
		key := strings.TrimSpace(vars[0])
		config.params[key] = strings.TrimSpace(vars[1])
	}
	log.Println("Template:", config.template)
	log.Println("Inkscape:", config.inkscape)
	log.Println("FFprobe:", config.ffprobe)
	log.Println("FFmpeg:", config.ffmpeg)
	log.Println("YoutubeUploader:", config.youtubeuploader)
	log.Println("Offset:", config.offsetSec)
	log.Println("Params:", config.params)
	log.Println("Args:", config.args)
	log.Println("ConfigDir:", config.configDir)
	log.Println("CacheDir:", config.cacheDir)
}
