package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"syscall"

	"github.com/masiulaniec/Dominator/imagepublishers/amipublisher"
	"github.com/masiulaniec/Dominator/lib/fsutil"
	libjson "github.com/masiulaniec/Dominator/lib/json"
	"github.com/masiulaniec/Dominator/lib/log"
	libtags "github.com/masiulaniec/Dominator/lib/tags"
)

const (
	filePerms = syscall.S_IRUSR | syscall.S_IWUSR | syscall.S_IRGRP |
		syscall.S_IROTH
)

func listUnusedImagesSubcommand(args []string, logger log.DebugLogger) {
	err := listUnusedImages(logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing unused images: %s\n", err)
		os.Exit(1)
	}
	logMemoryUsage(logger)
	os.Exit(0)
}

func listUnusedImages(logger log.DebugLogger) error {
	results, err := amipublisher.ListUnusedImages(targets, skipTargets,
		searchTags, excludeSearchTags, *minImageAge, logger)
	if err != nil {
		return err
	}
	if err := libjson.WriteWithIndent(os.Stdout, "    ", results); err != nil {
		return err
	}
	if *oldImageInstancesCsvFile != "" {
		err := writeInstancesCsv(*oldImageInstancesCsvFile,
			results.OldInstances)
		if err != nil {
			return err
		}
	}
	if *unusedImagesCsvFile != "" {
		err := writeUnusedImagesCsv(*unusedImagesCsvFile,
			results.UnusedImages)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeInstancesCsv(filename string,
	instances []amipublisher.Instance) error {
	file, err := fsutil.CreateRenamingWriter(filename, filePerms)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	// First find all the tag keys.
	tagKeysSet := make(map[string]struct{})
	for _, instance := range instances {
		for key := range instance.Tags {
			tagKeysSet[key] = struct{}{}
		}
	}
	tagKeysList := makeTagKeysList(tagKeysSet)
	header := []string{"Account", "Region", "AmiId", "InstanceId", "LaunchTime"}
	header = append(header, tagKeysList...)
	if err := writer.Write(header); err != nil {
		return err
	}
	for _, instance := range instances {
		record := []string{
			instance.AccountName,
			instance.Region,
			instance.AmiId,
			instance.InstanceId,
			instance.LaunchTime,
		}
		err := appendRecordAndWrite(writer, record, tagKeysList, instance.Tags)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeUnusedImagesCsv(filename string,
	images []amipublisher.Image) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, filePerms)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	// First find all the tag keys.
	tagKeysSet := make(map[string]struct{})
	for _, image := range images {
		for key := range image.Tags {
			tagKeysSet[key] = struct{}{}
		}
	}
	tagKeysList := makeTagKeysList(tagKeysSet)
	header := []string{
		"Account",
		"Region",
		"AmiId",
		"AmiName",
		"CreationDate",
		"Description",
		"Size",
	}
	header = append(header, tagKeysList...)
	if err := writer.Write(header); err != nil {
		return err
	}
	for _, image := range images {
		record := []string{
			image.AccountName,
			image.Region,
			image.AmiId,
			image.AmiName,
			image.CreationDate,
			image.Description,
			strconv.Itoa(int(image.Size)),
		}
		err := appendRecordAndWrite(writer, record, tagKeysList, image.Tags)
		if err != nil {
			return err
		}
	}
	return nil
}

func makeTagKeysList(tagKeysSet map[string]struct{}) []string {
	tagKeysList := make([]string, 0, len(tagKeysSet))
	for key := range tagKeysSet {
		tagKeysList = append(tagKeysList, key)
	}
	sort.Strings(tagKeysList)
	return tagKeysList
}

func appendRecordAndWrite(writer *csv.Writer, record []string,
	tagKeysList []string, tags libtags.Tags) error {
	for _, key := range tagKeysList {
		value := tags[key]
		record = append(record, value)
	}
	return writer.Write(record)
}

func deleteUnusedImagesSubcommand(args []string, logger log.DebugLogger) {
	err := deleteUnusedImages(logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting unused images: %s\n", err)
		os.Exit(1)
	}
	logMemoryUsage(logger)
	os.Exit(0)
}

func deleteUnusedImages(logger log.DebugLogger) error {
	results, err := amipublisher.DeleteUnusedImages(targets, skipTargets,
		searchTags, excludeSearchTags, *minImageAge, logger)
	if err != nil {
		return err
	}
	return libjson.WriteWithIndent(os.Stdout, "    ", results)
}
