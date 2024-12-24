# asdasd

## [postgis sql](https://postgis.net/workshops/postgis-intro/simple_sql.html)

* [PostGIS](https://postgis.net/docs/using_postgis_dbmanagement.html)
* [WIKI](https://en.wikipedia.org/wiki/Spatial_reference_system)
* [OSM](https://wiki.openstreetmap.org/wiki/Web_Mercator)
* [SRID](https://www.gaia-gis.it/fossil/libspatialite/wiki?name=GPX+tracks)

    `SRID 4326 (alias WGS 84)`
    `SRID 3857`

```sql

SELECT postgis_full_version();

CREATE TABLE geometries (name varchar, geom geometry);

INSERT INTO geometries VALUES
  ('Point', 'POINT(3123 2)'),
  ('Linestring', 'LINESTRING(0 0, 1 1, 2 1, 2 2)'),
  ('Polygon', 'POLYGON((0 0, 1 0, 1 1, 0 1, 0 0))'),
  ('PolygonWithHole', 'POLYGON((0 0, 10 0, 10 10, 0 10, 0 0),(1 1, 1 2, 2 2, 2 1, 1 1))'),
  ('Collection', 'GEOMETRYCOLLECTION(POINT(2 0),POLYGON((0 0, 1 0, 1 1, 0 1, 0 0)))');

-- // p1: 24.061874700710177, 42.0945152733475 
-- // p2: 24.061881490051746, 42.094441428780556 == 8.230228 ++ 8.230228

SELECT ST_DistanceSphere(
    ST_MakePoint(24.061874700710177, 42.0945152733475),
    ST_MakePoint(24.061881490051746, 42.094441428780556)
    );

-- // p1: 24.061878137290478, 42.09444251842797 
--    p2: 24.06187243759632, 42.094442350789905 
-- == 0.470658 ++ 12.386014
--    0.47065841
SELECT ST_DistanceSphere(
    ST_MakePoint(24.061878137290478, 42.09444251842797),
    ST_MakePoint(24.06187243759632, 42.094442350789905 )
    );

SELECT ST_DistanceSphere(
    ST_MakePoint(103.776047, 1.292149),
    ST_MakePoint(103.77607, 1.292212)
    );

SELECT ST_DistanceSphere(
    ST_MakePoint(103.776070, 1.292212),
    ST_MakePoint(103.77554,1.292406)
    );


```
