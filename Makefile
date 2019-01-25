default:
	make -C "./Santorini"
	make -C "./Santorini" copy

linux:
	make -C "./Santorini" linux	
	make -C "./Santorini" copy