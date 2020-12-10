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
  ID                  int64     `json:"id"`
	FileID              int64     `json:"file_id"`
	FileName            string    `json:"file_name"`
	FilePhysicalName    string    `json:"file_physical_name"`
	Datetime            time.Time `json:"datetime"`
	FileSizeOnDiskBytes int64     `json:"file_size_on_disk_bytes"`
	ReadBytesPerSec     float64   `json:"read_bytes_per_sec"`
	WriteBytesPerSec    float64   `json:"write_bytes_per_sec"`
  Count               int64     `json:"count"`
  MinReadBytesPerSec  float64   `json:"min_read_bytes_per_sec"`
  MaxReadBytesPerSec  float64   `json:"max_read_bytes_per_sec"`
  AvgReadBytesPerSec  float64   `json:"avg_read_bytes_per_sec"`
  MinWriteBytesPerSec float64   `json:"min_write_bytes_per_sec"`
  MaxWriteBytesPerSec float64   `json:"max_write_bytes_per_sec"`
  AvgWriteBytesPerSec float64   `json:"avg_write_bytes_per_sec"`
}

type Memory struct {
  Datetime          time.Time `json:"datetime"`
	PhysicalMemory    int       `json:"physical_memory"`
	UsedMemory        int       `json:"used_memory"`
  MinUsedMemory     int       `json:"min_used_memory"`
  MaxUsedMemory     int       `json:"max_used_memory"`
  AvgUsedMemory     int       `json:"avg_used_memory"`
	AvailableMemory   int       `json:"available_memory"`
	TotalPageFile     int       `json:"total_page_file"`
	UsedPageFile      int       `json:"used_page_file"`
	AvailablePageFile int       `json:"available_page_file"`
	SystemChace       int       `json:"system_chace"`
	SystemMemoryState string    `json:"system_memory_state"`
}

var dbFileIOs [][]db.DatabaseFileIO
var memorys []db.Memory

func createInstanceIO(dbFileIO []db.DatabaseFileIO) InstanceIO {
  if dbFileIOs == nil {
    dbFileIOs = make([][]db.DatabaseFileIO,0)
  }
  dbFileIOs = append(dbFileIOs,dbFileIO)
  if len(dbFileIOs) > int(websocketCount) {
    dbFileIOs = dbFileIOs[len(dbFileIOs) - int(websocketCount):]
  }

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

      var minRead, sumRead, maxRead float64
      var minWrite, sumWrite, maxWrite float64
      for j := range dbFileIOs {
        r := dbFileIOs[j][i].ReadBytesPerSec
        w := dbFileIOs[j][i].WriteBytesPerSec
        if j == 0 {
          minRead = r
          sumRead = r
          maxRead = r
          minWrite = w
          sumWrite = w
          maxWrite = w
        }else {
          minRead = min(minRead,r)
          maxRead = max(maxRead,r)
          sumRead = sumRead + r
          minWrite = min(minWrite,w)
          maxWrite = max(maxWrite,w)
          sumWrite = sumWrite + w
        }
      }
      f.MinReadBytesPerSec = minRead
      f.MaxReadBytesPerSec = maxRead
      f.AvgReadBytesPerSec = sumRead / float64(len(dbFileIOs))
      f.MinWriteBytesPerSec = minWrite
      f.MaxWriteBytesPerSec = maxWrite
      f.AvgWriteBytesPerSec = sumWrite / float64(len(dbFileIOs))

      databaseIO.Files = append(databaseIO.Files,f)
    }
    instanceIO.Databases = append(instanceIO.Databases,databaseIO)
  }

  return *instanceIO
}

func min(a, b float64) float64 {
  if a < b {
    return a
  }
  return b
}

func max(a, b float64) float64 {
  if a > b {
    return a
  }
  return b
}

