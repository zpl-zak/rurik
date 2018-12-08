/*
 * @Author: V4 Games
 * @Date: 2018-11-09 22:29:44
 * @Last Modified by: Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-11-09 22:55:14
 */

package core

import (
	"log"

	"github.com/asdine/storm"
)

// InitDatabase prepares the database connection
func InitDatabase() {
	CurrentSaveSystem.InitSaveSystem()
}

// ShutdownDatabase disposes all database data and closes all connections
func ShutdownDatabase() {
	CurrentSaveSystem.ShutdownSaveSystem()
}

func openDatabase(dbName string) *storm.DB {
	db, err := storm.Open(dbName)

	if err != nil {
		log.Fatalf("Could not open database file %s !\n", dbName)
		return nil
	}

	return db
}
