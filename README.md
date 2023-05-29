### Run

1. `docker build ./ -t genesis-education-task`
2. `docker run  -p 3000:80 genesis-education-task`
3. `curl http://localhost:3000/api/rate`

**Note**: SendEmails endpoint does not send emails since a SandGrid account was suspended because of API_KEY being exposed in this public repo.
