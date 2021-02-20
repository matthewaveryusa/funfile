package funfile

import (
	"fmt"
	"log"

	"crawshaw.io/sqlite"
	"golang.org/x/crypto/blake2b"
)

// FunFile test
type FunFile struct {
	c *sqlite.Conn
}

// Connect a FunFile to its database
func (f *FunFile) Connect(dbstring string, poolSize int) (err error) {
	f.c, err = sqlite.OpenConn(dbstring, 0)
	if err != nil {
		log.Fatal(err)
	}
	if f.c == nil {
		log.Fatal("conn is null")
		return
	}
	stmt, _, _ := f.c.PrepareTransient("CREATE TABLE IF NOT EXISTS files (path, file, sum);")
	if hasRow, err := stmt.Step(); err != nil {
		log.Fatal(err)
	} else {
		defer stmt.Finalize()
		if hasRow {
			log.Fatal("Create table returned rows...")
		}
	}
	return err
}

// Disconnect disconnects from the pool
func (f *FunFile) Disconnect() {
	if err := f.c.Close(); err != nil {
		log.Fatal("Disconnect ", err)
	}
}

// AddFile adds a file to the database
func (f *FunFile) AddFile(path string, file []byte) (err error) {
	stmt, err := f.c.Prepare("INSERT INTO files (path, file, sum) VALUES ($path, $file, $sum);")
	if err != nil {
		log.Fatal(err)
		return
	}
	stmt.SetText("$path", path)
	stmt.SetBytes("$file", file)
	sum := blake2b.Sum256(file)
	stmt.SetBytes("$sum", sum[:])
	if _, err = stmt.Step(); err != nil {
		log.Fatal("INSERT AddFile", err)
	}
	if err = stmt.Reset(); err != nil {
		log.Fatal("INSERT (reset) AddFile", err)
	}
	return err
}

// GetFile get a file
func (f *FunFile) GetFile(path string) (err error) {
	stmt, err := f.c.Prepare("SELECT path, file, sum FROM files WHERE path = $path;")
	if err != nil {
		return
	}
	stmt.SetText("$path", path)
	hasRows, err := stmt.Step()
	if err != nil {
		log.Fatal("INSERT AddFile", err)
	} else if !hasRows {
		fmt.Printf("GetFile returned no rows\n")
	} else {
		b := make([]byte, 1024)
		blen := stmt.GetBytes("file", b)
		var sum [blake2b.Size256]byte
		stmt.GetBytes("sum", sum[:])
		fmt.Printf("GetFile returned: %s, %x, %x\n", stmt.GetText("path"), b[:blen], sum)
	}
	if err = stmt.Reset(); err != nil {
		log.Fatal("GetFile (reset) AddFile", err)
	}
	return
}

// DeleteFile delete a file
func (*FunFile) DeleteFile(path *string) (err error) {
	return nil
}

// DeleteDir delete a dir
func (*FunFile) DeleteDir(path *string, recursive bool) (err error) {
	return nil
}

// ListFiles list files
func (*FunFile) ListFiles(path *string, recursive bool) (err error) {
	return nil
}

// ListDirs list dirs
func (*FunFile) ListDirs(path *string, recursive bool) (err error) {
	return nil
}

// ListFilesAndDirs list files and dirs
func (*FunFile) ListFilesAndDirs(path *string, recursive bool) (err error) {
	return nil
}
