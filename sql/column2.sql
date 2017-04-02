with x as (
  select * from pk_select
)
   select
   ' "'||table_name||'":  SQLStatements{Sqlstatement: `update '||table_name||' set '||
   string_agg(
     case when column_name like '%upd_date' then
column_name||'=now()  '
     else
column_name||'=:'||column_name end , ','
  order by  ordinal_position
) filter( where column_name not in (select column_name from x where table_name =  x.table_name )
 and column_name not like '%cr_uid' and column_name not like '%cr_date' )||
string_agg(
  ' where '||column_name||' =:'||column_name||' returning *'
, ','
order by  ordinal_position
) filter ( where column_name in (select column_name from x where table_name =  x.table_name )  )
||' `}, '
     from information_schema.columns x
   where table_name in ( select table_name from information_schema.tables
     where table_schema = 'public' and table_type ='BASE TABLE')
   group by table_name
   ;
