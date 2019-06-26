package v3

import "github.com/gomarkdown/markdown/ast"

func isTableBodyCell(tableCell *ast.TableCell) bool {
	row, _ := tableCell.Parent.(*ast.TableRow)
	_, isGrandParentTableBody := row.Parent.(*ast.TableBody)
	return isGrandParentTableBody
}
