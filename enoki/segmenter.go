package enoki

import (
	"bufio"
	"context"
	"log"
	"os/exec"
	"path"
	"path/filepath"
)

type Segmenter struct{}

// Segment is used to divide the provided source file into chunks
// and generate the manifest file.
func (s *Segmenter) Segment(ctx context.Context, p string) error {
	var (
		base     = filepath.Base(p)
		name     = base[:len(base)-len(filepath.Ext(base))]
		manifest = path.Join(filepath.Dir(p), name+".m3u8")
		segmentF = path.Join(filepath.Dir(p), name+"-%d.ts")
	)

	cmdS := []string{"/usr/bin/ffmpeg",
		"-i", p,
		"-codec", "copy",
		"-bsf", "h264_mp4toannexb",
		"-map", "0",
		"-f", "segment",
		"-segment_time", "2",
		"-segment_format", "mpegts",
		"-segment_list", manifest,
		"-segment_list_type", "m3u8",
		segmentF,
	}

	cmd := exec.Command(cmdS[0], cmdS[1:]...)

	stdout, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	in := bufio.NewScanner(stdout)

	for in.Scan() {
		log.Printf(in.Text())
	}

	if err := in.Err(); err != nil {
		log.Printf("error: %s", err)
	}
	return nil
}
