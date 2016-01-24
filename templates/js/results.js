$(document).ready(function() {

  function getData(met, unmet) {
    $.getJSON("/api/calc?" + (new Date()).getTime(), {
      "met": met,
      "unmet": unmet
    },
    function(data) {
      console.log(data);
      writeData(data);
    });
  }

  function writeData(data) {

    dataspot = $("#dataspot");
    html = "";

    console.log(data);

    for (var i = 0; i < data.Groups.length; i++) {
      html += "<table class='table-hover class-list-table table'>";
      for (var cls = 0; cls < data.Groups[i].Classes.length; cls++) {
        obj = data.Groups[i].Classes[cls];

        html += "<td>" + obj["Course"] + "</td><td>" + obj["Title"] + "</td><td>" + obj["Credits"] + "</td></tr>\n";
      }
      html += "</table>";
    }

    dataspot.html(html);

    $("#data").css("visibility", "visible");
    $("#loading").css("visibility", "hidden");
    $("#loading").css("display","none");
  }

  met = $("#met").val()
  unmet = $("#unmet").val()

  console.log(met);
  console.log(unmet);

  getData(met, unmet);

});
