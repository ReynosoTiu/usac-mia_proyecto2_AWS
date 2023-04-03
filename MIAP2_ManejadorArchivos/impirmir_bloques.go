package main

import "fmt"

func ImprimirInodo(inodo TABLA_INODOS) {
	fmt.Println("I_UID:", inodo.I_uid)
	fmt.Println("I_gid:", inodo.I_gid)
	fmt.Println("I_size:", inodo.I_size)
	fmt.Println("I_atime:", string(inodo.I_atime[:]))
	fmt.Println("I_actime:", string(inodo.I_ctime[:]))
	fmt.Println("I_mtime:", string(inodo.I_mtime[:]))
	for i := 0; i < 16; i++ {
		fmt.Println("I_block[", i, "]:", inodo.I_block[i])
	}

	fmt.Println("I_type:", string(inodo.I_type))
	fmt.Println("I_perm:", inodo.I_perm)

}

func ImprimirSuper(sp SUPER_BLOQUE) {
	fmt.Println("S_filesystem_type:", sp.S_filesystem_type)
	fmt.Println("S_inodes_count:", sp.S_inodes_count)
	fmt.Println("S_block_cout:", sp.S_blocks_count)
	fmt.Println("S_free_blocks_count:", sp.S_free_blocks_count)
	fmt.Println("S_free_inodes_count:", sp.S_free_inodes_count)
	fmt.Println("S_mtime:", aString(sp.S_mtime[:]))
	//fmt.Println("S_umtime:", aString(sp.S_umtime[:]))
	fmt.Println("S_mnt_count:", sp.S_mnt_count)
	fmt.Println("S_magic:", sp.S_magic)
	fmt.Println("S_inode_size:", sp.S_inode_size)
	fmt.Println("S_block_size:", sp.S_block_size)
	fmt.Println("S_first_ino:", sp.S_first_ino)
	fmt.Println("S_first_block:", sp.S_first_blo)
	fmt.Println("S_bm_inode_start:", sp.S_bm_inode_start)
	fmt.Println("S_bm_block_start:", sp.S_bm_block_start)
	fmt.Println("S_inode_start:", sp.S_inode_start)
	fmt.Println("S_block_start:", sp.S_block_start)

}
