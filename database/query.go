package database

const DatabaseFilesQuery string = `
	use master
	select
     ROW_NUMBER() OVER(ORDER BY  files.database_id,files.file_id) AS ID
		,db.name as DatabaseName
		,files.name as FileName
		,files.physical_name as FilePhysicalName
	from 
		sys.master_files as files
	inner join 
		sys.databases as db
	on 
		files.database_id = db.database_id
  `

const DatabaseFileIOQuery string = `
	use master
	select
     ROW_NUMBER() OVER(ORDER BY  files.database_id,files.file_id) AS ID
		,db.name as DatabaseName
		,files.name as FileName
		,files.physical_name as FilePhysicalName
		,io.size_on_disk_bytes as FileSizeOnDiskBytes
		,io.sample_ms as SampleMs
		,io.io_stall as IoStall
		,io.num_of_reads as NumOfReads
		,io.num_of_bytes_read as NumOfBytesRead
		,io.io_stall_read_ms as IoStallReadMs
		,io.num_of_writes as NumOfWrites
		,io.num_of_bytes_written as NumOfBytesWritten
		,io.io_stall_write_ms as IoStallWriteMs
	from 
		sys.dm_io_virtual_file_stats (null,null ) as io
	inner join 
		sys.databases as db
	on 
		io.database_id = db.database_id
	inner join 
		sys.master_files as files
	on 
		io.database_id = files.database_id
	and io.file_id = files.file_id
  `
