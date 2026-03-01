CREATE TABLE IF NOT EXISTS join_employees (
id INT,
name TEXT,
department_id INT,
PRIMARY KEY (id)
);

INSERT INTO join_employees (id, name, department_id) VALUES
(1, 'Alice', 10),
(2, 'Bob', 20),
(3, 'Charlie', 10),
(4, 'David', 30);

CREATE TABLE IF NOT EXISTS join_departments (
id INT,
department_name TEXT,
PRIMARY KEY (id)
);

INSERT INTO join_departments (id, department_name) VALUES
(10, 'HR'),
(20, 'Engineering'),
(30, 'Sales');

CREATE TABLE IF NOT EXISTS join_salaries (
	employee_id INT,
	salary FLOAT,
	PRIMARY KEY (employee_id)
);

INSERT INTO join_salaries (employee_id, salary) VALUES
(1, 60000.0),
(2, 80000.0),
(3, 75000.0),
(4, 50000.0);

CREATE TABLE IF NOT EXISTS join_locations (
	department_id INT,
	location TEXT,
	PRIMARY KEY (department_id)
);

INSERT INTO join_locations (department_id, location) VALUES
(10, 'New York'),
(20, 'San Francisco'),
(30, 'Chicago');