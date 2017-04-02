with x as (
  select * from pk_select
)
   select
   ' "'||table_name||'":  SQLStatements{Sqlstatement: `insert into '||table_name||'('||
   string_agg(
column_name  , ','
  order by  ordinal_position
)filter( where column_name not in (select column_name from x where table_name =  x.table_name ))||') values('||
  string_agg(
    case when column_name like '%cr_date' then
'now()'
    else
':'||column_name end , ','
 order by  ordinal_position
)filter( where column_name not in (select column_name from x where table_name =  x.table_name ))||') returning *`}, '
     from information_schema.columns x
   where table_name in ( select table_name from information_schema.tables
     where table_schema = 'public' and table_type ='BASE TABLE')
     and ordinal_position > 1
   group by table_name
   ;
