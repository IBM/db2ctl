package main

import "github.com/IBM/db2ctl/cmd"

func main() {
	cmd.Execute()

	// construct `sleep.sh` command
	// cmd := &exec.Cmd{
	// 	Path:   "./script.sh",
	// 	Args:   []string{"./script.sh", "3"},
	// 	Stdout: os.Stdout,
	// 	Stderr: os.Stdout,
	// }

	// // run `cmd` in background
	// cmd.Start()

	// // do something else
	// // for i := 1; i < 300000; i++ {
	// // 	fmt.Println(i)
	// // }

	// // wait `cmd` until it finishes
	// cmd.Wait()

}
