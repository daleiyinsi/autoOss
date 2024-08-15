// Package app
/**
 * Copyright Hangzhou Guangyin Network, Inc. All Rights Reserved.
 * Description of app
 * @name
 * @action
 * @time 15 8月 2024 13:24 周四
 * @update
 * @author li
 */
package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"oss/app/conf"
	"oss/app/storage"
)

type Monitor struct {
}

func (m *Monitor) UploadFile() {
	m.upload()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	var t *time.Timer
	go func() {
		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok {
					return
				}
				if t != nil {
					t.Stop()
				}
				t = time.AfterFunc(5*time.Second, func() {
					m.upload()
				})
			}
		}
	}()
	err = watcher.Add(conf.Settings.LocalPath)
	if err != nil {
		log.Fatal(err)
	}
	<-make(chan struct{})
}

var (
	ossStorage storage.OSSStorage
	once       sync.Once
)

func (*Monitor) getOss() storage.OSSStorage {
	once.Do(func() {
		s, err := storage.NewOSS()
		if err != nil {
			log.Fatalf("oss init error: %v", err)
		}
		ossStorage = s
	})
	return ossStorage
}

func (m *Monitor) upload() {
	oss := m.getOss()
	paths := oss.List(conf.Settings.Storage.Path)
	localPaths, err := listFilePaths(conf.Settings.LocalPath)
	if err != nil {
		log.Printf("list file error: %v", err)
	}
	var listNames []string
	for _, path := range paths {
		if filepath.Ext(path) == "" {
			continue
		}
		listNames = append(listNames, filepath.Base(path))
	}
	for _, path := range localPaths {
		name := filepath.Base(path)
		if !slices.Contains(listNames, name) {
			if err := oss.PutFromFile(fmt.Sprintf("%s%s", conf.Settings.Storage.Path, name), path); err != nil {
				log.Printf("upload file error: %v", err)
			}
		}
	}
}

func listFilePaths(dirname string) (paths []string, err error) {
	dirname = strings.TrimSuffix(dirname, string(os.PathSeparator))
	err = filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	return
}
