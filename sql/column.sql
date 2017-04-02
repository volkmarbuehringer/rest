


select
'type '||upper(substr(table_name,1,1))||substr(table_name,2)||' struct {'||
string_agg(
upper(substr(column_name,1,1))||
substr(column_name,2)||' '||
case data_type
when 'integer' then
case when is_nullable = 'YES' then
'*int'
else
'int'
end
when 'character varying' then
case when is_nullable = 'YES' then
'*string'
else
'string'
end
when 'character' then
case when is_nullable = 'YES' then
'*string'
else
'string'
end
when 'text' then
case when is_nullable = 'YES' then
'*string'
else
'string'
end

when 'numeric'  then
case
  when numeric_scale > 0 and is_nullable = 'YES' then '*float64'
  when numeric_scale > 0 and is_nullable = 'NO' then 'float64'
  when numeric_scale = 0 and is_nullable = 'NO' then 'int'
else
   '*string'
  end
when 'timestamp without time zone' then
case when is_nullable = 'YES' then '*string'
else
  'string'
  end
when 'date' then
  case when is_nullable = 'YES' then '*string'
  else
    'string'
    end
end
||' '||
'`json:"'||column_name||'" db:"'||column_name||'"`' ,CHR(10)
order by  ordinal_position
)||'}'||chr(10)||chr(10)
from information_schema.columns
where table_name in ( select table_name from information_schema.tables
  where table_schema = 'public' and table_type ='BASE TABLE')
group by table_name
;
