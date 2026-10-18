Maybe 4 - 5 tables - these are my initial tables IDEA

#Table 1  - Workflow Table
id - Workflow-ID
name - human readable name
description - what is this workflow itended to do
version - Workflow - version(not needed it's a little advanced scenwerio)
created-at - When this workflow was created
updated-at - when this was last updated

#Table 2 - Activity Table
id - Activity-ID
name - human readable name
description - what is this actiivity meant to do - also for human to read
input_schema - What is the input structure of this Acitvity (Like a dto)
output_schema - Whats the Output Structrue of this activity (Response Structure)
created-at
updated-at

#Table 3 Execution ID
id - Execution id
workflow-id - (What workflow's execution is this)
status - Here the status means whether its in running its completed or it failed
created-at
updated-at
completed-at(??) - to know when this was completed or failed

#Table 4 - State log table - The Most important one (I dont have definitive structrue of what this table should look like but i will try to keep it closest to it vision/final-working should look like)
execution-id - Execution-id
step-no - What step was completed
acticity-id - what acticity it is associated with
event - (it will consist of action performed during the execution, here are some example - i. initialized a xyz activity(INIT-XYZ) ii. Completed a XYZ activity (COMPL-XYZ iii, Completed a workflow(FINISHED))
state - What state were we left at when the step was completed (The values of variable or output etc)
timestamp - when was this step completed(This will also be used in advanced scenrio - its here just for future reference)

table 4 will be indexed on (execution-id, step-no) for better retrieval

Table 4 Is tentative will work upon it to refine it to make it more tailored to the work as i go onn
