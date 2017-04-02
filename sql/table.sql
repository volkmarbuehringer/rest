
select
'case "'||table_name||'":'||chr(10)||
  'm=[]interface{}{&[]dbstructs.'||upper(substr(table_name,1,1))||substr(table_name,2)||'{}}'||chr(10)
  ||'mm=[]interface{}{&dbstructs.'||upper(substr(table_name,1,1))||substr(table_name,2)||'{}}'||chr(10)
  ||chr(10)
  from information_schema.tables x
  where table_schema = 'public' and table_type ='BASE TABLE';

  create or replace view gugu as select los.*,loskunde.*,lt_ads.* ,1 as "$zahlung$zl_id", zahlung_los.* ,zahlung.*,1 as "$gewinn$g_id",gewinn.*,gewinnzahlung.* from los
  left outer join zahlung_los on zl_lid = l_id and zl_maxperiode > 500
  inner join loskunde on lk_lid = l_id
  inner join lt_ads on lk_k_adsid = ads_id
  left outer join gewinn on g_lid = l_Id
  left outer join gewinnzahlung on gz_gid = g_id
  left outer join zahlung on zl_zid = z_id

  create or replace view gugu as select los.*,1 as "$adresse$ads_id$lk_lid$l_id",ads_id, lk_lid,lt_ads.ads_name,ads_vorname ,1 as "$zahlung$zl_id$zl_lid$l_id", z.*,zahlung.*,1 as "$gewinn$g_id$g_zlid$zl_id",gewinn.*,gewinnzahlung.* from los
  left outer join zahlung_los z on zl_lid = l_id
  left outer join zahlung on z_id = zl_zid
    inner join loskunde on lk_lid = l_id
    inner join lt_ads on lk_k_adsid = ads_id
    left outer join gewinn on g_zlid = zl_Id
    left outer join gewinnzahlung on gz_gid = g_id



  create or replace view gugu2 as select /*+ first_rows */ los.*,1 as "$adresse$ads_id$lk_lid$l_id",ads_id, lk_lid,lt_ads.ads_name,ads_vorname ,1 as "$zahlung$zl_id$zl_lid$l_id", z.*,1 as "$gewinn$g_id$g_zlid$zl_id",gewinn.*,gewinnzahlung.* from los
  inner join zahlung_los z on zl_lid = l_id
    inner join loskunde on lk_lid = l_id
    inner join lt_ads on lk_k_adsid = ads_id
    left outer join gewinn on g_zlid = zl_Id
    left outer join gewinnzahlung on gz_gid = g_id
  where l_id in ( select l_id from los sample(0.1 ) where rownum < 10000 )

  create table tgugu2 as select * from gugu2
  alter table tgugu2 add primary key( l_id )

  select l_id from l
