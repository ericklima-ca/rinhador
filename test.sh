# curl -X POST http://localhost:3000/payments \
# -H "Content-Type: application/json" \
# -d '{
#     "correlationId": "4a7901b8-7d26-4d9d-aa19-4dc1c7cf60b3",
#     "amount": 1234.56
# }'

# GET /payments-summary?from=2020-07-10T12:34:56.000Z&to=2020-07-10T12:35:56.000Z

# HTTP 200 - Ok
# {
#     "default" : {
#         "totalRequests": 43236,
#         "totalAmount": 415542345.98
#     },
#     "fallback" : {
#         "totalRequests": 423545,
#         "totalAmount": 329347.34
#     }
# }

curl "http://localhost:3000/payments-summary?"