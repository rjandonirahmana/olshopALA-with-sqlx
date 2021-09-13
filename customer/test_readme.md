POST "/register"
fail used email 
{"meta":{"message":"failed to create account","code":200,"status":"Error 1054: Unknown column 'id' in 'field list'"},"data":null}

POST "/login"
success 
{"meta":{"message":"email not found","code":403,"status":"failed"},"data":{}}

--- FAIL: TestUpdatePhoneCustomer (0.01s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered]
	panic: runtime error: invalid memory address or nil pointer dereference
[signal 0xc0000005 code=0x0 addr=0x90 pc=0x106c0dd]