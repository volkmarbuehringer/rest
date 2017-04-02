select
'"'||table_name||'": SQLStatements{Sqlstatement: '||
'`select * from '||table_name||' order by '||
(select column_name from pk_select where table_name =  x.table_name )
||' limit :length offset :offset'||'`},'||chr(10)
  from information_schema.tables x
  where table_schema = 'public' and table_type ='BASE TABLE';
