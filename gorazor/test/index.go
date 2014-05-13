package cases

import (
	"bytes"
	"gorazor"
	"kp/models"
	"tpl/admin/layout"
)

func Index(users []*models.User, total int, limit int, offset int) string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n\n<h2 class=\"sub-header\">用户总数：")
	_buffer.WriteString(gorazor.HTMLEscape(gorazor.Itoa(total)))
	_buffer.WriteString("</h2>\n<div class=\"table-responsive\">\n	<table class=\"table table-striped\">\n		<thead>\n			<tr>\n				<th>名字</th>\n				<th>电邮</th>\n				<th>编辑</th>\n			</tr>\n		</thead>\n		<tbody>\n			")
	for _, u := range users {

		_buffer.WriteString("<tr>\n				<td>")
		_buffer.WriteString(gorazor.HTMLEscape(u.Name))
		_buffer.WriteString("</td>\n				<td>")
		_buffer.WriteString(gorazor.HTMLEscape(u.Email))
		_buffer.WriteString("</td>\n				<td><a href=\"/admin/user/edit?id=")
		_buffer.WriteString(gorazor.HTMLEscape(u.ID.Hex()))
		_buffer.WriteString("\">编辑</a></td>\n			</tr>")

	}
	_buffer.WriteString("\n		</tbody>\n	</table>\n</div>\n\n")

	title := func() string {
		var _buffer bytes.Buffer

		_buffer.WriteString("用户管理")

		return _buffer.String()
	}
	_buffer.WriteString("\n")

	return layout.Base(_buffer.String(), title())
}
