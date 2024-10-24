package service

import (
	"fmt"
	"mhxyHelper/internal/utils"
)

func buildVal(valMH, valRM float32) (float32, float32) {
	var err error
	if (valMH == 0 && valRM == 0) || (valMH != 0 && valRM != 0) {
		// log 无需转换
		return valMH, valRM
	}

	if valMH == 0 {
		valMH, err = utils.RM2MH(valRM)
		if err != nil {
			// TODO log
			fmt.Printf("RM2MH[ValRM: %f] err: %v\n", valRM, err)
			return valMH, valRM
		}
	}

	if valRM == 0 {
		valRM, err = utils.MH2RM(valMH)
		if err != nil {
			// TODO log
			fmt.Printf("MH2RM[ValMH: %f] err: %v\n", valMH, err)
			return valMH, valRM
		}
	}

	return valMH, valRM
}
