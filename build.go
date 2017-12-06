package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//https://stackoverflow.com/questions/17796043/embedding-text-file-into-compiled-executable

const assetsConfigFile = "src/assets.json"

func main() {
	// assetsDir := os.Args[1]
	// sourceFileName := os.Args[2]
	// constPrefix := os.Args[3]
	// assetExt := os.Args[4]

	assetsDefinitionStr, err := ioutil.ReadFile(assetsConfigFile)
	if err != nil {
		panic(err)
	}

	assets := make(map[string]string)

	json.Unmarshal([]byte(assetsDefinitionStr), &assets)

	for varname, assetPath := range assets {
		fmt.Println(varname, " = ", assetPath)
	}

	// fs, _ := ioutil.ReadDir(assetsDir)
	// out, _ := os.Create(sourceFileName)
	// out.Write([]byte("package main \n\nconst (\n"))
	// for _, f := range fs {
	// 	if strings.HasSuffix(f.Name(), assetExt) {
	// 		out.Write([]byte(strings.TrimSuffix(fmt.Sprintf(constPrefix, f.Name()), assetExt) + " = `"))
	// 		f, _ := os.Open(f.Name())
	// 		io.Copy(out, f)
	// 		out.Write([]byte("`\n"))
	// 	}
	// }
	// out.Write([]byte(")\n"))
}
