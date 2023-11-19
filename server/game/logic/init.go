package logic

import "sgserver/server/game/model/data"

func BeforeInit() {
	data.GetYield = RoleResService.GetYield

	data.GetUnion = RoleAttrService.GetUnion
}
