package main

import(
  "time"
	db "github.com/atoyr/SQLServerGo/database"
)

type InstanceIO struct {
	Datetime            time.Time    `json:"datetime"`
  Databases           []DatabaseIO `json:"databases"`
}

type DatabaseIO struct {
	Datetime            time.Time `json:"datetime"`
	DatabaseName        string    `json:"database_name"`
  ID                  int64     `json:"database_id"`
  Files               []FileIO  `json:"files"`
}

type FileIO struct {
  ID                  int64 `json:"id"`
	FileID              int64     `json:"file_id"`
	FileName            string    `json:"file_name"`
	FilePhysicalName    string    `json:"file_physical_name"`
	Datetime            time.Time `json:"datetime"`
	FileSizeOnDiskBytes int64     `json:"file_size_on_disk_bytes"`
	ReadBytesPerSec     float64   `json:"read_bytes_per_sec"`
	WriteBytesPerSec    float64   `json:"write_bytes_per_sec"`
}

func createInstanceIO(dbFileIO []db.DatabaseFileIO) InstanceIO {
  instanceIO := new(InstanceIO)
  if len(dbFileIO) > 0 {
    datetime := dbFileIO[0].Datetime
    instanceIO.Datetime = datetime
    instanceIO.Databases = make([]DatabaseIO,0)
    beforeDBName := ""
    var databaseIO DatabaseIO
    for i := range dbFileIO {
      if beforeDBName != dbFileIO[i].DatabaseName {
        if beforeDBName != "" {
          instanceIO.Databases = append(instanceIO.Databases,databaseIO)
        }
        databaseIO = DatabaseIO{}
        databaseIO.Datetime = datetime
        databaseIO.ID = dbFileIO[i].DatabaseID
        databaseIO.DatabaseName = dbFileIO[i].DatabaseName
        databaseIO.Files = make([]FileIO,0)
        beforeDBName = dbFileIO[i].DatabaseName
      }
      f := FileIO{}
      f.Datetime = datetime
      f.ID = dbFileIO[i].ID
      f.FileID = dbFileIO[i].FileID
      f.FileName = dbFileIO[i].FileName
      f.FilePhysicalName = dbFileIO[i].FilePhysicalName
      f.FileSizeOnDiskBytes = dbFileIO[i].FileSizeOnDiskBytes
      f.ReadBytesPerSec = dbFileIO[i].ReadBytesPerSec
      f.WriteBytesPerSec = dbFileIO[i].WriteBytesPerSec
      databaseIO.Files = append(databaseIO.Files,f)
    }
    instanceIO.Databases = append(instanceIO.Databases,databaseIO)
  }
  return *instanceIO
}
