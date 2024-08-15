// Package main
/**
 * Copyright Hangzhou Guangyin Network, Inc. All Rights Reserved.
 * Description of oss
 * @name
 * @action
 * @time 15 8月 2024 13:06 周四
 * @update
 * @author li
 */
package main

import (
	"log"

	"oss/app"
	"oss/app/conf"
)

func init() {
	err := conf.SetupSetting()
	if err != nil {
		log.Fatalf("init setting error: %v", err)
	}
}

func main() {
	m := app.Monitor{}
	m.UploadFile()
}
