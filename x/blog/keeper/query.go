package keeper

import (
	"github.com/Akagi201/planet/x/blog/types"
)

var _ types.QueryServer = Keeper{}
