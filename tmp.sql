SELECT
	DISTINCT ctu.table_name AS referenced_table_name
  FROM 
	information_schema.TABLE_CONSTRAINTS AS tc
	JOIN information_schema.CONSTRAINT_TABLE_USAGE AS ctu
	  ON tc.constraint_name = ctu.constraint_name
  WHERE tc.constraint_type = 'FOREIGN KEY'
	AND tc.table_name = 'child_example_table_2'