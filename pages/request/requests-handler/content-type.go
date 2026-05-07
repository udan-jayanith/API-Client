package requests_handler

type ContentType string

func (content_type ContentType) Parse() (t, sub_t string) {
	if content_type == "" {
		return
	}

	var i int
	for ; i < len(content_type) && content_type[i] != '/'; i++ {
	}

	if i == len(content_type) {
		t = string(content_type)
		return
	} else if content_type[i] == '/'{
		t = string(content_type[:i])
	}

	j := i
	for ; i < len(content_type) && content_type[i] != ';'; i++ {
	}

	if i == len(content_type){
		sub_t = string(content_type[j+1:i])
	}else if content_type[i] == ';'{
		sub_t = string(content_type[j+1:i])
	}
	
	return
}
