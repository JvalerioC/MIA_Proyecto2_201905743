package main

// se asignara un tamano en especifico al arreglo de bytes sino el tamano de mbr o cualquier estructura seria variable

// struct de tamano 45
type Partition struct {
	Part_status [1]byte  //tamano 1, valores activa:1, no activa:0
	Part_type   [1]byte  //tamano 1, valores P o E
	Part_fit    [1]byte  //tamano, valores B, F o W
	Part_start  [10]byte //tamano, valor maximo 1Gb = 1024 MB = 1073741824 bytes
	Part_size   [10]byte //tamano, valor maximo 1GB = 1073741824 bytes
	Part_name   [22]byte //tamano , le vamo a dar una longitud de 22 para redondear
}

// struct de tamano 215 bytes
type MBR struct {
	Mbr_tamano         [10]byte     //tamano, valor maximo 1GB = 1073741824 bytes
	Mbr_fecha_creacion [20]byte     //formato DD/MM/YYYY hh:mm:ss
	Mbr_dsk_signature  [4]byte      //tamano 10, numero aleatorio entre 0 y 9999
	Dsk_fit            [1]byte      //tamano 1, valores B, F, W
	Mbr_Partition      [4]Partition // tamano 180 bytes
}

// struct de tamano 55 bytes
type EBR struct {
	Part_status [1]byte  //tamano 1, valores activa:1, no activa:0
	Part_fit    [1]byte  //tamano, valores B, F o W
	Part_start  [10]byte //tamano, valor maximo 1Gb = 1073741824 bytes
	Part_size   [10]byte //tamano, valor maximo 1GB = 1073741824 bytes
	Part_next   [10]byte //tamano, valor maximo 1GB = 1073741824 bytes
	Part_name   [23]byte //tamano , le vamo a dar una longitud de 23 para redondear
}

// struct de tamano 130 bytes
type SuperBloque struct {
	S_filesystem_type   [1]byte  //tamano 1 byte, es ext2 o ext3
	S_inodes_count      [10]byte //tamano 10 bytes
	S_blocks_count      [10]byte //tamano 10 bytes
	S_free_blocks_count [10]byte //tamano 10 bytes
	S_free_inodes_count [10]byte //tamano 10 bytes
	S_mtime             [20]byte //tamano 20 bytes, //formato DD/MM/YYYY hh:mm:ss
	S_mnt_count         [1]byte  //tamano 1 byte,
	S_magic             [6]byte  //tamano 5 bytes,valor 6,2,1,6,7
	S_inode_size        [1]byte  //tamano 1 byte, valor del tamano del inodo
	S_block_size        [1]byte  //tamano 1 byte, valor del tamano del bloque
	S_first_ino         [10]byte //tamano 10 bytes, valor del primer inodo
	S_first_blo         [10]byte //tamano 10 bytes, valor del primer bloque
	S_bm_inode_start    [10]byte //tamano 10 bytes, valor del inicio del bitmap de inodos
	S_bm_block_start    [10]byte //tamano 10 bytes, valor del inicio del bitmap de bloques
	S_inode_start       [10]byte //tamano 10 bytes, valor del inicio de la tabla de inodos
	S_block_start       [10]byte //tamano 10 bytes, valor del inicio de la tabla de bloques

}

// struct de tamano 216 bytes, segun el programa
type Inodo struct {
	I_uid   [3]byte  //tamano 3 bytes, valor del id del usuario
	I_gid   [2]byte  //tamano 2 bytes, valor del id del grupo
	I_size  [10]byte //tamano 10 bytes, valor del tamano del archivo
	I_atime [20]byte //tamano 20 bytes, //formato DD/MM/YYYY hh:mm:ss
	I_ctime [20]byte //tamano 20 bytes, //formato DD/MM/YYYY hh:mm:ss
	I_mtime [20]byte //tamano 20 bytes, //formato DD/MM/YYYY hh:mm:ss
	I_block [16]int  //tamano 16 int, array de posiciones
	I_type  [1]byte  //tamano 1 bytes, valor 1:archivo, 0:carpeta
	I_perm  [3]byte  //tamano 3 bytes, valor de los permisos
}

// struct de tamano 30
type B_Contenido struct {
	B_name  [20]byte //tamano 12 bytes, nombre del archivo
	B_inodo [10]byte //tamano 10 bytes, valor del inodo
}

// struct de tamaño 120 bytes
type B_Carpeta struct {
	B_content [4]B_Contenido //tamano 120 bytes
}

// struct de tamano 64 bytes
type B_Archivo struct {
	B_content [64]byte //tamano 64 bytes
}

// struct para los montajes
type itemMount struct {
	Part   Partition //particion que fue montada, creo que solo neceesito el nombre y donde inicia
	Id     string    //id que identifica la montura
	Path   string    //para saber la ruta del disco, reportes
	Number int       //para saber que numero tiene asociado
}

// struct para el login
type Usuario struct {
	Iniciado  bool
	Admin     int
	User      string
	LoginItem itemMount
	Grupo     string
	Grupo_id  int
	User_id   int
}
