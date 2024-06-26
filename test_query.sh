#/bin/bash!
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"question": "why is the sun black", "history": ""}' \
  https://localhost:8080/query

  curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"question": "what kind of gas", "history": "CUSTOMER: why is the sun black? CAT: Meow, meow, thats a very interesting question! The sun isnt actually black. Its a giant ball of burning gas, and it appears yellow because of all the light it gives off.  Purrrr, maybe youre thinking about something else?  Perhaps you saw a picture of a solar eclipse, where the moon blocks out the suns light. But the sun is still there, just hidden for a little while."}' \
  http://localhost:8080/query