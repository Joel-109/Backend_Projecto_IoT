-- name: GetDishes :many
SELECT *
FROM Dish;

-- name: GetInvoices :many
SELECT *
FROM Invoices I, Orders O
WHERE I.id_order = O.id_order;

-- name: GetDishesInvoice :many
SELECT D.name,D.image,D.description,D.price,D.id_category, O.id_order, O.id_desk, O.id_status
FROM Invoices I, Dish D, Orders O
WHERE I.id_dish = D.id_dish AND
I.id_order = ?;

-- name: GetTotal :one
SELECT SUM(D.price)
FROM Invoices I, Dish D
WHERE D.id_dish = I.id_dish
AND I.id_order = ?;

-- name: GestOrderNumber :one
SELECT id_order
FROM Orders
WHERE id_desk = ?
ORDER BY id_order DESC;

-- name: GetStatus :many
SELECT *
FROM Status;

-- name: InsertOrder :exec
INSERT INTO Orders(id_desk,id_status)
VALUES(?,?);

-- name: InsertDish :exec
INSERT INTO Dish (description,name,image,price,id_category)
VALUES (?,?,?,?,?);

-- name: InsertInvoice :exec
INSERT INTO Invoices(id_order,id_dish)
VALUES (?,?);

-- name: InsertStatus :exec
INSERT INTO Status (id_status)
VALUES (?);

-- name: InsertCategory :exec
INSERT INTO Category (id_category)
VALUES (?);

-- name: UpdateOrderStatus :exec
UPDATE orders
SET id_status = ?
WHERE id_order = ?;

-- name: DeleteDish :exec
DELETE FROM Dish
WHERE id_dish = ?;

-- name: DeleteInvoiceDish :exec
DELETE FROM Invoices 
WHERE rowid IN (
  SELECT rowid FROM Invoices I
  WHERE I.id_dish = ? AND I.id_order = ?
  LIMIT 1
);
