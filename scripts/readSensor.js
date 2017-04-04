var lastTrigger = new Date("1800-01-01 12:12:12");

function ReadSensor() {
			$(function() {
				$.getJSON(
					"http://yourlocalhost.info:5555/read",
					function( data ) {
						$.each( data, function( key, val ) {
							switch(key){
								case 'temperature':
									$("#temperature").hide().html(val.toFixed(2)).fadeIn('fast');
									break;
								case 'moisture':
									$("#moisture").hide().html(val.toFixed(2)).fadeIn('fast');
									break;
								case 'luminosity':
									$("#luminosity").hide().html(val.toFixed(2)).fadeIn('fast');
									break;
								case 'triggered':
									if(val == true) {
										var now = new Date();
										if ( Math.floor((now - lastTrigger)/(1000*60*60)) > 9 ) {
											$("#trigger").show();
											setTimeout(function() { $("#trigger").hide(); }, 5000);
											lastTrigger = now;
										}
									}
									break;
								default:
									$("#test").hide().html("IN SWITCH").fadeIn('fast');
							}
						});
					}
				);
			});
		}
		
		setInterval(ReadSensor, 500);
