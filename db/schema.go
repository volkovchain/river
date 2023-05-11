package db

var Schema = `
CREATE TABLE IF NOT EXISTS employers (
                                         id INTEGER PRIMARY KEY AUTOINCREMENT,
                                         name TEXT NOT NULL,
                                         addr TEXT NOT NULL,
                                         amount_salary INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS payments (
                                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                                        salary_id INT,
                                        employee_id INT,
                                        amount INT,
                                        status VARCHAR(16),
                                        addr TEXT NOT NULL,
                                        error TEXT DEFAULT NULL,
                                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                        FOREIGN KEY(salary_id) REFERENCES salaries(id),
                                        FOREIGN KEY(employee_id) REFERENCES employers(id)
);

CREATE TABLE IF NOT EXISTS salaries (
                                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                                        status VARCHAR(16),
                                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)



`
