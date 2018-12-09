"use strict";

document.addEventListener("DOMContentLoaded", function(event) {
    
    $("#process").click(function() {
        $("#processing").removeClass("hidden");
        $("#processed-image").addClass("hidden");

        $.ajax({
            url: "http://localhost:8080/process_image", 
            type: "GET",
            accepts: "json",
            success: function(xhr) {
                var imageURL = xhr["imageURL"];
                console.log(imageURL);
                document.getElementById("processed-image").src="http://localhost:8080/processed_images/" + imageURL;
            },
            error: function(xhr){
                console.log(xhr)
                alert("An error occured: " + xhr.status + " " + xhr.statusText);
            }
        });

        $("#processing").addClass("hidden");
        $("#processed-image").removeClass("hidden");
        console.log("processing complete");
    })


});

$.urlParam = function(name){
	var results = new RegExp('[\?&]' + name + '=([^&#]*)').exec(window.location.href);
	return results[1] || 0;
}