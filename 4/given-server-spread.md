```
---------------------------------------------------------------------------------------------------
     spreadsheet server                                      +------- accountant starts
			      ||                             |            | 
	     |                ||          spread client <----+            |
	     |                ||            |                             |
	     |<-----------------------------| tcp connect                 |
	     |                ||            |                             |
	     |<~~~~~~~~~~~~~~~~~~~~~~~~~~~~~| sign-up name                |
	     |~~~~~~~~~~~~~~~~~~~~~~~~~~~~~>| receive internal name	  | 
	     |                ||            |                             |
	     |                ||            |<----------------------------| named spreadsheet
	     |                ||            |<----------------------------| set requests 
	     |                ||            |<----------------------------| named spreadsheet
	     |                ||            |<----------------------------| set requests 
	     |                ||            |                             | ...
	     |                ||            |<----------------------------| at request 
	     |                ||            |                             |
	     |<~~~~~~~~~~~~~~~~~~~~~~~~~~~~~| send batch 		  |
	     |~~~~~~~~~~~~~~~~~~~~~~~~~~~~~>| receive at response	  | 
	     |                ||            |---------------------------->| print at response 
	     |                ||            |                             | 
	     |                ||            |<----------------------------| named spreadsheet
	     |                ||            |<----------------------------| set requests 
	     |                ||            |<----------------------------| named spreadsheet
	     |                ||            |<----------------------------| set requests 
	     |                ||            |                             | ...
	     |                ||            |<----------------------------| at request 
	     |<~~~~~~~~~~~~~~~~~~~~~~~~~~~~~| send batch 		  |
	     |                ||            |                             | 

```

The message formats are as follows: 
	     
|  message		  |  format						 |
| ----------------------- | ---------------------------------------------------- |
| sign-up name	    	  | string 						 |
| internal name	    	  | string (this distinguishes all connections)          |
| 	       	    	  |  			 	       			 |
| batch                   | [(named-spreadsheet, or set-request,)*, at-request]  |
| 			  |   meaning a possibly empty repetition of 		 |
|			  |   named spreadsheets and set requests 		 |
|			  |   followed by _one_ at request			 |
|			  |   	       	  	   				 |
| 	       	    	  |  			 	       			 |
| at request              | ["at", string, number, number]			 |
| 	       	    	  |  where string is the name of the probed spreadsheet  |
| 	       	    	  |        the first number is the x coordinate		 |
| 	       	    	  |        the second number is the y coordinate	 |
| 	       	    	  |  			 	       			 |
| at response  	    	  | number		 	       			 |


All JSON values are well-formed and valid. See 
[Assignment 3](http://www.ccis.northeastern.edu/home/matthias/4500-f18/3.html) 
task 3 for _well-formed_ and _valid_. 
