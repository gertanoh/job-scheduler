# Job Scheduler System Design

## Problem Statement
Design a job scheduler that runs jobs at scheduled interval one-off or regularly
Need to define what is a job : a job is a go repo and the tasks are build and test

## Requirement
### Functional requirements

* User can submit a job or view a job
* User can see history of execution of all submitted jobs
* Jobs can have priority
* Job output shall be stored in the filestorage

### Non Functional requirement
* Scalable
* Highly available
* Reliability : job should run at least once and a job shall not be executed twice at the same time
* Durable : jobs information shall not be lost
