package database

const DatabasesQuery string = `
  use master
  SELECT 
     name
    ,database_id
    ,state_desc
    ,recovery_model_desc
  FROM sys.databases
`

const DatabaseFilesQuery string = `
	use master
	select
     ROW_NUMBER() OVER(ORDER BY  files.database_id,files.file_id) AS ID
		,db.name as DatabaseName
		,files.name as FileName
		,files.physical_name as FilePhysicalName
	from 
		sys.master_files as files with(nolock)
	inner join 
		sys.databases as db with(nolock)
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
		sys.databases as db with(nolock)
	on 
		io.database_id = db.database_id
	inner join 
		sys.master_files as files with(nolock)
	on 
		io.database_id = files.database_id
	and io.file_id = files.file_id
  `

const ServerInfoQuery string = `
    SELECT 
      SERVERPROPERTY('MachineName')
    , SERVERPROPERTY('InstanceName')
    , SERVERPROPERTY('ServerName')
    , SERVERPROPERTY('productversion')
    , SERVERPROPERTY('ProductMajorVersion')
    , SERVERPROPERTY('ProductLevel')
    , SERVERPROPERTY('Edition')
  `

const CpuUsedQuery string = `
  SELECT TOP 1
      record.value('(./Record/@id)[1]', 'int') AS record_id
    , record.value('(./Record/SchedulerMonitorEvent/SystemHealth/SystemIdle)[1]', 'int') AS [SystemIdle]
    , record.value('(./Record/SchedulerMonitorEvent/SystemHealth/ProcessUtilization)[1]', 'int') AS [SQLProcessUtilization]
    , [timestamp]
  FROM (
    SELECT [timestamp]
    , convert(XML, record) AS [record]
    FROM sys.dm_os_ring_buffers with(nolock)
    WHERE ring_buffer_type = N'RING_BUFFER_SCHEDULER_MONITOR'
    AND record LIKE '%<systemHealth>%'
  ) AS x
  ORDER BY timestamp desc
`
