/*
 * This file was last modified at 2024-06-09 20:34 by Victor N. Skurikhin.
 * test.go
 * $Id$
 */

package pkg1

import (
	"fmt"
	"os"
)

func main() {
	fmt.Print("Hello world")
	os.Clearenv()
	// формулируем ожидания: анализатор должен находить ошибку,
	// описанную в комментарии want
	os.Exit(0) // want "call os.Exit error"
}
