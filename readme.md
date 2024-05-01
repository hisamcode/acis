
user
id | avatar | name | email | password

wts mean what the fuck transaction is
id | name | detail
1  | in   | kredit
2  | out  | debit

categories
id | name | unicode
1  | Gas  | ⛽

transactions
id | user_id  | nominal | wts_id | category_id | date | detail
1  |     1    | 20.000  | 1      | 1           |      | 

misal
setiap insert transaction, increment category[3]
projects
id, projects, categories text[][gas, ⛽, 4], version(optimistic)
transactions
id, project_id, 

misal setiap hari, jam 1 malam bakal nge calculate, sync atau biar konsisten sama table transaction
dan kalo misal mau perform query, pakai transaction(sql) kalo gagal,gagal sekalian kalo misalkan di database yg sama.
kalo misalkan untuk total transacation nya di database lain, pakai distributed transaction, tapi karena kayak nya ribet jadi mending di database yang sama aja.
I think this is good, but for each transaction idk, how ?????
for example range transaction in this month will slow, bcs I must querying on transaction. or I can partition every month
yeah I think if the app have slowdown, I can partition every month. 
but overall, on total we have fast querying
total_transaction
id | user_id | remaining_total | yearly_total | monthly_total | daily_total | out_monthly bla bla | in_monthly bla bla
id, project_id, yearly, year
id, project_id, monthly, month-year
id, project_id, daily, date

tbl_yearly
id, project_id, total, year
tbl_monthly
id, project_id, total, montly
tbl_daily
id, project_id, total, daily

transaction in 20.000
update daily +20.000
update montly +20.000
update yearly +20.000


query requirement:
what total out this month ?
what total out this day ?
what total out this year ?

what total in this month ?
what total in this day ?
what total in this year ?

what total this day ?
what total this month ?
what total this year ?

what remaining total I have, all of them ?

transaction per month, day, year?


```
curl -X POST localhost:8000/signup -H "Content-Type: application/x-www-form-urlencoded" -d "name=test&email=test1@test.com&password=123123123&repeat_password=123123123" && windows-kill -SIGINT 13716
```