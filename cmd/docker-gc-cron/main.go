package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/robfig/cron/v3"
)

type Config struct {
	CronSchedule           string
	GracePeriodSeconds     int
	MinimumImagesToSave    int
	ForceContainerRemoval  bool
	ForceImageRemoval      bool
	CleanUpVolumes         bool
	DryRun                 bool
	ExcludeVolumesFile     string
	VolumeDeleteOnlyDriver string
}

func main() {
	config := loadConfig()

	if config.DryRun {
		log.Println("DRY RUN MODE: No containers or images will be removed")
	}

	c := cron.New()

	_, err := c.AddFunc(config.CronSchedule, func() {
		runGarbageCollection(config)
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	c.Start()
	log.Printf("Docker GC Cron started with schedule: %s", config.CronSchedule)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	c.Stop()
	log.Println("Docker GC Cron stopped")
}

func loadConfig() Config {
	cronSchedule := os.Getenv("CRON")
	if cronSchedule == "" {
		cronSchedule = "0 0 * * *"
	}

	gracePeriodSeconds := 3600
	if val := os.Getenv("GRACE_PERIOD_SECONDS"); val != "" {
		fmt.Sscanf(val, "%d", &gracePeriodSeconds)
	}

	minimumImagesToSave := 0
	if val := os.Getenv("MINIMUM_IMAGES_TO_SAVE"); val != "" {
		fmt.Sscanf(val, "%d", &minimumImagesToSave)
	}

	return Config{
		CronSchedule:           cronSchedule,
		GracePeriodSeconds:     gracePeriodSeconds,
		MinimumImagesToSave:    minimumImagesToSave,
		ForceContainerRemoval:  os.Getenv("FORCE_CONTAINER_REMOVAL") == "1",
		ForceImageRemoval:      os.Getenv("FORCE_IMAGE_REMOVAL") == "1",
		CleanUpVolumes:         os.Getenv("CLEAN_UP_VOLUMES") == "1",
		DryRun:                 os.Getenv("DRY_RUN") == "1",
		ExcludeVolumesFile:     os.Getenv("EXCLUDE_VOLUMES_IDS_FILE"),
		VolumeDeleteOnlyDriver: os.Getenv("VOLUME_DELETE_ONLY_DRIVER"),
	}
}

func runGarbageCollection(config Config) {
	log.Println("[Docker GC] Starting garbage collection")

	cleanContainers(config)
	cleanImages(config)

	if config.CleanUpVolumes {
		cleanVolumes(config)
	}

	log.Println("[Docker GC] Garbage collection completed")
}

func cleanContainers(config Config) {
	log.Println("Cleaning up exited containers")

	cmd := exec.Command("docker", "ps", "-a", "-q", "-f", "status=exited")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to list exited containers: %v", err)
		return
	}

	containerIDs := strings.TrimSpace(string(output))
	if containerIDs == "" {
		log.Println("No exited containers found")
		return
	}

	containerList := strings.Split(containerIDs, "\n")

	for _, containerID := range containerList {
		containerID = strings.TrimSpace(containerID)
		if containerID == "" {
			continue
		}

		log.Printf("Removing container: %s", containerID)

		if !config.DryRun {
			args := []string{"rm"}
			if config.ForceContainerRemoval {
				args = append(args, "-f")
			}
			args = append(args, containerID)

			rmCmd := exec.Command("docker", args...)
			if err := rmCmd.Run(); err != nil {
				log.Printf("Failed to remove container %s: %v", containerID, err)
			}
		}
	}
}

func cleanImages(config Config) {
	log.Println("Cleaning up unused images")

	cmd := exec.Command("docker", "images", "-q", "--filter", "dangling=true")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to list dangling images: %v", err)
		return
	}

	imageIDs := strings.TrimSpace(string(output))
	if imageIDs == "" {
		log.Println("No dangling images found")
		return
	}

	imageList := strings.Split(imageIDs, "\n")

	for _, imageID := range imageList {
		imageID = strings.TrimSpace(imageID)
		if imageID == "" {
			continue
		}

		log.Printf("Removing image: %s", imageID)

		if !config.DryRun {
			args := []string{"rmi"}
			if config.ForceImageRemoval {
				args = append(args, "-f")
			}
			args = append(args, imageID)

			rmCmd := exec.Command("docker", args...)
			if err := rmCmd.Run(); err != nil {
				log.Printf("Failed to remove image %s: %v", imageID, err)
			}
		}
	}
}

func cleanVolumes(config Config) {
	log.Println("Cleaning up dangling volumes")

	cmd := exec.Command("docker", "volume", "ls", "-qf", "dangling=true")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to list dangling volumes: %v", err)
		return
	}

	volumeIDs := strings.TrimSpace(string(output))
	if volumeIDs == "" {
		log.Println("No dangling volumes found")
		return
	}

	volumeList := strings.Split(volumeIDs, "\n")

	excludedVolumes := loadExcludedVolumes(config.ExcludeVolumesFile)

	for _, volumeID := range volumeList {
		volumeID = strings.TrimSpace(volumeID)
		if volumeID == "" {
			continue
		}

		if excludedVolumes[volumeID] {
			log.Printf("Skipping excluded volume: %s", volumeID)
			continue
		}

		log.Printf("Removing volume: %s", volumeID)

		if !config.DryRun {
			rmCmd := exec.Command("docker", "volume", "rm", volumeID)
			if err := rmCmd.Run(); err != nil {
				log.Printf("Failed to remove volume %s: %v", volumeID, err)
			}
		}
	}
}

func loadExcludedVolumes(filename string) map[string]bool {
	excluded := make(map[string]bool)

	if filename == "" {
		return excluded
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return excluded
	}

	lines := string(content)
	for _, line := range strings.Split(lines, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			excluded[line] = true
		}
	}

	return excluded
}
