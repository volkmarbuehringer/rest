package db

import (
	"fmt"
	"strings"
)

const BindVar string = "$"
const EqBindVar string = "=$"
const dbTimestamp string = "now()"
const dbschema = "public"

func dbSequenzer(tab string) string {
	return ""
}

const sqlallcols string = `select column_name,case when data_type = 'numeric' then
case when numeric_scale > 0 then column_name||'::varchar' else
column_name||'::bigint' end else column_name end as column_type
from information_schema.columns
where table_name =$1 and table_schema = '` + dbschema + `' order by  ordinal_position `

const sqlLimit string = `select %s from ` + dbschema + `.%s where 1 = 1 %s order by %s limit $1 offset $2`
const sqlID string = `select %s from ` + dbschema + `.%s where %s %s order by %s`
const sqlInsert string = `insert into` + dbschema + `.%s(%s)values(%s) %s`
const sqlUpdate string = `update ` + dbschema + `.%s set %s where %s %s`

func (m *SQLStatement) ReturnClause() string {
	return fmt.Sprintf(" returning %s", strings.Join(m.ColumnsAllType, ","))
}

const sqlalltabs string = `SELECT
  tc.table_name,c.column_name
 FROM
 information_schema.table_constraints
 tc JOIN
 information_schema.constraint_column_usage
AS
 ccu USING
(constraint_schema,
 constraint_name)
JOIN
 information_schema.columns
AS
 c ON
 c.table_schema
= tc.constraint_schema
AND tc.table_name
= c.table_name
AND ccu.column_name
= c.column_name
where
 constraint_type =
'PRIMARY KEY'
 and tc.table_schema ='` + dbschema + `'
union all
select c.table_name,column_name
from information_schema.views v
inner join information_schema.columns c on v.table_name = c.table_name and ordinal_position = 1
where v.table_schema ='` + dbschema + `' and c.table_schema = '` + dbschema + `'`
