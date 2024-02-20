# Job Scheduler System Design

## Problem Statement
Design a job scheduler that runs jobs at scheduled interval one-off or regularly.
For this project scope, a Job is a CI/CD build and test project.

## Requirement
### Functional requirements

* User can submit a job or view a job
* User can see history of execution of all submitted jobs
* Jobs can have priority
* Job output shall be stored in the filestorage
* Job must be run at least once

### Non Functional requirement
* Scalable
* Highly available
* Reliability : job should run at least once and a job shall not be executed twice at the same time
* Durable : jobs information shall not be lost


## Design

The system is configured to perform task as early as possible, with a resolution of a minute.
The output of the execution are stored in cloud storage, S3.

### API routes
- /api/v1/login : POST, login is handled with auth0.
- /api/v1/submit_job : POST, data is a yaml file
- /api/v1/retrieve_history_execution/job_id : GET, retrieve json history of execution
- /api/v1/last_execution_logs/job_id
- /api/v1/job_status/job_id

### Database Design


### Notes
A job is a go package that is build and test, The client will submit the job via an API call with yaml spec (i.e url, task (build, test), format, ) and the output will be store in s3.

Only working on the backend. 
Going with these tech stack : 
Go(echo), Docker API, AWS RDS (postgres), DynamoDB, S3
For authentification, let's do auth0

