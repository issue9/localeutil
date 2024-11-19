// SPDX-FileCopyrightText: 2024 caixw
//
// SPDX-License-Identifier: MIT

package locale

type X struct {
	F1 string `comment:"x_1"`
	F2 string `comment:"x_2"`
}

type Y struct {
	F1 string `comment:"y_1"`
	F2 *X     `comment:"y_2"`
}
