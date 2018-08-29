package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var maxLevels int = -1
var currentLevels int
var rootCmd = &cobra.Command{
	Use:   "tree",
	Short: "a tree clone",
	Long:  `shows a tree of files and directories`,
	Run: func(cmd *cobra.Command, args []string) {
		dirArgs := []string{"."} //argumentos recebem uma lista de strings
		if len(args) > 0 {       //se o tamanho dos args for maior que um,
			dirArgs = args
		}

		for _, arg := range dirArgs { //passa por todos os args
			str, err := tree(arg, "")
			if err != nil {
				log.Printf("tree %s: %v\n", arg, err)
			}
			fmt.Print(str)
		}
	},
}

func main() {
	rootCmd.PersistentFlags().IntVar(&maxLevels, "levels", -1, "Max number of levels")
	rootCmd.Execute()
}

func tree(root, indent string) (string, error) { //criamos uma recursao para fazer a "arvore" - indent (indentation)
	fi, err := os.Stat(root) //retorna informacoes
	if err != nil {
		return "", fmt.Errorf("could not stat %s: %v", root, err)
	}

	text := fmt.Sprintln(fi.Name(), "[", ByteCountSI(fi.Size()), "]") //"printa" o nome e o tamanho
	if !fi.IsDir() {                                                  //se nao for um diretorio, nao tem mais nada o que fazer
		return text, nil
	}

	fis, err := ioutil.ReadDir(root) // ReadDir reads the directory named by dirname and returns a list of directory entries sorted by filename
	if err != nil {
		return text, fmt.Errorf("could not read dir %s: %v", root, err)
	}

	if fi.Name() == ".." {
		currentLevels--
	} else {
		currentLevels++
	}
	if maxLevels != -1 && currentLevels > maxLevels {
		return text, nil
	}

	var names []string //criou isso pq quando era o ultimo, dava erro
	for _, fi := range fis {
		if fi.Name()[0] != '.' { //se n for ., adiciona o nome
			names = append(names, fi.Name())
		}
	}

	for i, name := range names { //
		add := "│  "           // sem isso, ficavam separados. Só ficava └── ou ├──, nao tinha juncao
		if i == len(names)-1 { //se for o ultimo, n precisa ├──, se n fica errado
			text += fmt.Sprintf(indent + "└──") //printa a posicao, e para isso precisa saber o indent
			add = "   "                         // se for o ultimo, nao vai printar │ , e sim "  "(espaco)
		} else {
			text += fmt.Sprintf(indent + "├──") //se for um diretorio tem q ter esse formato

		}

		text2, err := tree(filepath.Join(root, name), indent+add)
		if err != nil { //indent+add para criar os "espacos"
			return text, err
		}
		text += text2
	}

	return text, nil

}

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit { //se o numero for menor que 1000, escreve só com B
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp]) //dependendo do valor da divisao, pega um valor do slice/array
}
