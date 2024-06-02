package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// Define the paths for Videos and Audios subfolders
	videosFolder := filepath.Join(cwd, "Videos")
	audiosFolder := filepath.Join(cwd, "Audios")

	// Ensure the Videos folder exists
	if _, err := os.Stat(videosFolder); os.IsNotExist(err) {
		fmt.Println("Videos folder does not exist:", videosFolder)
		return
	}

	// Ensure the Audios folder exists
	if _, err := os.Stat(audiosFolder); os.IsNotExist(err) {
		err := os.Mkdir(audiosFolder, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating Audios folder:", err)
			return
		}
	}

	// Convert MKV to MP4
	err = filepath.Walk(videosFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".mkv") {
			mp4Path := strings.TrimSuffix(path, ".mkv") + ".mp4"
			cmd := exec.Command("ffmpeg", "-i", path, "-c", "copy", mp4Path)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				return fmt.Errorf("error converting %s to mp4: %v", path, err)
			}
			err = os.Remove(path)
			if err != nil {
				return fmt.Errorf("error removing original mkv file %s: %v", path, err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error during MKV to MP4 conversion:", err)
		return
	}

	// Extract audio from MP4 to MP3
	err = filepath.Walk(videosFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".mp4") {
			mp3Path := filepath.Join(audiosFolder, strings.TrimSuffix(info.Name(), ".mp4")+".mp3")
			cmd := exec.Command("ffmpeg", "-i", path, "-q:a", "0", "-map", "a", mp3Path)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				return fmt.Errorf("error extracting audio from %s to mp3: %v", path, err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error during audio extraction:", err)
	}
}
