// package main

// import (
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"os"
// 	"strings"
// )

//https://stackoverflow.com/questions/17796043/embedding-text-file-into-compiled-executable

// func main(assetsDir string, sourceFileName string, constPrefix string, assetExt string) {
// 	fs, _ := ioutil.ReadDir(assetsDir)
// 	out, _ := os.Create(sourceFileName)
// 	out.Write([]byte("package main \n\nconst (\n"))
// 	for _, f := range fs {
// 		if strings.HasSuffix(f.Name(), assetExt) {
// 			out.Write([]byte(strings.TrimSuffix(fmt.Sprintf(constPrefix, f.Name()), assetExt) + " = `"))
// 			f, _ := os.Open(f.Name())
// 			io.Copy(out, f)
// 			out.Write([]byte("`\n"))
// 		}
// 	}
// 	out.Write([]byte(")\n"))
// }
