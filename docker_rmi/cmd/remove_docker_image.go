package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type dockerImageDate struct {
	Unit  string
	Value int
}

type dockerImageInfo struct {
	Repository string
	Tag        string
	ImageID    string
	Created    dockerImageDate
	Size       string
}

var rootCmd = &cobra.Command{
	Use:     "remove",
	Example: "",
	Version: "v0.1.0",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func Exec() {
	images := getDockerImages()
	i := make([]string, 0)
	for _, v := range images {
		i = append(i, v.ImageID)
	}
	removeImages(i)
}

func getDockerImages() []*dockerImageInfo {
	out, err := exec.Command("docker", "images").Output()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	res := strings.Split(string(out), "\n")
	images := make([]*dockerImageInfo, 0)
	reg, err := regexp.Compile("(\\S*)\\s*(\\S*)\\s*(\\S*)\\s*(\\S*)\\s*(\\S*)\\s*(\\S*)\\s*(\\S*)")
	if err != nil {
		fmt.Println("compile regexp failed with err: ", err)
		return nil
	}
	for i, s := range res {
		// first line is title, skip it
		if i == 0 {
			continue
		}
		matches := reg.FindStringSubmatch(s)
		if len(matches) != 8 {
			fmt.Println("regexp match failed, string: ", s)
		}
		if matches[3] == "" {
			continue
		}
		info := &dockerImageInfo{
			Repository: matches[1],
			Tag:        matches[2],
			ImageID:    matches[3],
			Size:       matches[7],
		}
		v, err := strconv.Atoi(matches[4])
		if err != nil {
			info.Created = dockerImageDate{
				Unit:  matches[5],
				Value: v,
			}
		}
		images = append(images, info)
	}
	return images
}

func removeImages(images []string) {
	for _, v := range images {
		_, err := exec.Command("docker", "rmi", v).Output()
		if err != nil {
			fmt.Println(fmt.Sprintf("remove docker image: %s error: %v", v, err))
		}
	}
}
