package userconfig

import (
	"fmt"

	"github.com/ItsHisoka17/Helix/internal/analyzer/userconfig/types"
)

func InitializeConfig(*types.DefaulConfig) {
	getUserInput()
}

func getUserInput() *types.DefaulConfig {

	fmt.Printf("UserConfig | Initializng Config...")
	return &types.DefaulConfig{}
}
