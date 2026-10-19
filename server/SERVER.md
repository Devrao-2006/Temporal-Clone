# Whats Server Should do??  

- It Should register any new workflow, activity, listen to creation of a execution, listen to recieving state log  
- Simply to Give data to database  
- When in execution it should give the workflow data or activity data to the matchng server  

## Apis  - for user and worker (there is one outgoing request to matching server)

(PUT) workflow/create  
(PUT) workflow/update  
(GET) workflow  
(GET) worfklow/:id  

(PUT) activity/create  
(PUT) activity/update  
(GET) activity  
(GET) activity/:id  

(PUT) execution/create  
(PUT) execution/update //will be used by user as well as worker  
(GET) execution  
(GET) execution/:id  //to know the status as well  

### Outgoing Request (To Matching Server)

When an execution starts or when a step completes and triggers the next step, the Server dispatches a task to the Matching Server so workers can pick it up.

(POST) task/dispatch