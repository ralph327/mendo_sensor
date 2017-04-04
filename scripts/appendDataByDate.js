function underscore(id) {
	return id.replace(/ /g, "_");
}

function appendDataByDate(date, format) {		
	if(document.getElementById(underscore(date) + "_" + format) == null){
		var div = document.createElement('div');
			$.getJSON(
				"http://yourlocalhost.info:5555/appendDataByDate",
				{date: date, format: format},
				function( data ) {
					div.id = underscore(date) + "_" + format;
					div.className = format;
					div.innerHTML = data.html;					
				}
			);
		
		document.getElementById(underscore(date)).appendChild(div);
	}else{
		$("#" + underscore(date) + "_" + format).toggle();
	}
}

function getChildFormat(format) {
	var returnFormat = "dteMonthYear";
	$.ajax({
		url: "http://yourlocalhost.info:5555/getChildFormat",
		async: false,
		data: {format: format},
		success: function( data ) {
			returnFormat = data.format;				
		},
		dataType:"json"
	});
	
	return returnFormat;
}
