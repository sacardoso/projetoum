package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	args := []string{"."} //argumentos recebem uma lista de strings
	if len(os.Args) > 1 { //se o tamanho dos args for maior que um,
		args = os.Args[1:] //vamos usa-los
	}

	for _, arg := range args { //pssa por todos os args
		err := tree(arg, "")
		if err != nil {
			log.Printf("tree %s: %v\n", arg, err)
		}
	}

}

func tree(root, indent string) error { //criamos uma recursao para fazer a "arvore" - indent (indentation)
	fi, err := os.Stat(root) //retorna infor
	if err != nil {
		return fmt.Errorf("could not stat %s: %v", root, err)
	}

	fmt.Println(fi.Name(), "[", ByteCountSI(fi.Size()), "]") //"printa" o nome e o tamanho
	if !fi.IsDir() {                                         //se nao for um diretorio, nao tem mais nada o que fazer
		return nil
	}

	fis, err := ioutil.ReadDir(root) // ReadDir reads the directory named by dirname and returns a list of directory entries sorted by filename
	if err != nil {
		return fmt.Errorf("could not read dir %s: %v", root, err)
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
			fmt.Printf(indent + "└──") //printa a posicao, e para isso precisa saber o indent
			add = "   "                // se for o ultimo, nao vai printar │ , e sim "  "(espaco)
		} else {
			fmt.Printf(indent + "├──") //se for um diretorio tem q ter esse formato

		}

		if err := tree(filepath.Join(root, name), indent+add); err != nil { //indent+add para criar os "espacos"
			return err
		}
	}

	return nil

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
